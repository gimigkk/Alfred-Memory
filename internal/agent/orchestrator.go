package agent

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/gimigkk/Alfred-Memory/assets/prompts"
	"github.com/gimigkk/Alfred-Memory/internal/embed"
	"github.com/gimigkk/Alfred-Memory/internal/ladybug"
	"github.com/gimigkk/Alfred-Memory/internal/llm"
)

type Orchestrator struct {
	LLM    *llm.RouterClient
	Embed  *embed.GeminiClient
	DBConn *ladybug.Connection
}

func NewOrchestrator(llm *llm.RouterClient, embed *embed.GeminiClient, dbConn *ladybug.Connection) *Orchestrator {
	return &Orchestrator{
		LLM:    llm,
		Embed:  embed,
		DBConn: dbConn,
	}
}

func (o *Orchestrator) RunAgenticIngestion(runID string, transcript string, dryRun bool) (*LinkingOutput, error) {
	log.Printf("\n\033[36mStarting Agentic Ingestion for block: %s\033[0m", runID)

	// ==========================================
	// PHASE 1: INITIALIZATION
	// ==========================================
	state := newIngestionState()
	tools := GetIngestionTools()

	// ==========================================
	// PHASE 2: TOOL DELEGATION
	// ==========================================
	executor := func(name, args string) (string, error) {
		if name == "query_rag" {
			return o.handleQueryRag(args, state)
		} else if name == "extract_transcript_manifest" {
			return o.handleExtractManifest(args, transcript, state)
		} else if name == "commit_mutations" {
			return o.handleCommitMutations(args, state)
		}
		return "", fmt.Errorf("unknown tool: %s", name)
	}

	log.Println("\n\033[33m[AGENT] Starting Investigation Loop...\033[0m")
	systemPrompt := prompts.BuildDiscoveryPrompt()

	log.Printf("\033[90m--- SYSTEM PROMPT ---\n%s\n---------------------\033[0m\n", systemPrompt)
	log.Printf("\033[36mInitiating Agentic Ingestion Loop...\033[0m")

	// ==========================================
	// PHASE 3: STATE INTERCEPTOR
	// ==========================================
	interceptor := func(history *[]llm.Message, lastThought string) {
		if strings.Contains(lastThought, "[REQUEST_SCHEMA]") {
			if !state.HasQueriedVault {
				*history = append(*history, llm.Message{Role: "user", Content: "ERROR: You cannot request the schema yet. You must use the query_rag tool to verify the extracted entities against the vault first."})
			} else if !state.SchemaInjected {
				state.SchemaInjected = true
				newHistory := make([]llm.Message, 0, len(*history))
				for _, m := range *history {
					if !strings.HasPrefix(m.Content, "[SYSTEM_INJECTION_SKILL_COMMIT]") {
						newHistory = append(newHistory, m)
					}
				}
				*history = newHistory

				injectionContent := "[SYSTEM_INJECTION_SKILL_COMMIT]\nYou have completed the discovery phase. You must now apply the following Schema Constraints to commit your findings:\n\n" + prompts.BuildCommitPrompt()
				log.Printf("\n\033[90m--- SYSTEM INJECTION (SKILL COMMIT) ---\n%s\n---------------------------------------\033[0m\n", injectionContent)

				*history = append(*history, llm.Message{
					Role:    "user",
					Content: injectionContent,
				})
			}
		}
	}

	// ==========================================
	// PHASE 4: AGENT EXECUTION
	// ==========================================
	mutationsJSON, err := o.LLM.GenerateAgentic(systemPrompt, transcript, tools, executor, interceptor)
	if err != nil {
		return nil, fmt.Errorf("agent loop failed: %w", err)
	}

	log.Printf("\033[90m--- AGENT FINAL COMMIT ---\n%s\n------------------\033[0m\n", mutationsJSON)

	// ==========================================
	// PHASE 5: POST-PROCESSING
	// ==========================================

	// Parse the final mutations JSON
	// Since commit_mutations takes an argument with {"mutations": [...]}, we parse it like LinkingOutput
	var linkOut LinkingOutput
	if err := json.Unmarshal([]byte(mutationsJSON), &linkOut); err != nil {
		return nil, fmt.Errorf("failed to parse final mutations: %w", err)
	}

	log.Printf("\033[32m✔ Agent generated %d mutations.\033[0m Executing against DB...\n\n", len(linkOut.Mutations))

	o.remapTempIDs(linkOut.Mutations)
	o.executeAndVerifyEdges(linkOut.Mutations, transcript, state, dryRun)

	linkOut.ManifestAccounting = buildManifestAccounting(state.ExtractedManifestLines, linkOut.Mutations)
	log.Printf("\033[32m✔ Manifest accounting built for %d lines.\033[0m", len(linkOut.ManifestAccounting))

	return &linkOut, nil
}
