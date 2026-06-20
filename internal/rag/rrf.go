package rag

import (
	"sort"
)

// RRF performs Reciprocal Rank Fusion on vector search ranks and PageRank scores.
// vectorHits maps NodeID -> rank (1 is best)
// pageRankHits maps NodeID -> PageRank score (higher is better)
func RRF(vectorHits map[string]int, pageRankHits map[string]float64, k int) []string {
	type prHit struct {
		id    string
		score float64
	}
	var prList []prHit
	for id, score := range pageRankHits {
		prList = append(prList, prHit{id: id, score: score})
	}
	
	// Sort PR hits descending
	sort.Slice(prList, func(i, j int) bool {
		return prList[i].score > prList[j].score
	})

	prRanks := make(map[string]int)
	for i, hit := range prList {
		prRanks[hit.id] = i + 1
	}

	scores := make(map[string]float64)
	
	// Add vector scores
	for id, rank := range vectorHits {
		scores[id] += 1.0 / float64(k+rank)
	}
	
	// Add PageRank scores
	for id, rank := range prRanks {
		scores[id] += 1.0 / float64(k+rank)
	}

	type rrfHit struct {
		id    string
		score float64
	}
	var finalHits []rrfHit
	for id, score := range scores {
		finalHits = append(finalHits, rrfHit{id: id, score: score})
	}

	// Sort final by RRF score descending
	sort.Slice(finalHits, func(i, j int) bool {
		return finalHits[i].score > finalHits[j].score
	})

	var result []string
	for _, hit := range finalHits {
		result = append(result, hit.id)
	}

	return result
}
