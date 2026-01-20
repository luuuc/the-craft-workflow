package reviewer

import (
	"errors"
	"fmt"
)

// Reviewer name constants.
const (
	NameCouncil = "Council"
	NameAI      = "AI"
	NameNone    = "None"

	// CLI flag values (lowercase).
	FlagCouncil = "council"
	FlagAI      = "ai"
	FlagNone    = "none"
)

// ReviewRequest contains the context for a review.
type ReviewRequest struct {
	Intent string
	Notes  []string
}

// ReviewResponse contains the review output.
type ReviewResponse struct {
	Content  string // The review text
	Reviewer string // e.g., "AI", "Council", "None"
}

// Reviewer can review workflow intent.
type Reviewer interface {
	// Name returns the reviewer identifier.
	Name() string

	// Available returns true if this reviewer can be used.
	Available() bool

	// Review examines the intent and returns feedback.
	Review(req ReviewRequest) (ReviewResponse, error)
}

// GetBestReviewer returns the highest-priority available reviewer.
// Priority: Council > AI > Null
func GetBestReviewer() Reviewer {
	reviewers := []Reviewer{
		&CouncilReviewer{},
		&AIReviewer{},
		&NullReviewer{},
	}

	for _, r := range reviewers {
		if r.Available() {
			return r
		}
	}

	return &NullReviewer{} // Should never reach here
}

// GetReviewer returns a specific reviewer by name.
func GetReviewer(name string) (Reviewer, error) {
	switch name {
	case FlagCouncil:
		r := &CouncilReviewer{}
		if !r.Available() {
			return nil, errors.New("council not found in PATH")
		}
		return r, nil
	case FlagAI:
		r := &AIReviewer{}
		if !r.Available() {
			return nil, errors.New("CRAFT_AI_API_KEY not set")
		}
		return r, nil
	case FlagNone:
		return &NullReviewer{}, nil
	default:
		return nil, fmt.Errorf("unknown reviewer: %s", name)
	}
}
