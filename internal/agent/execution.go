package agent

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"
)

func (o *Orchestrator) remapTempIDs(mutations []Mutation) {
	// --- ID REMAPPING PHASE ---
	// Swap temp_ IDs for real UUIDs (using UnixNano) before execution
	idMap := make(map[string]string)
	for i, m := range mutations {
		if m.NodeID == "" {
			m.NodeID = fmt.Sprintf("temp_%s_%d", m.NodeType, i)
			mutations[i].NodeID = m.NodeID
		}

		if m.Operation == "CREATE_NODE" && strings.HasPrefix(m.NodeID, "temp_") {
			// Preserve the semantic intent of the temp ID for human readability
			readablePart := strings.TrimPrefix(m.NodeID, "temp_")
			shortHash := fmt.Sprintf("%x", time.Now().UnixNano()+int64(i))
			if len(shortHash) > 6 {
				shortHash = shortHash[len(shortHash)-6:]
			}
			newID := fmt.Sprintf("%s_%s", readablePart, shortHash)

			idMap[m.NodeID] = newID
			mutations[i].NodeID = newID
		}
	}

	// Update all edges to point to the permanent IDs
	for i := range mutations {
		for j, edge := range mutations[i].AddEdges {
			if mappedID, exists := idMap[edge.TargetNodeID]; exists {
				mutations[i].AddEdges[j].TargetNodeID = mappedID
			}
		}
	}
}

