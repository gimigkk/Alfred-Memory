package agent

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/gimigkk/Alfred-Memory/internal/rag"
)

func (o *Orchestrator) handleQueryRag(args string, state *ingestionState) (string, error) {
	var parsed map[string]any
	if err := json.Unmarshal([]byte(args), &parsed); err != nil {
		return "", err
	}

	var rawQueries []string
	if qRaw, ok := parsed["queries"].([]any); ok {
		for _, q := range qRaw {
			rawQueries = append(rawQueries, fmt.Sprint(q))
		}
	}
	var targetSpeakers []string
	if tsRaw, ok := parsed["target_speakers"].([]any); ok {
		for _, ts := range tsRaw {
			targetSpeakers = append(targetSpeakers, strings.TrimSpace(strings.ToLower(fmt.Sprint(ts))))
		}
	}

	if len(targetSpeakers) != len(rawQueries) {
		return "", fmt.Errorf("ERROR: target_speakers length (%d) must exactly match queries length (%d). You MUST provide a parallel array. If a query is not meant to resolve a speaker, use an empty string '' in the target_speakers array.", len(targetSpeakers), len(rawQueries))
	}

	var queries []string
	var validTargetSpeakers []string

	for i, qStr := range rawQueries {
		speaker := targetSpeakers[i]
		cleanQuery := strings.TrimSpace(strings.ToLower(qStr))

		if speaker != "" {
			// Validate speaker exists in manifest
			foundSpeaker := false
			for _, s := range state.ManifestSpeakers {
				if s == speaker {
					foundSpeaker = true
					break
				}
			}
			if !foundSpeaker {
				return "", fmt.Errorf("ERROR: target_speaker '%s' does not match any speaker in the manifest. If this is a manifest speaker, you MUST fix the typo. If it is not a speaker, pass an empty string '' for this element in target_speakers.", speaker)
			}
			state.QueryAttempts[speaker] = true
			queries = append(queries, qStr)
			validTargetSpeakers = append(validTargetSpeakers, speaker)
			state.ExecutedQueries[cleanQuery] = true
			log.Printf("\033[90m   [Audit] Query: \"%s\" | Target: \"%s\"\033[0m", qStr, speaker)
		} else {
			if len(cleanQuery) >= 2 {
				queries = append(queries, qStr)
				validTargetSpeakers = append(validTargetSpeakers, "")
				state.ExecutedQueries[cleanQuery] = true
				log.Printf("\033[90m   [Audit] Query: \"%s\" | Target: NONE\033[0m", qStr)
			}
		}
	}

	if len(queries) == 0 {
		return "", fmt.Errorf("ERROR: All provided queries were too short or empty. You must provide genuine search terms.")
	}

	results := make(map[string]any)
	var allFoundNames []string

	// BATCH EMBEDDING: Fetch all vectors in a single API call! Consumes only 1 RPM.
	vecs, err := o.Embed.GetVectors(queries)
	if err != nil {
		// If embedding fails, return error for the tool
		return "", fmt.Errorf("batch embedding failed: %w", err)
	}

	for i, query := range queries {
		vec := vecs[i]
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
				state.ValidToolNodeIDs[n.ID] = true
				state.ValidToolNodeTypes[n.ID] = n.NodeType
				state.ValidToolNodeContent[n.ID] = n.Content
			}
		}
	}

	log.Printf("\033[33m▶ [AGENT ACTION]\033[0m Called query_rag for %d queries: %v", len(queries), queries)
	foundAny := false
	for i, query := range queries {
		speaker := validTargetSpeakers[i]

		if qRes, ok := results[query].(*rag.Subgraph); ok && len(qRes.Nodes) > 0 {
			var nodeNames []string
			for _, n := range qRes.Nodes {
				nodeNames = append(nodeNames, fmt.Sprintf("%s (%s)", n.ID, n.NodeType))
			}
			log.Printf("   └─ Result: \"%s\" → %v", query, nodeNames)
			foundAny = true

			if speaker != "" {
				// Monotonic Upgrade
				if state.ResolvedSpeakers[speaker] != "EXISTING" {
					state.ResolvedSpeakers[speaker] = "EXISTING"
				}
			}
		}
	}
	if !foundAny {
		log.Printf("   └─ Result: No nodes matched any of the queries.")
	}

	subBytes, _ := json.Marshal(results)

	// Save tool results for substring verification later
	state.LastToolResults += string(subBytes) + " "

	return string(subBytes), nil
}

