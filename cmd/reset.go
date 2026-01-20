package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"craft/internal/workflow"
)

// Reset abandons the current workflow.
func Reset(args []string) int {
	// Check for --force flag
	force := false
	for _, arg := range args {
		if arg == "--force" || arg == "-f" {
			force = true
			break
		}
	}

	if !workflow.Exists() {
		fmt.Println("No workflow to reset.")
		return 0
	}

	w, err := workflow.Load()
	if err != nil {
		// File exists but can't be parsed - still allow reset
		if !force {
			fmt.Print("Abandon corrupted workflow? [y/N] ")
			if !confirmReset() {
				fmt.Println("Cancelled.")
				return 0
			}
		}
	} else if !force {
		fmt.Printf("Abandon workflow \"%s\"? [y/N] ", w.Intent)
		if !confirmReset() {
			fmt.Println("Cancelled.")
			return 0
		}
	}

	if err := workflow.Delete(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return 1
	}

	fmt.Println("Workflow abandoned.")
	return 0
}

func confirmReset() bool {
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return false
	}
	input = strings.TrimSpace(strings.ToLower(input))
	return input == "y" || input == "yes"
}
