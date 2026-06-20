package agent

import (
	"fmt"
	"strings"
	"testing"
)

// Helper to run the speaker coverage check from orchestrator.go
func runSpeakerCoverageCheck(extractedSpeakers []string, validToolNodeContent, validToolNodeTypes, batchCreatedContent, batchCreatedNodeTypes map[string]string, mutRaw []any) error {
	for _, speaker := range extractedSpeakers {
		if strings.HasPrefix(speaker, "_62_") {
			continue
		}
		speakerRepresented := false

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
					speakerTokens := tokenize(speaker)
					if isSubslice(speakerTokens, tokenize(id)) {
						return true
					}
					contentStr := validToolNodeContent[id]
					if contentStr == "" {
						contentStr = batchCreatedContent[id]
					}
					if isSubslice(speakerTokens, tokenize(contentStr)) {
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
			return fmt.Errorf("ERROR: You extracted lines from speaker '%s' but never represented them in any mutation. If they do not own a task, they MUST receive a MENTIONED_IN edge per Rule 2.", speaker)
		}
	}
	return nil
}

// Helper to run the duplicate Person guard from orchestrator.go
func runDuplicatePersonGuard(validToolNodeContent, validToolNodeTypes, batchCreatedContent, batchCreatedNodeTypes map[string]string, mutRaw []any) error {
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
						return err
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
						return err
					}
				}
			}
		}
	}
	return nil
}

// 1. Test Speaker Coverage Bug: Verify that the correct UPDATE_NODE for "rafid_harsyah" (resolving to person_rafid) passes validation now.
func TestSpeakerCoverageBug(t *testing.T) {
	extractedSpeakers := []string{"rafid_harsyah"}

	validToolNodeTypes := map[string]string{
		"person_rafid": "Person",
	}
	validToolNodeContent := map[string]string{
		"person_rafid": "Name: Rafid Harsyah, Aliases: Rapit, Rafid",
	}

	batchCreatedNodeTypes := make(map[string]string)
	batchCreatedContent := make(map[string]string)

	mutRaw := []any{
		map[string]any{
			"node_id":   "person_rafid",
			"operation": "UPDATE_NODE",
			"add_edges": []any{
				map[string]any{
					"target_node_id": "temp_event_gobak_sodor",
					"rel_type":      "MENTIONED_IN",
					"evidence_refs": []any{
						map[string]any{"line_index": 7, "quote": "Ongkeh"},
					},
				},
			},
		},
	}

	err := runSpeakerCoverageCheck(extractedSpeakers, validToolNodeContent, validToolNodeTypes, batchCreatedContent, batchCreatedNodeTypes, mutRaw)
	if err != nil {
		t.Fatalf("Expected speaker coverage check to pass, but got: %v", err)
	}
}

// 2. Test Duplicate Person Guard: Verify that attempting to CREATE_NODE a duplicate Person node is caught and rejected.
func TestDuplicatePersonGuard(t *testing.T) {
	// Vault contains person_rafid
	validToolNodeTypes := map[string]string{
		"person_rafid": "Person",
	}
	validToolNodeContent := map[string]string{
		"person_rafid": "Name: Rafid Harsyah, Aliases: Rapit, Rafid",
	}

	// Batch created content for new node (which is a duplicate)
	batchCreatedNodeTypes := map[string]string{
		"temp_person_rafid_speaker": "Person",
	}
	batchCreatedContent := map[string]string{
		"temp_person_rafid_speaker": "Name: Rafid Harsyah, Aliases: rafid_harsyah",
	}

	mutRaw := []any{
		map[string]any{
			"node_id":   "temp_person_rafid_speaker",
			"node_type": "Person",
			"operation": "CREATE_NODE",
			"properties": map[string]any{
				"name":    "Rafid Harsyah",
				"aliases": []any{"rafid_harsyah"},
			},
		},
	}

	err := runDuplicatePersonGuard(validToolNodeContent, validToolNodeTypes, batchCreatedContent, batchCreatedNodeTypes, mutRaw)
	if err == nil {
		t.Fatal("Expected duplicate person guard to fail, but it passed!")
	}
	t.Logf("Correctly rejected duplicate Person creation: %v", err)
}

