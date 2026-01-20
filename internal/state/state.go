package state

import "fmt"

// State represents a workflow state.
type State string

const (
	Thinking State = "thinking"
	Shaping  State = "shaping"
	Building State = "building"
	Shipped  State = "shipped"
)

// Valid returns true if s is a recognized state.
func (s State) Valid() bool {
	switch s {
	case Thinking, Shaping, Building, Shipped:
		return true
	default:
		return false
	}
}

// String returns the state as a string.
func (s State) String() string {
	return string(s)
}

// ValidateTransition checks if a transition from one state to another is allowed.
// Transition rules:
//   - thinking → shaping (via accept)
//   - thinking → building (via accept --skip-shaping)
//   - shaping → building (via approve)
//   - building → shipped (via ship)
func ValidateTransition(from, to State) error {
	if !from.Valid() {
		return fmt.Errorf("invalid current state: %s", from)
	}
	if !to.Valid() {
		return fmt.Errorf("invalid target state: %s", to)
	}

	switch from {
	case Thinking:
		if to == Shaping || to == Building {
			return nil
		}
		return fmt.Errorf("invalid transition: cannot go from %s to %s", from, to)
	case Shaping:
		if to == Building {
			return nil
		}
		return fmt.Errorf("invalid transition: cannot go from %s to %s", from, to)
	case Building:
		if to == Shipped {
			return nil
		}
		return fmt.Errorf("invalid transition: cannot go from %s to %s", from, to)
	case Shipped:
		return fmt.Errorf("invalid transition: %s is a terminal state", from)
	default:
		return fmt.Errorf("unknown state: %s", from)
	}
}

// NextValidActions returns the actions available from the given state.
func NextValidActions(current State) []string {
	switch current {
	case Thinking:
		return []string{"accept", "accept --skip-shaping", "reject", "reset"}
	case Shaping:
		return []string{"shape", "approve", "revise", "reset"}
	case Building:
		return []string{"ship", "reset"}
	case Shipped:
		return []string{"reset"}
	default:
		return nil
	}
}
