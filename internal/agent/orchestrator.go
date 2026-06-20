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
	"github.com/gimigkk/Alfred-Memory/internal/rag"
	"github.com/gimigkk/Alfred-Memory/internal/waha"
)

type EdgeMutation struct {
	RelType      string `json:"rel_type"`
	TargetNodeID string `json:"target_node_id"`
	Context      string `json:"context,omitempty"`
	Role         string `json:"role,omitempty"`
	Descriptor   string `json:"descriptor,omitempty"`
}

type Mutation struct {
	Operation  string                 `json:"operation"` // CREATE_NODE or UPDATE_NODE
	NodeType   string                 `json:"node_type,omitempty"`
	NodeID     string                 `json:"node_id,omitempty"`
	Properties map[string]interface{} `json:"properties"`
	AddEdges   []EdgeMutation         `json:"add_edges,omitempty"`
}

type LinkingOutput struct {
	Mutations []Mutation `json:"mutations"`
}

type Orchestrator struct {
	LLM    *llm.GroqClient
	Embed  *embed.GeminiClient
	DBConn *ladybug.Connection
}

func NewOrchestrator(llm *llm.GroqClient, embed *embed.GeminiClient, dbConn *ladybug.Connection) *Orchestrator {
	return &Orchestrator{
		LLM:    llm,
		Embed:  embed,
		DBConn: dbConn,
	}
}

func (o *Orchestrator) RunAgenticIngestion(block *waha.ConversationBlock) error {
	log.Printf("\n\033[36mStarting Agentic Ingestion for block: %s\033[0m", block.ID)
	transcript := block.FormatTranscript()

	tools := []llm.ToolDef{
		{
			Type: "function",
			Function: llm.FunctionDef{
				Name:        "query_rag",
				Description: "Search the knowledge vault for relevant nodes using semantic search.",
				Parameters: map[string]any{
					"type":       "object",
					"properties": map[string]any{
						"query": map[string]any{
							"type":        "string",
							"description": "The search query (e.g. 'Bahlil' or 'Friday event').",
						},
					},
					"required": []string{"query"},
				},
			},
		},
		{
			Type: "function",
			Function: llm.FunctionDef{
				Name:        "commit_mutations",
				Description: "Commit the final graph mutations to the vault once all entities are resolved.",
				Parameters: map[string]any{
					"type":       "object",
					"properties": map[string]any{
						"mutations": map[string]any{
							"type":  "array",
							"items": map[string]any{
								"type":       "object",
								"properties": map[string]any{
									"operation":  map[string]any{"type": "string", "enum": []string{"CREATE_NODE", "UPDATE_NODE"}},
									"node_type":  map[string]any{"type": "string"},
									"node_id":    map[string]any{"type": "string"},
									"properties": map[string]any{"type": "object"},
									"add_edges":  map[string]any{
										"type":  "array",
										"items": map[string]any{
											"type":       "object",
											"properties": map[string]any{
												"rel_type":       map[string]any{"type": "string"},
												"target_node_id": map[string]any{"type": "string"},
												"context":        map[string]any{"type": "string"},
											},
										},
									},
								},
								"required": []string{"operation", "node_id"},
							},
						},
					},
					"required": []string{"mutations"},
				},
			},
		},
	}

	executor := func(name, args string) (string, error) {
		if name == "query_rag" {
			var parsed map[string]string
			if err := json.Unmarshal([]byte(args), &parsed); err != nil {
				return "", err
			}
			query := parsed["query"]
			
			log.Printf("\033[33m▶ [AGENT ACTION]\033[0m Called query_rag(\"%s\")", query)
			sub, err := rag.QueryRAG(o.DBConn, o.Embed, query, 3, 60)
			if err != nil {
				return "", err
			}
			
			if len(sub.Nodes) == 0 {
				log.Printf("   └─ Result: No nodes found.")
				return "No results found.", nil
			}
			
			// Mock DB returns everything. Filter it here so it behaves like a real search.
			var filteredNodes []rag.Node
			lowerQuery := strings.ToLower(query)
			for _, n := range sub.Nodes {
				if strings.Contains(strings.ToLower(n.ID), lowerQuery) || strings.Contains(strings.ToLower(n.Content), lowerQuery) {
					filteredNodes = append(filteredNodes, n)
				}
			}
			sub.Nodes = filteredNodes

			if len(sub.Nodes) == 0 {
				log.Printf("   └─ Result: No nodes matched the query '%s'.", query)
				return "No results found.", nil
			}
			
			var foundNames []string
			for _, n := range sub.Nodes {
				foundNames = append(foundNames, fmt.Sprintf("%s (%s)", n.ID, n.NodeType))
			}
			log.Printf("   └─ Result: Found %d nodes: %v", len(sub.Nodes), foundNames)
			subBytes, _ := json.Marshal(sub)
			return string(subBytes), nil
		}
		return "", fmt.Errorf("unknown tool: %s", name)
	}

	log.Println("\n\033[33m[AGENT] Starting Investigation Loop...\033[0m")
	log.Printf("\033[90m--- SYSTEM PROMPT ---\n%s\n---------------------\033[0m\n", prompts.IngestionAgentPrompt)
	log.Printf("\033[90m--- USER PROMPT ---\n%s\n-------------------\033[0m\n", transcript)

	mutationsJSON, err := o.LLM.GenerateAgentic(prompts.IngestionAgentPrompt, transcript, tools, executor)
	if err != nil {
		return fmt.Errorf("agent loop failed: %w", err)
	}

	log.Printf("\033[90m--- AGENT FINAL COMMIT ---\n%s\n------------------\033[0m\n", mutationsJSON)

	// Parse the final mutations JSON
	// Since commit_mutations takes an argument with {"mutations": [...]}, we parse it like LinkingOutput
	var linkOut LinkingOutput
	if err := json.Unmarshal([]byte(mutationsJSON), &linkOut); err != nil {
		return fmt.Errorf("failed to parse final mutations: %w", err)
	}

	log.Printf("\033[32m✔ Agent generated %d mutations.\033[0m Executing against DB...\n\n", len(linkOut.Mutations))

	for i, m := range linkOut.Mutations {
		content, _ := m.Properties["content"].(string)

		// Fallback for missing node ID
		if m.NodeID == "" {
			m.NodeID = fmt.Sprintf("temp_%s_%d", m.NodeType, i)
		}

		color := "\033[35m" // Magenta for CREATE
		if m.Operation == "UPDATE_NODE" {
			color = "\033[36m" // Cyan for UPDATE
		}
		log.Printf("%s🔨 Mutation: [%s] %s (ID: %s)\033[0m", color, m.Operation, m.NodeType, m.NodeID)

		if m.Operation == "CREATE_NODE" || m.Operation == "UPDATE_NODE" {
			ladybug.AddMockNode(m.NodeID, m.NodeType, content)
		}

		if content != "" {
			log.Printf("   ├─ Content: \033[37m%s\033[0m", content)
		}
		for k, v := range m.Properties {
			if k != "content" {
				log.Printf("   ├─ %s: \033[37m%v\033[0m", k, v)
			}
		}
		for _, e := range m.AddEdges {
			log.Printf("   └─ Add Edge: \033[33m%s\033[0m -> \033[32m%s\033[0m", e.RelType, e.TargetNodeID)
			ladybug.AddMockEdge(m.NodeID, e.TargetNodeID, e.RelType)
		}
		log.Println() // Add blank line between mutations
	}

	return nil
}
