package cmd

import (
	"fmt"
	"os"
	"strings"

	"craft/internal/reviewer"
	"craft/internal/state"
	"craft/internal/workflow"
)

// Think displays the current workflow state for deliberation.
func Think(args []string) int {
	w, err := workflow.Load()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error: No workflow found. Run 'craft start' to begin.")
		return 1
	}

	// Parse --review flag
	reviewFlag, reviewerName := parseReviewFlag(args)

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

	// Handle review if requested
	if reviewFlag {
		fmt.Println()
		return runReview(w, reviewerName)
	}

	return 0
}

// parseReviewFlag extracts --review or --review=X from args.
// Returns (hasFlag, reviewerName). If reviewerName is empty, auto-detect.
func parseReviewFlag(args []string) (bool, string) {
	for _, arg := range args {
		if arg == "--review" {
			return true, ""
		}
		if strings.HasPrefix(arg, "--review=") {
			return true, strings.TrimPrefix(arg, "--review=")
		}
	}
	return false, ""
}

// runReview invokes the appropriate reviewer and displays output.
func runReview(w *workflow.Workflow, reviewerName string) int {
	var rev reviewer.Reviewer
	var err error

	if reviewerName == "" {
		rev = reviewer.GetBestReviewer()
		fmt.Printf("Reviewing with %s...\n\n", rev.Name())
	} else {
		rev, err = reviewer.GetReviewer(reviewerName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s\n", err)
			return 1
		}
		fmt.Printf("Reviewing with %s...\n\n", rev.Name())
	}

	req := reviewer.ReviewRequest{
		Intent: w.Intent,
		Notes:  w.Notes,
	}

	resp, err := rev.Review(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		return 1
	}

	fmt.Println(resp.Content)
	return 0
}
