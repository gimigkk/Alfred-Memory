package ladybug

import "strings"

type Database struct{}
type Connection struct{}

func NewDatabase(path string) (*Database, error) {
	return &Database{}, nil
}

func (db *Database) Close() {}

func NewConnection(db *Database) (*Connection, error) {
	return &Connection{}, nil
}

func (c *Connection) Close() {}

type QueryResult struct {
	rows [][]any
	idx  int
}

var mockNodes = [][]any{
	{"person_bahlil", "Person", "Name: Bahlil, Aliases: Bro, Bahlil", map[string]any{"name": "Bahlil", "aliases": []any{"Bro", "Bahlil"}, "needs_clarification": false}},
	{"person_rafid", "Person", "Name: Rafid Harsyah, Aliases: Rapit, Rafid", map[string]any{"name": "Rafid Harsyah", "aliases": []any{"Rapit", "Rafid"}, "needs_clarification": false}},
	{"person_rafif", "Person", "Name: Rafif Ilmany, Aliases: Pip, Rafif, rafif_ilmany_ieee25", map[string]any{"name": "Rafif Ilmany", "aliases": []any{"Pip", "Rafif", "rafif_ilmany_ieee25"}, "needs_clarification": false}},
	{"person_rezonaldo", "Person", "Name: Rezonaldo, Aliases: Jon, Rezonaldo, rezonaldo_ieee__, VP, Vice President of External", map[string]any{"name": "Rezonaldo", "aliases": []any{"Jon", "Rezonaldo", "rezonaldo_ieee__", "VP", "Vice President of External"}, "needs_clarification": false}},
	{"person_apta", "Person", "Name: Apta, Aliases: Apta, apta_ieee25", map[string]any{"name": "Apta", "aliases": []any{"Apta", "apta_ieee25"}, "needs_clarification": false}},
	{"person_gilang", "Person", "Name: Gilang Muhamad W, Aliases: Gilang, Lang, Gilang Muhamad, m3_117_gilang_muhamad_w, M3-117_Gilang Muhamad W, You, THE USER", map[string]any{"name": "Gilang Muhamad W", "aliases": []any{"Gilang", "Lang", "Gilang Muhamad", "m3_117_gilang_muhamad_w", "M3-117_Gilang Muhamad W", "You", "THE USER"}, "needs_clarification": false}},
	{"person_jeslyn", "Person", "Name: Jeslyn, Aliases: Jes, Jeslyn, jeslyn_ieee", map[string]any{"name": "Jeslyn", "aliases": []any{"Jes", "Jeslyn", "jeslyn_ieee"}, "needs_clarification": false}},
	{"person_naufal", "Person", "Name: Naufal, Aliases: Naufal, Opal, m_naufal_ieee__", map[string]any{"name": "Naufal", "aliases": []any{"Naufal", "Opal", "m_naufal_ieee__"}, "needs_clarification": false}},
	{"person_rendi", "Person", "Name: Rendi Ramadana, Aliases: Ren, Rendi", map[string]any{"name": "Rendi Ramadana", "aliases": []any{"Ren", "Rendi"}, "needs_clarification": false}},
	{"person_nadine", "Person", "Name: Nadine, Aliases: Din, Nadine, nadine_ieee26", map[string]any{"name": "Nadine", "aliases": []any{"Din", "Nadine", "nadine_ieee26"}, "needs_clarification": false}},
	{"person_clint", "Person", "Name: Clint, Aliases: Clint", map[string]any{"name": "Clint", "aliases": []any{"Clint"}, "needs_clarification": false}},
	{"event_dpp", "Event", "Event presentasi design DPP hari Jumat", map[string]any{"content": "Event presentasi design DPP hari Jumat", "status": "planned", "needs_clarification": false}},
}

var mockEdges = [][]any{
	{"person_bahlil", "event_dpp", "PARTICIPATES_IN", "Bahlil is the key person for the DPP event"},
}

var initialMockNodes [][]any
var initialMockEdges [][]any

func init() {
	initialMockNodes = append([][]any(nil), mockNodes...)
	initialMockEdges = append([][]any(nil), mockEdges...)
}

func ResetMock() {
	mockNodes = append([][]any(nil), initialMockNodes...)
	mockEdges = append([][]any(nil), initialMockEdges...)
}

