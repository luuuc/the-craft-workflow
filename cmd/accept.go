package cmd

import (
	"fmt"
	"os"
	"strings"

	"craft/internal/state"
	"craft/internal/workflow"
)

// Accept confirms alignment and advances from thinking to shaping (or building with --skip-shaping).
func Accept(args []string) int {
	w, err := workflow.Load()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error: No workflow found. Run 'craft start' to begin.")
		return 1
	}

	if w.State != state.Thinking {
		fmt.Fprintf(os.Stderr, "Error: Invalid transition. Current state: %s\n", w.State)
		return 1
	}

	// Check for --skip-shaping flag
	skipShaping := false
	var filteredArgs []string
	for _, arg := range args {
		if arg == "--skip-shaping" {
			skipShaping = true
		} else {
			filteredArgs = append(filteredArgs, arg)
		}
	}

	// Get optional note for history
	var note string
	if len(filteredArgs) > 0 {
		note = strings.Join(filteredArgs, " ")
		note = strings.Trim(note, "\"'")
		note = strings.TrimSpace(note)
		w.AddNote(note)
	}

	// Determine target state
	targetState := state.Shaping
	if skipShaping {
		targetState = state.Building
	}

	if err := w.TransitionWithNote(targetState, note); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return 1
	}

	if err := w.Save(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return 1
	}

	if skipShaping {
		fmt.Println("Intent frozen. State: building")
	} else {
		fmt.Println("Intent frozen. State: shaping")
		fmt.Println()
		fmt.Println("Structure your work, then run `craft approve` to start building.")
		fmt.Println("Or run `craft shape --generate` for AI assistance.")
	}
	return 0
}
