package main

import (
	"fmt"
	"strings"
)

// Paste mock logic here
var mockNodes = [][]any{}
var mockEdges = [][]any{}

func findStringEnd(s string) int {
	for i := 0; i < len(s); i++ {
		if s[i] == '\'' {
			escaped := false
			j := i - 1
			for j >= 0 && s[j] == '\\' {
				escaped = !escaped
				j--
			}
			if !escaped {
				return i
			}
		}
	}
	return -1
}

func parseCypherProps(propStr string) map[string]any {
	props := make(map[string]any)
	inKey := true
	inStr := false
	inArr := false
	escaped := false
	currKey := ""
	currVal := ""
	for i := 0; i < len(propStr); i++ {
		c := propStr[i]
		if inKey {
			if c == ' ' { continue }
			if c == ':' { inKey = false; continue }
			currKey += string(c)
		} else {
			if c == '\\' && !escaped { escaped = true; currVal += string(c); continue }
			if c == '\'' && !escaped { inStr = !inStr }
			if c == '[' && !inStr { inArr = true }
			if c == ']' && !inStr { inArr = false }
			if c == ',' && !inStr && !inArr {
				props[strings.TrimSpace(currKey)] = cleanCypherVal(currVal)
				currKey = ""
				currVal = ""
				inKey = true
				escaped = false
				continue
			}
			currVal += string(c)
			escaped = false
		}
	}
	if currKey != "" { props[strings.TrimSpace(currKey)] = cleanCypherVal(currVal) }
	return props
}

func cleanCypherVal(val string) any {
	val = strings.TrimSpace(val)
	if strings.HasPrefix(val, "'") && strings.HasSuffix(val, "'") {
		s := val[1 : len(val)-1]
		s = strings.ReplaceAll(s, "\\'", "'")
		s = strings.ReplaceAll(s, "\\\\", "\\")
		return s
	}
	if val == "true" { return true }
	if val == "false" { return false }
	return val
}

func ExecuteQuery(query string) {
	if strings.HasPrefix(query, "CREATE (n:") {
		nodeType := ""
		if typeStart := strings.Index(query, "CREATE (n:"); typeStart != -1 {
			typeEnd := strings.Index(query[typeStart+10:], " ")
			if typeEnd != -1 { nodeType = query[typeStart+10 : typeStart+10+typeEnd] }
		}
		propsStr := ""
		if propStart := strings.Index(query, "{"); propStart != -1 {
			propEnd := strings.LastIndex(query, "}")
			if propEnd != -1 && propEnd > propStart { propsStr = query[propStart+1 : propEnd] }
		}
		props := parseCypherProps(propsStr)
		id, _ := props["id"].(string)
		content, _ := props["content"].(string)
		mockNodes = append(mockNodes, []any{id, nodeType, content, props})

	} else if strings.HasPrefix(query, "MATCH (a), (b) WHERE a.id =") {
		source := ""
		target := ""
		rel := ""
		if sStart := strings.Index(query, "a.id = '"); sStart != -1 {
			sEnd := findStringEnd(query[sStart+8:])
			if sEnd != -1 { source = query[sStart+8 : sStart+8+sEnd] }
		}
		if tStart := strings.Index(query, "b.id = '"); tStart != -1 {
			tEnd := findStringEnd(query[tStart+8:])
			if tEnd != -1 { target = query[tStart+8 : tStart+8+tEnd] }
		}
		if rStart := strings.Index(query, "CREATE (a)-[r:"); rStart != -1 {
			rEnd := strings.Index(query[rStart+14:], " ")
			rEnd2 := strings.Index(query[rStart+14:], "]")
			if rEnd == -1 || (rEnd2 != -1 && rEnd2 < rEnd) { rEnd = rEnd2 }
			if rEnd != -1 { rel = query[rStart+14 : rStart+14+rEnd] }
		}
		source = strings.ReplaceAll(source, "\\'", "'")
		source = strings.ReplaceAll(source, "\\\\", "\\")
		target = strings.ReplaceAll(target, "\\'", "'")
		target = strings.ReplaceAll(target, "\\\\", "\\")
		mockEdges = append(mockEdges, []any{source, target, rel, ""})
	}
}

