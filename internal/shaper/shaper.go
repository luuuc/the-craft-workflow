package shaper

// Shaper name constants.
const (
	NameShapeCLI = "ShapeCLI"
	NameAI       = "AI"
	NameManual   = "Manual"
)

// ShapeRequest contains context for structure generation.
type ShapeRequest struct {
	Intent string
	Notes  []string
}

// ShapeResult contains the generated structure.
type ShapeResult struct {
	PitchPath string   // Path to generated pitch
	CardPaths []string // Paths to generated cards
	Shaper    string   // e.g., "AI", "ShapeCLI", "Manual"
}

// Shaper can generate project structure from intent.
type Shaper interface {
	// Name returns the shaper identifier.
	Name() string

	// Available returns true if this shaper can be used.
	Available() bool

	// Shape generates structure and writes files.
	Shape(req ShapeRequest) (ShapeResult, error)
}

// GetBestShaper returns the highest-priority available shaper.
// Priority: ShapeCLI > AI > Manual (nil)
func GetBestShaper() Shaper {
	shapers := []Shaper{
		&ShapeCLIShaper{},
		&AIShaper{},
	}

	for _, s := range shapers {
		if s.Available() {
			return s
		}
	}

	return nil // No shaper available, manual mode
}
