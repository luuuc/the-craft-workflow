package workflow

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"craft/internal/state"
)

const (
	CraftDir      = ".craft"
	WorkflowFile  = "workflow.md"
	SchemaVersion = 1
)

// Workflow represents a craft workflow.
type Workflow struct {
	State         state.State
	SchemaVersion int
	Checksum      string
	Intent        string
	Notes         []string
}

// Path returns the full path to the workflow file.
func Path() string {
	return filepath.Join(CraftDir, WorkflowFile)
}

// Exists returns true if a workflow file exists.
func Exists() bool {
	_, err := os.Stat(Path())
	return err == nil
}

// EnsureDir creates the .craft directory if it doesn't exist.
func EnsureDir() error {
	return os.MkdirAll(CraftDir, 0755)
}

// Load reads and parses the workflow file.
func Load() (*Workflow, error) {
	data, err := os.ReadFile(Path())
	if err != nil {
		if os.IsNotExist(err) {
			return nil, errors.New("no workflow found")
		}
		return nil, fmt.Errorf("failed to read workflow: %w", err)
	}
	return Parse(data)
}

// Parse parses workflow content from bytes.
func Parse(data []byte) (*Workflow, error) {
	content := string(data)

	// Extract front matter
	if !strings.HasPrefix(content, "---\n") {
		return nil, errors.New("invalid workflow file: missing front matter")
	}

	parts := strings.SplitN(content[4:], "\n---\n", 2)
	if len(parts) != 2 {
		return nil, errors.New("invalid workflow file: malformed front matter")
	}

	frontMatter := parts[0]
	body := parts[1]

	w := &Workflow{
		SchemaVersion: SchemaVersion,
	}

	// Parse front matter
	for _, line := range strings.Split(frontMatter, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		kv := strings.SplitN(line, ":", 2)
		if len(kv) != 2 {
			continue
		}
		key := strings.TrimSpace(kv[0])
		value := strings.TrimSpace(kv[1])

		switch key {
		case "state":
			w.State = state.State(value)
		case "schema_version":
			fmt.Sscanf(value, "%d", &w.SchemaVersion)
		case "checksum":
			w.Checksum = value
		}
	}

	if !w.State.Valid() {
		return nil, fmt.Errorf("invalid workflow state: %s", w.State)
	}

	// Parse body
	w.Intent, w.Notes = parseBody(body)

	return w, nil
}

func parseBody(body string) (intent string, notes []string) {
	lines := strings.Split(body, "\n")
	inIntent := false
	inNotes := false

	var intentLines []string

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		if trimmed == "# Intent" {
			inIntent = true
			inNotes = false
			continue
		}
		if trimmed == "## Notes" {
			inIntent = false
			inNotes = true
			continue
		}

		if inIntent && trimmed != "" {
			intentLines = append(intentLines, trimmed)
		}
		if inNotes && strings.HasPrefix(trimmed, "- ") {
			notes = append(notes, strings.TrimPrefix(trimmed, "- "))
		}
	}

	intent = strings.Join(intentLines, " ")
	return intent, notes
}

// Save writes the workflow to disk atomically.
func (w *Workflow) Save() error {
	if err := EnsureDir(); err != nil {
		return fmt.Errorf("failed to create .craft directory: %w", err)
	}

	content := w.Format()

	// Write to temp file first
	tmpPath := Path() + ".tmp"
	if err := os.WriteFile(tmpPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write workflow: %w", err)
	}

	// Atomic rename
	if err := os.Rename(tmpPath, Path()); err != nil {
		os.Remove(tmpPath)
		return fmt.Errorf("failed to save workflow: %w", err)
	}

	return nil
}

// formatNotes returns notes formatted for the workflow file.
func (w *Workflow) formatNotes() string {
	if len(w.Notes) == 0 {
		return "(none)"
	}
	var noteLines []string
	for _, n := range w.Notes {
		noteLines = append(noteLines, "- "+n)
	}
	return strings.Join(noteLines, "\n")
}

// contentForChecksum returns the content used for checksum computation.
// This excludes the checksum field itself to allow verification.
func (w *Workflow) contentForChecksum() string {
	return fmt.Sprintf(`---
state: %s
schema_version: %d
---

# Intent
%s

## Notes
%s
`, w.State, w.SchemaVersion, w.Intent, w.formatNotes())
}

// Format returns the workflow as a formatted string with checksum.
func (w *Workflow) Format() string {
	content := w.contentForChecksum()
	checksum := ComputeChecksum([]byte(content))
	w.Checksum = checksum

	return fmt.Sprintf(`---
state: %s
schema_version: %d
checksum: %s
---

# Intent
%s

## Notes
%s
`, w.State, w.SchemaVersion, checksum, w.Intent, w.formatNotes())
}

// ComputeChecksum generates a SHA-256 checksum (first 8 hex chars).
func ComputeChecksum(content []byte) string {
	hash := sha256.Sum256(content)
	return hex.EncodeToString(hash[:])[:8]
}

// ValidateChecksum verifies the stored checksum matches the content.
func (w *Workflow) ValidateChecksum() error {
	expected := ComputeChecksum([]byte(w.contentForChecksum()))
	if w.Checksum != expected {
		return errors.New("checksum mismatch: workflow file may have been modified externally")
	}
	return nil
}

// Delete removes the workflow file.
func Delete() error {
	err := os.Remove(Path())
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete workflow: %w", err)
	}

	// Try to remove .craft dir if empty
	os.Remove(CraftDir)

	return nil
}

// New creates a new workflow with the given intent.
func New(intent string) *Workflow {
	return &Workflow{
		State:         state.Thinking,
		SchemaVersion: SchemaVersion,
		Intent:        intent,
		Notes:         nil,
	}
}

// AddNote appends a note to the workflow.
func (w *Workflow) AddNote(note string) {
	note = strings.TrimSpace(note)
	note = strings.Trim(note, "\"'")
	if note != "" {
		w.Notes = append(w.Notes, note)
	}
}

// Transition validates and performs a state transition.
func (w *Workflow) Transition(to state.State) error {
	if err := state.ValidateTransition(w.State, to); err != nil {
		return err
	}
	w.State = to
	return nil
}
