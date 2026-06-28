package agent

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gimigkk/Alfred-Memory/assets/prompts"
	"github.com/gimigkk/Alfred-Memory/internal/llm"
)

type ChatEvent struct {
	Type    string      `json:"type"` // "thought", "tool_call", "tool_result", "message", "error"
	Content string      `json:"content,omitempty"`
	Tool    string      `json:"tool,omitempty"`
	Args    interface{} `json:"args,omitempty"`
	Result  interface{} `json:"result,omitempty"`
}

func (o *Orchestrator) RunChatAgent(message string, history []llm.Message, emitEvent func(ChatEvent)) {
	log.Printf("\n\033[36mStarting Chat Agent\033[0m")
	log.Printf("\033[1;34m[USER]\033[0m %s\n", message)

	currentTime := time.Now().In(time.FixedZone("WIB", 7*3600)).Format("Monday, 02 Jan 2006 15:04:05 WIB")
	systemPrompt := prompts.BuildChatPrompt(currentTime, o.OwnerID)
	
	// Create chat tools
	tools := GetChatTools()

	log.Printf("\033[90m--- [SKILL CHAT] injected ---\033[0m\n")

	// Wrap emitEvent to also log locally
	emit := func(e ChatEvent) {
		switch e.Type {
		case "tool_call":
			argsJSON, _ := json.MarshalIndent(e.Args, "", "  ")
			log.Printf("\033[33m[CHAT] 🛠️ Calling %s with args:\n%s\033[0m", e.Tool, string(argsJSON))
		case "thought":
			log.Printf("\033[90m[CHAT] 🧠 %s\033[0m", e.Content)
		case "error":
			log.Printf("\033[31m[CHAT ERROR] %s\033[0m", e.Content)
		}
		emitEvent(e)
	}

	executor := func(name, args string) (string, error) {
		// Parse args for UI
		var parsedArgs interface{}
		json.Unmarshal([]byte(args), &parsedArgs)
		
		emit(ChatEvent{
			Type: "tool_call",
			Tool: name,
			Args: parsedArgs,
		})

		var result string
		var err error

		switch name {
		case "query_rag":
			// We can reuse a modified version of handleQueryRag, but we need to pass a state or modify it to not need ingestion state
			// For chat, we can just call it directly.
			result, err = o.handleChatQueryRag(args)
		case "query_node_history":
			result, err = o.handleQueryNodeHistory(args)
		case "ask_user_for_hint":
			result, err = o.handleAskUserForHint(args)
		case "commit_chat_mutations":
			result, err = o.handleCommitChatMutations(args)
		case "upsert_reminder":
			result, err = o.handleUpsertReminder(args)
		case "check_reminders":
			result, err = o.handleCheckReminders(args)
		default:
			err = fmt.Errorf("unknown tool: %s", name)
		}

		if err != nil {
			emit(ChatEvent{
				Type:    "error",
				Content: fmt.Sprintf("Tool %s failed: %v", name, err),
			})
			return "", err
		}

		// For UI, we want to send the result, but maybe truncated if it's too big
		var parsedResult interface{}
		if json.Unmarshal([]byte(result), &parsedResult) != nil {
			parsedResult = result
		}
		
		emit(ChatEvent{
			Type:   "tool_result",
			Tool:   name,
			Result: parsedResult,
		})

		// If ask_user_for_hint is called, we return a special error to break the loop
		if name == "ask_user_for_hint" {
			if argsMap, ok := parsedArgs.(map[string]interface{}); ok {
				if q, ok := argsMap["question"].(string); ok {
					emit(ChatEvent{
						Type:    "message",
						Content: q,
					})
				}
			}
			return result, fmt.Errorf("YIELD_TO_USER")
		}

		return result, nil
	}

	err := o.runChatLoop(systemPrompt, message, history, tools, executor, emit)
	if err != nil && err.Error() != "YIELD_TO_USER" {
		emit(ChatEvent{
			Type:    "error",
			Content: fmt.Sprintf("Agent loop failed: %v", err),
		})
	}
}

