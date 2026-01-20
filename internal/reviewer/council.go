package reviewer

import (
	"fmt"
	"os/exec"
	"strings"
)

// CouncilReviewer invokes council-cli for multi-perspective review.
type CouncilReviewer struct{}

func (r *CouncilReviewer) Name() string {
	return NameCouncil
}

func (r *CouncilReviewer) Available() bool {
	_, err := exec.LookPath("council")
	return err == nil
}

func (r *CouncilReviewer) Review(req ReviewRequest) (ReviewResponse, error) {
	// Build the review input
	var input strings.Builder
	input.WriteString(req.Intent)
	if len(req.Notes) > 0 {
		input.WriteString("\n\nNotes:\n")
		for _, note := range req.Notes {
			input.WriteString("- " + note + "\n")
		}
	}

	cmd := exec.Command("council", "review", input.String())
	output, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			stderr := string(exitErr.Stderr)
			if strings.Contains(stderr, "unknown command") {
				return ReviewResponse{}, fmt.Errorf("council review command not available (council-cli may need updating)")
			}
			return ReviewResponse{}, fmt.Errorf("council failed: %s", stderr)
		}
		return ReviewResponse{}, fmt.Errorf("council failed: %w", err)
	}

	return ReviewResponse{
		Content:  strings.TrimSpace(string(output)),
		Reviewer: NameCouncil,
	}, nil
}
