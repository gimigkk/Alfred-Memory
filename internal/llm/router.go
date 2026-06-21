package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

type RouterClient struct {
	GeminiKey string
	GroqKey   string
	Models    []string
	Cooldowns map[string]time.Time
}

type FunctionDef struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Parameters  map[string]any `json:"parameters"`
}

type ToolDef struct {
	Type     string      `json:"type"`
	Function FunctionDef `json:"function"`
}

func NewRouterClient(geminiKey, groqKey string) *RouterClient {
	client := &RouterClient{
		GeminiKey: geminiKey,
		GroqKey:   groqKey,
		Cooldowns: make(map[string]time.Time),
	}
	client.fetchAvailableModels()
	return client
}

func (c *RouterClient) fetchAvailableModels() {
	c.Models = []string{}

	// Fetch Gemini Models
	geminiURL := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models?key=%s", c.GeminiKey)
	respG, err := http.Get(geminiURL)
	if err == nil {
		defer respG.Body.Close()
		var resG struct {
			Models []struct {
				Name string `json:"name"`
			} `json:"models"`
		}
		if json.NewDecoder(respG.Body).Decode(&resG) == nil {
			var gemini3 []string
			var others []string
			for _, m := range resG.Models {
				// gemini models return as "models/gemini-..."
				// we want to skip embeddings, tts
				name := m.Name
				if strings.Contains(name, "gemini") && !strings.Contains(name, "embedding") && !strings.Contains(name, "tts") && !strings.Contains(name, "vision") {
					short := strings.TrimPrefix(name, "models/")
					if strings.Contains(short, "gemini-3") {
						gemini3 = append(gemini3, "gemini/"+short)
					} else {
						others = append(others, "gemini/"+short)
					}
				}
			}
			c.Models = append(c.Models, gemini3...)
			c.Models = append(c.Models, others...)
		}
	}

	// Fetch Groq Models
	reqGroq, _ := http.NewRequest("GET", "https://api.groq.com/openai/v1/models", nil)
	reqGroq.Header.Set("Authorization", "Bearer "+c.GroqKey)
	client := &http.Client{}
	respGroq, err := client.Do(reqGroq)
	if err == nil {
		defer respGroq.Body.Close()
		var resGroq struct {
			Data []struct {
				ID string `json:"id"`
			} `json:"data"`
		}
		if json.NewDecoder(respGroq.Body).Decode(&resGroq) == nil {
			for _, m := range resGroq.Data {
				name := m.ID
				if !strings.Contains(name, "whisper") && !strings.Contains(name, "guard") && !strings.Contains(name, "safeguard") && !strings.Contains(name, "orpheus") {
					c.Models = append(c.Models, "groq/"+name)
				}
			}
		}
	}
	log.Printf("Dynamically loaded %d models for Agentic Loop.", len(c.Models))
}

type Message struct {
	Role    string
	Content string
}

func cleanJSON(s string) string {
	s = strings.TrimSpace(s)
	if strings.HasPrefix(s, "```json") {
		s = strings.TrimPrefix(s, "```json")
		s = strings.TrimSuffix(s, "```")
	} else if strings.HasPrefix(s, "```") {
		s = strings.TrimPrefix(s, "```")
		s = strings.TrimSuffix(s, "```")
	}
	return strings.TrimSpace(s)
}

func stripSpaces(s string) string {
	var sb strings.Builder
	for _, r := range s {
		if r != ' ' && r != '\t' && r != '\n' && r != '\r' {
			sb.WriteRune(r)
		}
	}
	return sb.String()
}


