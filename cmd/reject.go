package cmd

import (
	"fmt"
	"os"
	"strings"

	"craft/internal/state"
	"craft/internal/workflow"
)

// Reject records a concern and stays in thinking state.
func Reject(args []string) int {
	w, err := workflow.Load()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error: No workflow found. Run 'craft start' to begin.")
		return 1
	}

	if w.State != state.Thinking {
		fmt.Fprintf(os.Stderr, "Error: Invalid transition. Current state: %s\n", w.State)
		return 1
	}

	// Add note (required for reject to be meaningful)
	if len(args) > 0 {
		note := strings.Join(args, " ")
		w.AddNote(note)
	}

	if err := w.Save(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return 1
	}

	fmt.Println("Concern recorded. State: thinking")
	return 0
}
