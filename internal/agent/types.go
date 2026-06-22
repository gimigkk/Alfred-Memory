package agent

type EvidenceRef struct {
	Quote     string `json:"quote"`
	LineIndex int    `json:"line_index"`
}

type EdgeMutation struct {
	RelType      string        `json:"rel_type"`
	TargetNodeID string        `json:"target_node_id"`
	EvidenceRefs []EvidenceRef `json:"evidence_refs,omitempty"`
}

type Mutation struct {
	Operation  string                 `json:"operation"` // CREATE_NODE or UPDATE_NODE
	NodeType   string                 `json:"node_type,omitempty"`
	NodeID     string                 `json:"node_id,omitempty"`
	Properties map[string]interface{} `json:"properties"`
	AddEdges   []EdgeMutation         `json:"add_edges,omitempty"`
}

// ManifestItem is the final, post-execution accounting record for a single
// transcript line. It is constructed programmatically by the orchestrator
// after mutations have been validated and executed — never supplied by the LLM —
// so it is guaranteed to reflect what actually survived, not what the model claimed.
type ManifestItem struct {
	Line          string `json:"line"`
	Speaker       string `json:"speaker,omitempty"`
	ActionTaken   string `json:"action_taken"`
	SkippedReason string `json:"skipped_reason,omitempty"`
}

// ExtractedManifestLine captures one line as reported by the LLM's
// extract_transcript_manifest call. Shape and SkippedReason are the model's own
// characterization at extraction time, captured before any mutation exists —
// this is the only point in the pipeline where the model's stated *reason* for
// skipping a line (as opposed to the bare fact that it was skipped) is available,
// so it must be preserved here rather than re-derived later.
type ExtractedManifestLine struct {
	Speaker       string `json:"speaker"`
	Line          string `json:"line"`
	Shape         string `json:"shape"`
	SkippedReason string `json:"skipped_reason"`
}

type LinkingOutput struct {
	ManifestAccounting []ManifestItem `json:"manifest_accounting"`
	Mutations          []Mutation     `json:"mutations"`
}
