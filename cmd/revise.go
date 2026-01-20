package cmd

import (
	"fmt"
	"os"
	"strings"

	"craft/internal/state"
	"craft/internal/workflow"
)

// Revise records a concern during shaping without advancing state.
func Revise(args []string) int {
	w, err := workflow.Load()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error: No workflow found. Run 'craft start' to begin.")
		return 1
	}

	if w.State != state.Shaping {
		fmt.Fprintln(os.Stderr, "Error: Invalid state. Revise only works during shaping.")
		return 1
	}

	// Require note argument
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "Error: Note required. Usage: craft revise \"note\"")
		return 1
	}

	note := strings.Join(args, " ")
	note = strings.Trim(note, "\"'")
	note = strings.TrimSpace(note)

	if note == "" {
		fmt.Fprintln(os.Stderr, "Error: Note required. Usage: craft revise \"note\"")
		return 1
	}

	// Add note with [revise] prefix
	w.AddNote("[revise] " + note)

	if err := w.Save(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return 1
	}

	fmt.Println("Concern recorded. State: shaping")
	return 0
}
