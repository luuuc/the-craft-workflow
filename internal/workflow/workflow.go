package workflow

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"craft/internal/state"
)

const (
	CraftDir      = ".craft"
	WorkflowFile  = "workflow.md"
	SchemaVersion = 3
)

// YAML front matter keys
const (
	keyState         = "state"
	keySchemaVersion = "schema_version"
	keyChecksum      = "checksum"
	keyStartedAt     = "started_at"
	keyHistory       = "history:"
	keyHistoryState  = "- state:"
	keyHistoryAt     = "at:"
	keyHistoryNote   = "note:"
)

// HistoryEntry records a state transition with timestamp and optional note.
type HistoryEntry struct {
	State string
	At    time.Time
	Note  string
}

// Workflow represents a craft workflow.
type Workflow struct {
	State         state.State
	SchemaVersion int
	Checksum      string
	StartedAt     time.Time
	History       []HistoryEntry
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

	w, err := Parse(data)
	if err != nil {
		return nil, err
	}

	// Handle v1 migration in Load (not Parse) since it may need filesystem access
	if w.SchemaVersion < 2 && len(w.History) == 0 {
		w.synthesizeV1History(Path())
	}

	return w, nil
}

// synthesizeV1History creates initial history for v1 workflows being migrated.
func (w *Workflow) synthesizeV1History(filePath string) {
	if w.StartedAt.IsZero() {
		if info, err := os.Stat(filePath); err == nil {
			w.StartedAt = info.ModTime().UTC()
		} else {
			w.StartedAt = time.Now().UTC()
		}
	}
	w.History = []HistoryEntry{{
		State: string(w.State),
		At:    w.StartedAt,
	}}
}

// Parse parses workflow content from bytes. This is a pure function with no side effects.
func Parse(data []byte) (*Workflow, error) {
	content := string(data)

	frontMatter, body, err := extractFrontMatter(content)
	if err != nil {
		return nil, err
	}

	w := &Workflow{
		SchemaVersion: 1, // Default to v1, will be overwritten if present
	}

	// Parse front matter fields and history
	parseFrontMatter(frontMatter, w)

	if !w.State.Valid() {
		return nil, fmt.Errorf("invalid workflow state: %s", w.State)
	}

	// Parse body
	w.Intent, w.Notes = parseBody(body)

	return w, nil
}

// extractFrontMatter splits content into front matter and body.
func extractFrontMatter(content string) (frontMatter, body string, err error) {
	if !strings.HasPrefix(content, "---\n") {
		return "", "", errors.New("invalid workflow file: missing front matter")
	}

	parts := strings.SplitN(content[4:], "\n---\n", 2)
	if len(parts) != 2 {
		return "", "", errors.New("invalid workflow file: malformed front matter")
	}

	return parts[0], parts[1], nil
}

// parseFrontMatter parses YAML front matter into the workflow struct.
func parseFrontMatter(frontMatter string, w *Workflow) {
	lines := strings.Split(frontMatter, "\n")
	inHistory := false
	var currentEntry *HistoryEntry

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Check for history array start
		if trimmed == keyHistory {
			inHistory = true
			continue
		}

		// Handle history entries
		if inHistory {
			if trimmed == "" {
				continue
			}

			// New history entry starts with "- state:"
			if strings.HasPrefix(trimmed, keyHistoryState) {
				if currentEntry != nil {
					w.History = append(w.History, *currentEntry)
				}
				currentEntry = &HistoryEntry{
					State: strings.TrimSpace(strings.TrimPrefix(trimmed, keyHistoryState)),
				}
				continue
			}

			// If line doesn't start with whitespace, we're done with history
			if !strings.HasPrefix(line, "  ") && !strings.HasPrefix(line, "\t") {
				if currentEntry != nil {
					w.History = append(w.History, *currentEntry)
					currentEntry = nil
				}
				inHistory = false
				// Fall through to process this line as a regular key
			} else if currentEntry != nil {
				parseHistoryEntryField(trimmed, currentEntry)
				continue
			}
		}

		// Regular key-value parsing
		if trimmed == "" {
			continue
		}
		kv := strings.SplitN(trimmed, ":", 2)
		if len(kv) != 2 {
			continue
		}
		key := strings.TrimSpace(kv[0])
		value := strings.TrimSpace(kv[1])

		switch key {
		case keyState:
			w.State = state.State(value)
		case keySchemaVersion:
			fmt.Sscanf(value, "%d", &w.SchemaVersion)
		case keyChecksum:
			w.Checksum = value
		case keyStartedAt:
			if t, err := time.Parse(time.RFC3339, value); err == nil {
				w.StartedAt = t
			}
		}
	}

	// Don't forget last history entry
	if currentEntry != nil {
		w.History = append(w.History, *currentEntry)
	}
}

