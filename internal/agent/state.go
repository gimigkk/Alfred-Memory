package agent

type ingestionState struct {
	HasExtractedManifest   bool
	SchemaInjected         bool
	ExtractedManifestLines []ExtractedManifestLine
	LastToolResults        string
	ValidToolNodeIDs       map[string]bool
	ValidToolNodeTypes     map[string]string
	ValidToolNodeContent   map[string]string
	ManifestSpeakers       []string
	QueryAttempts          map[string]bool
	ResolvedSpeakers       map[string]string
	ExecutedQueries        map[string]bool
}

func newIngestionState() *ingestionState {
	return &ingestionState{
		ValidToolNodeIDs:     make(map[string]bool),
		ValidToolNodeTypes:   make(map[string]string),
		ValidToolNodeContent: make(map[string]string),
		ManifestSpeakers:     make([]string, 0),
		QueryAttempts:        make(map[string]bool),
		ResolvedSpeakers:     make(map[string]string),
		ExecutedQueries:      make(map[string]bool),
	}
}
