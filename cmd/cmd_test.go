package cmd

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"craft/internal/workflow"
)

func setupTest(t *testing.T) (cleanup func()) {
	t.Helper()
	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	os.Chdir(tmpDir)
	return func() {
		os.Chdir(origDir)
	}
}

func TestStartSuccess(t *testing.T) {
	cleanup := setupTest(t)
	defer cleanup()

	code := Start([]string{"Add rate limiting"})
	if code != 0 {
		t.Errorf("Start() = %d, want 0", code)
	}

	if !workflow.Exists() {
		t.Error("Workflow file should exist after start")
	}

	w, err := workflow.Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if w.Intent != "Add rate limiting" {
		t.Errorf("Intent = %q, want %q", w.Intent, "Add rate limiting")
	}
}

func TestStartWithQuotes(t *testing.T) {
	cleanup := setupTest(t)
	defer cleanup()

	code := Start([]string{`"Quoted intent"`})
	if code != 0 {
		t.Errorf("Start() = %d, want 0", code)
	}

	w, _ := workflow.Load()
	if w.Intent != "Quoted intent" {
		t.Errorf("Intent = %q, want %q", w.Intent, "Quoted intent")
	}
}

func TestStartNoIntent(t *testing.T) {
	cleanup := setupTest(t)
	defer cleanup()

	code := Start([]string{})
	if code != 1 {
		t.Errorf("Start() = %d, want 1", code)
	}
}

func TestStartEmptyIntent(t *testing.T) {
	cleanup := setupTest(t)
	defer cleanup()

	code := Start([]string{""})
	if code != 1 {
		t.Errorf("Start() = %d, want 1", code)
	}
}

func TestStartExisting(t *testing.T) {
	cleanup := setupTest(t)
	defer cleanup()

	Start([]string{"First"})
	code := Start([]string{"Second"})
	if code != 1 {
		t.Errorf("Start() = %d, want 1 for existing workflow", code)
	}
}

func TestThinkNoWorkflow(t *testing.T) {
	cleanup := setupTest(t)
	defer cleanup()

	code := Think(nil)
	if code != 1 {
		t.Errorf("Think() = %d, want 1", code)
	}
}

func TestThinkWithWorkflow(t *testing.T) {
	cleanup := setupTest(t)
	defer cleanup()

	Start([]string{"Test intent"})
	code := Think(nil)
	if code != 0 {
		t.Errorf("Think() = %d, want 0", code)
	}
}

func TestAcceptFromThinking(t *testing.T) {
	cleanup := setupTest(t)
	defer cleanup()

	Start([]string{"Test"})
	code := Accept(nil)
	if code != 0 {
		t.Errorf("Accept() = %d, want 0", code)
	}

	w, _ := workflow.Load()
	if w.State != "building" {
		t.Errorf("State = %s, want building", w.State)
	}
}

func TestAcceptWithNote(t *testing.T) {
	cleanup := setupTest(t)
	defer cleanup()

	Start([]string{"Test"})
	Accept([]string{"Decided to use token bucket"})

	w, _ := workflow.Load()
	if len(w.Notes) != 1 || w.Notes[0] != "Decided to use token bucket" {
		t.Errorf("Notes = %v, want [Decided to use token bucket]", w.Notes)
	}
}

func TestAcceptFromBuilding(t *testing.T) {
	cleanup := setupTest(t)
	defer cleanup()

	Start([]string{"Test"})
	Accept(nil)
	code := Accept(nil)
	if code != 1 {
		t.Errorf("Accept() from building = %d, want 1", code)
	}
}

func TestRejectFromThinking(t *testing.T) {
	cleanup := setupTest(t)
	defer cleanup()

	Start([]string{"Test"})
	code := Reject([]string{"Need more thought"})
	if code != 0 {
		t.Errorf("Reject() = %d, want 0", code)
	}

	w, _ := workflow.Load()
	if w.State != "thinking" {
		t.Errorf("State = %s, want thinking", w.State)
	}
	if len(w.Notes) != 1 {
		t.Errorf("Notes count = %d, want 1", len(w.Notes))
	}
}

func TestRejectFromBuilding(t *testing.T) {
	cleanup := setupTest(t)
	defer cleanup()

	Start([]string{"Test"})
	Accept(nil)
	code := Reject([]string{"Too late"})
	if code != 1 {
		t.Errorf("Reject() from building = %d, want 1", code)
	}
}

func TestShipFromBuilding(t *testing.T) {
	cleanup := setupTest(t)
	defer cleanup()

	Start([]string{"Test"})
	Accept(nil)
	code := Ship(nil)
	if code != 0 {
		t.Errorf("Ship() = %d, want 0", code)
	}

	w, _ := workflow.Load()
	if w.State != "shipped" {
		t.Errorf("State = %s, want shipped", w.State)
	}
}

func TestShipFromThinking(t *testing.T) {
	cleanup := setupTest(t)
	defer cleanup()

	Start([]string{"Test"})
	code := Ship(nil)
	if code != 1 {
		t.Errorf("Ship() from thinking = %d, want 1", code)
	}
}

