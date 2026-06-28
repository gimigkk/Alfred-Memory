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
func BuildChatPrompt(currentTime, ownerID string) string {
	return fmt.Sprintf("%s\n\n%s\n\n%s\n\n### SYSTEM CLOCK\nThe current system time is: %s\n\n### USER IDENTITY\nThe user you are interacting with has the Node ID: %s. You MUST use this ID when linking their Tasks/Events.", CorePersona, CoreSchema, SkillChat, currentTime, ownerID)
}
