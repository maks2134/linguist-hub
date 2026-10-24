package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/makskozlov/linguist-hub/internal/constants"
	"github.com/makskozlov/linguist-hub/internal/models"
)

// LLMService talks to an OpenAI-compatible API (OpenRouter by default).
type LLMService struct {
	baseURL    string
	apiKey     string
	model      string
	httpClient *http.Client
}

func NewLLMService(apiKey, baseURL, model string) (*LLMService, error) {
	if strings.TrimSpace(apiKey) == "" {
		return nil, fmt.Errorf("LLM_API_KEY / OPENROUTER_API_KEY is not set")
	}
	baseURL = strings.TrimRight(strings.TrimSpace(baseURL), "/")
	if baseURL == "" {
		baseURL = constants.DefaultLLMBaseURL
	}
	if !strings.HasSuffix(baseURL, "/v1") {
		baseURL += "/v1"
	}
	if strings.TrimSpace(model) == "" {
		model = constants.DefaultLLMModel
	}
	slog.Info("llm client initialized", "base_url", baseURL, "model", model)
	return &LLMService{
		baseURL: baseURL,
		apiKey:  apiKey,
		model:   model,
		httpClient: &http.Client{
			Timeout: 90 * time.Second,
		},
	}, nil
}

func (s *LLMService) Close() error { return nil }

func (s *LLMService) Model() string { return s.model }

type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatRequest struct {
	Model       string        `json:"model"`
	Messages    []chatMessage `json:"messages"`
	Temperature float64       `json:"temperature"`
	MaxTokens   int           `json:"max_tokens,omitempty"`
}

type chatResponse struct {
	Choices []struct {
		Message chatMessage `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

func (s *LLMService) chat(ctx context.Context, systemPrompt string, history []models.ChatMessage, userMessage string) (string, error) {
	msgs := make([]chatMessage, 0, len(history)+2)
	if strings.TrimSpace(systemPrompt) != "" {
		msgs = append(msgs, chatMessage{Role: "system", Content: systemPrompt})
	}
	for _, m := range history {
		role := "user"
		if m.Role == "assistant" || m.Role == "model" {
			role = "assistant"
		}
		msgs = append(msgs, chatMessage{Role: role, Content: m.Content})
	}
	msgs = append(msgs, chatMessage{Role: "user", Content: userMessage})

	body, err := json.Marshal(chatRequest{
		Model:       s.model,
		Messages:    msgs,
		Temperature: 0.4,
		MaxTokens:   2048,
	})
	if err != nil {
		return "", fmt.Errorf("marshal chat request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.baseURL+"/chat/completions", bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("create chat request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.apiKey)
	req.Header.Set("HTTP-Referer", "https://github.com/maks2134/linguist-hub")
	req.Header.Set("X-Title", "Linguist Hub")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("llm http: %w", err)
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return "", fmt.Errorf("read llm response: %w", err)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("llm http %d: %s", resp.StatusCode, truncate(string(raw), 400))
	}

	var parsed chatResponse
	if err := json.Unmarshal(raw, &parsed); err != nil {
		return "", fmt.Errorf("decode llm response: %w", err)
	}
	if parsed.Error != nil && parsed.Error.Message != "" {
		return "", fmt.Errorf("llm: %s", parsed.Error.Message)
	}
	if len(parsed.Choices) == 0 {
		return "", fmt.Errorf("llm returned no choices")
	}
	return strings.TrimSpace(parsed.Choices[0].Message.Content), nil
}

func (s *LLMService) DiagnosticChat(ctx context.Context, systemPrompt string, history []models.ChatMessage, userMessage string) (string, error) {
	return s.chat(ctx, systemPrompt, history, userMessage)
}

func (s *LLMService) TutorChat(ctx context.Context, systemPrompt string, history []models.ChatMessage, userMessage string) (string, error) {
	return s.chat(ctx, systemPrompt, history, userMessage)
}

func (s *LLMService) GeneratePlan(ctx context.Context, targetLanguage, nativeLanguage, level string) ([]byte, error) {
	prompt := GetPlanGenerationPrompt(targetLanguage, nativeLanguage, level)
	content, err := s.chat(ctx, "You return only valid JSON arrays. No markdown fences.", nil, prompt)
	if err != nil {
		return nil, fmt.Errorf("failed to generate plan: %w", err)
	}

	content = strings.TrimSpace(content)
	content = strings.TrimPrefix(content, "```json")
	content = strings.TrimPrefix(content, "```")
	content = strings.TrimSuffix(content, "```")
	content = strings.TrimSpace(content)

	var modules []models.Module
	if err := json.Unmarshal([]byte(content), &modules); err != nil {
		return nil, fmt.Errorf("failed to parse plan JSON: %w, content: %s", err, truncate(content, 300))
	}
	return []byte(content), nil
}

// TranscribeBlob is not supported via OpenRouter text chat; use browser STT (FEATURE_SERVER_STT=false).
func (s *LLMService) TranscribeBlob(ctx context.Context, data []byte, mimeType, locale string) (string, error) {
	_ = ctx
	_ = data
	_ = mimeType
	_ = locale
	return "", fmt.Errorf("server STT is not available with OpenRouter; use client-side speech recognition")
}

type LevelAssessment struct {
	Level      string   `json:"level"`
	Weaknesses []string `json:"weaknesses"`
}

type OnboardingResult struct {
	Phase      string   `json:"phase"`
	Level      string   `json:"level"`
	Goals      []string `json:"goals"`
	Weaknesses []string `json:"weaknesses"`
}

func (s *LLMService) ExtractOnboardingResult(response string) (*OnboardingResult, error) {
	response = strings.TrimSpace(response)
	startIdx := strings.LastIndex(response, "{")
	endIdx := strings.LastIndex(response, "}")
	if startIdx == -1 || endIdx == -1 || startIdx >= endIdx {
		return nil, nil
	}
	var result OnboardingResult
	if err := json.Unmarshal([]byte(response[startIdx:endIdx+1]), &result); err != nil {
		return nil, nil
	}
	if result.Phase != "completed" {
		return nil, nil
	}
	if !validLevel(result.Level) {
		return nil, nil
	}
	return &result, nil
}

func (s *LLMService) ExtractLevelFromResponse(response string) (*LevelAssessment, error) {
	response = strings.TrimSpace(response)
	startIdx := strings.Index(response, "{")
	endIdx := strings.LastIndex(response, "}")
	if startIdx == -1 || endIdx == -1 || startIdx >= endIdx {
		return nil, nil
	}
	var assessment LevelAssessment
	if err := json.Unmarshal([]byte(response[startIdx:endIdx+1]), &assessment); err != nil {
		return nil, nil
	}
	if !validLevel(assessment.Level) {
		return nil, nil
	}
	return &assessment, nil
}

func validLevel(level string) bool {
	switch level {
	case "A1", "A2", "B1", "B2", "C1", "C2":
		return true
	default:
		return false
	}
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "…"
}
