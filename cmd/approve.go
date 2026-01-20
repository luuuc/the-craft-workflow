package cmd

import (
	"fmt"
	"os"

	"craft/internal/state"
	"craft/internal/structure"
	"craft/internal/workflow"
)

// Approve approves the structure and advances from shaping to building.
func Approve(args []string) int {
	w, err := workflow.Load()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error: No workflow found. Run 'craft start' to begin.")
		return 1
	}

	if w.State != state.Shaping {
		fmt.Fprintf(os.Stderr, "Error: Invalid transition. Current state: %s\n", w.State)
		if w.State == state.Thinking {
			fmt.Fprintln(os.Stderr, "Run `craft accept` first.")
		}
		return 1
	}

	// Check that pitch.md exists
	if !structure.HasPitch() {
		fmt.Fprintln(os.Stderr, "Error: No structure found. Create .craft/pitch.md or run `craft shape --generate`.")
		return 1
	}

	if err := w.TransitionWithNote(state.Building, ""); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return 1
	}

	if err := w.Save(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return 1
	}

	fmt.Println("Structure approved. State: building")
	return 0
}