func (c *RouterClient) GenerateAgentic(systemPrompt string, userPrompt string, tools []ToolDef, executor func(name, args string) (string, error)) (string, error) {
	toolsJSON, _ := json.MarshalIndent(tools, "", "  ")
	fullSystemPrompt := systemPrompt + "\n\nYou MUST respond ONLY with a JSON object. You must ALWAYS include your reasoning in a 'thought' field. To call a tool, return:\n{\"thought\": \"your reasoning...\", \"tool_name\": \"...\", \"arguments\": {...}}\n\nTools available:\n" + string(toolsJSON)

	history := []Message{
		{Role: "user", Content: userPrompt},
	}

	var lastCommitArgs string
	var lastCommitErr string

	for step := 0; step < 30; step++ {
		var content string
		var success bool
		var lastErr string

		// Use the dynamically loaded models
		for _, modelRef := range c.Models {
			if time.Now().Before(c.Cooldowns[modelRef]) {
				continue
			}

			parts := strings.SplitN(modelRef, "/", 2)
			provider := parts[0]
			model := parts[1]

			var err error
			if provider == "gemini" {
				content, err = c.callGemini(model, fullSystemPrompt, history)
			} else {
				content, err = c.callGroq(model, fullSystemPrompt, history)
			}

			if err != nil {
				if strings.Contains(err.Error(), "429") || strings.Contains(err.Error(), "500") || strings.Contains(err.Error(), "502") || strings.Contains(err.Error(), "503") || strings.Contains(err.Error(), "504") {
					log.Printf("Model %s hit rate limit or server error. Placed on 1 minute cooldown. Trying next...", modelRef)
					c.Cooldowns[modelRef] = time.Now().Add(1 * time.Minute)
					lastErr = err.Error()
					continue
				}
				log.Printf("Model %s failed: %s. Trying next...", modelRef, err.Error())
				lastErr = err.Error()
				continue
			}

			success = true
			break
		}

		if !success {
			return "", fmt.Errorf("all models failed. last error: %s", lastErr)
		}

		history = append(history, Message{Role: "assistant", Content: content})

		cleanContent := cleanJSON(content)
		var parseAttempt struct {
			Thought   string                 `json:"thought"`
			ToolName  string                 `json:"tool_name"`
			Arguments map[string]interface{} `json:"arguments"`
			Mutations []interface{}          `json:"mutations"`
		}

		if err := json.Unmarshal([]byte(cleanContent), &parseAttempt); err != nil {
			history = append(history, Message{Role: "user", Content: "Error: You must output a valid JSON object. Do not wrap in markdown blocks if unsupported."})
			continue
		}

		if parseAttempt.Thought != "" {
			log.Printf("\n\033[90m[AGENT THOUGHT]\033[0m %s\n", parseAttempt.Thought)
		}

		toolName := parseAttempt.ToolName
		if toolName == "" && len(parseAttempt.Mutations) > 0 {
			toolName = "commit_mutations"
		}

		if toolName != "" {
			argsJSON, _ := json.Marshal(parseAttempt.Arguments)
			argsStr := string(argsJSON)
			toolResult, err := executor(toolName, argsStr)
			if err != nil {
				errMsg := err.Error()
				toolResult = fmt.Errorf("tool error: %v", err).Error()
				if toolName == "commit_mutations" {
					normCurrent := stripSpaces(argsStr)
					normLast := stripSpaces(lastCommitArgs)
					if normCurrent == normLast && lastCommitErr == errMsg {
						toolResult += "\n\nWARNING: Your last two attempts were nearly identical and both failed for the same reason. Re-read the error carefully — repeating the same mutation will not work. If you believe your mutation is correct, the validator itself may have a bug; in that case, try an alternative representation (e.g., a different but still rule-compliant edge) rather than repeating the exact same payload."
					}
					lastCommitArgs = argsStr
					lastCommitErr = errMsg
				}
			} else {
				if toolName == "commit_mutations" {
					return argsStr, nil
				}
			}
			history = append(history, Message{Role: "user", Content: "Tool Result:\n" + toolResult})
		} else {
			history = append(history, Message{Role: "user", Content: "Error: No tool_name specified."})
		}
	}

	return "", fmt.Errorf("agentic loop reached max steps")
}

func (c *RouterClient) callGemini(model, systemPrompt string, history []Message) (string, error) {
	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent?key=%s", model, c.GeminiKey)

	type Part struct {
		Text string `json:"text"`
	}
	type Content struct {
		Role  string `json:"role"`
		Parts []Part `json:"parts"`
	}

	var contents []Content
	for _, m := range history {
		role := "user"
		if m.Role == "assistant" {
			role = "model"
		}
		contents = append(contents, Content{
			Role:  role,
			Parts: []Part{{Text: m.Content}},
		})
	}

	reqBody := map[string]interface{}{
		"systemInstruction": map[string]interface{}{
			"parts": []map[string]interface{}{
				{"text": systemPrompt},
			},
		},
		"contents": contents,
		"generationConfig": map[string]interface{}{
			"temperature":      0.1,
			"responseMimeType": "application/json",
		},
	}

	jsonData, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var result struct {
		Candidates []struct {
			Content struct {
				Parts []struct {
					Text string `json:"text"`
				} `json:"parts"`
			} `json:"content"`
		} `json:"candidates"`
	}

	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return "", err
	}

	if len(result.Candidates) == 0 || len(result.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("no candidates returned from gemini")
	}

	return result.Candidates[0].Content.Parts[0].Text, nil
}

func (c *RouterClient) callGroq(model, systemPrompt string, history []Message) (string, error) {
	url := "https://api.groq.com/openai/v1/chat/completions"

	messages := []map[string]string{
		{"role": "system", "content": systemPrompt},
	}
	for _, m := range history {
		messages = append(messages, map[string]string{
			"role":    m.Role,
			"content": m.Content,
		})
	}

	reqBody := map[string]interface{}{
		"model":       model,
		"messages":    messages,
		"temperature": 0.1,
	}
	if !strings.HasPrefix(model, "openai/") {
		reqBody["response_format"] = map[string]string{"type": "json_object"}
	}

	jsonData, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	req.Header.Set("Authorization", "Bearer "+c.GroqKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return "", err
	}

	if len(result.Choices) == 0 {
		return "", fmt.Errorf("no choices returned from groq")
	}

	return result.Choices[0].Message.Content, nil
}