func (o *Orchestrator) handleDeclareNewSpeaker(args string, state *ingestionState) (string, error) {
	var parsed map[string]any
	if err := json.Unmarshal([]byte(args), &parsed); err != nil {
		return "", err
	}
	targetSpeaker, ok := parsed["target_speaker"].(string)
	if !ok || strings.TrimSpace(targetSpeaker) == "" {
		return "", fmt.Errorf("ERROR: target_speaker is required.")
	}
	targetSpeaker = strings.ToLower(strings.TrimSpace(targetSpeaker))

	if state.ResolvedSpeakers[targetSpeaker] == "EXISTING" {
		return "", fmt.Errorf("ERROR: target_speaker '%s' is already resolved as an EXISTING entity. You cannot downgrade them to NEW.", targetSpeaker)
	}
	if !state.QueryAttempts[targetSpeaker] {
		return "", fmt.Errorf("ERROR: You cannot declare '%s' as new without attempting to query them via query_rag first.", targetSpeaker)
	}

	state.ResolvedSpeakers[targetSpeaker] = "NEW"
	log.Printf("\033[33m▶ [AGENT ACTION]\033[0m Called declare_new_speaker for '%s'. Entity marked as NEW.", targetSpeaker)
	return fmt.Sprintf("Success: '%s' is confirmed as a new entity.", targetSpeaker), nil
}

func (o *Orchestrator) handleExtractManifest(args string, transcript string, state *ingestionState) (string, error) {
	state.HasExtractedManifest = true
	var parsed map[string][]map[string]any
	if err := json.Unmarshal([]byte(args), &parsed); err != nil {
		return "", err
	}
	lines := parsed["extracted_lines"]

	state.ExtractedManifestLines = []ExtractedManifestLine{}
	for _, mLine := range lines {
		var ml ExtractedManifestLine
		if lineStr, ok := mLine["line"].(string); ok {
			ml.Line = lineStr
		}
		if speakerStr, ok := mLine["speaker"].(string); ok {
			ml.Speaker = speakerStr
			speakerNorm := strings.ToLower(strings.TrimSpace(speakerStr))
			state.ManifestSpeakers = append(state.ManifestSpeakers, speakerNorm)
		}
		if shapeStr, ok := mLine["shape"].(string); ok {
			ml.Shape = shapeStr
		}
		if reasonStr, ok := mLine["skipped_reason"].(string); ok {
			ml.SkippedReason = reasonStr
		}
		state.ExtractedManifestLines = append(state.ExtractedManifestLines, ml)
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
		state.HasExtractedManifest = false // Force them to do it again properly
		errMsg := fmt.Sprintf("ERROR: Your manifest is incomplete. You dropped the following candidate lines:\n%s\nYou MUST include them in the manifest (use skipped_reason if you don't plan to act on them).", strings.Join(missing, "\n"))
		log.Printf("\033[31m[AGENT FAILED EXTRACTION]\033[0m Missing %d lines.", len(missing))
		return errMsg, nil
	}

	log.Printf("\033[33m▶ [AGENT ACTION]\033[0m Called extract_transcript_manifest successfully. Extracted %d lines.", len(lines))
	return "Manifest accepted. You may now query_rag or commit_mutations.", nil
}

