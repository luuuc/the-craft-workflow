package cmd

import (
	"fmt"
	"os"

	"craft/internal/shaper"
	"craft/internal/state"
	"craft/internal/structure"
	"craft/internal/workflow"
)

// Shape displays shaping status or generates structure with --generate flag.
func Shape(args []string) int {
	w, err := workflow.Load()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error: No workflow found. Run 'craft start' to begin.")
		return 1
	}

	if w.State != state.Shaping {
		fmt.Fprintf(os.Stderr, "Error: Shape only works in shaping state. Current state: %s\n", w.State)
		return 1
	}

	// Check for --generate flag
	generate := false
	for _, arg := range args {
		if arg == "--generate" {
			generate = true
			break
		}
	}

	if generate {
		return generateStructure(w)
	}

	return showShapingStatus(w)
}

func showShapingStatus(w *workflow.Workflow) int {
	pitch, cards, err := structure.ListStructure()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error listing structure: %v\n", err)
		return 1
	}

	fmt.Printf("Shaping: %s\n", w.Intent)

	if pitch == "" && len(cards) == 0 {
		fmt.Println("Structure: (none)")
		fmt.Println("Next: craft shape --generate OR create .craft/pitch.md")
	} else {
		fmt.Println("Structure:")
		if pitch != "" {
			fmt.Printf("  %s\n", pitch)
		}
		for _, c := range cards {
			fmt.Printf("  %s\n", c)
		}
		fmt.Println("Next: craft approve")
	}

	return 0
}

func generateStructure(w *workflow.Workflow) int {
	s := shaper.GetBestShaper()

	if s == nil || s.Name() == shaper.NameManual {
		fmt.Println("No shaper available. Create .craft/pitch.md manually.")
		return 0
	}

	fmt.Printf("Generating via %s...\n", s.Name())

	req := shaper.ShapeRequest{
		Intent: w.Intent,
		Notes:  w.Notes,
	}

	result, err := s.Shape(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return 1
	}

	fmt.Println("Created:")
	if result.PitchPath != "" {
		fmt.Printf("  %s\n", result.PitchPath)
	}
	for _, card := range result.CardPaths {
		fmt.Printf("  %s\n", card)
	}
	fmt.Println("Next: craft approve")
	return 0
}
