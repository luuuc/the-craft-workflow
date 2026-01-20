package main

import (
	"os"
	"testing"
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

func TestRunHelp(t *testing.T) {
	tests := []struct {
		args []string
	}{
		{[]string{}},
		{[]string{"--help"}},
		{[]string{"-h"}},
		{[]string{"help"}},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			code := run(tt.args)
			if code != 0 {
				t.Errorf("run(%v) = %d, want 0", tt.args, code)
			}
		})
	}
}

func TestRunVersion(t *testing.T) {
	tests := [][]string{
		{"--version"},
		{"-v"},
	}

	for _, args := range tests {
		code := run(args)
		if code != 0 {
			t.Errorf("run(%v) = %d, want 0", args, code)
		}
	}
}

func TestRunUnknownCommand(t *testing.T) {
	code := run([]string{"unknown"})
	if code != 1 {
		t.Errorf("run(unknown) = %d, want 1", code)
	}
}

func TestRunCommands(t *testing.T) {
	cleanup := setupTest(t)
	defer cleanup()

	// Test start
	code := run([]string{"start", "Test intent"})
	if code != 0 {
		t.Errorf("start = %d, want 0", code)
	}

	// Test think
	code = run([]string{"think"})
	if code != 0 {
		t.Errorf("think = %d, want 0", code)
	}

	// Test status
	code = run([]string{"status"})
	if code != 0 {
		t.Errorf("status = %d, want 0", code)
	}

	// Test reject
	code = run([]string{"reject", "Need more thought"})
	if code != 0 {
		t.Errorf("reject = %d, want 0", code)
	}

	// Test accept
	code = run([]string{"accept"})
	if code != 0 {
		t.Errorf("accept = %d, want 0", code)
	}

	// Test ship
	code = run([]string{"ship"})
	if code != 0 {
		t.Errorf("ship = %d, want 0", code)
	}

	// Test reset
	code = run([]string{"reset", "--force"})
	if code != 0 {
		t.Errorf("reset = %d, want 0", code)
	}
}
