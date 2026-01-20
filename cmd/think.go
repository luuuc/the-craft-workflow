package cmd

import (
	"fmt"
	"os"
	"strings"

	"craft/internal/state"
	"craft/internal/workflow"
)

// Think displays the current workflow state for deliberation.
func Think(_ []string) int {
	w, err := workflow.Load()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error: No workflow found. Run 'craft start' to begin.")
		return 1
	}

	fmt.Println("# Intent")
	fmt.Println(w.Intent)
	fmt.Println()

	fmt.Println("## Notes")
	if len(w.Notes) == 0 {
		fmt.Println("(none)")
	} else {
		for _, note := range w.Notes {
			fmt.Printf("- %s\n", note)
		}
	}
	fmt.Println()

	fmt.Printf("State: %s\n", w.State)
	actions := state.NextValidActions(w.State)
	fmt.Printf("Actions: %s\n", strings.Join(actions, ", "))

	return 0
}
