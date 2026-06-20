package rag

type Subgraph struct {
	Nodes []Node
	Edges []Edge
}

type Node struct {
	ID       string
	NodeType string
	Content  string
}

type Edge struct {
	FromID  string
	ToID    string
	RelType string
	Context string
}
