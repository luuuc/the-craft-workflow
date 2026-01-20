package workflow

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"craft/internal/state"
)

func TestNew(t *testing.T) {
	w := New("Add rate limiting")
	if w.State != state.Thinking {
		t.Errorf("New().State = %v, want %v", w.State, state.Thinking)
	}
	if w.Intent != "Add rate limiting" {
		t.Errorf("New().Intent = %v, want %q", w.Intent, "Add rate limiting")
	}
	if w.SchemaVersion != SchemaVersion {
		t.Errorf("New().SchemaVersion = %v, want %v", w.SchemaVersion, SchemaVersion)
	}
	if w.StartedAt.IsZero() {
		t.Error("New().StartedAt should not be zero")
	}
	if len(w.History) != 1 {
		t.Errorf("New().History len = %d, want 1", len(w.History))
	}
	if w.History[0].State != string(state.Thinking) {
		t.Errorf("New().History[0].State = %v, want %v", w.History[0].State, state.Thinking)
	}
}

func TestWorkflowFormat(t *testing.T) {
	w := New("Test intent")
	formatted := w.Format()

	if !strings.Contains(formatted, "state: thinking") {
		t.Error("Format() missing state")
	}
	if !strings.Contains(formatted, "schema_version: 3") {
		t.Error("Format() missing schema_version")
	}
	if !strings.Contains(formatted, "checksum:") {
		t.Error("Format() missing checksum")
	}
	if !strings.Contains(formatted, "started_at:") {
		t.Error("Format() missing started_at")
	}
	if !strings.Contains(formatted, "history:") {
		t.Error("Format() missing history")
	}
	if !strings.Contains(formatted, "# Intent") {
		t.Error("Format() missing Intent header")
	}
	if !strings.Contains(formatted, "Test intent") {
		t.Error("Format() missing intent content")
	}
}

func TestWorkflowRoundTrip(t *testing.T) {
	w := New("Round trip test")
	w.AddNote("First note")
	w.AddNote("Second note")

	formatted := w.Format()
	parsed, err := Parse([]byte(formatted))
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	if parsed.State != w.State {
		t.Errorf("State = %v, want %v", parsed.State, w.State)
	}
	if parsed.Intent != w.Intent {
		t.Errorf("Intent = %v, want %v", parsed.Intent, w.Intent)
	}
	if len(parsed.Notes) != len(w.Notes) {
		t.Errorf("Notes len = %v, want %v", len(parsed.Notes), len(w.Notes))
	}
	for i, note := range parsed.Notes {
		if note != w.Notes[i] {
			t.Errorf("Notes[%d] = %v, want %v", i, note, w.Notes[i])
		}
	}
}

func TestComputeChecksum(t *testing.T) {
	content := []byte("test content")
	checksum := ComputeChecksum(content)

	if len(checksum) != 8 {
		t.Errorf("Checksum length = %d, want 8", len(checksum))
	}

	// Same content should produce same checksum
	checksum2 := ComputeChecksum(content)
	if checksum != checksum2 {
		t.Errorf("Checksum mismatch for same content: %s != %s", checksum, checksum2)
	}

	// Different content should produce different checksum
	checksum3 := ComputeChecksum([]byte("different content"))
	if checksum == checksum3 {
		t.Error("Different content produced same checksum")
	}
}

func TestValidateChecksum(t *testing.T) {
	w := New("Test")
	w.Format() // This sets the checksum

	if err := w.ValidateChecksum(); err != nil {
		t.Errorf("ValidateChecksum() error = %v", err)
	}

	// Tamper with the checksum
	w.Checksum = "00000000"
	if err := w.ValidateChecksum(); err == nil {
		t.Error("ValidateChecksum() should fail with wrong checksum")
	}
}

func TestTransition(t *testing.T) {
	tests := []struct {
		name    string
		from    state.State
		to      state.State
		wantErr bool
	}{
		{"thinking to building", state.Thinking, state.Building, false},
		{"building to shipped", state.Building, state.Shipped, false},
		{"thinking to shipped", state.Thinking, state.Shipped, true},
		{"building to thinking", state.Building, state.Thinking, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &Workflow{State: tt.from}
			err := w.Transition(tt.to)
			if (err != nil) != tt.wantErr {
				t.Errorf("Transition() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil && w.State != tt.to {
				t.Errorf("State = %v, want %v", w.State, tt.to)
			}
		})
	}
}