func (o *Orchestrator) runChatLoop(systemPrompt, message string, history []llm.Message, tools []llm.ToolDef, executor func(string, string) (string, error), emit func(ChatEvent)) error {
	fullHistory := append([]llm.Message(nil), history...)
	if len(fullHistory) == 0 || fullHistory[len(fullHistory)-1].Content != message {
		fullHistory = append(fullHistory, llm.Message{Role: "user", Content: message})
	}
	
	toolsJSON, _ := json.MarshalIndent(tools, "", "  ")
	fullSystemPrompt := systemPrompt + "\n\nYou MUST respond ONLY with a JSON object. You must ALWAYS include your reasoning in a 'thought' field. To call a tool, return:\n{\"thought\": \"your reasoning...\", \"tool_name\": \"...\", \"arguments\": {...}}\nTo reply to the user with a final answer, return:\n{\"thought\": \"...\", \"final_answer\": \"...\"}\n\nTools available:\n" + string(toolsJSON)

	var lastToolName string
	var lastToolArgs string
	var sameToolCount int

	for step := 0; step < 15; step++ {
		// Use the primary model for chat, e.g. the first one in the list
		var content string
		var err error
		var success bool
		var lastErr string
		for _, modelRef := range o.LLM.Models {
			if time.Now().Before(o.LLM.Cooldowns[modelRef]) {
				continue
			}

			parts := strings.SplitN(modelRef, "/", 2)
			provider := parts[0]
			model := parts[1]

			if provider == "gemini" {
				content, err = o.LLM.CallGeminiRaw(model, fullSystemPrompt, fullHistory)
			} else {
				content, err = o.LLM.CallGroqRaw(model, fullSystemPrompt, fullHistory)
			}

			if err != nil {
				if strings.Contains(err.Error(), "429") || strings.Contains(err.Error(), "500") || strings.Contains(err.Error(), "502") || strings.Contains(err.Error(), "503") || strings.Contains(err.Error(), "504") {
					log.Printf("\033[31m[API FALLBACK] Model %s hit rate limit/error. Placed on 1 min cooldown. Rotating...\033[0m", modelRef)
					o.LLM.Cooldowns[modelRef] = time.Now().Add(1 * time.Minute)
					lastErr = err.Error()
					continue
				}
				log.Printf("Model %s failed in chat: %v", modelRef, err)
				lastErr = err.Error()
				continue
			}

			success = true
			break
		}

		if !success {
			log.Printf("All models failed or on cooldown. Sleeping 15s and retrying step... (last error: %s)", lastErr)
			time.Sleep(15 * time.Second)
			step--
			continue
		}

		fullHistory = append(fullHistory, llm.Message{Role: "assistant", Content: content})

		cleanContent := cleanJSON(content)
		var parseAttempt struct {
			Thought     string                 `json:"thought"`
			ToolName    string                 `json:"tool_name"`
			Arguments   map[string]interface{} `json:"arguments"`
			FinalAnswer string                 `json:"final_answer"`
			Answer      string                 `json:"answer"`
			Response    string                 `json:"response"`
			Message     string                 `json:"message"`
		}

		if err := json.Unmarshal([]byte(cleanContent), &parseAttempt); err != nil {
			fullHistory = append(fullHistory, llm.Message{Role: "user", Content: "Error: You must output a valid JSON object."})
			continue
		}

		tName := strings.ToLower(strings.TrimSpace(parseAttempt.ToolName))
		if tName == "none" || tName == "null" || tName == "false" || tName == "n/a" || tName == "no_tool" {
			parseAttempt.ToolName = ""
		}

		finalAns := parseAttempt.FinalAnswer
		if finalAns == "" {
			finalAns = parseAttempt.Answer
		}
		if finalAns == "" {
			finalAns = parseAttempt.Response
		}
		if finalAns == "" {
			finalAns = parseAttempt.Message
		}
		if finalAns == "" && parseAttempt.ToolName == "" && parseAttempt.Thought != "" {
			// Fallback: If the model stubbornly refuses to output final_answer or a tool, yield its thought to break the loop.
			finalAns = parseAttempt.Thought
		}

		if parseAttempt.Thought != "" && finalAns != parseAttempt.Thought {
			emit(ChatEvent{
				Type:    "thought",
				Content: parseAttempt.Thought,
			})
		} else if parseAttempt.ToolName != "" {
			emit(ChatEvent{
				Type:    "thought",
				Content: "Decided to invoke " + parseAttempt.ToolName,
			})
		} else if finalAns != "" {
			emit(ChatEvent{
				Type:    "thought",
				Content: "Preparing final response...",
			})
		}

		if finalAns != "" {
			emit(ChatEvent{
				Type:    "message",
				Content: finalAns,
			})
			return nil
		}

		if parseAttempt.ToolName != "" {
			argsJSON, _ := json.Marshal(parseAttempt.Arguments)
			argsStr := string(argsJSON)

			if parseAttempt.ToolName == lastToolName && argsStr == lastToolArgs {
				sameToolCount++
				if sameToolCount >= 3 {
					emit(ChatEvent{Type: "error", Content: "Agent got stuck in a repetitive loop. Halting."})
					return fmt.Errorf("agent repetitive loop")
				}
			} else {
				lastToolName = parseAttempt.ToolName
				lastToolArgs = argsStr
				sameToolCount = 0
			}

			toolResult, err := executor(parseAttempt.ToolName, argsStr)
			if err != nil {
				if err.Error() == "YIELD_TO_USER" {
					return err
				}
				toolResult = fmt.Sprintf("Error: %v", err)
			}
			
			if sameToolCount > 0 {
				toolResult += "\n\nWARNING: You just repeated the exact same tool call and arguments as your previous step. Repeating identical actions will not work. You must change your approach or provide a final_answer."
			}

			fullHistory = append(fullHistory, llm.Message{Role: "user", Content: "Tool Result:\n" + toolResult})
		} else {
			fullHistory = append(fullHistory, llm.Message{Role: "user", Content: "Error: No tool_name or final_answer specified. You MUST provide either 'tool_name' to invoke a tool, or 'final_answer' to respond to the user. Do not just provide a thought."})
		}
	}
	
	return fmt.Errorf("max steps reached")
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
