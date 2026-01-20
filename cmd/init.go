package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"craft/internal/templates"
)

// Init copies integration templates to the current project.
func Init(args []string) int {
	flags := parseInitFlags(args)

	if !flags.claude && !flags.cursor && !flags.all {
		fmt.Fprintln(os.Stderr, "Error: No template specified.")
		fmt.Fprintln(os.Stderr, "Usage: craft init [--claude] [--cursor] [--all]")
		return 1
	}

	if flags.all {
		flags.claude = true
		flags.cursor = true
	}

	var created, skipped []string

	if flags.claude {
		c, s := installClaudeTemplates()
		created = append(created, c...)
		skipped = append(skipped, s...)
	}

	if flags.cursor {
		c, s := installCursorTemplates()
		created = append(created, c...)
		skipped = append(skipped, s...)
	}

	for _, f := range created {
		fmt.Printf("Created: %s\n", f)
	}
	for _, f := range skipped {
		fmt.Printf("Skipped: %s (already exists)\n", f)
	}

	if len(created) == 0 && len(skipped) > 0 {
		fmt.Println("\nAll templates already exist.")
	}

	return 0
}

type initFlags struct {
	claude bool
	cursor bool
	all    bool
}

func parseInitFlags(args []string) initFlags {
	var flags initFlags
	for _, arg := range args {
		switch arg {
		case "--claude":
			flags.claude = true
		case "--cursor":
			flags.cursor = true
		case "--all":
			flags.all = true
		}
	}
	return flags
}

func installClaudeTemplates() (created, skipped []string) {
	// CLAUDE.md
	claudeMD, err := templates.ClaudeTemplate()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading CLAUDE.md template: %v\n", err)
		return
	}
	if writeIfNotExists("CLAUDE.md", claudeMD) {
		created = append(created, "CLAUDE.md")
	} else {
		skipped = append(skipped, "CLAUDE.md")
	}

	// .claude/commands/craft.md
	craftCmd, err := templates.CraftCommand()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading craft.md template: %v\n", err)
		return
	}
	cmdPath := filepath.Join(".claude", "commands", "craft.md")
	if writeIfNotExists(cmdPath, craftCmd) {
		created = append(created, cmdPath)
	} else {
		skipped = append(skipped, cmdPath)
	}

	return
}

func installCursorTemplates() (created, skipped []string) {
	cursorRules, err := templates.CursorRules()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading cursor-rules.md template: %v\n", err)
		return
	}
	if writeIfNotExists(".cursorrules", cursorRules) {
		created = append(created, ".cursorrules")
	} else {
		skipped = append(skipped, ".cursorrules")
	}
	return
}

func writeIfNotExists(path string, content []byte) bool {
	if _, err := os.Stat(path); err == nil {
		return false // file exists
	}

	dir := filepath.Dir(path)
	if dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			fmt.Fprintf(os.Stderr, "Error creating directory %s: %v\n", dir, err)
			return false
		}
	}

	if err := os.WriteFile(path, content, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing %s: %v\n", path, err)
		return false
	}
	return true
}