// 3. Test Legitimate Person Creation: Assert that a new short-named legitimate person (e.g. "Opal" where no Opal exists in Person nodes) is accepted.
func TestLegitimatePersonCreation(t *testing.T) {
	// Vault contains person_rafid and event_dpp
	validToolNodeTypes := map[string]string{
		"person_rafid": "Person",
		"event_dpp":    "Event",
	}
	validToolNodeContent := map[string]string{
		"person_rafid": "Name: Rafid Harsyah, Aliases: Rapit, Rafid",
		"event_dpp":    "Event presentasi design DPP hari Jumat", // design, dpp, jumat exist in vault
	}

	// Creating a genuinely new person named "Opal" (4 characters)
	batchCreatedNodeTypes := map[string]string{
		"temp_person_opal": "Person",
	}
	batchCreatedContent := map[string]string{
		"temp_person_opal": "Name: Opal, Aliases: Opal",
	}

	mutRaw := []any{
		map[string]any{
			"node_id":   "temp_person_opal",
			"node_type": "Person",
			"operation": "CREATE_NODE",
			"properties": map[string]any{
				"name":    "Opal",
				"aliases": []any{"Opal"},
			},
		},
	}

	err := runDuplicatePersonGuard(validToolNodeContent, validToolNodeTypes, batchCreatedContent, batchCreatedNodeTypes, mutRaw)
	if err != nil {
		t.Fatalf("Expected new legitimate person creation to pass, but got: %v", err)
	}
}

// 4. Test Single Token Collision:
// - Case A: Person named "Jumat" (5 chars) matching "Jumat" inside an Event should pass (Person-only filtering).
// - Case B: Person named "Rapit" (5 chars) matching "Rapit" inside another Person's aliases should fail.
func TestSingleTokenCollision(t *testing.T) {
	// Vault contains person_rafid and event_dpp
	validToolNodeTypes := map[string]string{
		"person_rafid": "Person",
		"event_dpp":    "Event",
	}
	validToolNodeContent := map[string]string{
		"person_rafid": "Name: Rafid Harsyah, Aliases: Rapit, Rafid",
		"event_dpp":    "Event presentasi design DPP hari Jumat", // jumat exists in Event
	}

	// Case A: Create a Person named "Jumat" (5 chars)
	mutRawA := []any{
		map[string]any{
			"node_id":   "temp_person_jumat",
			"node_type": "Person",
			"operation": "CREATE_NODE",
			"properties": map[string]any{
				"name":    "Jumat",
				"aliases": []any{"Jumat"},
			},
		},
	}
	batchCreatedNodeTypesA := map[string]string{"temp_person_jumat": "Person"}
	batchCreatedContentA := map[string]string{"temp_person_jumat": "Name: Jumat, Aliases: Jumat"}

	errA := runDuplicatePersonGuard(validToolNodeContent, validToolNodeTypes, batchCreatedContentA, batchCreatedNodeTypesA, mutRawA)
	if errA != nil {
		t.Fatalf("Case A failed: Person named 'Jumat' was incorrectly rejected: %v", errA)
	}
	t.Log("Case A passed: Type filtering correctly ignored non-Person token matches.")

	// Case B: Create a Person named "Rapit" (5 chars) which matches person_rafid's alias
	mutRawB := []any{
		map[string]any{
			"node_id":   "temp_person_rapit",
			"node_type": "Person",
			"operation": "CREATE_NODE",
			"properties": map[string]any{
				"name":    "Rapit",
				"aliases": []any{"Rapit"},
			},
		},
	}
	batchCreatedNodeTypesB := map[string]string{"temp_person_rapit": "Person"}
	batchCreatedContentB := map[string]string{"temp_person_rapit": "Name: Rapit, Aliases: Rapit"}

	errB := runDuplicatePersonGuard(validToolNodeContent, validToolNodeTypes, batchCreatedContentB, batchCreatedNodeTypesB, mutRawB)
	if errB == nil {
		t.Fatal("Case B failed: Expected duplicate Person 'Rapit' to be rejected, but it passed!")
	}
	t.Logf("Case B passed: Correctly rejected duplicate Person 'Rapit': %v", errB)
}
