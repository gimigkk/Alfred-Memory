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
	{"person_bahlil", "Person", "Name: Bahlil, Aliases: Bro, Bahlil"},
	{"person_rafid", "Person", "Name: Rafid Harsyah, Aliases: Rapit, Rafid"},
	{"person_rafif", "Person", "Name: Rafif Ilmany, Aliases: Pip, Rafif, rafif_ilmany_ieee25"},
	{"person_rezonaldo", "Person", "Name: Rezonaldo, Aliases: Jon, Rezonaldo, rezonaldo_ieee__, VP, Vice President of External"},
	{"person_apta", "Person", "Name: Apta, Aliases: Apta, apta_ieee25"},
	{"person_gilang", "Person", "Name: Gilang Muhamad W, Aliases: Gilang, Lang, Gilang Muhamad, m3_117_gilang_muhamad_w, M3-117_Gilang Muhamad W, You, THE USER"},
	{"person_jeslyn", "Person", "Name: Jeslyn, Aliases: Jes, Jeslyn, jeslyn_ieee"},
	{"person_naufal", "Person", "Name: Naufal, Aliases: Naufal, Opal"},
	{"person_rendi", "Person", "Name: Rendi Ramadana, Aliases: Ren, Rendi"},
	{"person_nadine", "Person", "Name: Nadine, Aliases: Din, Nadine"},
	{"person_clint", "Person", "Name: Clint, Aliases: Clint"},
	{"event_dpp", "Event", "Event presentasi design DPP hari Jumat"},
}

var mockEdges = [][]any{
	{"person_bahlil", "event_dpp", "PARTICIPATES_IN", "Bahlil is the key person for the DPP event"},
}

// Helper to inject agent mutations into the mock DB state
func AddMockNode(id, nodeType, content string) {
	mockNodes = append(mockNodes, []any{id, nodeType, content})
}

func AddMockEdge(from, to, relType string) {
	mockEdges = append(mockEdges, []any{from, to, relType, ""})
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
