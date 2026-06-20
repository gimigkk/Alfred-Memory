package rag

import (
	"fmt"
	"log"

	ladybug "github.com/gimigkk/Alfred-Memory/internal/ladybug"
	"github.com/gimigkk/Alfred-Memory/internal/embed"
)

func QueryRAG(conn *ladybug.Connection, gemini *embed.GeminiClient, query string, topK int, rrfK int) (*Subgraph, error) {
	// 1. Embed query
	vec, err := gemini.GetVector(query)
	if err != nil {
		return nil, fmt.Errorf("failed to embed query: %w", err)
	}

	vecStr := "["
	for i, v := range vec {
		if i > 0 {
			vecStr += ", "
		}
		vecStr += fmt.Sprintf("%f", v)
	}
	vecStr += "]"

	// 2. Vector Search (HNSW)
	vecQuery := fmt.Sprintf(`
		CALL QUERY_VECTOR_INDEX('doc_idx', %s, %d) 
		YIELD node 
		RETURN node.id
	`, vecStr, topK)

	res, err := conn.Query(vecQuery)
	if err != nil {
		return nil, fmt.Errorf("vector search failed: %w", err)
	}

	vectorHits := make(map[string]int)
	rank := 1
	for res.HasNext() {
		tuple := res.GetNext()
		idStr, ok := tuple[0].(string)
		if ok {
			vectorHits[idStr] = rank
			rank++
		}
	}
	res.Close()

	if len(vectorHits) == 0 {
		return &Subgraph{}, nil
	}

	// 3. Graph Expand (1-hop)
	idsList := "["
	first := true
	for id := range vectorHits {
		if !first {
			idsList += ", "
		}
		idsList += fmt.Sprintf("'%s'", id)
		first = false
	}
	idsList += "]"

	expandQuery := fmt.Sprintf(`
		MATCH (n)-[e]-(m)
		WHERE n.id IN %s
		RETURN m.id
	`, idsList)

	resExpand, err := conn.Query(expandQuery)
	if err != nil {
		return nil, fmt.Errorf("expand failed: %w", err)
	}

	expandedIds := make(map[string]bool)
	for id := range vectorHits {
		expandedIds[id] = true
	}

	for resExpand.HasNext() {
		tuple := resExpand.GetNext()
		mID, ok := tuple[0].(string)
		if ok {
			expandedIds[mID] = true
		}
	}
	resExpand.Close()

	// 4. PageRank on expanded set
	pageRankHits := make(map[string]float64)
	
	prQuery := "CALL project_graph('g', 'Document', 'LINKS_TO') CALL pagerank('g') YIELD node, rank RETURN node.id, rank"
	resPR, err := conn.Query(prQuery)
	if err == nil {
		for resPR.HasNext() {
			tuple := resPR.GetNext()
			id, ok1 := tuple[0].(string)
			val, ok2 := tuple[1].(float64)
			if ok1 && ok2 && expandedIds[id] {
				pageRankHits[id] = val
			}
		}
		resPR.Close()
	} else {
		log.Printf("Warning: PageRank failed or unsupported, using dummy scores. Error: %v", err)
		for id := range expandedIds {
			pageRankHits[id] = 1.0
		}
	}

	// 5. RRF Fusion
	finalIDs := RRF(vectorHits, pageRankHits, rrfK)

	// Fetch full nodes and edges for subgraph
	subgraph := &Subgraph{}
	if len(finalIDs) > 0 {
		fIDsList := "["
		for i, id := range finalIDs {
			if i > 0 {
				fIDsList += ", "
			}
			fIDsList += fmt.Sprintf("'%s'", id)
		}
		fIDsList += "]"

		nodeQuery := fmt.Sprintf("MATCH (n) WHERE n.id IN %s RETURN n.id, n.type, n.content", fIDsList)
		resNodes, err := conn.Query(nodeQuery)
		if err == nil {
			for resNodes.HasNext() {
				tuple := resNodes.GetNext()
				subgraph.Nodes = append(subgraph.Nodes, Node{
					ID:       tuple[0].(string),
					NodeType: tuple[1].(string),
					Content:  tuple[2].(string),
				})
			}
			resNodes.Close()
		}

		edgeQuery := fmt.Sprintf("MATCH (n)-[e]->(m) WHERE n.id IN %s AND m.id IN %s RETURN n.id, m.id, label(e), e.context", fIDsList, fIDsList)
		resEdges, err := conn.Query(edgeQuery)
		if err == nil {
			for resEdges.HasNext() {
				tuple := resEdges.GetNext()
				subgraph.Edges = append(subgraph.Edges, Edge{
					FromID:  tuple[0].(string),
					ToID:    tuple[1].(string),
					RelType: tuple[2].(string),
					Context: tuple[3].(string),
				})
			}
			resEdges.Close()
		}
	}

	return subgraph, nil
}
