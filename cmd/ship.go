package cmd

import (
	"fmt"
	"os"

	"craft/internal/state"
	"craft/internal/workflow"
)

// Ship finalizes the workflow.
func Ship(_ []string) int {
	w, err := workflow.Load()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error: No workflow found. Run 'craft start' to begin.")
		return 1
	}

	if w.State == state.Shipped {
		fmt.Fprintln(os.Stderr, "Error: Invalid transition. Current state: shipped")
		fmt.Fprintln(os.Stderr, "Workflow already complete.")
		return 1
	}

	if w.State != state.Building {
		fmt.Fprintf(os.Stderr, "Error: Invalid transition. Current state: %s\n", w.State)
		fmt.Fprintln(os.Stderr, "Must accept before shipping.")
		return 1
	}

	if err := w.Transition(state.Shipped); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return 1
	}

	if err := w.Save(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return 1
	}

	fmt.Println("Workflow complete. State: shipped")
	fmt.Println()
	fmt.Printf("Intent: %s\n", w.Intent)
	return 0
}
