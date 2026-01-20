package workflow

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

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
}

func TestWorkflowFormat(t *testing.T) {
	w := New("Test intent")
	formatted := w.Format()

	if !strings.Contains(formatted, "state: thinking") {
		t.Error("Format() missing state")
	}
	if !strings.Contains(formatted, "schema_version: 1") {
		t.Error("Format() missing schema_version")
	}
	if !strings.Contains(formatted, "checksum:") {
		t.Error("Format() missing checksum")
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
