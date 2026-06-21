package agent

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/gimigkk/Alfred-Memory/assets/prompts"
	"github.com/gimigkk/Alfred-Memory/internal/embed"
	"github.com/gimigkk/Alfred-Memory/internal/ladybug"
	"github.com/gimigkk/Alfred-Memory/internal/llm"
	"github.com/gimigkk/Alfred-Memory/internal/rag"
)

type EvidenceRef struct {
	Quote     string `json:"quote"`
	LineIndex int    `json:"line_index"`
}

type EdgeMutation struct {
	RelType      string         `json:"rel_type"`
	TargetNodeID string         `json:"target_node_id"`
	EvidenceRefs []EvidenceRef  `json:"evidence_refs,omitempty"`
}

type Mutation struct {
	Operation  string                 `json:"operation"` // CREATE_NODE or UPDATE_NODE
	NodeType   string                 `json:"node_type,omitempty"`
	NodeID     string                 `json:"node_id,omitempty"`
	Properties map[string]interface{} `json:"properties"`
	AddEdges   []EdgeMutation         `json:"add_edges,omitempty"`
}

type ManifestItem struct {
	Line        string `json:"line"`
	ActionTaken string `json:"action_taken"`
}

type LinkingOutput struct {
	ManifestAccounting []ManifestItem `json:"manifest_accounting"`
	Mutations          []Mutation     `json:"mutations"`
}

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

	tools := []llm.ToolDef{
		{
			Type: "function",
			Function: llm.FunctionDef{
				Name:        "extract_transcript_manifest",
				Description: "Extract EVERY SINGLE LINE from the transcript sequentially. You MUST NOT skip any lines.",
				Parameters: map[string]any{
					"type": "object",
					"properties": map[string]any{
						"extracted_lines": map[string]any{
							"type": "array",
							"items": map[string]any{
								"type": "object",
								"properties": map[string]any{
									"speaker":        map[string]any{"type": "string"},
									"line":           map[string]any{"type": "string"},
									"shape":          map[string]any{"type": "string", "enum": []string{"directive", "confirmation", "question", "commitment", "mention"}},
									"skipped_reason": map[string]any{"type": "string"},
								},
								"required": []string{"speaker", "line", "shape"},
							},
						},
					},
					"required": []string{"extracted_lines"},
				},
			},
		},
		{
			Type: "function",
			Function: llm.FunctionDef{
				Name:        "query_rag",
				Description: "Search the knowledge vault for relevant nodes using semantic search.",
				Parameters: map[string]any{
					"type":       "object",
					"properties": map[string]any{
						"queries": map[string]any{
							"type":        "array",
							"items":       map[string]any{"type": "string"},
							"description": "An array of search queries (e.g. ['Bahlil', 'Friday event', 'Rafid']).",
						},
					},
					"required": []string{"queries"},
				},
			},
		},
		{
			Type: "function",
			Function: llm.FunctionDef{
				Name:        "commit_mutations",
				Description: "Commit the final graph mutations to the vault once all entities are resolved. YOU MUST CALL extract_transcript_manifest FIRST.",
				Parameters: map[string]any{
					"type":       "object",
					"properties": map[string]any{
						"manifest_accounting": map[string]any{
							"type": "array",
							"description": "You MUST provide an accounting for every single line you extracted in the manifest. Did you act on it? Or skip it?",
							"items": map[string]any{
								"type": "object",
								"properties": map[string]any{
									"line": map[string]any{"type": "string"},
									"action_taken": map[string]any{"type": "string", "description": "e.g., CREATE_TASK, UPDATE_EDGE, SKIP"},
									"skipped_reason": map[string]any{"type": "string"},
								},
								"required": []string{"line", "action_taken"},
							},
						},
						"mutations": map[string]any{
							"type":  "array",
							"items": map[string]any{
								"type":       "object",
								"properties": map[string]any{
									"operation":  map[string]any{"type": "string", "enum": []string{"CREATE_NODE", "UPDATE_NODE"}},
									"node_type":  map[string]any{"type": "string"},
									"node_id":    map[string]any{"type": "string"},
									"properties": map[string]any{
										"type": "object",
										"properties": map[string]any{
											"content": map[string]any{"type": "string"},
											"status": map[string]any{"type": "string"},
											"needs_clarification": map[string]any{"type": "boolean"},
											"clarification_basis": map[string]any{
												"type": "string",
												"description": "REQUIRED for Task/Event/Insight. Explain your deduction based solely on what this transcript says about this entity — ignore the content or confidence of any other node in this same mutation set.",
											},
										},
									},
									"add_edges":  map[string]any{
										"type":  "array",
										"items": map[string]any{
											"type":       "object",
											"properties": map[string]any{
												"rel_type":       map[string]any{"type": "string"},
												"target_node_id": map[string]any{"type": "string"},
												"evidence_refs": map[string]any{
													"type": "array",
													"items": map[string]any{
														"type": "object",
														"properties": map[string]any{
															"quote": map[string]any{"type": "string"},
															"line_index": map[string]any{"type": "integer"},
														},
														"required": []string{"quote", "line_index"},
													},
												},
											},
											"required": []string{"rel_type", "target_node_id", "evidence_refs"},
										},
									},
								},
								"required": []string{"operation", "node_id"},
							},
						},
					},
					"required": []string{"manifest_accounting", "mutations"},
				},
			},
		},
	}

	hasExtractedManifest := false
	var extractedManifestLines []string
	var lastToolResults string
	validToolNodeIDs := make(map[string]bool)
	validToolNodeTypes := make(map[string]string)
	validToolNodeContent := make(map[string]string)
	var extractedSpeakers []string
	var queriedTerms []string


	executor := func(name, args string) (string, error) {
		if name == "query_rag" {
			var parsed map[string][]string
			if err := json.Unmarshal([]byte(args), &parsed); err != nil {
				return "", err
			}
			queries := parsed["queries"]
			
			results := make(map[string]any)
			var allFoundNames []string
			
			// BATCH EMBEDDING: Fetch all vectors in a single API call! Consumes only 1 RPM.
			vecs, err := o.Embed.GetVectors(queries)
			if err != nil {
				// If embedding fails, return error for the tool
				return "", fmt.Errorf("batch embedding failed: %w", err)
			}

			// Using package-level normalizeStr
			for i, query := range queries {
				vec := vecs[i]
				queriedTerms = append(queriedTerms, normalizeStr(query))
				sub, err := rag.QueryRAG(o.DBConn, vec, query, 3, 60)
				if err != nil {
					results[query] = map[string]string{"error": err.Error()}
					continue
				}
				
				// Mock DB returns everything. Filter it here so it behaves like a real search.
				var filteredNodes []rag.Node
				normQuery := normalizeStr(query)
				for _, n := range sub.Nodes {
					if strings.Contains(normalizeStr(n.ID), normQuery) || strings.Contains(normalizeStr(n.Content), normQuery) {
						filteredNodes = append(filteredNodes, n)
					}
				}
				sub.Nodes = filteredNodes

				if len(sub.Nodes) == 0 {
					results[query] = "NO_MATCH: no evidence found connecting these entities — do not create a link based on this query."
				} else {
					results[query] = sub
					for _, n := range sub.Nodes {
						allFoundNames = append(allFoundNames, fmt.Sprintf("%s (%s)", n.ID, n.NodeType))
						validToolNodeIDs[n.ID] = true
						validToolNodeTypes[n.ID] = n.NodeType
						validToolNodeContent[n.ID] = n.Content
					}
				}
			}
			
			log.Printf("\033[33m▶ [AGENT ACTION]\033[0m Called query_rag for %d queries: %v", len(queries), queries)
			foundAny := false
			for _, query := range queries {
				if qRes, ok := results[query].(*rag.Subgraph); ok && len(qRes.Nodes) > 0 {
					var nodeNames []string
					for _, n := range qRes.Nodes {
						nodeNames = append(nodeNames, fmt.Sprintf("%s (%s)", n.ID, n.NodeType))
					}
					log.Printf("   └─ Result: \"%s\" → %v", query, nodeNames)
					foundAny = true
				}
			}
			if !foundAny {
				log.Printf("   └─ Result: No nodes matched any of the queries.")
			}

			subBytes, _ := json.Marshal(results)
			
			// Save tool results for substring verification later
			lastToolResults += string(subBytes) + " "
			
			return string(subBytes), nil
		} else if name == "extract_transcript_manifest" {
			hasExtractedManifest = true
			var parsed map[string][]map[string]any
			if err := json.Unmarshal([]byte(args), &parsed); err != nil {
				return "", err
			}
			lines := parsed["extracted_lines"]
			
			extractedManifestLines = []string{}
			for _, mLine := range lines {
				if lineStr, ok := mLine["line"].(string); ok {
					extractedManifestLines = append(extractedManifestLines, lineStr)
				}
				if speakerStr, ok := mLine["speaker"].(string); ok {
					extractedSpeakers = append(extractedSpeakers, speakerStr)
				}
			}
			
			// HEURISTIC VALIDATION: check if the agent dropped obvious candidate lines
			// We look for "@", "?", "Ada", "Iya", "Oke", "tolong" in the raw transcript.
			var missing []string
			transcriptLines := strings.Split(transcript, "\n")
			for _, tLine := range transcriptLines {
				lowerLine := strings.ToLower(tLine)
				if strings.Contains(tLine, "@") || strings.Contains(tLine, "?") || strings.Contains(lowerLine, "ada") || strings.Contains(lowerLine, "iya") || strings.Contains(lowerLine, "oke") || strings.Contains(lowerLine, "tolong") {
					// Check if this line is in the manifest
					found := false
					for _, mLine := range lines {
						if mStr, ok := mLine["line"].(string); ok && strings.Contains(tLine, mStr) {
							found = true
							break
						}
					}
					if !found {
						missing = append(missing, tLine)
					}
				}
			}

			if len(missing) > 0 {
				hasExtractedManifest = false // Force them to do it again properly
				errMsg := fmt.Sprintf("ERROR: Your manifest is incomplete. You dropped the following candidate lines:\n%s\nYou MUST include them in the manifest (use skipped_reason if you don't plan to act on them).", strings.Join(missing, "\n"))
				log.Printf("\033[31m[AGENT FAILED EXTRACTION]\033[0m Missing %d lines.", len(missing))
				return errMsg, nil
			}

			log.Printf("\033[33m▶ [AGENT ACTION]\033[0m Called extract_transcript_manifest successfully. Extracted %d lines.", len(lines))
			return "Manifest accepted. You may now query_rag or commit_mutations.", nil
		} else if name == "commit_mutations" {
			if !hasExtractedManifest {
				return "", fmt.Errorf("ERROR: You are strictly forbidden from committing mutations until you have successfully called extract_transcript_manifest and passed validation.")
			}
			
			var parsed map[string]any
			if err := json.Unmarshal([]byte(args), &parsed); err != nil {
				return "", err
			}
			
			// Verify manifest accounting
			accRaw, ok := parsed["manifest_accounting"].([]any)
			if !ok {
				return "", fmt.Errorf("ERROR: Missing or invalid manifest_accounting array. You must provide an accounting for every single line you extracted.")
			}
			
			var accountedLines []string
			hasCreateTaskInAccounting := false

			for _, acc := range accRaw {
				if accMap, ok := acc.(map[string]any); ok {
					if line, ok := accMap["line"].(string); ok {
						accountedLines = append(accountedLines, line)
					}
					if actionTaken, ok := accMap["action_taken"].(string); ok {
						act := strings.ToUpper(actionTaken)
						if strings.Contains(act, "CREATE_TASK") || strings.Contains(act, "TASK") {
							hasCreateTaskInAccounting = true
						}
					}
				}
			}
			
			for _, extLine := range extractedManifestLines {
				found := false
				for _, accLine := range accountedLines {
					if strings.TrimSpace(extLine) == strings.TrimSpace(accLine) {
						found = true
						break
					}
				}
				if !found {
					return "", fmt.Errorf("ERROR: Your manifest_accounting is missing line: '%s'. You must account for ALL extracted lines.", extLine)
				}
			}
			
			// Verify clarification_basis, cross-reference, and target node IDs
			hasTaskMutation := false
			batchCreatedIDs := make(map[string]bool)
			batchCreatedNodeTypes := make(map[string]string)
			batchCreatedContent := make(map[string]string)
			
			// Build list of all quotes used anywhere
			var allQuotes []string

			if mutRaw, ok := parsed["mutations"].([]any); ok {
				// First pass to collect created IDs
				for i, m := range mutRaw {
					if mItem, ok := m.(map[string]any); ok {
						op, _ := mItem["operation"].(string)
						if op == "CREATE_NODE" {
							nodeID, _ := mItem["node_id"].(string)
							nodeType, _ := mItem["node_type"].(string)
							
							var contentStr string
							if props, ok := mItem["properties"].(map[string]any); ok {
								if c, ok := props["content"].(string); ok {
									contentStr += " " + c
								}
								if n, ok := props["name"].(string); ok {
									contentStr += " " + n
								}
								if p, ok := props["phone_number"].(string); ok {
									contentStr += " " + p
								}
								if aliases, ok := props["aliases"].([]any); ok {
									for _, a := range aliases {
										contentStr += " " + fmt.Sprint(a)
									}
								}
							}

							if nodeID != "" {
								batchCreatedIDs[nodeID] = true
								batchCreatedNodeTypes[nodeID] = nodeType
								batchCreatedContent[nodeID] = contentStr
							} else {
								impliedID := fmt.Sprintf("temp_%s_%d", nodeType, i)
								batchCreatedIDs[impliedID] = true
								batchCreatedNodeTypes[impliedID] = nodeType
								batchCreatedContent[impliedID] = contentStr
							}
						}
						
						// Collect quotes from this mutation
						if edgesRaw, ok := mItem["add_edges"].([]any); ok {
							for _, eRaw := range edgesRaw {
								if edgeMap, ok := eRaw.(map[string]any); ok {
									if refsRaw, ok := edgeMap["evidence_refs"].([]any); ok {
										for _, refRaw := range refsRaw {
											if refMap, ok := refRaw.(map[string]any); ok {
												if quote, ok := refMap["quote"].(string); ok {
													allQuotes = append(allQuotes, quote)
												}
											}
										}
									}
								}
							}
						}
					}
				}

				// Check manifest accounting
				if accRaw, ok := parsed["manifest_accounting"].([]any); ok {
					for _, aRaw := range accRaw {
						if accMap, ok := aRaw.(map[string]any); ok {
							act, _ := accMap["action_taken"].(string)
							line, _ := accMap["line"].(string)
							if act == "UPDATE_TASK" || act == "UPDATE_EDGE" {
								foundQuote := false
								for _, q := range allQuotes {
									if strings.Contains(line, q) || strings.Contains(q, line) {
										foundQuote = true
										break
									}
								}
								if !foundQuote {
									return "", fmt.Errorf("ERROR: You claimed action_taken '%s' for line '%s', but this line was never used as a quote in any evidence_refs.", act, line)
								}
							}
						}
					}
				}

				// Check User Resolution (Rule 16) coverage
				for _, speaker := range extractedSpeakers {
					if strings.HasPrefix(speaker, "_62_") { continue }
					speakerNorm := normalizeStr(speaker)
					queried := false
					for _, qt := range queriedTerms {
						if strings.Contains(qt, speakerNorm) || strings.Contains(speakerNorm, qt) {
							queried = true
							break
						}
					}
					if !queried {
						return "", fmt.Errorf("ERROR: speaker '%s' appears in your extracted manifest but was never queried via query_rag. You must query for every speaker, including 'THE USER' or 'You', before committing mutations.", speaker)
					}
				}

				// Check participant coverage
				for _, speaker := range extractedSpeakers {
					if strings.HasPrefix(speaker, "_62_") { continue }
					speakerRepresented := false
					
					// Check if this speaker matches any person node in mutations
					for _, m := range mutRaw {
						if mItem, ok := m.(map[string]any); ok {
							nodeID, _ := mItem["node_id"].(string)
							ntype, _ := mItem["node_type"].(string)
							if ntype == "" {
								ntype = validToolNodeTypes[nodeID]
								if ntype == "" {
									ntype = batchCreatedNodeTypes[nodeID]
								}
							}
							
							checkID := func(id string, typ string) bool {
								if typ != "Person" { return false }
								speakerTokens := filterNoiseTokens(tokenize(speaker))
								if len(speakerTokens) == 0 {
									return false // nothing meaningful left to match on
								}
								if hasTokenOverlap(speakerTokens, tokenize(id)) { return true }
								contentStr := validToolNodeContent[id]
								if contentStr == "" { contentStr = batchCreatedContent[id] }
								if hasTokenOverlap(speakerTokens, tokenize(contentStr)) { return true }
								return false
							}

							if checkID(nodeID, ntype) {
								speakerRepresented = true
								break
							}
							
							if edgesRaw, ok := mItem["add_edges"].([]any); ok {
								for _, eRaw := range edgesRaw {
									if edgeMap, ok := eRaw.(map[string]any); ok {
										targetID, _ := edgeMap["target_node_id"].(string)
										ttype := validToolNodeTypes[targetID]
										if ttype == "" { ttype = batchCreatedNodeTypes[targetID] }
										if checkID(targetID, ttype) {
											speakerRepresented = true
											break
										}
									}
								}
							}
							if speakerRepresented { break }
						}
					}
					
					if !speakerRepresented {
						return "", fmt.Errorf("ERROR: You extracted lines from speaker '%s' but never represented them in any mutation. If they do not own a task, they MUST receive a MENTIONED_IN edge per Rule 2.", speaker)
					}
				}

				// Second pass to check target IDs, basis, and directional edges
				for _, m := range mutRaw {
					if mItem, ok := m.(map[string]any); ok {
						op, _ := mItem["operation"].(string)
						nodeID, _ := mItem["node_id"].(string)
						
						sourceType, _ := mItem["node_type"].(string)
						if sourceType == "" {
							sourceType = validToolNodeTypes[nodeID]
							if sourceType == "" {
								sourceType = batchCreatedNodeTypes[nodeID]
							}
						}

						if op == "CREATE_NODE" && sourceType == "Person" {
							var newNamesAndAliases []string
							if props, ok := mItem["properties"].(map[string]any); ok {
								if name, ok := props["name"].(string); ok && strings.TrimSpace(name) != "" {
									newNamesAndAliases = append(newNamesAndAliases, name)
								}
								if aliases, ok := props["aliases"].([]any); ok {
									for _, a := range aliases {
										if aStr := fmt.Sprint(a); strings.TrimSpace(aStr) != "" {
											newNamesAndAliases = append(newNamesAndAliases, aStr)
										}
									}
								}
							}

							checkDuplicate := func(compareID, compareContent string) error {
								for _, item := range newNamesAndAliases {
									normItem := normalizeStr(item)
									if len(normItem) < 4 {
										continue // skip short names/aliases
									}
									itemTokens := tokenize(item)
									
									// Use order-independent isTokenSubset matching
									if isTokenSubset(itemTokens, tokenize(compareID)) || isTokenSubset(itemTokens, tokenize(compareContent)) {
										return fmt.Errorf("ERROR: You are attempting to CREATE_NODE for Person '%s' (alias/name '%s') which already exists in the vault/batch as '%s'. You must use UPDATE_NODE instead per Rule 4.", nodeID, item, compareID)
									}
								}
								return nil
							}

							// Check against vault RAG nodes (Person type only)
							for valID, valContent := range validToolNodeContent {
								valType := validToolNodeTypes[valID]
								if valType != "Person" {
									continue
								}
								if err := checkDuplicate(valID, valContent); err != nil {
									return "", err
								}
							}

							// Check against other batch created nodes (Person type only)
							for batchID, batchContent := range batchCreatedContent {
								if batchID == nodeID {
									continue // Skip self
								}
								batchType := batchCreatedNodeTypes[batchID]
								if batchType != "Person" {
									continue
								}
								if err := checkDuplicate(batchID, batchContent); err != nil {
									return "", err
								}
							}
						}

						if op == "UPDATE_NODE" {
							if nodeID != "" && !validToolNodeIDs[nodeID] {
								return "", fmt.Errorf("ERROR: You are attempting to UPDATE_NODE '%s', but this node was never returned by query_rag. If this is a new entity, use CREATE_NODE instead.", nodeID)
							}
						}

						if sourceType == "Task" {
							hasTaskMutation = true
						}

						// Target Node ID verification
						if edgesRaw, ok := mItem["add_edges"].([]any); ok {
							for _, eRaw := range edgesRaw {
								if edgeMap, ok := eRaw.(map[string]any); ok {
									targetID, _ := edgeMap["target_node_id"].(string)
									relType, _ := edgeMap["rel_type"].(string)
									
									var targetType string

									if targetID != "" {
										if strings.HasPrefix(targetID, "temp_") {
											if !batchCreatedIDs[targetID] {
												return "", fmt.Errorf("ERROR: You used target_node_id '%s' which starts with temp_, but you never created this node in the current mutations batch.", targetID)
											}
											targetType = batchCreatedNodeTypes[targetID]
										} else {
											if !validToolNodeIDs[targetID] && !batchCreatedIDs[targetID] {
												return "", fmt.Errorf("ERROR: You used target_node_id '%s' but this node was never returned by query_rag, nor created in this batch. You must query for it or create it.", targetID)
											}
											targetType = validToolNodeTypes[targetID]
											if targetType == "" {
												targetType = batchCreatedNodeTypes[targetID]
											}
										}
									}

									// Directional edge validation
									if sourceType != "" && targetType != "" {
										if relType == "ASSIGNED_TO" && (sourceType != "Person" || targetType != "Task") {
											return "", fmt.Errorf("ERROR: ASSIGNED_TO edges MUST originate from a Person and point to a Task. Found %s -> ASSIGNED_TO -> %s.", sourceType, targetType)
										}
										if relType == "PART_OF" && (sourceType != "Task" || targetType != "Event") {
											return "", fmt.Errorf("ERROR: PART_OF edges MUST originate from a Task and point to an Event. Found %s -> PART_OF -> %s.", sourceType, targetType)
										}
										if relType == "HAS_ROLE" && (sourceType != "Person" || targetType != "Event") {
											return "", fmt.Errorf("ERROR: HAS_ROLE edges MUST originate from a Person and point to an Event. Found %s -> HAS_ROLE -> %s.", sourceType, targetType)
										}
										if relType == "MENTIONED_IN" && sourceType != "Person" {
											return "", fmt.Errorf("ERROR: MENTIONED_IN edges MUST originate from a Person. Found %s -> MENTIONED_IN.", sourceType)
										}
										if relType == "MENTIONED_IN" && targetType != "Task" && targetType != "Event" {
											return "", fmt.Errorf("ERROR: MENTIONED_IN edges MUST point to a Task or Event. Found pointing to %s.", targetType)
										}
									}

								}
							}
						}
						
						if sourceType == "Task" || sourceType == "Event" || sourceType == "Insight" {
							props, ok := mItem["properties"].(map[string]any)
							if !ok {
								continue // Will fail schema validation anyway
							}
							cb, ok := props["clarification_basis"].(string)
							if !ok || strings.TrimSpace(cb) == "" {
								return "", fmt.Errorf("ERROR: %s node MUST have a non-empty clarification_basis explaining your deduction.", sourceType)
							}
							if len(strings.TrimSpace(cb)) < 15 {
								return "", fmt.Errorf("ERROR: clarification_basis on %s is too short ('%s'). You must explain your deduction referencing Who/What/When/Where explicitly.", sourceType, cb)
							}
						}
					}
				}
			}

			if hasCreateTaskInAccounting && !hasTaskMutation {
				return "", fmt.Errorf("ERROR: You claimed to CREATE_TASK or act on a Task in manifest_accounting, but no Task mutation exists in the final mutations array.")
			}

			return "OK", nil
		}
		return "", fmt.Errorf("unknown tool: %s", name)
	}

	log.Println("\n\033[33m[AGENT] Starting Investigation Loop...\033[0m")
	log.Printf("\033[90m--- SYSTEM PROMPT ---\n%s\n---------------------\033[0m\n", prompts.IngestionAgentPrompt)
	log.Printf("\033[90m--- USER PROMPT ---\n%s\n-------------------\033[0m\n", transcript)

	mutationsJSON, err := o.LLM.GenerateAgentic(prompts.IngestionAgentPrompt, transcript, tools, executor)
	if err != nil {
		return nil, fmt.Errorf("agent loop failed: %w", err)
	}

	log.Printf("\033[90m--- AGENT FINAL COMMIT ---\n%s\n------------------\033[0m\n", mutationsJSON)

	// Parse the final mutations JSON
	// Since commit_mutations takes an argument with {"mutations": [...]}, we parse it like LinkingOutput
	var linkOut LinkingOutput
	if err := json.Unmarshal([]byte(mutationsJSON), &linkOut); err != nil {
		return nil, fmt.Errorf("failed to parse final mutations: %w", err)
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

		if m.Operation == "CREATE_NODE" || (m.Operation == "UPDATE_NODE" && m.NodeType != "" && content != "") {
		if !dryRun {
				ladybug.AddMockNode(m.NodeID, m.NodeType, content)
		}
		}

		if content != "" {
			log.Printf("   ├─ Content: \033[37m%s\033[0m", content)
		}
		for k, v := range m.Properties {
			if k != "content" {
				log.Printf("   ├─ %s: \033[37m%v\033[0m", k, v)
			}
		}
		transcriptLines := strings.Split(transcript, "\n")
		var survivingEdges []EdgeMutation

		for _, e := range m.AddEdges {
			// Substring verification for evidence refs
			validCount := 0
			var failedQuotes []string

			for _, ref := range e.EvidenceRefs {
				quote := strings.TrimSpace(ref.Quote)
				
				// Length and token heuristic check
				if len(e.EvidenceRefs) < 2 {
					if len(quote) < 5 {
						failedQuotes = append(failedQuotes, fmt.Sprintf("'%s' (too short)", quote))
						continue
					}
					if len(quote) < 8 && !strings.Contains(quote, " ") {
						failedQuotes = append(failedQuotes, fmt.Sprintf("'%s' (single token too short)", quote))
						continue
					}
				}

				// The quote must exist either in the specific transcript line or in the tool results
				passed := false
				if ref.LineIndex >= 0 && ref.LineIndex < len(transcriptLines) {
					if strings.Contains(transcriptLines[ref.LineIndex], ref.Quote) {
						passed = true
					}
				}
				
				if !passed && strings.Contains(lastToolResults, ref.Quote) {
					passed = true
				}
				
				if passed {
					validCount++
				} else {
					failedQuotes = append(failedQuotes, fmt.Sprintf("'%s' (not found at line %d)", quote, ref.LineIndex))
				}
			}

			if validCount == 0 {
				log.Printf("   \033[31m└─ [REJECTED EDGE]\033[0m %s -> %s (Failed all %d refs: %s)", e.RelType, e.TargetNodeID, len(e.EvidenceRefs), strings.Join(failedQuotes, ", "))
				continue
			}
			
			survivingEdges = append(survivingEdges, e)

			if len(failedQuotes) > 0 {
				log.Printf("   └─ Add Edge: \033[33m%s\033[0m -> \033[32m%s\033[0m (Verified %d/%d refs, %d failed: %s)", e.RelType, e.TargetNodeID, validCount, len(e.EvidenceRefs), len(failedQuotes), strings.Join(failedQuotes, ", "))
			} else {
				log.Printf("   └─ Add Edge: \033[33m%s\033[0m -> \033[32m%s\033[0m (Verified %d/%d refs)", e.RelType, e.TargetNodeID, validCount, len(e.EvidenceRefs))
			}
				if !dryRun {
					ladybug.AddMockEdge(m.NodeID, e.TargetNodeID, e.RelType)
				}
		}
		
		linkOut.Mutations[i].AddEdges = survivingEdges
		
		log.Println() // Add blank line between mutations
	}

	return &linkOut, nil
}

