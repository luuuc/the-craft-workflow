package cmd

import (
	"fmt"
	"os"
	"strings"

	"craft/internal/state"
	"craft/internal/workflow"
)

// Accept confirms alignment and advances from thinking to building.
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

	// Get optional note for history
	var note string
	if len(args) > 0 {
		note = strings.Join(args, " ")
		note = strings.Trim(note, "\"'")
		note = strings.TrimSpace(note)
		w.AddNote(note)
	}

	if err := w.TransitionWithNote(state.Building, note); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return 1
	}

	if err := w.Save(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return 1
	}

	fmt.Println("Intent frozen. State: building")
	return 0
}