func TestShipFromShipped(t *testing.T) {
	cleanup := setupTest(t)
	defer cleanup()

	Start([]string{"Test"})
	Accept(nil)
	Ship(nil)
	code := Ship(nil)
	if code != 1 {
		t.Errorf("Ship() from shipped = %d, want 1", code)
	}
}

func TestStatusNoWorkflow(t *testing.T) {
	cleanup := setupTest(t)
	defer cleanup()

	code := Status(nil)
	if code != 0 {
		t.Errorf("Status() = %d, want 0", code)
	}
}

func TestStatusWithWorkflow(t *testing.T) {
	cleanup := setupTest(t)
	defer cleanup()

	Start([]string{"Test intent"})
	code := Status(nil)
	if code != 0 {
		t.Errorf("Status() = %d, want 0", code)
	}
}

func TestResetForce(t *testing.T) {
	cleanup := setupTest(t)
	defer cleanup()

	Start([]string{"Test"})
	code := Reset([]string{"--force"})
	if code != 0 {
		t.Errorf("Reset(--force) = %d, want 0", code)
	}

	if workflow.Exists() {
		t.Error("Workflow should not exist after reset")
	}
}

func TestResetNoWorkflow(t *testing.T) {
	cleanup := setupTest(t)
	defer cleanup()

	code := Reset([]string{"--force"})
	if code != 0 {
		t.Errorf("Reset() with no workflow = %d, want 0", code)
	}
}

func TestFullWorkflow(t *testing.T) {
	cleanup := setupTest(t)
	defer cleanup()

	// Start
	if code := Start([]string{"Add rate limiting to API"}); code != 0 {
		t.Fatalf("Start() = %d", code)
	}

	// Think and reject
	if code := Think(nil); code != 0 {
		t.Fatalf("Think() = %d", code)
	}
	if code := Reject([]string{"Need to consider rate limit headers"}); code != 0 {
		t.Fatalf("Reject() = %d", code)
	}

	// Think and accept
	if code := Think(nil); code != 0 {
		t.Fatalf("Think() = %d", code)
	}
	if code := Accept([]string{"Decided on token bucket algorithm"}); code != 0 {
		t.Fatalf("Accept() = %d", code)
	}

	// Ship
	if code := Ship(nil); code != 0 {
		t.Fatalf("Ship() = %d", code)
	}

	// Verify final state
	w, _ := workflow.Load()
	if w.State != "shipped" {
		t.Errorf("Final state = %s, want shipped", w.State)
	}
	if len(w.Notes) != 2 {
		t.Errorf("Notes count = %d, want 2", len(w.Notes))
	}

	// Reset
	if code := Reset([]string{"--force"}); code != 0 {
		t.Fatalf("Reset() = %d", code)
	}
	if workflow.Exists() {
		t.Error("Workflow should not exist after reset")
	}
}

func TestConfirm(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"y\n", true},
		{"Y\n", true},
		{"yes\n", true},
		{"YES\n", true},
		{"Yes\n", true},
		{"n\n", false},
		{"N\n", false},
		{"no\n", false},
		{"\n", false},
		{"maybe\n", false},
		{"", false}, // EOF
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			reader := strings.NewReader(tt.input)
			got := confirm(reader)
			if got != tt.want {
				t.Errorf("confirm(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestResetInteractiveYes(t *testing.T) {
	cleanup := setupTest(t)
	defer cleanup()

	Start([]string{"Test"})

	// Simulate user typing "y"
	oldStdin := stdinReader
	stdinReader = strings.NewReader("y\n")
	defer func() { stdinReader = oldStdin }()

	code := Reset(nil)
	if code != 0 {
		t.Errorf("Reset() with 'y' = %d, want 0", code)
	}

	if workflow.Exists() {
		t.Error("Workflow should not exist after confirmed reset")
	}
}

func TestResetInteractiveNo(t *testing.T) {
	cleanup := setupTest(t)
	defer cleanup()

	Start([]string{"Test"})

	// Simulate user typing "n"
	oldStdin := stdinReader
	stdinReader = strings.NewReader("n\n")
	defer func() { stdinReader = oldStdin }()

	code := Reset(nil)
	if code != 0 {
		t.Errorf("Reset() with 'n' = %d, want 0", code)
	}

	if !workflow.Exists() {
		t.Error("Workflow should still exist after cancelled reset")
	}
}

func TestStatusChecksumWarning(t *testing.T) {
	cleanup := setupTest(t)
	defer cleanup()

	Start([]string{"Test checksum"})

	// Tamper with the checksum by writing directly to the file
	content := `---
state: thinking
schema_version: 1
checksum: 00000000
---

# Intent
Test checksum

## Notes
(none)
`
	os.WriteFile(".craft/workflow.md", []byte(content), 0644)

	// Capture stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	code := Status(nil)

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	if code != 0 {
		t.Errorf("Status() = %d, want 0", code)
	}

	if !strings.Contains(output, "Warning") {
		t.Error("Status() should display checksum warning for tampered file")
	}
}
