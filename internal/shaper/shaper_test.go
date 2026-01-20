package shaper

import (
	"bytes"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"craft/internal/structure"
)

func setupTest(t *testing.T) func() {
	t.Helper()
	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	os.Chdir(tmpDir)
	return func() {
		os.Chdir(origDir)
	}
}

func TestAIShaperAvailable(t *testing.T) {
	cleanup := setupTest(t)
	defer cleanup()

	s := &AIShaper{}

	// Not available without API key
	os.Unsetenv(envAPIKey)
	if s.Available() {
		t.Error("Available() = true, want false (no API key)")
	}

	// Available with API key
	os.Setenv(envAPIKey, "test-key")
	defer os.Unsetenv(envAPIKey)
	if !s.Available() {
		t.Error("Available() = false, want true (API key set)")
	}
}

func TestAIShaperName(t *testing.T) {
	s := &AIShaper{}
	if got := s.Name(); got != NameAI {
		t.Errorf("Name() = %q, want %q", got, NameAI)
	}
}

func TestShapeCLIShaperName(t *testing.T) {
	s := &ShapeCLIShaper{}
	if got := s.Name(); got != NameShapeCLI {
		t.Errorf("Name() = %q, want %q", got, NameShapeCLI)
	}
}

func TestGetBestShaper(t *testing.T) {
	cleanup := setupTest(t)
	defer cleanup()

	// No shapers available
	os.Unsetenv(envAPIKey)
	s := GetBestShaper()
	if s != nil {
		t.Errorf("GetBestShaper() = %v, want nil (no shapers available)", s)
	}

	// AI available
	os.Setenv(envAPIKey, "test-key")
	defer os.Unsetenv(envAPIKey)
	s = GetBestShaper()
	if s == nil || s.Name() != NameAI {
		t.Errorf("GetBestShaper() = %v, want AI shaper", s)
	}
}

// MockHTTPClient for testing AI shaper
type MockHTTPClient struct {
	Response *http.Response
	Err      error
}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return m.Response, m.Err
}

func TestAIShaperShape(t *testing.T) {
	cleanup := setupTest(t)
	defer cleanup()

	os.Setenv(envAPIKey, "test-key")
	defer os.Unsetenv(envAPIKey)

	// Mock response for pitch
	pitchResponse := `# Pitch: Test Feature

## Problem
This is a test problem.

## Solution
This is the solution.

## Scope

### In Scope
- Feature A

### Out of Scope
- Feature B

## Tasks
- [ ] Task 1
- [ ] Task 2`

	// Mock response for cards
	cardsResponse := `===CARD===
# Card: First Card

## Summary
First card summary.

## Tasks
- [ ] Task 1

## Acceptance Criteria
- Done when task 1 is complete
===END===

===CARD===
# Card: Second Card

## Summary
Second card summary.

## Tasks
- [ ] Task 2

## Acceptance Criteria
- Done when task 2 is complete
===END===`

	callCount := 0
	mockClient := &MockHTTPClient{}

	// We need to mock multiple calls
	s := &AIShaper{
		Client: &MockSequentialClient{
			Responses: []*http.Response{
				{
					StatusCode: 200,
					Body:       io.NopCloser(bytes.NewBufferString(`{"choices":[{"message":{"content":"` + escapeJSON(pitchResponse) + `"}}]}`)),
				},
				{
					StatusCode: 200,
					Body:       io.NopCloser(bytes.NewBufferString(`{"choices":[{"message":{"content":"` + escapeJSON(cardsResponse) + `"}}]}`)),
				},
			},
			CallCount: &callCount,
		},
	}

	req := ShapeRequest{
		Intent: "Test feature",
		Notes:  []string{"Note 1"},
	}

	result, err := s.Shape(req)
	if err != nil {
		t.Fatalf("Shape() error = %v", err)
	}

	if result.Shaper != NameAI {
		t.Errorf("Shaper = %q, want %q", result.Shaper, NameAI)
	}

	if result.PitchPath != structure.PitchPath() {
		t.Errorf("PitchPath = %q, want %q", result.PitchPath, structure.PitchPath())
	}

	if len(result.CardPaths) != 2 {
		t.Errorf("CardPaths count = %d, want 2", len(result.CardPaths))
	}

	// Verify files were created
	if !structure.HasPitch() {
		t.Error("Pitch file should exist")
	}

	cards, _ := structure.ListCards()
	if len(cards) != 2 {
		t.Errorf("Cards count = %d, want 2", len(cards))
	}

	_ = mockClient // silence unused warning
}

