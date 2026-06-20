package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type GroqClient struct {
	APIKey string
}

func NewGroqClient(apiKey string) *GroqClient {
	return &GroqClient{APIKey: apiKey}
}

// GenerateJSON sends a prompt to Groq forcing a JSON object response
func (c *GroqClient) GenerateJSON(systemPrompt, userPrompt string) (string, error) {
	url := "https://api.groq.com/openai/v1/chat/completions"

	reqBody := map[string]interface{}{
		"model":           "llama-3.3-70b-versatile",
		"response_format": map[string]string{"type": "json_object"},
		"messages": []map[string]string{
			{"role": "system", "content": systemPrompt},
			{"role": "user", "content": userPrompt},
		},
		"temperature": 0.0,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+c.APIKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("groq api error (status %d): %s", resp.StatusCode, string(bodyBytes))
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

// Tool definitions for OpenAI format
type ToolDef struct {
	Type     string      `json:"type"`
	Function FunctionDef `json:"function"`
}

type FunctionDef struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Parameters  map[string]any `json:"parameters"`
}

type ToolCall struct {
	ID       string `json:"id"`
	Type     string `json:"type"`
	Function struct {
		Name      string `json:"name"`
		Arguments string `json:"arguments"`
	} `json:"function"`
}

type Message struct {
	Role       string     `json:"role"`
	Content    string     `json:"content"`
	Name       string     `json:"name,omitempty"`
	ToolCalls  []ToolCall `json:"tool_calls,omitempty"`
	ToolCallID string     `json:"tool_call_id,omitempty"`
}

type AgentResult struct {
	Mutations []map[string]any
}

// GenerateAgentic runs a ReAct loop using pure JSON (bypassing native tool bugs)
func (c *GroqClient) GenerateAgentic(systemPrompt string, userPrompt string, tools []ToolDef, executor func(name, args string) (string, error)) (string, error) {
	url := "https://api.groq.com/openai/v1/chat/completions"

	// Inject tool schemas into the system prompt
	toolsJSON, _ := json.MarshalIndent(tools, "", "  ")
	fullSystemPrompt := systemPrompt + "\n\nYou MUST respond ONLY with a JSON object. You must ALWAYS include your reasoning in a 'thought' field. To call a tool, return:\n{\"thought\": \"your reasoning...\", \"tool_name\": \"...\", \"arguments\": {...}}\n\nTools available:\n" + string(toolsJSON)

	messages := []map[string]string{
		{"role": "system", "content": fullSystemPrompt},
		{"role": "user", "content": userPrompt},
	}

	for step := 0; step < 10; step++ {
		var resp *http.Response
		var bodyBytes []byte
		var success bool
		var lastErr string

		models := []string{
			"llama-3.3-70b-versatile",
			"llama-3.1-8b-instant",
			"llama3-8b-8192",
		}

		for _, model := range models {
			reqBody := map[string]interface{}{
				"model":           model,
				"messages":        messages,
				"response_format": map[string]string{"type": "json_object"},
				"temperature":     0.1,
			}

			jsonData, _ := json.Marshal(reqBody)
			req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
			req.Header.Set("Authorization", "Bearer "+c.APIKey)
			req.Header.Set("Content-Type", "application/json")

			client := &http.Client{}
			var err error
			resp, err = client.Do(req)
			if err != nil {
				lastErr = err.Error()
				continue
			}

			bodyBytes, _ = io.ReadAll(resp.Body)
			resp.Body.Close()

			if resp.StatusCode == http.StatusOK {
				success = true
				break
			}
			
			lastErr = fmt.Sprintf("status %d: %s", resp.StatusCode, string(bodyBytes))

			if resp.StatusCode == http.StatusTooManyRequests {
				log.Printf("Rate limit hit for %s, falling back to next model...", model)
				continue
			}

			// If it's a 400 Bad Request (like invalid JSON prompt), failing over won't help
			break
		}

		if !success {
			return "", fmt.Errorf("groq api error (all models failed): %s", lastErr)
		}

		var result struct {
			Choices []struct {
				Message struct {
					Role    string `json:"role"`
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

		content := result.Choices[0].Message.Content
		messages = append(messages, map[string]string{"role": "assistant", "content": content})

		var parseAttempt struct {
			Thought   string                 `json:"thought"`
			ToolName  string                 `json:"tool_name"`
			Arguments map[string]interface{} `json:"arguments"`
			Mutations []interface{}          `json:"mutations"`
		}

		if err := json.Unmarshal([]byte(content), &parseAttempt); err != nil {
			messages = append(messages, map[string]string{"role": "user", "content": "Error: You must output a valid JSON object."})
			continue
		}

		// Print the thought for terminal visibility
		if parseAttempt.Thought != "" {
			log.Printf("\n\033[90m[AGENT THOUGHT]\033[0m %s\n", parseAttempt.Thought)
		}

		// Support both standard tool_name format or direct mutations output
		toolName := parseAttempt.ToolName
		if toolName == "" && len(parseAttempt.Mutations) > 0 {
			toolName = "commit_mutations"
		}

		if toolName == "commit_mutations" {
			argsJSON, _ := json.Marshal(parseAttempt.Arguments)
			// Return just the arguments JSON string so the caller can unmarshal it into LinkingOutput
			return string(argsJSON), nil
		} else if toolName != "" {
			// Execute tool
			argsJSON, _ := json.Marshal(parseAttempt.Arguments)
			toolResult, err := executor(toolName, string(argsJSON))
			if err != nil {
				toolResult = fmt.Errorf("tool error: %v", err).Error()
			}
			messages = append(messages, map[string]string{"role": "user", "content": "Tool Result:\n" + toolResult})
		} else {
			messages = append(messages, map[string]string{"role": "user", "content": "Error: No tool_name specified."})
		}
	}

	return "", fmt.Errorf("agentic loop reached max steps")
}
