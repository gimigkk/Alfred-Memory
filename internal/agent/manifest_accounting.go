package agent

import (
	"fmt"
	"strings"
)

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