func TestAddNote(t *testing.T) {
	w := New("Test")

	w.AddNote("Note 1")
	if len(w.Notes) != 1 || w.Notes[0] != "Note 1" {
		t.Errorf("AddNote() = %v, want [Note 1]", w.Notes)
	}

	w.AddNote("  Note 2  ")
	if len(w.Notes) != 2 || w.Notes[1] != "Note 2" {
		t.Errorf("AddNote() with spaces = %v, want [Note 1, Note 2]", w.Notes)
	}

	w.AddNote(`"Quoted note"`)
	if len(w.Notes) != 3 || w.Notes[2] != "Quoted note" {
		t.Errorf("AddNote() with quotes = %v", w.Notes)
	}

	w.AddNote("")
	if len(w.Notes) != 3 {
		t.Error("AddNote() should ignore empty notes")
	}
}

func TestSaveAndLoad(t *testing.T) {
	// Use temp directory
	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(origDir)

	w := New("Save and load test")
	w.AddNote("Test note")

	if err := w.Save(); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	// Check file exists
	if _, err := os.Stat(filepath.Join(CraftDir, WorkflowFile)); err != nil {
		t.Errorf("Workflow file not created: %v", err)
	}

	loaded, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if loaded.Intent != w.Intent {
		t.Errorf("Loaded Intent = %v, want %v", loaded.Intent, w.Intent)
	}
	if len(loaded.Notes) != len(w.Notes) {
		t.Errorf("Loaded Notes = %v, want %v", loaded.Notes, w.Notes)
	}
}

func TestDelete(t *testing.T) {
	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(origDir)

	w := New("Delete test")
	w.Save()

	if err := Delete(); err != nil {
		t.Errorf("Delete() error = %v", err)
	}

	if Exists() {
		t.Error("Workflow should not exist after Delete()")
	}
}

func TestLoadMissing(t *testing.T) {
	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(origDir)

	_, err := Load()
	if err == nil {
		t.Error("Load() should error when no workflow exists")
	}
}

func TestParseInvalid(t *testing.T) {
	tests := []struct {
		name    string
		content string
	}{
		{"no front matter", "# Intent\nTest"},
		{"malformed front matter", "---\nstate: thinking\n# Intent\nTest"},
		{"invalid state", "---\nstate: invalid\n---\n# Intent\nTest"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Parse([]byte(tt.content))
			if err == nil {
				t.Error("Parse() should error for invalid content")
			}
		})
	}
}

func TestHistoryRoundTrip(t *testing.T) {
	w := New("History test")
	w.TransitionWithNote(state.Building, "Ready to build")

	formatted := w.Format()
	parsed, err := Parse([]byte(formatted))
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	if len(parsed.History) != 2 {
		t.Fatalf("History len = %d, want 2", len(parsed.History))
	}

	if parsed.History[0].State != "thinking" {
		t.Errorf("History[0].State = %v, want thinking", parsed.History[0].State)
	}
	if parsed.History[1].State != "building" {
		t.Errorf("History[1].State = %v, want building", parsed.History[1].State)
	}
	if parsed.History[1].Note != "Ready to build" {
		t.Errorf("History[1].Note = %v, want 'Ready to build'", parsed.History[1].Note)
	}
}

func TestTransitionWithNote(t *testing.T) {
	w := New("Test")
	initialHistoryLen := len(w.History)

	err := w.TransitionWithNote(state.Building, "My note")
	if err != nil {
		t.Fatalf("TransitionWithNote() error = %v", err)
	}

	if w.State != state.Building {
		t.Errorf("State = %v, want building", w.State)
	}

	if len(w.History) != initialHistoryLen+1 {
		t.Errorf("History len = %d, want %d", len(w.History), initialHistoryLen+1)
	}

	lastEntry := w.History[len(w.History)-1]
	if lastEntry.State != "building" {
		t.Errorf("Last history entry state = %v, want building", lastEntry.State)
	}
	if lastEntry.Note != "My note" {
		t.Errorf("Last history entry note = %v, want 'My note'", lastEntry.Note)
	}
}

