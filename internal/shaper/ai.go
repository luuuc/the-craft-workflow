package shaper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"craft/internal/structure"
)

const (
	envAPIKey      = "CRAFT_AI_API_KEY"
	envModel       = "CRAFT_AI_MODEL"
	envBaseURL     = "CRAFT_AI_BASE_URL"
	defaultModel   = "gpt-4o-mini"
	defaultBaseURL = "https://api.openai.com/v1"
)

// Pre-compiled regexes for card parsing.
var (
	cardRegex  = regexp.MustCompile(`(?s)===CARD===\s*(.+?)\s*===END===`)
	titleRegex = regexp.MustCompile(`(?m)^#\s*Card:\s*(.+)$`)
	slugRegex  = regexp.MustCompile(`[^a-z0-9]+`)
)

// HTTPClient interface for testability.
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// AIShaper uses an OpenAI-compatible API to generate structure.
type AIShaper struct {
	Client HTTPClient // Optional; uses http.DefaultClient if nil
}

func (s *AIShaper) Name() string {
	return NameAI
}

func (s *AIShaper) Available() bool {
	return os.Getenv(envAPIKey) != ""
}

func (s *AIShaper) Shape(req ShapeRequest) (ShapeResult, error) {
	apiKey := os.Getenv(envAPIKey)
	if apiKey == "" {
		return ShapeResult{}, fmt.Errorf("%s not set", envAPIKey)
	}

	model := os.Getenv(envModel)
	if model == "" {
		model = defaultModel
	}

	baseURL := os.Getenv(envBaseURL)
	if baseURL == "" {
		baseURL = defaultBaseURL
	}

	client := s.Client
	if client == nil {
		client = &http.Client{Timeout: 120 * time.Second}
	}

	// Ensure structure directory exists
	if err := structure.EnsureStructureDir(); err != nil {
		return ShapeResult{}, fmt.Errorf("failed to create structure dir: %w", err)
	}

	// Generate pitch
	pitchPrompt := buildPitchPrompt(req)
	pitchContent, err := callAPI(client, baseURL, apiKey, model, pitchPrompt)
	if err != nil {
		return ShapeResult{}, fmt.Errorf("failed to generate pitch: %w", err)
	}

	// Write pitch file
	pitchPath := structure.PitchPath()
	if err := os.WriteFile(pitchPath, []byte(pitchContent), 0644); err != nil {
		return ShapeResult{}, fmt.Errorf("failed to write pitch: %w", err)
	}

	// Generate cards - cleanup pitch on failure
	cardsPrompt := buildCardsPrompt(req, pitchContent)
	cardsContent, err := callAPI(client, baseURL, apiKey, model, cardsPrompt)
	if err != nil {
		os.Remove(pitchPath) // Cleanup partial state
		return ShapeResult{}, fmt.Errorf("failed to generate cards: %w", err)
	}

	// Parse and write card files - cleanup pitch on failure
	cardPaths, err := parseAndWriteCards(cardsContent)
	if err != nil {
		os.Remove(pitchPath) // Cleanup partial state
		return ShapeResult{}, fmt.Errorf("failed to write cards: %w", err)
	}

	return ShapeResult{
		PitchPath: pitchPath,
		CardPaths: cardPaths,
		Shaper:    NameAI,
	}, nil
}

func buildPitchPrompt(req ShapeRequest) string {
	var sb strings.Builder
	sb.WriteString("Generate a pitch document for this software feature.\n\n")
	sb.WriteString(fmt.Sprintf("Intent: %s\n\n", req.Intent))

	if len(req.Notes) > 0 {
		sb.WriteString("Notes:\n")
		for _, note := range req.Notes {
			sb.WriteString(fmt.Sprintf("- %s\n", note))
		}
		sb.WriteString("\n")
	}

	sb.WriteString(`Format the pitch EXACTLY like this (use markdown):

# Pitch: [Title]

## Problem
[What problem does this solve? 2-3 sentences]

## Solution
[How will it be solved? Be specific about approach]

## Scope

### In Scope
- [Bullet points of what's included]

### Out of Scope
- [Bullet points of what's NOT included]

## Tasks
- [ ] [High-level task 1]
- [ ] [High-level task 2]
- [ ] [High-level task 3]

Keep it concise. Focus on clarity, not length.`)

	return sb.String()
}

func buildCardsPrompt(req ShapeRequest, pitchContent string) string {
	var sb strings.Builder
	sb.WriteString("Break down this pitch into implementation cards.\n\n")
	sb.WriteString("Pitch:\n")
	sb.WriteString(pitchContent)
	sb.WriteString("\n\n")

	sb.WriteString(`Generate 2-5 cards. Each card should be a focused, completable unit of work.

Output format - use this EXACT structure with === as separator:

===CARD===
# Card: [Descriptive Title]

## Summary
[One sentence describing what this card accomplishes]

## Tasks
- [ ] Task 1
- [ ] Task 2
- [ ] Task 3

## Acceptance Criteria
- [How do we know it's done?]
===END===

Repeat the ===CARD=== ... ===END=== block for each card.
Number the cards implicitly by order (first card = 01, second = 02, etc).
Keep cards focused - if a card has more than 5 tasks, split it.`)

	return sb.String()
}

func parseAndWriteCards(content string) ([]string, error) {
	// Parse cards between ===CARD=== and ===END=== markers
	matches := cardRegex.FindAllStringSubmatch(content, -1)

	var cardPaths []string
	for i, match := range matches {
		if len(match) < 2 {
			continue
		}

		cardContent := strings.TrimSpace(match[1])

		// Extract title for filename
		titleMatch := titleRegex.FindStringSubmatch(cardContent)
		title := "untitled"
		if len(titleMatch) >= 2 {
			title = strings.TrimSpace(titleMatch[1])
		}

		// Create slug from title
		slug := strings.ToLower(title)
		slug = slugRegex.ReplaceAllString(slug, "-")
		slug = strings.Trim(slug, "-")
		if len(slug) > 40 {
			slug = slug[:40]
		}

		filename := fmt.Sprintf("%02d-%s.md", i+1, slug)
		cardPath := fmt.Sprintf("%s/%s", structure.CardsDirPath(), filename)

		if err := os.WriteFile(cardPath, []byte(cardContent), 0644); err != nil {
			return nil, err
		}

		cardPaths = append(cardPaths, cardPath)
	}

	return cardPaths, nil
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