func normalizeStr(s string) string {
	reg := regexp.MustCompile("[^a-zA-Z0-9]+")
	return strings.ToLower(reg.ReplaceAllString(s, ""))
}

func tokenize(s string) []string {
	reg := regexp.MustCompile("[^a-zA-Z0-9]+")
	parts := reg.Split(s, -1)
	var tokens []string
	for _, p := range parts {
		p = strings.ToLower(p)
		if p != "" {
			tokens = append(tokens, p)
		}
	}
	return tokens
}

func isSubslice(sub, main []string) bool {
	if len(sub) == 0 {
		return false
	}
	if len(sub) > len(main) {
		return false
	}
	for i := 0; i <= len(main)-len(sub); i++ {
		match := true
		for j := 0; j < len(sub); j++ {
			if main[i+j] != sub[j] {
				match = false
				break
			}
		}
		if match {
			return true
		}
	}
	return false
}

func isTokenSubset(sub, main []string) bool {
	if len(sub) == 0 {
		return false
	}
	mainSet := make(map[string]bool, len(main))
	for _, t := range main {
		mainSet[t] = true
	}
	for _, t := range sub {
		if !mainSet[t] {
			return false
		}
	}
	return true
}

func hasTokenOverlap(sub, main []string) bool {
	if len(sub) == 0 {
		return false
	}
	mainSet := make(map[string]bool, len(main))
	for _, t := range main {
		mainSet[t] = true
	}
	for _, t := range sub {
		if mainSet[t] {
			return true
		}
	}
	return false
}

func filterNoiseTokens(tokens []string) []string {
	var filtered []string
	numericRe := regexp.MustCompile(`^[0-9]+$`)
	for _, t := range tokens {
		if numericRe.MatchString(t) {
			continue // drop purely numeric cohort years / IDs
		}
		// TODO: "ieee" is dataset-specific; consider a configurable noise-word list if Alfred is used beyond this org
		if t == "ieee" {
			continue // drop cohort/org suffix
		}
		if len(t) < 3 {
			continue // drop short tokens like "m", "w"
		}
		filtered = append(filtered, t)
	}
	if len(filtered) == 0 {
		return tokens // filtering would remove everything — fall back to unfiltered
	}
	return filtered
}
