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
					speakerTokens := filterNoiseTokens(tokenize(speaker))
					if len(speakerTokens) == 0 {
						return false
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
					"rel_type":       "MENTIONED_IN",
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

// 5. TestSpeakerCoverageHandleWithSuffix: Verify that a speaker handle with noise suffixes (e.g. "nadine_ieee26") matches the clean person node "person_nadine".
func TestSpeakerCoverageHandleWithSuffix(t *testing.T) {
	extractedSpeakers := []string{"nadine_ieee26"}

	validToolNodeTypes := map[string]string{
		"person_nadine": "Person",
	}
	validToolNodeContent := map[string]string{
		"person_nadine": "Name: Nadine, Aliases: Din, Nadine",
	}

	batchCreatedNodeTypes := make(map[string]string)
	batchCreatedContent := make(map[string]string)

	mutRaw := []any{
		map[string]any{
			"node_id":   "person_nadine",
			"operation": "UPDATE_NODE",
			"add_edges": []any{
				map[string]any{
					"target_node_id": "temp_event_gobak_sodor",
					"rel_type":       "MENTIONED_IN",
					"evidence_refs": []any{
						map[string]any{"line_index": 12, "quote": "Iya"},
					},
				},
			},
		},
	}

	err := runSpeakerCoverageCheck(extractedSpeakers, validToolNodeContent, validToolNodeTypes, batchCreatedContent, batchCreatedNodeTypes, mutRaw)
	if err != nil {
		t.Fatalf("Expected speaker coverage check to pass for handle with suffix, but got: %v", err)
	}
}

// 6. TestSpeakerCoverageNoiseFalseMatch: Verify that speaker "jon_smith_ieee99" does NOT match an unrelated Person node that only matches on the noise token ("ieee") or suffix.
func TestSpeakerCoverageNoiseFalseMatch(t *testing.T) {
	// Nadine is the only speaker, but the update mutation is on an unrelated Person node (person_apta) that only shares the "ieee" token
	extractedSpeakers := []string{"jon_smith_ieee99"}

	validToolNodeTypes := map[string]string{
		"person_apta": "Person",
	}
	validToolNodeContent := map[string]string{
		"person_apta": "Name: Apta, Aliases: Apta, apta_ieee25", // "ieee" is present
	}

	batchCreatedNodeTypes := make(map[string]string)
	batchCreatedContent := make(map[string]string)

	mutRaw := []any{
		map[string]any{
			"node_id":   "person_apta",
			"operation": "UPDATE_NODE",
			"add_edges": []any{
				map[string]any{
					"target_node_id": "temp_event_gobak_sodor",
					"rel_type":       "MENTIONED_IN",
					"evidence_refs": []any{
						map[string]any{"line_index": 12, "quote": "Iya"},
					},
				},
			},
		},
	}

	err := runSpeakerCoverageCheck(extractedSpeakers, validToolNodeContent, validToolNodeTypes, batchCreatedContent, batchCreatedNodeTypes, mutRaw)
	if err == nil {
		t.Fatal("Expected speaker coverage check to fail due to noise-only match, but it passed!")
	}
	t.Logf("Correctly rejected false noise match: %v", err)
}

// 7. TestSpeakerCoverageShortNameOnly: Verify that fallback to unfiltered tokens allows matching speaker "jo_99" (where "jo" is 2 chars, under length floor) to a node containing "Jo".
func TestSpeakerCoverageShortNameOnly(t *testing.T) {
	extractedSpeakers := []string{"jo_99"}

	validToolNodeTypes := map[string]string{
		"person_jo": "Person",
	}
	validToolNodeContent := map[string]string{
		"person_jo": "Name: Jo, Aliases: Jo",
	}

	batchCreatedNodeTypes := make(map[string]string)
	batchCreatedContent := make(map[string]string)

	mutRaw := []any{
		map[string]any{
			"node_id":   "person_jo",
			"operation": "UPDATE_NODE",
			"add_edges": []any{
				map[string]any{
					"target_node_id": "temp_event_gobak_sodor",
					"rel_type":       "MENTIONED_IN",
					"evidence_refs": []any{
						map[string]any{"line_index": 12, "quote": "Iya"},
					},
				},
			},
		},
	}

	err := runSpeakerCoverageCheck(extractedSpeakers, validToolNodeContent, validToolNodeTypes, batchCreatedContent, batchCreatedNodeTypes, mutRaw)
	if err != nil {
		t.Fatalf("Expected speaker coverage check to pass for short name fallback, but got: %v", err)
	}
}

// --- buildManifestAccounting tests ---
//
// These tests cover the programmatic ledger builder that replaces the
// LLM-supplied manifest_accounting. The central property under test is that
// the ledger is derived strictly from the *surviving* mutation set (post edge
// validation), never from what the model claimed — so a rejected edge must
// never appear as an accounted-for action, and a SKIP'd line must carry
// forward whatever skipped_reason the model gave at extraction time.

// 8. TestManifestAccounting_EvidenceRefMatch: a line directly quoted in a
// surviving edge's evidence_refs should be reported with that edge's rel_type.
func TestManifestAccounting_EvidenceRefMatch(t *testing.T) {
	extracted := []ExtractedManifestLine{
		{Speaker: "nadine_ieee26", Line: "btw @Rendi Ramadana IEEE²⁵ lu 27rb ke gopay gua yaahh", Shape: "directive"},
	}
	mutations := []Mutation{
		{
			Operation:  "UPDATE_NODE",
			NodeID:     "person_rendi",
			Properties: map[string]any{},
			AddEdges: []EdgeMutation{
				{
					RelType:      "MENTIONED_IN",
					TargetNodeID: "task_bayar_gopay_rendi",
					EvidenceRefs: []EvidenceRef{
						{Quote: "lu 27rb ke gopay gua yaahh", LineIndex: 12},
					},
				},
			},
		},
	}

	ledger := buildManifestAccounting(extracted, mutations)
	if len(ledger) != 1 {
		t.Fatalf("expected 1 ledger entry, got %d", len(ledger))
	}
	if ledger[0].ActionTaken != "MENTIONED_IN" {
		t.Fatalf("expected action_taken MENTIONED_IN, got %q", ledger[0].ActionTaken)
	}
}

// 9. TestManifestAccounting_RejectedEdgeIsSkip: this is the direct regression
// test for the "Ongkeh" bug. An edge whose only evidence_ref failed validation
// (e.g. the executor's single-token-too-short heuristic) must never survive
// into the mutations passed to buildManifestAccounting — simulating that here
// by simply not including the rejected edge — and the corresponding line must
// be reported as SKIP, not as the edge's rel_type, even though the model
// originally intended to cite it.
func TestManifestAccounting_RejectedEdgeIsSkip(t *testing.T) {
	extracted := []ExtractedManifestLine{
		{Speaker: "rafid_harsyah", Line: "Ongkeh", Shape: "mention"},
	}
	// Simulates post-execution state: the MENTIONED_IN edge citing "Ongkeh" was
	// rejected by the edge-survival loop in RunAgenticIngestion (single token,
	// too short) and so never made it into this mutation's AddEdges.
	mutations := []Mutation{
		{
			Operation:  "UPDATE_NODE",
			NodeID:     "person_rafid",
			Properties: map[string]any{},
			AddEdges:   nil, // rejected edge removed
		},
	}

	ledger := buildManifestAccounting(extracted, mutations)
	if len(ledger) != 1 {
		t.Fatalf("expected 1 ledger entry, got %d", len(ledger))
	}
	if ledger[0].ActionTaken != "SKIP" {
		t.Fatalf("expected rejected-edge line to be reported as SKIP, got %q — a rejected edge must never appear as an accounted-for action", ledger[0].ActionTaken)
	}
}

// 10. TestManifestAccounting_SkippedReasonPreserved: a line with no surviving
// mutation, but for which the model gave a skipped_reason at extraction time,
// should carry that reason forward into the final ledger.
func TestManifestAccounting_SkippedReasonPreserved(t *testing.T) {
	extracted := []ExtractedManifestLine{
		{
			Speaker:       "apta_ieee25",
			Line:          "AI ini",
			Shape:         "mention",
			SkippedReason: "Ambiguous reaction/banter, no resolvable referent or entity.",
		},
	}
	mutations := []Mutation{} // nothing references this line

	ledger := buildManifestAccounting(extracted, mutations)
	if len(ledger) != 1 {
		t.Fatalf("expected 1 ledger entry, got %d", len(ledger))
	}
	if ledger[0].ActionTaken != "SKIP" {
		t.Fatalf("expected SKIP, got %q", ledger[0].ActionTaken)
	}
	if ledger[0].SkippedReason != "Ambiguous reaction/banter, no resolvable referent or entity." {
		t.Fatalf("expected skipped_reason to be preserved from extraction, got %q", ledger[0].SkippedReason)
	}
}

// 11. TestManifestAccounting_FabricatedEventNotLinked: regression test for the
// "es gobak sodor" bug. A single ambiguous question with no surviving mutation
// referencing it must be reported as SKIP — confirming the ledger does not
// independently invent a link to a node just because the line discusses an
// ambiguous-sounding term.
func TestManifestAccounting_FabricatedEventNotLinked(t *testing.T) {
	extracted := []ExtractedManifestLine{
		{Speaker: "jeslyn_ieee", Line: "es gobak sodor apaan dh", Shape: "question"},
	}
	// Correct post-fix behavior: no Event node was created for this line at all,
	// so no mutation in the surviving set references it.
	mutations := []Mutation{
		{
			Operation: "UPDATE_NODE",
			NodeID:    "person_jeslyn",
			AddEdges:  nil,
		},
	}

	ledger := buildManifestAccounting(extracted, mutations)
	if ledger[0].ActionTaken != "SKIP" {
		t.Fatalf("expected ambiguous question with no real referent to be SKIP, got %q", ledger[0].ActionTaken)
	}
}

// 12. TestManifestAccounting_PartialContentParaphraseMatch: a Task's `content`
// is written as an Indonesian paraphrase, not a verbatim copy, of the line that
// motivated it. The ledger must still attribute the line to that Task's
// CREATE_NODE via a meaningful (>=2 token) contiguous overlap, not require an
// exact substring match.
func TestManifestAccounting_PartialContentParaphraseMatch(t *testing.T) {
	extracted := []ExtractedManifestLine{
		{Speaker: "nadine_ieee26", Line: "guyssss ini yg belom bayar siapaa?? baru aqila, rapid, dan rapip", Shape: "question"},
	}
	mutations := []Mutation{
		{
			Operation: "CREATE_NODE",
			NodeType:  "Task",
			NodeID:    "temp_task_belum_bayar",
			Properties: map[string]any{
				"content": "Beberapa anggota belum bayar: Aqila, Rapid, Rapip disebut belum melunasi.",
			},
		},
	}

	ledger := buildManifestAccounting(extracted, mutations)
	if ledger[0].ActionTaken != "CREATE_NODE" {
		t.Fatalf("expected paraphrased Task content to tier-2 match as CREATE_NODE, got %q", ledger[0].ActionTaken)
	}
}

// 13. TestManifestAccounting_NoFalseSingleTokenMatch: a line sharing only one
// generic token with an unrelated node's content must NOT be matched — guards
// against reintroducing the single-token false-positive class this pipeline
// has already had to fix elsewhere (filterNoiseTokens, TestSingleTokenCollision).
func TestManifestAccounting_NoFalseSingleTokenMatch(t *testing.T) {
	extracted := []ExtractedManifestLine{
		{Speaker: "m_naufal_ieee__", Line: "pikm", Shape: "mention"},
	}
	mutations := []Mutation{
		{
			Operation: "CREATE_NODE",
			NodeType:  "Task",
			NodeID:    "temp_task_unrelated",
			Properties: map[string]any{
				// Shares no real tokens with "pikm" — sanity check only.
				"content": "Beberapa anggota belum bayar iuran kegiatan.",
			},
		},
	}

	ledger := buildManifestAccounting(extracted, mutations)
	if ledger[0].ActionTaken != "SKIP" {
		t.Fatalf("expected unrelated short token to remain SKIP, got %q", ledger[0].ActionTaken)
	}
}
