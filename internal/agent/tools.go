package agent

import "github.com/gimigkk/Alfred-Memory/internal/llm"

var ingestionTools []llm.ToolDef

func init() {
	ingestionTools = []llm.ToolDef{
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
						"target_speakers": map[string]any{
							"type":        "array",
							"items":       map[string]any{"type": "string"},
							"description": "REQUIRED whenever this call is used to resolve manifest speakers. Must exactly match the length of 'queries'. For each query, provide the exact speaker label it resolves, or an empty string '' if that specific query is not for a speaker resolution. A length mismatch against 'queries' rejects the entire call.",
						},
					},
					"required": []string{"queries"},
				},
			},
		},
		{
			Type: "function",
			Function: llm.FunctionDef{
				Name:        "declare_new_speaker",
				Description: "Declare that a speaker from the manifest is a new entity that does not exist in the vault. You MUST attempt to search for them via query_rag first.",
				Parameters: map[string]any{
					"type": "object",
					"properties": map[string]any{
						"target_speaker": map[string]any{
							"type":        "string",
							"description": "The EXACT speaker label from the manifest.",
						},
					},
					"required": []string{"target_speaker"},
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
											"rag_verification_query": map[string]any{
												"type":        "string",
												"description": "REQUIRED for CREATE_NODE on Person, Event, Project, or Circle. The exact string you queried via query_rag to verify this entity didn't exist.",
											},
											"content":             map[string]any{"type": "string"},
											"status":              map[string]any{"type": "string"},
											"verbatim":            map[string]any{"type": "string"},
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
}

func GetIngestionTools() []llm.ToolDef {
	return ingestionTools
}