func QueryObligations(idsStr string) [][]any {
	rows := [][]any{}
	speakerIDs := make(map[string]bool)
	for _, id := range strings.Split(idsStr, ",") {
		id = strings.TrimSpace(id)
		id = strings.Trim(id, "'")
		if id != "" { speakerIDs[id] = true }
	}
	for _, edge := range mockEdges {
		source := edge[0].(string)
		target := edge[1].(string)
		rel := edge[2].(string)
		var speakerID, nodeID string
		if speakerIDs[source] {
			speakerID = source
			nodeID = target
		} else if speakerIDs[target] {
			speakerID = target
			nodeID = source
		}
		if speakerID != "" {
			for _, n := range mockNodes {
				if n[0].(string) == nodeID {
					props, ok := n[3].(map[string]any)
					if ok {
						if nc, ok := props["needs_clarification"].(bool); ok && nc {
							content, _ := props["content"].(string)
							cb, _ := props["clarification_basis"].(string)
							nodeType := n[1].(string)
							rows = append(rows, []any{nodeID, nodeType, content, cb, rel})
						}
					}
					break
				}
			}
		}
	}
	return rows
}

func main() {
	// Block 2: Payment
	ExecuteQuery(`CREATE (n:Event {id: 'event_payment_coord_245785', clarification_basis: 'Apa tujuan pembayaran ini? Siapa saja yang diwajibkan membayar?', content: 'Koordinasi pembayaran antar anggota IEEE²⁵.', needs_clarification: true, rag_verification_query: 'Koordinasi Pembayaran', verbatim: '...'})`)
	ExecuteQuery(`CREATE (n:Task {id: 'task_transfer_apta_245a57', clarification_basis: 'Berapa nominal transfer?', content: 'Transfer dana ke rekening BCA Apta.', needs_clarification: true, rag_verification_query: 'Transfer BCA Apta', verbatim: '...'})`)
	ExecuteQuery(`CREATE (n:Task {id: 'task_transfer_jeslyn_245b82', clarification_basis: 'Siapa yang bertanggung jawab?', content: 'Transfer dana ke rekening Seabank Jeslyn.', needs_clarification: true, rag_verification_query: 'Transfer Seabank Jeslyn', verbatim: '...'})`)
	
	ExecuteQuery(`MATCH (a), (b) WHERE a.id = 'task_transfer_apta_245a57' AND b.id = 'event_payment_coord_245785' CREATE (a)-[r:PART_OF {evidence_refs: '[{...}]'}]->(b)`)
	ExecuteQuery(`MATCH (a), (b) WHERE a.id = 'task_transfer_jeslyn_245b82' AND b.id = 'event_payment_coord_245785' CREATE (a)-[r:PART_OF {evidence_refs: '[{...}]'}]->(b)`)
	ExecuteQuery(`MATCH (a), (b) WHERE a.id = 'person_naufal' AND b.id = 'task_transfer_apta_245a57' CREATE (a)-[r:ASSIGNED_TO {evidence_refs: '[{...}]'}]->(b)`)
	ExecuteQuery(`MATCH (a), (b) WHERE a.id = 'person_apta' AND b.id = 'task_transfer_apta_245a57' CREATE (a)-[r:MENTIONED_IN {evidence_refs: '[{...}]'}]->(b)`)
	ExecuteQuery(`MATCH (a), (b) WHERE a.id = 'person_jeslyn' AND b.id = 'task_transfer_jeslyn_245b82' CREATE (a)-[r:MENTIONED_IN {evidence_refs: '[{...}]'}]->(b)`)
	ExecuteQuery(`MATCH (a), (b) WHERE a.id = 'person_rafid' AND b.id = 'event_payment_coord_245785' CREATE (a)-[r:MENTIONED_IN {evidence_refs: '[{...}]'}]->(b)`)

	fmt.Printf("Mock Nodes: %d, Edges: %d\n", len(mockNodes), len(mockEdges))
	
	res := QueryObligations("'person_nadine', 'person_rendi', 'person_apta', 'person_rafid', 'person_jeslyn', 'person_naufal'")
	fmt.Printf("Found %d obligations\n", len(res))
	for _, r := range res {
		fmt.Printf("Obligation: %v\n", r)
	}
}