func TestRecordTransition(t *testing.T) {
	w := New("Test")
	initialLen := len(w.History)

	w.RecordTransition("A note")

	if len(w.History) != initialLen+1 {
		t.Errorf("History len = %d, want %d", len(w.History), initialLen+1)
	}

	lastEntry := w.History[len(w.History)-1]
	if lastEntry.Note != "A note" {
		t.Errorf("Note = %v, want 'A note'", lastEntry.Note)
	}
	if lastEntry.At.IsZero() {
		t.Error("At should not be zero")
	}
}

func TestMigrateV1ToLatest(t *testing.T) {
	// Create a v1 workflow content
	v1Content := `---
state: thinking
schema_version: 1
checksum: abcd1234
---

# Intent
V1 workflow

## Notes
- A note
`

	parsed, err := Parse([]byte(v1Content))
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	if parsed.SchemaVersion != 1 {
		t.Errorf("Parsed schema_version = %d, want 1", parsed.SchemaVersion)
	}

	// Use temp directory for save
	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(origDir)

	// Save should migrate to latest version
	if err := parsed.Save(); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	if parsed.SchemaVersion != SchemaVersion {
		t.Errorf("After migration schema_version = %d, want %d", parsed.SchemaVersion, SchemaVersion)
	}

	if len(parsed.History) == 0 {
		t.Error("Migration should synthesize history")
	}

	if parsed.StartedAt.IsZero() {
		t.Error("Migration should set started_at")
	}
}

func TestParseHistoryWithNote(t *testing.T) {
	content := `---
state: building
schema_version: 2
checksum: abcd1234
started_at: 2024-01-15T10:30:00Z
history:
  - state: thinking
    at: 2024-01-15T10:30:00Z
  - state: building
    at: 2024-01-15T14:20:00Z
    note: "Decided on token bucket"
---

# Intent
Add rate limiting

## Notes
- Note 1
`

	w, err := Parse([]byte(content))
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	expectedTime, _ := time.Parse(time.RFC3339, "2024-01-15T10:30:00Z")
	if !w.StartedAt.Equal(expectedTime) {
		t.Errorf("StartedAt = %v, want %v", w.StartedAt, expectedTime)
	}

	if len(w.History) != 2 {
		t.Fatalf("History len = %d, want 2", len(w.History))
	}

	if w.History[0].State != "thinking" {
		t.Errorf("History[0].State = %v, want thinking", w.History[0].State)
	}

	if w.History[1].Note != "Decided on token bucket" {
		t.Errorf("History[1].Note = %v, want 'Decided on token bucket'", w.History[1].Note)
	}
}

func TestQuoteEscapingInNotes(t *testing.T) {
	w := New("Quote test")
	noteWithQuotes := `Note with "quotes" inside`
	w.TransitionWithNote(state.Building, noteWithQuotes)

	formatted := w.Format()

	// Verify quotes are escaped in the formatted output
	if !strings.Contains(formatted, `\"quotes\"`) {
		t.Errorf("Format() should escape quotes in notes, got:\n%s", formatted)
	}

	// Verify round-trip preserves the note
	parsed, err := Parse([]byte(formatted))
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	if len(parsed.History) < 2 {
		t.Fatalf("History len = %d, want at least 2", len(parsed.History))
	}

	if parsed.History[1].Note != noteWithQuotes {
		t.Errorf("Round-trip note = %q, want %q", parsed.History[1].Note, noteWithQuotes)
	}
}

func TestQuoteEscapingEdgeCases(t *testing.T) {
	tests := []struct {
		name string
		note string
	}{
		{"single quote", `He said "hello"`},
		{"multiple quotes", `"start" and "end"`},
		{"adjacent quotes", `""double""`},
		{"escaped already", `already \"escaped\"`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := New("Test")
			w.TransitionWithNote(state.Building, tt.note)

			formatted := w.Format()
			parsed, err := Parse([]byte(formatted))
			if err != nil {
				t.Fatalf("Parse() error = %v", err)
			}

			lastEntry := parsed.History[len(parsed.History)-1]
			if lastEntry.Note != tt.note {
				t.Errorf("Round-trip note = %q, want %q", lastEntry.Note, tt.note)
			}
		})
	}
}
