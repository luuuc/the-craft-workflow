package cmd

import (
	"fmt"
	"strings"

	"craft/internal/state"
	"craft/internal/workflow"
)

// Status displays the current workflow state and valid actions.
func Status(_ []string) int {
	w, err := workflow.Load()
	if err != nil {
		fmt.Println("No workflow found. Run 'craft start' to begin.")
		return 0
	}

	// Check checksum and warn if mismatch
	if err := w.ValidateChecksum(); err != nil {
		fmt.Println("Warning: Workflow file modified externally. State may be inconsistent.")
		fmt.Println()
	}

	fmt.Printf("State: %s\n", w.State)
	fmt.Printf("Intent: %s\n", w.Intent)
	fmt.Println()

	fmt.Println("Notes:")
	if len(w.Notes) == 0 {
		fmt.Println("(none)")
	} else {
		for _, note := range w.Notes {
			fmt.Printf("- %s\n", note)
		}
	}
	fmt.Println()

	actions := state.NextValidActions(w.State)
	fmt.Printf("Actions: %s\n", strings.Join(actions, ", "))

	return 0
}
