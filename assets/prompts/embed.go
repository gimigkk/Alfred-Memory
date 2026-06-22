package prompts

import (
	_ "embed"
	"fmt"
)

//go:embed core_persona.md
var CorePersona string

//go:embed core_schema.md
var CoreSchema string

//go:embed skill_discovery.md
var SkillDiscovery string

//go:embed skill_commit.md
var SkillCommit string

//go:embed skill_chat.md
var SkillChat string

// BuildDiscoveryPrompt concatenates the persona and the discovery rules
func BuildDiscoveryPrompt() string {
	return fmt.Sprintf("%s\n\n%s", CorePersona, SkillDiscovery)
}

// BuildCommitPrompt concatenates the topology schema and the commit mapping rules
func BuildCommitPrompt() string {
	return fmt.Sprintf("%s\n\n%s", CoreSchema, SkillCommit)
}

// BuildChatPrompt deterministically concatenates the modules for the Chat pipeline
func BuildChatPrompt() string {
	return fmt.Sprintf("%s\n\n%s\n\n%s", CorePersona, CoreSchema, SkillChat)
}
