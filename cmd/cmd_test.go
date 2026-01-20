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
	if w.State != "shaping" {
		t.Errorf("State = %s, want shaping", w.State)
	}
}

func TestAcceptSkipShaping(t *testing.T) {
	cleanup := setupTest(t)
	defer cleanup()

	Start([]string{"Test"})
	code := Accept([]string{"--skip-shaping"})
	if code != 0 {
		t.Errorf("Accept(--skip-shaping) = %d, want 0", code)
	}

	w, _ := workflow.Load()
	if w.State != "building" {
		t.Errorf("State = %s, want building", w.State)
	}
}

func TestAcceptSkipShapingWithNote(t *testing.T) {
	cleanup := setupTest(t)
	defer cleanup()

	Start([]string{"Test"})
	code := Accept([]string{"--skip-shaping", "Quick fix"})
	if code != 0 {
		t.Errorf("Accept(--skip-shaping, note) = %d, want 0", code)
	}

	w, _ := workflow.Load()
	if w.State != "building" {
		t.Errorf("State = %s, want building", w.State)
	}
	if len(w.Notes) != 1 || w.Notes[0] != "Quick fix" {
		t.Errorf("Notes = %v, want [Quick fix]", w.Notes)
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

func TestAcceptFromShaping(t *testing.T) {
	cleanup := setupTest(t)
	defer cleanup()

	Start([]string{"Test"})
	Accept(nil) // Now in shaping
	code := Accept(nil)
	if code != 1 {
		t.Errorf("Accept() from shaping = %d, want 1", code)
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

func TestRejectFromShaping(t *testing.T) {
	cleanup := setupTest(t)
	defer cleanup()

	Start([]string{"Test"})
	Accept(nil) // Now in shaping
	code := Reject([]string{"Too late"})
	if code != 1 {
		t.Errorf("Reject() from shaping = %d, want 1", code)
	}
}

func TestShipFromBuilding(t *testing.T) {
	cleanup := setupTest(t)
	defer cleanup()

	Start([]string{"Test"})
	Accept([]string{"--skip-shaping"})
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
	Accept([]string{"--skip-shaping"})
	Ship(nil)
	code := Ship(nil)
	if code != 1 {
		t.Errorf("Ship() from shipped = %d, want 1", code)
	}
}

func TestShipFromShaping(t *testing.T) {
	cleanup := setupTest(t)
	defer cleanup()

	Start([]string{"Test"})
	Accept(nil) // Now in shaping
	code := Ship(nil)
	if code != 1 {
		t.Errorf("Ship() from shaping = %d, want 1", code)
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

func TestFullWorkflowWithSkipShaping(t *testing.T) {
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

	// Think and accept (skip shaping)
	if code := Think(nil); code != 0 {
		t.Fatalf("Think() = %d", code)
	}
	if code := Accept([]string{"--skip-shaping", "Decided on token bucket algorithm"}); code != 0 {
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

func TestInitNoFlags(t *testing.T) {
	cleanup := setupTest(t)
	defer cleanup()

	code := Init(nil)
	if code != 1 {
		t.Errorf("Init() with no flags = %d, want 1", code)
	}
}

func TestInitClaude(t *testing.T) {
	cleanup := setupTest(t)
	defer cleanup()

	code := Init([]string{"--claude"})
	if code != 0 {
		t.Errorf("Init(--claude) = %d, want 0", code)
	}

	if _, err := os.Stat("CLAUDE.md"); os.IsNotExist(err) {
		t.Error("CLAUDE.md should exist after init --claude")
	}
	if _, err := os.Stat(".claude/commands/craft.md"); os.IsNotExist(err) {
		t.Error(".claude/commands/craft.md should exist after init --claude")
	}
}

func TestInitCursor(t *testing.T) {
	cleanup := setupTest(t)
	defer cleanup()

	code := Init([]string{"--cursor"})
	if code != 0 {
		t.Errorf("Init(--cursor) = %d, want 0", code)
	}

	if _, err := os.Stat(".cursorrules"); os.IsNotExist(err) {
		t.Error(".cursorrules should exist after init --cursor")
	}
}

func TestInitAll(t *testing.T) {
	cleanup := setupTest(t)
	defer cleanup()

	code := Init([]string{"--all"})
	if code != 0 {
		t.Errorf("Init(--all) = %d, want 0", code)
	}

	if _, err := os.Stat("CLAUDE.md"); os.IsNotExist(err) {
		t.Error("CLAUDE.md should exist after init --all")
	}
	if _, err := os.Stat(".claude/commands/craft.md"); os.IsNotExist(err) {
		t.Error(".claude/commands/craft.md should exist after init --all")
	}
	if _, err := os.Stat(".cursorrules"); os.IsNotExist(err) {
		t.Error(".cursorrules should exist after init --all")
	}
}

func TestInitSkipsExisting(t *testing.T) {
	cleanup := setupTest(t)
	defer cleanup()

	// Create CLAUDE.md with custom content
	os.WriteFile("CLAUDE.md", []byte("custom content"), 0644)

	code := Init([]string{"--claude"})
	if code != 0 {
		t.Errorf("Init(--claude) = %d, want 0", code)
	}

	// Verify custom content preserved
	content, _ := os.ReadFile("CLAUDE.md")
	if string(content) != "custom content" {
		t.Error("Init should not overwrite existing CLAUDE.md")
	}

	// But .claude/commands/craft.md should be created
	if _, err := os.Stat(".claude/commands/craft.md"); os.IsNotExist(err) {
		t.Error(".claude/commands/craft.md should exist")
	}
}

// Tests for new shaping phase commands

func TestShapeFromShaping(t *testing.T) {
	cleanup := setupTest(t)
	defer cleanup()

	Start([]string{"Test"})
	Accept(nil) // Now in shaping

	code := Shape(nil)
	if code != 0 {
		t.Errorf("Shape() = %d, want 0", code)
	}
}

func TestShapeFromThinking(t *testing.T) {
	cleanup := setupTest(t)
	defer cleanup()

	Start([]string{"Test"})
	code := Shape(nil)
	if code != 1 {
		t.Errorf("Shape() from thinking = %d, want 1", code)
	}
}

func TestApproveFromShaping(t *testing.T) {
	cleanup := setupTest(t)
	defer cleanup()

	Start([]string{"Test"})
	Accept(nil) // Now in shaping

	// Create pitch file
	os.MkdirAll(".craft", 0755)
	os.WriteFile(".craft/pitch.md", []byte("# Pitch"), 0644)

	code := Approve(nil)
	if code != 0 {
		t.Errorf("Approve() = %d, want 0", code)
	}

	w, _ := workflow.Load()
	if w.State != "building" {
		t.Errorf("State = %s, want building", w.State)
	}
}

func TestApproveNoPitch(t *testing.T) {
	cleanup := setupTest(t)
	defer cleanup()

	Start([]string{"Test"})
	Accept(nil) // Now in shaping

	code := Approve(nil)
	if code != 1 {
		t.Errorf("Approve() without pitch = %d, want 1", code)
	}
}

func TestApproveFromThinking(t *testing.T) {
	cleanup := setupTest(t)
	defer cleanup()

	Start([]string{"Test"})
	code := Approve(nil)
	if code != 1 {
		t.Errorf("Approve() from thinking = %d, want 1", code)
	}
}

func TestReviseFromShaping(t *testing.T) {
	cleanup := setupTest(t)
	defer cleanup()

	Start([]string{"Test"})
	Accept(nil) // Now in shaping

	code := Revise([]string{"Need to rethink this"})
	if code != 0 {
		t.Errorf("Revise() = %d, want 0", code)
	}

	w, _ := workflow.Load()
	if w.State != "shaping" {
		t.Errorf("State = %s, want shaping", w.State)
	}
	if len(w.Notes) != 1 {
		t.Errorf("Notes count = %d, want 1", len(w.Notes))
	}
	if !strings.Contains(w.Notes[0], "[revise]") {
		t.Errorf("Note should contain [revise] prefix, got %q", w.Notes[0])
	}
}

func TestReviseNoNote(t *testing.T) {
	cleanup := setupTest(t)
	defer cleanup()

	Start([]string{"Test"})
	Accept(nil) // Now in shaping

	code := Revise(nil)
	if code != 1 {
		t.Errorf("Revise() without note = %d, want 1", code)
	}
}

func TestReviseFromThinking(t *testing.T) {
	cleanup := setupTest(t)
	defer cleanup()

	Start([]string{"Test"})
	code := Revise([]string{"Some note"})
	if code != 1 {
		t.Errorf("Revise() from thinking = %d, want 1", code)
	}
}

func TestFullWorkflowWithShaping(t *testing.T) {
	cleanup := setupTest(t)
	defer cleanup()

	// Start
	if code := Start([]string{"Add rate limiting"}); code != 0 {
		t.Fatalf("Start() = %d", code)
	}

	// Think and accept (goes to shaping)
	if code := Accept([]string{"Token bucket algorithm"}); code != 0 {
		t.Fatalf("Accept() = %d", code)
	}

	w, _ := workflow.Load()
	if w.State != "shaping" {
		t.Fatalf("State after accept = %s, want shaping", w.State)
	}

	// Create pitch manually (simulating manual shaping)
	os.WriteFile(".craft/pitch.md", []byte("# Pitch: Rate Limiting"), 0644)

	// Approve (goes to building)
	if code := Approve(nil); code != 0 {
		t.Fatalf("Approve() = %d", code)
	}

	w, _ = workflow.Load()
	if w.State != "building" {
		t.Fatalf("State after approve = %s, want building", w.State)
	}

	// Ship
	if code := Ship(nil); code != 0 {
		t.Fatalf("Ship() = %d", code)
	}

	w, _ = workflow.Load()
	if w.State != "shipped" {
		t.Errorf("Final state = %s, want shipped", w.State)
	}
}
