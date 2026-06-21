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
	RelType      string        `json:"rel_type"`
	TargetNodeID string        `json:"target_node_id"`
	EvidenceRefs []EvidenceRef `json:"evidence_refs,omitempty"`
}

type Mutation struct {
	Operation  string                 `json:"operation"` // CREATE_NODE or UPDATE_NODE
	NodeType   string                 `json:"node_type,omitempty"`
	NodeID     string                 `json:"node_id,omitempty"`
	Properties map[string]interface{} `json:"properties"`
	AddEdges   []EdgeMutation         `json:"add_edges,omitempty"`
}

// ManifestItem is the final, post-execution accounting record for a single
// transcript line. It is constructed programmatically by the orchestrator
// after mutations have been validated and executed — never supplied by the LLM —
// so it is guaranteed to reflect what actually survived, not what the model claimed.
type ManifestItem struct {
	Line          string `json:"line"`
	Speaker       string `json:"speaker,omitempty"`
	ActionTaken   string `json:"action_taken"`
	SkippedReason string `json:"skipped_reason,omitempty"`
}

// ExtractedManifestLine captures one line as reported by the LLM's
// extract_transcript_manifest call. Shape and SkippedReason are the model's own
// characterization at extraction time, captured before any mutation exists —
// this is the only point in the pipeline where the model's stated *reason* for
// skipping a line (as opposed to the bare fact that it was skipped) is available,
// so it must be preserved here rather than re-derived later.
type ExtractedManifestLine struct {
	Speaker       string `json:"speaker"`
	Line          string `json:"line"`
	Shape         string `json:"shape"`
	SkippedReason string `json:"skipped_reason"`
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
					"type": "object",
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
				Description: "Commit the final graph mutations to the vault once all entities are resolved. This is all-or-nothing: if it returns an error, the entire batch is rejected. You must resubmit all mutations in your next attempt. YOU MUST CALL extract_transcript_manifest FIRST.",
				Parameters: map[string]any{
					"type": "object",
					"properties": map[string]any{
						"mutations": map[string]any{
							"type": "array",
							"items": map[string]any{
								"type": "object",
								"properties": map[string]any{
									"operation": map[string]any{"type": "string", "enum": []string{"CREATE_NODE", "UPDATE_NODE"}},
									"node_type": map[string]any{"type": "string"},
									"node_id":   map[string]any{"type": "string"},
									"properties": map[string]any{
										"type": "object",
										"properties": map[string]any{
											"content":             map[string]any{"type": "string"},
											"status":              map[string]any{"type": "string"},
											"needs_clarification": map[string]any{"type": "boolean"},
											"clarification_basis": map[string]any{
												"type":        "string",
												"description": "REQUIRED for Task/Event/Insight. Explain your deduction based solely on what this transcript says about this entity — ignore the content or confidence of any other node in this same mutation set.",
											},
										},
									},
									"add_edges": map[string]any{
										"type": "array",
										"items": map[string]any{
											"type": "object",
											"properties": map[string]any{
												"rel_type":       map[string]any{"type": "string"},
												"target_node_id": map[string]any{"type": "string"},
												"evidence_refs": map[string]any{
													"type": "array",
													"items": map[string]any{
														"type": "object",
														"properties": map[string]any{
															"quote":      map[string]any{"type": "string"},
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
					"required": []string{"mutations"},
				},
			},
		},
	}

	hasExtractedManifest := false
	var extractedManifestLines []ExtractedManifestLine
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
				queryTokens := filterNoiseTokens(tokenize(query))
				
				for _, n := range sub.Nodes {
					// Semantic match heuristic: strict substring OR token overlap (simulating fuzzy vector search)
					if strings.Contains(normalizeStr(n.ID), normQuery) || 
					   strings.Contains(normalizeStr(n.Content), normQuery) || 
					   hasTokenOverlap(queryTokens, tokenize(n.Content)) {
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

			extractedManifestLines = []ExtractedManifestLine{}
			for _, mLine := range lines {
				var ml ExtractedManifestLine
				if lineStr, ok := mLine["line"].(string); ok {
					ml.Line = lineStr
				}
				if speakerStr, ok := mLine["speaker"].(string); ok {
					ml.Speaker = speakerStr
					extractedSpeakers = append(extractedSpeakers, speakerStr)
				}
				if shapeStr, ok := mLine["shape"].(string); ok {
					ml.Shape = shapeStr
				}
				if reasonStr, ok := mLine["skipped_reason"].(string); ok {
					ml.SkippedReason = reasonStr
				}
				extractedManifestLines = append(extractedManifestLines, ml)
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

				// Check User Resolution (Rule 16) coverage
				for _, speaker := range extractedSpeakers {
					if strings.HasPrefix(speaker, "_62_") {
						continue
					}
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
					if strings.HasPrefix(speaker, "_62_") {
						continue
					}
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
								if typ != "Person" {
									return false
								}
								speakerTokens := filterNoiseTokens(tokenize(speaker))
								if len(speakerTokens) == 0 {
									return false // nothing meaningful left to match on
								}
								if hasTokenOverlap(speakerTokens, tokenize(id)) {
									return true
								}
								contentStr := validToolNodeContent[id]
								if contentStr == "" {
									contentStr = batchCreatedContent[id]
								}
								if hasTokenOverlap(speakerTokens, tokenize(contentStr)) {
									return true
								}
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
										if ttype == "" {
											ttype = batchCreatedNodeTypes[targetID]
										}
										if checkID(targetID, ttype) {
											speakerRepresented = true
											break
										}
									}
								}
							}
							if speakerRepresented {
								break
							}
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

			if hasTaskMutation {
				log.Printf("\033[90m   (commit includes at least one Task mutation)\033[0m")
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

		// --- CLARITY GUARD ---
		if cb, ok := m.Properties["clarification_basis"].(string); ok {
			cbLower := strings.ToLower(cb)
			if strings.Contains(cbLower, "unknown") {
				
				if nc, hasNc := m.Properties["needs_clarification"].(bool); hasNc && !nc {
					m.Properties["needs_clarification"] = true
					log.Printf("   \033[33m[GUARD]\033[0m Overriding needs_clarification to TRUE due to missing context in basis.")
				}
			}
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

				// Removed strict length heuristics. We will just check if the quote exists.

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

				// Fallback: If LLM messed up the line index, search the whole transcript
				if !passed && strings.Contains(transcript, ref.Quote) {
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

	// Build the manifest accounting ledger programmatically, from the mutations
	// that actually survived validation and edge-quote verification above. This
	// guarantees 100% consistency between the ledger and what was truly applied —
	// the LLM no longer asserts this; see buildManifestAccounting for the matching logic.
	linkOut.ManifestAccounting = buildManifestAccounting(extractedManifestLines, linkOut.Mutations)
	log.Printf("\033[32m✔ Manifest accounting built for %d lines.\033[0m", len(linkOut.ManifestAccounting))

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

// longestCommonRun returns the length of the longest contiguous run of tokens
// shared between a and b, preserving order (classic longest-common-substring,
// applied to token sequences rather than characters). Used where a single shared
// token would be too weak a signal but exact whole-sequence containment (isSubslice)
// would be too strict — e.g. matching a transcript line against a paraphrased
// node content where only a meaningful phrase, not the whole line, recurs verbatim.
func longestCommonRun(a, b []string) int {
	if len(a) == 0 || len(b) == 0 {
		return 0
	}
	prev := make([]int, len(b)+1)
	best := 0
	for i := 1; i <= len(a); i++ {
		curr := make([]int, len(b)+1)
		for j := 1; j <= len(b); j++ {
			if a[i-1] == b[j-1] {
				curr[j] = prev[j-1] + 1
				if curr[j] > best {
					best = curr[j]
				}
			}
		}
		prev = curr
	}
	return best
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

// buildManifestAccounting constructs the final manifest ledger programmatically,
// from the mutations that actually survived validation and edge-quote verification
// (i.e. the post-execution state of `mutations`, where rejected edges have already
// been filtered out of each mutation's AddEdges). It never trusts an LLM-supplied
// claim about what happened to a line — every action_taken here is derived by
// re-scanning the surviving mutation set for evidence that actually cites the line.
//
// Matching priority per line (first match wins):
//  1. The line is (sub)matched by a quote in a surviving edge's evidence_refs.
//     This is the strongest signal: evidence_refs are already verified, in the
//     execution loop, to be real substrings of a real transcript line. The
//     action_taken is set to that edge's rel_type, scoped to the specific
//     source/target node pair, since the same line may support multiple edges.
//  2. The line is not cited by any edge, but its content tokens form a contiguous
//     subsequence inside a surviving node's own content/name/aliases (i.e. the
//     line plausibly seeded that node's properties directly, e.g. a Task whose
//     `content` restates the line). Token-subsequence matching (not single-token
//     overlap) is used deliberately — single shared tokens are exactly the false
//     positive class this pipeline has previously had to guard against, and a
//     manifest ledger that reuses that weak match would misattribute a SKIP'd
//     line to an unrelated CREATE_NODE.
//  3. Neither — the line was not used by any surviving mutation. Reported as
//     SKIP. If the model gave a skipped_reason at extraction time (before any
//     mutation existed, so before it had any incentive to rationalize a
//     post-hoc decision), that reason is preserved; otherwise left blank.
func buildManifestAccounting(extractedLines []ExtractedManifestLine, mutations []Mutation) []ManifestItem {
	type quoteHit struct {
		relType  string
		nodeID   string
		targetID string
	}

	// Index every surviving edge's evidence quotes for tier-1 matching.
	var edgeHits []quoteHit
	edgeQuoteText := make(map[int]string) // index into edgeHits -> quote text used (for substring test)
	hitIdx := 0
	for _, m := range mutations {
		for _, e := range m.AddEdges {
			for _, ref := range e.EvidenceRefs {
				q := strings.TrimSpace(ref.Quote)
				if q == "" {
					continue
				}
				edgeHits = append(edgeHits, quoteHit{
					relType:  e.RelType,
					nodeID:   m.NodeID,
					targetID: e.TargetNodeID,
				})
				edgeQuoteText[hitIdx] = q
				hitIdx++
			}
		}
	}

	// Index every surviving node's content/name/aliases for tier-2 matching.
	type nodeContentEntry struct {
		nodeID    string
		operation string
		tokens    []string
	}
	var nodeEntries []nodeContentEntry
	for _, m := range mutations {
		var contentStr string
		if c, ok := m.Properties["content"].(string); ok {
			contentStr += " " + c
		}
		if n, ok := m.Properties["name"].(string); ok {
			contentStr += " " + n
		}
		if aliases, ok := m.Properties["aliases"].([]any); ok {
			for _, a := range aliases {
				contentStr += " " + fmt.Sprint(a)
			}
		}
		if strings.TrimSpace(contentStr) == "" {
			continue
		}
		nodeEntries = append(nodeEntries, nodeContentEntry{
			nodeID:    m.NodeID,
			operation: m.Operation,
			tokens:    tokenize(contentStr),
		})
	}

	ledger := make([]ManifestItem, 0, len(extractedLines))

	for _, ml := range extractedLines {
		line := strings.TrimSpace(ml.Line)
		item := ManifestItem{
			Line:    ml.Line,
			Speaker: ml.Speaker,
		}

		if line == "" {
			item.ActionTaken = "SKIP"
			item.SkippedReason = ml.SkippedReason
			ledger = append(ledger, item)
			continue
		}

		// Tier 1: matched by a surviving edge's evidence_ref.
		matched := false
		for i, hit := range edgeHits {
			q := edgeQuoteText[i]
			if q == "" {
				continue
			}
			if strings.Contains(line, q) || strings.Contains(q, line) {
				item.ActionTaken = hit.relType
				matched = true
				break
			}
		}

		// Tier 2: matched by a surviving node's own content (CREATE_NODE/UPDATE_NODE
		// whose properties plausibly draw on this line). Node `content` is written
		// as an Indonesian narrative paraphrase per the prompt's storage rules, not
		// a verbatim copy of the line, so we cannot require the whole line as a
		// contiguous subsequence inside the node content (that would almost never
		// fire). Instead we require the longest common contiguous token run between
		// the two to be at least 2 tokens — long enough to rule out a single shared
		// word (the false-positive class this pipeline already guards against
		// elsewhere), short enough to still catch genuine paraphrase overlap.
		if !matched {
			lineTokens := tokenize(line)
			for _, ne := range nodeEntries {
				if longestCommonRun(lineTokens, ne.tokens) >= 2 {
					item.ActionTaken = ne.operation
					matched = true
					break
				}
			}
		}

		if !matched {
			item.ActionTaken = "SKIP"
			item.SkippedReason = ml.SkippedReason
		}

		ledger = append(ledger, item)
	}

	return ledger
}