func (c *Connection) Query(query string) (*QueryResult, error) {
	var rows [][]any

	if strings.Contains(query, "rank RETURN node.id, rank") {
		// Mock PageRank hits
		rows = [][]any{}
		for _, n := range mockNodes {
			rows = append(rows, []any{n[0], 1.0})
		}
	} else if strings.Contains(query, "RETURN n.id, n.type, n.content") {
		rows = mockNodes
	} else if strings.Contains(query, "MATCH (n) RETURN n.id, label(n), n.content") || strings.Contains(query, "MATCH (n) RETURN n.id") {
		// Vault fetch all nodes
		rows = mockNodes
	} else if strings.Contains(query, "MATCH (a)-[r]->(b) RETURN a.id, b.id, label(r)") {
		// Vault fetch all edges
		rows = mockEdges
	} else if strings.Contains(query, "RETURN n.id, m.id, label(e), e.context") {
		rows = mockEdges // simplified
	} else if strings.Contains(query, "RETURN node.id") || strings.Contains(query, "RETURN m.id") || strings.Contains(query, "RETURN n.id LIMIT 1") {
		// Mock Vector Search hits - return all node IDs so the agent has full context
		rows = [][]any{}
		for _, n := range mockNodes {
			rows = append(rows, []any{n[0]})
		}
	} else if strings.Contains(query, "RETURN d") {
		// Used in seed.go to check if exists, return empty so it seeds
		rows = [][]any{}
	} else if strings.HasPrefix(query, "CREATE (n:") {
		// CREATE (n:Type {id: '...', ...})
		nodeType := ""

		if typeStart := strings.Index(query, "CREATE (n:"); typeStart != -1 {
			typeEnd := strings.Index(query[typeStart+10:], " ")
			if typeEnd != -1 {
				nodeType = query[typeStart+10 : typeStart+10+typeEnd]
			}
		}
		propsStr := ""
		if propStart := strings.Index(query, "{"); propStart != -1 {
			propEnd := strings.LastIndex(query, "}")
			if propEnd != -1 && propEnd > propStart {
				propsStr = query[propStart+1 : propEnd]
			}
		}

		props := parseCypherProps(propsStr)

		id, _ := props["id"].(string)
		content, _ := props["content"].(string)
		mockNodes = append(mockNodes, []any{id, nodeType, content, props})
		rows = [][]any{}

	} else if strings.HasPrefix(query, "MATCH (n) WHERE n.id =") {
		// SET update
		id := ""
		if idStart := strings.Index(query, "n.id = '"); idStart != -1 {
			idEnd := findStringEnd(query[idStart+8:])
			if idEnd != -1 {
				id = query[idStart+8 : idStart+8+idEnd]
			}
		}
		content := ""
		if contentStart := strings.Index(query, "n.content = '"); contentStart != -1 {
			contentEnd := findStringEnd(query[contentStart+13:])
			if contentEnd != -1 {
				content = query[contentStart+13 : contentStart+13+contentEnd]
			}
		}
		content = strings.ReplaceAll(content, "\\'", "'")
		content = strings.ReplaceAll(content, "\\\\", "\\")

		for i, n := range mockNodes {
			if n[0].(string) == id {
				existingProps, _ := n[3].(map[string]any)
				if existingProps == nil {
					existingProps = make(map[string]any)
				}
				if content != "" {
					existingProps["content"] = content
					mockNodes[i][2] = content
				}
				break
			}
		}
		rows = [][]any{}

	} else if strings.HasPrefix(query, "MATCH (a), (b) WHERE a.id =") {
		// Edge creation
		source := ""
		target := ""
		rel := ""

		if sStart := strings.Index(query, "a.id = '"); sStart != -1 {
			sEnd := findStringEnd(query[sStart+8:])
			if sEnd != -1 {
				source = query[sStart+8 : sStart+8+sEnd]
			}
		}
		if tStart := strings.Index(query, "b.id = '"); tStart != -1 {
			tEnd := findStringEnd(query[tStart+8:])
			if tEnd != -1 {
				target = query[tStart+8 : tStart+8+tEnd]
			}
		}
		if rStart := strings.Index(query, "CREATE (a)-[r:"); rStart != -1 {
			rEnd := strings.Index(query[rStart+14:], " ")
			rEnd2 := strings.Index(query[rStart+14:], "]")
			if rEnd == -1 || (rEnd2 != -1 && rEnd2 < rEnd) {
				rEnd = rEnd2
			}
			if rEnd != -1 {
				rel = query[rStart+14 : rStart+14+rEnd]
			}
		}

		source = strings.ReplaceAll(source, "\\'", "'")
		source = strings.ReplaceAll(source, "\\\\", "\\")
		target = strings.ReplaceAll(target, "\\'", "'")
		target = strings.ReplaceAll(target, "\\\\", "\\")

		mockEdges = append(mockEdges, []any{source, target, rel, ""})
		rows = [][]any{}

	} else {
		// default insert/create
		rows = [][]any{}
	}

	return &QueryResult{
		rows: rows,
		idx:  0,
	}, nil
}

func (r *QueryResult) Close() {}

func (r *QueryResult) HasNext() bool {
	return r.idx < len(r.rows)
}

func (r *QueryResult) GetNext() []any {
	row := r.rows[r.idx]
	r.idx++
	return row
}

func findStringEnd(s string) int {
	for i := 0; i < len(s); i++ {
		if s[i] == '\'' {
			// Check if it's escaped
			escaped := false
			// Count preceding backslashes
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
			if c == ' ' {
				continue
			}
			if c == ':' {
				inKey = false
				continue
			}
			currKey += string(c)
		} else {
			if c == '\\' && !escaped {
				escaped = true
				currVal += string(c)
				continue
			}

			if c == '\'' && !escaped {
				inStr = !inStr
			}
			if c == '[' && !inStr {
				inArr = true
			}
			if c == ']' && !inStr {
				inArr = false
			}

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
	if currKey != "" {
		props[strings.TrimSpace(currKey)] = cleanCypherVal(currVal)
	}
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
	if val == "true" {
		return true
	}
	if val == "false" {
		return false
	}
	if strings.HasPrefix(val, "[") && strings.HasSuffix(val, "]") {
		return val
	}
	return val
}