// parseHistoryEntryField parses a single field within a history entry.
func parseHistoryEntryField(line string, entry *HistoryEntry) {
	if strings.HasPrefix(line, keyHistoryAt) {
		timeStr := strings.TrimSpace(strings.TrimPrefix(line, keyHistoryAt))
		if t, err := time.Parse(time.RFC3339, timeStr); err == nil {
			entry.At = t
		}
		return
	}
	if strings.HasPrefix(line, keyHistoryNote) {
		note := strings.TrimSpace(strings.TrimPrefix(line, keyHistoryNote))
		// Remove only the outer quotes (not all quotes like Trim does)
		if len(note) >= 2 && note[0] == '"' && note[len(note)-1] == '"' {
			note = note[1 : len(note)-1]
		}
		// Unescape: order matters - unescape backslashes first, then quotes
		note = strings.ReplaceAll(note, `\\`, "\x00") // temp placeholder
		note = strings.ReplaceAll(note, `\"`, `"`)
		note = strings.ReplaceAll(note, "\x00", `\`)
		entry.Note = note
		return
	}
	if strings.HasPrefix(line, keyState+":") {
		entry.State = strings.TrimSpace(strings.TrimPrefix(line, keyState+":"))
	}
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

	// Migrate schema if needed
	w.migrateSchema()

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

// migrateSchema upgrades older schema versions to the current version.
// Called by Save() after Load() has already synthesized history.
func (w *Workflow) migrateSchema() {
	if w.SchemaVersion >= SchemaVersion {
		return
	}

	// V1 to V2+: Ensure StartedAt and History are set
	if w.SchemaVersion < 2 {
		if w.StartedAt.IsZero() {
			w.StartedAt = time.Now().UTC()
		}

		if len(w.History) == 0 {
			w.History = []HistoryEntry{{
				State: string(w.State),
				At:    w.StartedAt,
			}}
		}
	}

	// V2 to V3: No structural changes, just new shaping state support
	// Existing workflows continue to work - shaping is only for new workflows

	w.SchemaVersion = SchemaVersion
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

// formatHistory returns history formatted for the workflow file.
func (w *Workflow) formatHistory() string {
	if len(w.History) == 0 {
		return ""
	}
	var lines []string
	for _, h := range w.History {
		entry := fmt.Sprintf("  - %s: %s\n    %s %s", keyState, h.State, keyHistoryAt, h.At.Format(time.RFC3339))
		if h.Note != "" {
			// Escape for YAML: backslashes first, then quotes
			escapedNote := strings.ReplaceAll(h.Note, `\`, `\\`)
			escapedNote = strings.ReplaceAll(escapedNote, `"`, `\"`)
			entry += fmt.Sprintf("\n    %s \"%s\"", keyHistoryNote, escapedNote)
		}
		lines = append(lines, entry)
	}
	return keyHistory + "\n" + strings.Join(lines, "\n")
}

// formatWorkflowContent formats the workflow without the checksum field.
func (w *Workflow) formatWorkflowContent() string {
	historySection := w.formatHistory()
	if historySection != "" {
		historySection = "\n" + historySection
	}
	return fmt.Sprintf(`---
%s: %s
%s: %d
%s: %s%s
---

# Intent
%s

## Notes
%s
`, keyState, w.State, keySchemaVersion, w.SchemaVersion, keyStartedAt, w.StartedAt.Format(time.RFC3339), historySection, w.Intent, w.formatNotes())
}

// contentForChecksum returns the content used for checksum computation.
// This excludes the checksum field itself to allow verification.
func (w *Workflow) contentForChecksum() string {
	return w.formatWorkflowContent()
}

// Format returns the workflow as a formatted string with checksum.
func (w *Workflow) Format() string {
	checksumContent := w.formatWorkflowContent()
	checksum := ComputeChecksum([]byte(checksumContent))
	w.Checksum = checksum

	historySection := w.formatHistory()
	if historySection != "" {
		historySection = "\n" + historySection
	}

	return fmt.Sprintf(`---
%s: %s
%s: %d
%s: %s
%s: %s%s
---

# Intent
%s

## Notes
%s
`, keyState, w.State, keySchemaVersion, w.SchemaVersion, keyChecksum, checksum, keyStartedAt, w.StartedAt.Format(time.RFC3339), historySection, w.Intent, w.formatNotes())
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
	now := time.Now().UTC()
	return &Workflow{
		State:         state.Thinking,
		SchemaVersion: SchemaVersion,
		StartedAt:     now,
		History: []HistoryEntry{{
			State: string(state.Thinking),
			At:    now,
		}},
		Intent: intent,
		Notes:  nil,
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
	return w.TransitionWithNote(to, "")
}

// TransitionWithNote validates and performs a state transition with an optional note.
func (w *Workflow) TransitionWithNote(to state.State, note string) error {
	if err := state.ValidateTransition(w.State, to); err != nil {
		return err
	}
	w.State = to
	w.RecordTransition(note)
	return nil
}

// RecordTransition adds a history entry for the current state.
func (w *Workflow) RecordTransition(note string) {
	w.History = append(w.History, HistoryEntry{
		State: string(w.State),
		At:    time.Now().UTC(),
		Note:  note,
	})
}
