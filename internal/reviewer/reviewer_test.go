package reviewer

import (
	"bytes"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"
)

func TestNullReviewer_Name(t *testing.T) {
	r := &NullReviewer{}
	if r.Name() != "None" {
		t.Errorf("expected Name() = 'None', got '%s'", r.Name())
	}
}

func TestNullReviewer_Available(t *testing.T) {
	r := &NullReviewer{}
	if !r.Available() {
		t.Error("NullReviewer should always be available")
	}
}

func TestNullReviewer_Review(t *testing.T) {
	r := &NullReviewer{}
	req := ReviewRequest{Intent: "Test intent"}
	resp, err := r.Review(req)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if resp.Reviewer != "None" {
		t.Errorf("expected Reviewer = 'None', got '%s'", resp.Reviewer)
	}
	if resp.Content == "" {
		t.Error("expected non-empty content")
	}
}

func TestAIReviewer_Name(t *testing.T) {
	r := &AIReviewer{}
	if r.Name() != "AI" {
		t.Errorf("expected Name() = 'AI', got '%s'", r.Name())
	}
}

func TestAIReviewer_Available_NoKey(t *testing.T) {
	os.Unsetenv("CRAFT_AI_API_KEY")
	r := &AIReviewer{}
	if r.Available() {
		t.Error("AIReviewer should not be available without API key")
	}
}

func TestAIReviewer_Available_WithKey(t *testing.T) {
	os.Setenv("CRAFT_AI_API_KEY", "test-key")
	defer os.Unsetenv("CRAFT_AI_API_KEY")
	r := &AIReviewer{}
	if !r.Available() {
		t.Error("AIReviewer should be available with API key")
	}
}

func TestCouncilReviewer_Name(t *testing.T) {
	r := &CouncilReviewer{}
	if r.Name() != "Council" {
		t.Errorf("expected Name() = 'Council', got '%s'", r.Name())
	}
}

func TestGetBestReviewer_NullFallback(t *testing.T) {
	// Ensure no AI key is set
	os.Unsetenv("CRAFT_AI_API_KEY")
	// Council likely not in PATH during tests

	r := GetBestReviewer()
	// Should fall back to NullReviewer (or AI/Council if available)
	if r == nil {
		t.Error("GetBestReviewer should never return nil")
	}
}

func TestGetReviewer_Unknown(t *testing.T) {
	_, err := GetReviewer("unknown")
	if err == nil {
		t.Error("expected error for unknown reviewer")
	}
}

func TestGetReviewer_None(t *testing.T) {
	r, err := GetReviewer("none")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if r.Name() != "None" {
		t.Errorf("expected NullReviewer, got '%s'", r.Name())
	}
}

func TestGetReviewer_AI_NotAvailable(t *testing.T) {
	os.Unsetenv("CRAFT_AI_API_KEY")
	_, err := GetReviewer("ai")
	if err == nil {
		t.Error("expected error when AI not available")
	}
}

func TestGetReviewer_AI_Available(t *testing.T) {
	os.Setenv("CRAFT_AI_API_KEY", "test-key")
	defer os.Unsetenv("CRAFT_AI_API_KEY")
	r, err := GetReviewer("ai")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if r.Name() != "AI" {
		t.Errorf("expected AIReviewer, got '%s'", r.Name())
	}
}

func TestBuildPrompt(t *testing.T) {
	req := ReviewRequest{
		Intent: "Add rate limiting",
		Notes:  []string{"Token bucket", "Per-user"},
	}
	prompt := buildPrompt(req)

	if prompt == "" {
		t.Error("expected non-empty prompt")
	}
	if !strings.Contains(prompt, "Add rate limiting") {
		t.Error("prompt should contain intent")
	}
	if !strings.Contains(prompt, "Token bucket") {
		t.Error("prompt should contain notes")
	}
}

func TestBuildPrompt_NoNotes(t *testing.T) {
	req := ReviewRequest{
		Intent: "Test intent",
		Notes:  nil,
	}
	prompt := buildPrompt(req)

	if !strings.Contains(prompt, "(none)") {
		t.Error("prompt should indicate no notes")
	}
}

// mockHTTPClient implements HTTPClient for testing.
type mockHTTPClient struct {
	Response *http.Response
	Err      error
}

func (m *mockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return m.Response, m.Err
}

func TestAIReviewer_Review_WithMock(t *testing.T) {
	os.Setenv("CRAFT_AI_API_KEY", "test-key")
	defer os.Unsetenv("CRAFT_AI_API_KEY")

	mockResp := `{"choices":[{"message":{"content":"Test review content"}}]}`
	mock := &mockHTTPClient{
		Response: &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewBufferString(mockResp)),
		},
	}

	r := &AIReviewer{Client: mock}
	req := ReviewRequest{Intent: "Test intent"}
	resp, err := r.Review(req)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if resp.Content != "Test review content" {
		t.Errorf("expected 'Test review content', got '%s'", resp.Content)
	}
	if resp.Reviewer != "AI" {
		t.Errorf("expected Reviewer = 'AI', got '%s'", resp.Reviewer)
	}
}

func TestAIReviewer_Review_APIError(t *testing.T) {
	os.Setenv("CRAFT_AI_API_KEY", "test-key")
	defer os.Unsetenv("CRAFT_AI_API_KEY")

	mockResp := `{"error":{"message":"Invalid API key"}}`
	mock := &mockHTTPClient{
		Response: &http.Response{
			StatusCode: 401,
			Body:       io.NopCloser(bytes.NewBufferString(mockResp)),
		},
	}

	r := &AIReviewer{Client: mock}
	req := ReviewRequest{Intent: "Test intent"}
	_, err := r.Review(req)

	if err == nil {
		t.Error("expected error for API error response")
	}
	if !strings.Contains(err.Error(), "Invalid API key") {
		t.Errorf("expected error to contain 'Invalid API key', got '%s'", err.Error())
	}
}