func (o *Orchestrator) handleQuerySpeakerObligations(args string, state *ingestionState) (string, error) {
	var parsed map[string]any
	if err := json.Unmarshal([]byte(args), &parsed); err != nil {
		return "", err
	}

	var speakerIDs []string
	if raw, ok := parsed["speaker_ids"].([]any); ok {
		for _, s := range raw {
			speakerIDs = append(speakerIDs, fmt.Sprint(s))
		}
	}
	if len(speakerIDs) == 0 {
		return "", fmt.Errorf("ERROR: speaker_ids array is required and must not be empty")
	}

	// Validate all IDs were returned by previous query_rag calls
	for _, id := range speakerIDs {
		if !state.ValidToolNodeIDs[id] {
			return "", fmt.Errorf("ERROR: speaker_id '%s' was never returned by query_rag. You can only query obligations for resolved speakers", id)
		}
	}

	state.HasQueriedObligations = true

	// Build Cypher IN list
	idsList := "["
	for i, id := range speakerIDs {
		if i > 0 {
			idsList += ", "
		}
		idsList += fmt.Sprintf("'%s'", escapeCypher(id))
	}
	idsList += "]"

	// Direct graph traversal: find all unclarified Task/Event nodes and Circle memberships connected to these speakers.
	// NOTE: Circle memberships are intentionally returned here even under Layer 1 (Mention Capture) rules,
	// because the agent needs this visibility to link to existing Circles per Rule 13's corroboration path.
	// We query per-speaker because LadybugDB does not reliably return p.id as a 6th column.
	type obligation struct {
		NodeID             string   `json:"node_id"`
		NodeType           string   `json:"node_type"`
		Content            string   `json:"content"`
		ClarificationBasis string   `json:"clarification_basis"`
		EdgeType           string   `json:"edge_type"`
		ConnectedSpeakers  []string `json:"connected_speakers"`
	}

	var obligations []*obligation
	nodeMap := make(map[string]*obligation)

	for _, spkID := range speakerIDs {
		query := fmt.Sprintf(`
			MATCH (p)-[e]-(t)
			WHERE p.id = '%s'
			  AND (t.needs_clarification = true OR (label(t) = 'Circle' AND label(e) = 'MEMBER_OF'))
			RETURN t.id, label(t), t.content, t.clarification_basis, label(e)
		`, escapeCypher(spkID))

		res, err := o.DBConn.Query(query)
		if err != nil {
			log.Printf("   \033[33m[WARNING]\033[0m query_speaker_obligations query failed for %s: %v. Skipping.", spkID, err)
			continue
		}

		for res.HasNext() {
			tuple := res.GetNext()
			nodeID, _ := tuple[0].(string)
			nodeType, _ := tuple[1].(string)
			content, _ := tuple[2].(string)
			clarBasis, _ := tuple[3].(string)
			edgeType, _ := tuple[4].(string)

			if existing, ok := nodeMap[nodeID]; ok {
				found := false
				for _, s := range existing.ConnectedSpeakers {
					if s == spkID {
						found = true
						break
					}
				}
				if !found {
					existing.ConnectedSpeakers = append(existing.ConnectedSpeakers, spkID)
				}
			} else {
				ob := &obligation{
					NodeID:             nodeID,
					NodeType:           nodeType,
					Content:            content,
					ClarificationBasis: clarBasis,
					EdgeType:           edgeType,
					ConnectedSpeakers:  []string{spkID},
				}
				obligations = append(obligations, ob)
				nodeMap[nodeID] = ob

				// Register these nodes as valid for UPDATE_NODE
				state.ValidToolNodeIDs[nodeID] = true
				state.ValidToolNodeTypes[nodeID] = nodeType
				state.ValidToolNodeContent[nodeID] = content
			}
		}
		res.Close()
	}

	log.Printf("\033[33m▶ [AGENT ACTION]\033[0m Called query_speaker_obligations for %d speakers. Found %d unclarified nodes.", len(speakerIDs), len(obligations))
	for _, ob := range obligations {
		log.Printf("   └─ %s (%s): %s", ob.NodeID, ob.NodeType, truncate(ob.Content, 80))
	}
	if len(obligations) == 0 {
		log.Printf("   └─ No unclarified obligations found.")
	}

	result, _ := json.Marshal(obligations)
	state.LastToolResults += string(result) + " "
	return string(result), nil
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

func (o *Orchestrator) handleCommitMutations(args string, state *ingestionState) (string, error) {
	if !state.HasExtractedManifest {
		return "", fmt.Errorf("ERROR: You are strictly forbidden from committing mutations until you have successfully called extract_transcript_manifest and passed validation.")
	}
	if !state.SchemaInjected {
		return "", fmt.Errorf("ERROR: You are not authorized to commit yet. You must output the thought [REQUEST_SCHEMA] to receive the graph mapping rules first.")
	}

	// Strict schema validation first
	var rootPayload struct {
		Mutations []Mutation `json:"mutations"`
	}
	dec := json.NewDecoder(strings.NewReader(args))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&rootPayload); err != nil {
		return "", fmt.Errorf("JSON Schema Error: %v. You likely hallucinated an unknown field (e.g. nesting 'mutations' inside a mutation object instead of directly placing 'add_edges' in it). Stick strictly to the required schema.", err)
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

					if nodeType == "Circle" {
						return "", fmt.Errorf("ERROR: You are attempting to CREATE_NODE for a Circle. Inline Circle creation is currently STRICTLY FORBIDDEN (see skill_commit.md Rule 15). You MUST capture the reference in the group_mentions array of the relevant Task/Event instead.")
					}

					var ragQuery string
					var contentStr string
					if props, ok := mItem["properties"].(map[string]any); ok {
						if c, ok := props["content"].(string); ok {
							contentStr += " " + c
						}
						if r, ok := props["rag_verification_query"].(string); ok {
							ragQuery = r
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

					if nodeType == "Person" || nodeType == "Event" || nodeType == "Project" || nodeType == "Circle" {
						if ragQuery == "" {
							return "", fmt.Errorf("ERROR: You are attempting to create a new %s node, but you did not provide a rag_verification_query.", nodeType)
						}

						cleanQuery := strings.TrimSpace(strings.ToLower(ragQuery))
						if !state.ExecutedQueries[cleanQuery] {
							return "", fmt.Errorf("ERROR: You are attempting to create a new %s node, but our logs show you never ran the query '%s'. You must execute query_rag with this exact string to verify the entity before creating it.", nodeType, ragQuery)
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
		for _, speaker := range state.ManifestSpeakers {
			if _, ok := state.ResolvedSpeakers[speaker]; !ok {
				return "", fmt.Errorf("ERROR: speaker '%s' appears in your extracted manifest but was never resolved. You must query for every speaker, including 'THE USER' or 'You', before committing mutations.", speaker)
			}
		}

		// Check participant coverage
		for _, speaker := range state.ManifestSpeakers {
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
						ntype = state.ValidToolNodeTypes[nodeID]
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
						contentStr := state.ValidToolNodeContent[id]
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
								ttype := state.ValidToolNodeTypes[targetID]
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
					sourceType = state.ValidToolNodeTypes[nodeID]
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
					for valID, valContent := range state.ValidToolNodeContent {
						valType := state.ValidToolNodeTypes[valID]
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
					if nodeID != "" && !state.ValidToolNodeIDs[nodeID] {
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
									if !state.ValidToolNodeIDs[targetID] && !batchCreatedIDs[targetID] {
										return "", fmt.Errorf("ERROR: You used target_node_id '%s' but this node was never returned by query_rag, nor created in this batch. You must query for it or create it.", targetID)
									}
									targetType = state.ValidToolNodeTypes[targetID]
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

					// Content validation
					contentStr, hasContent := props["content"].(string)
					if hasContent && len(strings.TrimSpace(contentStr)) < 40 {
						return "", fmt.Errorf("ERROR: content on %s is too short ('%s'). It MUST be a highly descriptive, verbose narrative containing ALL known facts (Who, What, When, Where, Why), ensuring ZERO DATA LOSS. Do not just write a short title.", sourceType, contentStr)
					}

					// Clarification Basis validation
					nc, hasNc := props["needs_clarification"].(bool)
					cb, hasCb := props["clarification_basis"].(string)

					if !hasNc || nc {
						if !hasCb || strings.TrimSpace(cb) == "" {
							return "", fmt.Errorf("ERROR: %s node MUST have a non-empty clarification_basis explaining what specific details are missing.", sourceType)
						}
						if len(strings.TrimSpace(cb)) < 15 {
							return "", fmt.Errorf("ERROR: clarification_basis on %s is too short ('%s'). You must ask specific questions about Who/What/When/Where explicitly.", sourceType, cb)
						}
					} else {
						// If needs_clarification is false, cb should be empty
						if hasCb && strings.TrimSpace(cb) != "" {
							return "", fmt.Errorf("ERROR: clarification_basis on %s MUST be empty if needs_clarification is false. All facts belong in the content field.", sourceType)
						}
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
