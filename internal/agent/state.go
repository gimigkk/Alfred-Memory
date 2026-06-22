package agent

type ingestionState struct {
	HasExtractedManifest   bool
	HasQueriedVault        bool
	SchemaInjected         bool
	ExtractedManifestLines []ExtractedManifestLine
	LastToolResults        string
	ValidToolNodeIDs       map[string]bool
	ValidToolNodeTypes     map[string]string
	ValidToolNodeContent   map[string]string
	ExtractedSpeakers      []string
	QueriedTerms           []string
}

func newIngestionState() *ingestionState {
	return &ingestionState{
		ValidToolNodeIDs:     make(map[string]bool),
		ValidToolNodeTypes:   make(map[string]string),
		ValidToolNodeContent: make(map[string]string),
		ExtractedSpeakers:    make([]string, 0),
		QueriedTerms:         make([]string, 0),
	}
}
