package templates

import "embed"

//go:embed files/*
var embedded embed.FS

// ClaudeTemplate returns the CLAUDE.md template content.
func ClaudeTemplate() ([]byte, error) {
	return embedded.ReadFile("files/CLAUDE.md")
}

// CraftCommand returns the craft.md slash command content.
func CraftCommand() ([]byte, error) {
	return embedded.ReadFile("files/craft.md")
}

// CursorRules returns the cursor-rules.md template content.
func CursorRules() ([]byte, error) {
	return embedded.ReadFile("files/cursor-rules.md")
}

// IntegrationDoc returns the INTEGRATION.md template content.
func IntegrationDoc() ([]byte, error) {
	return embedded.ReadFile("files/INTEGRATION.md")
}
