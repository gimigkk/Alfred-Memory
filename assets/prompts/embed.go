package prompts

import (
	_ "embed"
)

//go:embed ingestion_agent.md
var IngestionAgentPrompt string

// cache bust 5