// MockSequentialClient returns different responses for sequential calls
type MockSequentialClient struct {
	Responses []*http.Response
	CallCount *int
}

func (m *MockSequentialClient) Do(req *http.Request) (*http.Response, error) {
	idx := *m.CallCount
	*m.CallCount++
	if idx < len(m.Responses) {
		return m.Responses[idx], nil
	}
	return nil, io.EOF
}

func escapeJSON(s string) string {
	// Simple JSON string escaping
	s = bytes.NewBuffer([]byte(s)).String()
	var result bytes.Buffer
	for _, c := range s {
		switch c {
		case '"':
			result.WriteString(`\"`)
		case '\\':
			result.WriteString(`\\`)
		case '\n':
			result.WriteString(`\n`)
		case '\r':
			result.WriteString(`\r`)
		case '\t':
			result.WriteString(`\t`)
		default:
			result.WriteRune(c)
		}
	}
	return result.String()
}

func TestParseAndWriteCards(t *testing.T) {
	cleanup := setupTest(t)
	defer cleanup()

	structure.EnsureStructureDir()

	content := `===CARD===
# Card: Feature Implementation

## Summary
Implement the feature.

## Tasks
- [ ] Task 1

## Acceptance Criteria
- Feature works
===END===

===CARD===
# Card: Add Tests

## Summary
Add test coverage.

## Tasks
- [ ] Write tests

## Acceptance Criteria
- Tests pass
===END===`

	paths, err := parseAndWriteCards(content)
	if err != nil {
		t.Fatalf("parseAndWriteCards() error = %v", err)
	}

	if len(paths) != 2 {
		t.Fatalf("parseAndWriteCards() = %d paths, want 2", len(paths))
	}

	// Check filenames are numbered correctly
	if !filepath.IsAbs(paths[0]) {
		expected := filepath.Join(structure.CardsDirPath(), "01-feature-implementation.md")
		if paths[0] != expected {
			t.Errorf("First card path = %q, want %q", paths[0], expected)
		}
	}

	// Verify files exist
	for _, p := range paths {
		if _, err := os.Stat(p); os.IsNotExist(err) {
			t.Errorf("Card file %q should exist", p)
		}
	}
}

func TestBuildPitchPrompt(t *testing.T) {
	req := ShapeRequest{
		Intent: "Add rate limiting",
		Notes:  []string{"Token bucket", "Per-user limits"},
	}

	prompt := buildPitchPrompt(req)

	if !bytes.Contains([]byte(prompt), []byte("Add rate limiting")) {
		t.Error("Prompt should contain intent")
	}

	if !bytes.Contains([]byte(prompt), []byte("Token bucket")) {
		t.Error("Prompt should contain notes")
	}
}

func TestBuildCardsPrompt(t *testing.T) {
	req := ShapeRequest{
		Intent: "Test feature",
	}
	pitchContent := "# Pitch: Test\n\n## Problem\nTest problem"

	prompt := buildCardsPrompt(req, pitchContent)

	if !bytes.Contains([]byte(prompt), []byte("Test problem")) {
		t.Error("Prompt should contain pitch content")
	}

	if !bytes.Contains([]byte(prompt), []byte("===CARD===")) {
		t.Error("Prompt should describe card format")
	}
}
