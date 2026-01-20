package cmd

import (
	"fmt"
	"strings"

	"craft/internal/display"
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

	// Show started_at with relative time
	if !w.StartedAt.IsZero() {
		fmt.Printf("Started: %s (%s)\n", w.StartedAt.Local().Format("2006-01-02 15:04"), display.RelativeTime(w.StartedAt))
	}
	fmt.Println()

	// Show history timeline
	if len(w.History) > 0 {
		fmt.Println("History:")
		for _, h := range w.History {
			timeStr := h.At.Local().Format("15:04")
			if h.Note != "" {
				fmt.Printf("  %s %s \"%s\"\n", timeStr, h.State, h.Note)
			} else {
				fmt.Printf("  %s %s\n", timeStr, h.State)
			}
		}
		fmt.Println()
	}

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