func (o *Orchestrator) executeAndVerifyEdges(mutations []Mutation, transcript string, state *ingestionState, dryRun bool) {
	// --- TITLE GUARD ---
	// Proactively fix LLM hallucinations where it puts the title inside brackets in the content field
	for i := range mutations {
		m := &mutations[i]
		log.Printf("   \033[33m[DEBUG TITLE GUARD]\033[0m Processing mutation: %s", m.Operation)
		if m.Operation == "CREATE_NODE" || m.Operation == "UPDATE_NODE" {
			title, _ := m.Properties["title"].(string)
			content, _ := m.Properties["content"].(string)
			log.Printf("   \033[33m[DEBUG TITLE GUARD]\033[0m ID: %s, Title: '%s', Content: '%s'", m.NodeID, title, content)

			
			if title == "" && strings.HasPrefix(content, "[") {
				idx := strings.Index(content, "]")
				if idx != -1 {
					m.Properties["title"] = content[1:idx]
					
					// Strip the prefix and " - " if it exists
					remainder := strings.TrimSpace(content[idx+1:])
					if strings.HasPrefix(remainder, "- ") {
						remainder = strings.TrimSpace(remainder[2:])
					}
					m.Properties["content"] = remainder
					log.Printf("   \033[33m[TITLE GUARD]\033[0m Extracted title '%s' from content", m.Properties["title"])
				}
			}
		}
	}

	// --- BATCH EMBEDDING PHASE ---
	var embedTexts []string
	var embedMutations []*Mutation

	for i := range mutations {
		m := &mutations[i]
		if m.Operation == "CREATE_NODE" && m.NodeType != "Person" {
			title, _ := m.Properties["title"].(string)
			content, _ := m.Properties["content"].(string)
			log.Printf("   \033[90m[EMBED TRACE]\033[0m CREATE_NODE: %s, title: '%s', content: '%s'", m.NodeType, title, content)
			if title != "" || content != "" {
				embedTexts = append(embedTexts, title+" - "+content)
				embedMutations = append(embedMutations, m)
			}
		} else if m.Operation == "UPDATE_NODE" {
			title, hasTitle := m.Properties["title"].(string)
			content, hasContent := m.Properties["content"].(string)
			log.Printf("   \033[90m[EMBED TRACE]\033[0m UPDATE_NODE: %s, hasTitle: %v, hasContent: %v", m.NodeID, hasTitle, hasContent)

			if hasTitle || hasContent {
				// We need both to create the new embedding. Fetch missing from DB.
				if !hasTitle || !hasContent {
					query := fmt.Sprintf("MATCH (n) WHERE n.id = '%s' RETURN n.title, n.content", escapeCypher(m.NodeID))
					res, err := o.DBConn.Query(query)
					if err == nil && res.HasNext() {
						tuple := res.GetNext()
						if !hasTitle {
							if t, ok := tuple[0].(string); ok {
								title = t
							}
						}
						if !hasContent {
							if c, ok := tuple[1].(string); ok {
								content = c
							}
						}
					}
					if res != nil {
						res.Close()
					}
				}
				embedTexts = append(embedTexts, title+" - "+content)
				embedMutations = append(embedMutations, m)
			}
		}
	}

	if len(embedTexts) > 0 && !dryRun {
		vecs, err := o.Embed.GetVectors(embedTexts)
		if err != nil {
			log.Printf("   \033[31m[EMBED ERROR]\033[0m Failed to batch embed %d nodes: %v", len(embedTexts), err)
		} else {
			for i, m := range embedMutations {
				if i < len(vecs) {
					m.Properties["embedding"] = vecs[i]
				}
			}
			log.Printf("   \033[32m[EMBED]\033[0m Successfully batched generated %d embeddings.", len(vecs))
		}
	}

	transcriptLines := strings.Split(transcript, "\n")

	for i, m := range mutations {
		content, _ := m.Properties["content"].(string)

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

		if m.Operation == "CREATE_NODE" || m.Operation == "UPDATE_NODE" {
			if !dryRun {
				var query string
				if m.Operation == "CREATE_NODE" {
					query = buildCreateCypher(m)
				} else {
					query = buildUpdateCypher(m)
				}
				if query != "" {
					log.Printf("   \033[90m[CYPHER]\033[0m %s", query)
					_, err := o.DBConn.Query(query)
					if err != nil {
						log.Printf("   \033[31m[DB ERROR]\033[0m Failed to execute %s: %v", m.Operation, err)
					}
				}
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

		var survivingEdges []EdgeMutation

		for _, e := range m.AddEdges {
			// Substring verification for evidence refs
			validCount := 0
			var failedQuotes []string

			for _, ref := range e.EvidenceRefs {
				quote := strings.TrimSpace(ref.Quote)

				// The quote must exist either in the specific transcript line or in the tool results
				passed := false
				if ref.LineIndex >= 0 && ref.LineIndex < len(transcriptLines) {
					if strings.Contains(transcriptLines[ref.LineIndex], ref.Quote) {
						passed = true
					}
				}

				if !passed && strings.Contains(state.LastToolResults, ref.Quote) {
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
				query := buildEdgeCypher(m.NodeID, e)
				log.Printf("   \033[90m[CYPHER]\033[0m %s", query)
				_, err := o.DBConn.Query(query)
				if err != nil {
					log.Printf("   \033[31m[DB ERROR]\033[0m Failed to create edge %s: %v", e.RelType, err)
				}
			}
		}

		mutations[i].AddEdges = survivingEdges

		log.Println() // Add blank line between mutations
	}
}

func escapeCypher(val string) string {
	val = strings.ReplaceAll(val, "\\", "\\\\")
	return strings.ReplaceAll(val, "'", "\\'")
}

func buildCreateCypher(m Mutation) string {
	var props []string
	props = append(props, fmt.Sprintf("id: '%s'", escapeCypher(m.NodeID)))

	for k, v := range m.Properties {
		if k == "id" {
			continue // Prevent LLM from overriding the root NodeID
		}
		switch typed := v.(type) {
		case string:
			props = append(props, fmt.Sprintf("%s: '%s'", k, escapeCypher(typed)))
		case bool:
			props = append(props, fmt.Sprintf("%s: %t", k, typed))
		case []any:
			if k == "group_mentions" {
				b, _ := json.Marshal(typed)
				props = append(props, fmt.Sprintf("%s: '%s'", k, escapeCypher(string(b))))
				continue
			}
			var list []string
			for _, item := range typed {
				list = append(list, fmt.Sprintf("'%s'", escapeCypher(fmt.Sprint(item))))
			}
			props = append(props, fmt.Sprintf("%s: [%s]", k, strings.Join(list, ", ")))
		case []float32:
			var list []string
			for _, item := range typed {
				list = append(list, fmt.Sprintf("%f", item))
			}
			props = append(props, fmt.Sprintf("%s: [%s]", k, strings.Join(list, ", ")))
		default:
			b, _ := json.Marshal(v)
			props = append(props, fmt.Sprintf("%s: '%s'", k, escapeCypher(string(b))))
		}
	}

	return fmt.Sprintf("CREATE (n:%s {%s})", m.NodeType, strings.Join(props, ", "))
}

func buildUpdateCypher(m Mutation) string {
	var sets []string

	if _, hasContent := m.Properties["content"]; hasContent {
		timestamp := time.Now().Format("2006-01-02 15:04")
		sets = append(sets, fmt.Sprintf("n.history = ['%s - ' + COALESCE(n.content, '')] + COALESCE(n.history, [])", timestamp))
	}

	for k, v := range m.Properties {
		if k == "id" {
			continue // Prevent LLM from modifying the root NodeID via properties
		}
		switch typed := v.(type) {
		case string:
			sets = append(sets, fmt.Sprintf("n.%s = '%s'", k, escapeCypher(typed)))
		case bool:
			sets = append(sets, fmt.Sprintf("n.%s = %t", k, typed))
		case []any:
			if k == "group_mentions" {
				b, _ := json.Marshal(typed)
				sets = append(sets, fmt.Sprintf("n.%s = '%s'", k, escapeCypher(string(b))))
				continue
			}
			var list []string
			for _, item := range typed {
				list = append(list, fmt.Sprintf("'%s'", escapeCypher(fmt.Sprint(item))))
			}
			sets = append(sets, fmt.Sprintf("n.%s = [%s]", k, strings.Join(list, ", ")))
		case []float32:
			var list []string
			for _, item := range typed {
				list = append(list, fmt.Sprintf("%f", item))
			}
			sets = append(sets, fmt.Sprintf("n.%s = [%s]", k, strings.Join(list, ", ")))
		default:
			b, _ := json.Marshal(v)
			sets = append(sets, fmt.Sprintf("n.%s = '%s'", k, escapeCypher(string(b))))
		}
	}

	if len(sets) == 0 {
		return ""
	}
	return fmt.Sprintf("MATCH (n) WHERE n.id = '%s' SET %s", escapeCypher(m.NodeID), strings.Join(sets, ", "))
}

func buildEdgeCypher(sourceID string, e EdgeMutation) string {
	props := ""
	if len(e.EvidenceRefs) > 0 {
		b, _ := json.Marshal(e.EvidenceRefs)
		props = fmt.Sprintf(" {evidence_refs: '%s'}", escapeCypher(string(b)))
	}

	return fmt.Sprintf("MATCH (a), (b) WHERE a.id = '%s' AND b.id = '%s' CREATE (a)-[r:%s%s]->(b)",
		escapeCypher(sourceID), escapeCypher(e.TargetNodeID), e.RelType, props)
}
