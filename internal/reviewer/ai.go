package reviewer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	envAPIKey      = "CRAFT_AI_API_KEY"
	envModel       = "CRAFT_AI_MODEL"
	envBaseURL     = "CRAFT_AI_BASE_URL"
	defaultModel   = "gpt-4o-mini"
	defaultBaseURL = "https://api.openai.com/v1"
)

// HTTPClient interface for testability.
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// AIReviewer uses an OpenAI-compatible API to review intent.
type AIReviewer struct {
	Client HTTPClient // Optional; uses http.DefaultClient if nil
}

func (r *AIReviewer) Name() string {
	return NameAI
}

func (r *AIReviewer) Available() bool {
	return os.Getenv(envAPIKey) != ""
}

func (r *AIReviewer) Review(req ReviewRequest) (ReviewResponse, error) {
	apiKey := os.Getenv(envAPIKey)
	if apiKey == "" {
		return ReviewResponse{}, fmt.Errorf("%s not set", envAPIKey)
	}

	model := os.Getenv(envModel)
	if model == "" {
		model = defaultModel
	}

	baseURL := os.Getenv(envBaseURL)
	if baseURL == "" {
		baseURL = defaultBaseURL
	}

	client := r.Client
	if client == nil {
		client = &http.Client{Timeout: 60 * time.Second}
	}

	prompt := buildPrompt(req)
	content, err := callAPI(client, baseURL, apiKey, model, prompt)
	if err != nil {
		return ReviewResponse{}, fmt.Errorf("AI review failed: %w", err)
	}

	return ReviewResponse{
		Content:  content,
		Reviewer: NameAI,
	}, nil
}

func buildPrompt(req ReviewRequest) string {
	var sb strings.Builder
	sb.WriteString("You are reviewing a software development intent before implementation begins.\n\n")
	sb.WriteString(fmt.Sprintf("Intent: %s\n\n", req.Intent))

	sb.WriteString("Notes so far:\n")
	if len(req.Notes) == 0 {
		sb.WriteString("(none)\n")
	} else {
		for _, note := range req.Notes {
			sb.WriteString(fmt.Sprintf("- %s\n", note))
		}
	}

	sb.WriteString("\nPlease review this intent and provide:\n")
	sb.WriteString("1. Clarifying questions the developer should consider\n")
	sb.WriteString("2. Potential concerns or risks\n")
	sb.WriteString("3. Suggestions for scope refinement\n\n")
	sb.WriteString("Be direct and constructive. Focus on helping the developer think clearly.")

	return sb.String()
}

type chatRequest struct {
	Model    string        `json:"model"`
	Messages []chatMessage `json:"messages"`
}

type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

func callAPI(client HTTPClient, baseURL, apiKey, model, prompt string) (string, error) {
	reqBody := chatRequest{
		Model: model,
		Messages: []chatMessage{
			{Role: "user", Content: prompt},
		},
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	url := strings.TrimSuffix(baseURL, "/") + "/chat/completions"
	httpReq, err := http.NewRequest("POST", url, bytes.NewReader(jsonBody))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+apiKey)

	resp, err := client.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	var chatResp chatResponse
	if err := json.NewDecoder(resp.Body).Decode(&chatResp); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	if chatResp.Error != nil {
		return "", fmt.Errorf("API error: %s", chatResp.Error.Message)
	}

	if len(chatResp.Choices) == 0 {
		return "", fmt.Errorf("no response from AI")
	}

	return chatResp.Choices[0].Message.Content, nil
}
