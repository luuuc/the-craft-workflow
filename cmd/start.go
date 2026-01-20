package cmd

import (
	"fmt"
	"os"
	"strings"

	"craft/internal/workflow"
)

// Start begins a new workflow with the given intent.
func Start(args []string) int {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "Error: Intent required. Usage: craft start \"<intent>\"")
		return 1
	}

	intent := strings.Join(args, " ")
	intent = strings.Trim(intent, "\"'")
	intent = strings.TrimSpace(intent)

	if intent == "" {
		fmt.Fprintln(os.Stderr, "Error: Intent cannot be empty. Usage: craft start \"<intent>\"")
		return 1
	}

	if workflow.Exists() {
		fmt.Fprintln(os.Stderr, "Error: Workflow already exists. Run 'craft reset' to abandon.")
		return 1
	}

	w := workflow.New(intent)
	if err := w.Save(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return 1
	}

	fmt.Println("Workflow started. State: thinking")
	return 0
}
