package prompts

import (
	_ "embed"
	"fmt"
)

//go:embed core_persona.md
var CorePersona string

//go:embed core_schema.md
var CoreSchema string

//go:embed skill_ingestion.md
var SkillIngestion string

//go:embed skill_chat.md
var SkillChat string

// BuildIngestionPrompt deterministically concatenates the modules for the Ingestion pipeline
func BuildIngestionPrompt() string {
	return fmt.Sprintf("%s\n\n%s\n\n%s", CorePersona, CoreSchema, SkillIngestion)
}

// BuildChatPrompt deterministically concatenates the modules for the Chat pipeline
func BuildChatPrompt() string {
	return fmt.Sprintf("%s\n\n%s\n\n%s", CorePersona, CoreSchema, SkillChat)
}
