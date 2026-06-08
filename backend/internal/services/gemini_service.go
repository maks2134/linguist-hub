package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"github.com/makskozlov/linguist-hub/internal/constants"
	"github.com/makskozlov/linguist-hub/internal/models"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
)

type GeminiService struct {
	client *genai.Client
	model  *genai.GenerativeModel
}

func NewGeminiService(apiKey string) (*GeminiService, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("GEMINI_API_KEY is not set")
	}

	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("failed to create Gemini client: %w", err)
	}

	model := client.GenerativeModel(constants.GeminiModel)
	slog.Info("gemini model initialized", "model", model)

	return &GeminiService{
		client: client,
		model:  model,
	}, nil
}

func logGeminiError(op string, err error) {
	var apiErr *googleapi.Error
	if errors.As(err, &apiErr) {
		slog.Error("gemini operation failed", "operation", op, "http_code", apiErr.Code, "error", apiErr.Message)
	} else {
		slog.Error("gemini operation failed", "operation", op, "error", err)
	}
}

func (s *GeminiService) Close() error {
	return s.client.Close()
}

func (s *GeminiService) modelWithInstruction(systemPrompt string) *genai.GenerativeModel {
	model := *s.model
	model.SystemInstruction = &genai.Content{
		Parts: []genai.Part{genai.Text(systemPrompt)},
	}
	return &model
}

func (s *GeminiService) DiagnosticChat(ctx context.Context, systemPrompt string, history []models.ChatMessage, userMessage string) (string, error) {
	model := s.modelWithInstruction(systemPrompt)

	chat := model.StartChat()
	if len(history) > 0 {
		chat.History = s.convertHistory(history)
	}

	resp, err := chat.SendMessage(ctx, genai.Text(userMessage))
	if err != nil {
		logGeminiError("DiagnosticChat", err)
		return "", fmt.Errorf("failed to send message (check API key and quota): %w", err)
	}

	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("no response from Gemini API")
	}

	return fmt.Sprintf("%v", resp.Candidates[0].Content.Parts[0]), nil
}

func (s *GeminiService) GeneratePlan(ctx context.Context, targetLanguage, nativeLanguage, level string) ([]byte, error) {
	prompt := GetPlanGenerationPrompt(targetLanguage, nativeLanguage, level)

	resp, err := s.model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		logGeminiError("GeneratePlan", err)
		return nil, fmt.Errorf("failed to generate plan: %w", err)
	}

	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return nil, fmt.Errorf("no response from Gemini")
	}

	content := fmt.Sprintf("%v", resp.Candidates[0].Content.Parts[0])
	content = strings.TrimSpace(content)
	content = strings.Trim(content, "`")
	content = strings.TrimPrefix(content, "json")
	content = strings.TrimSpace(content)

	var modules []models.Module
	if err := json.Unmarshal([]byte(content), &modules); err != nil {
		return nil, fmt.Errorf("failed to parse plan JSON: %w, content: %s", err, content)
	}

	return []byte(content), nil
}

func (s *GeminiService) TutorChat(ctx context.Context, systemPrompt string, history []models.ChatMessage, userMessage string) (string, error) {
	model := s.modelWithInstruction(systemPrompt)

	chat := model.StartChat()
	if len(history) > 0 {
		chat.History = s.convertHistory(history)
	}

	resp, err := chat.SendMessage(ctx, genai.Text(userMessage))
	if err != nil {
		logGeminiError("TutorChat", err)
		return "", fmt.Errorf("failed to send message: %w", err)
	}

	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("no response from Gemini")
	}

	return fmt.Sprintf("%v", resp.Candidates[0].Content.Parts[0]), nil
}

func (s *GeminiService) TranscribeBlob(ctx context.Context, data []byte, mimeType, locale string) (string, error) {
	if len(data) == 0 {
		return "", fmt.Errorf("empty audio data")
	}
	prompt := fmt.Sprintf("Transcribe the following audio accurately in locale %s. Return only the transcript text, no commentary.", locale)
	model := s.modelWithInstruction(prompt)
	resp, err := model.GenerateContent(ctx,
		genai.Blob{MIMEType: mimeType, Data: data},
		genai.Text("Transcribe."),
	)
	if err != nil {
		logGeminiError("TranscribeBlob", err)
		return "", fmt.Errorf("transcription failed: %w", err)
	}
	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("no transcription from Gemini")
	}
	return strings.TrimSpace(fmt.Sprintf("%v", resp.Candidates[0].Content.Parts[0])), nil
}

func (s *GeminiService) convertHistory(messages []models.ChatMessage) []*genai.Content {
	history := make([]*genai.Content, 0, len(messages))
	for _, msg := range messages {
		role := "user"
		if msg.Role == "assistant" {
			role = "model"
		}
		history = append(history, &genai.Content{
			Parts: []genai.Part{genai.Text(msg.Content)},
			Role:  role,
		})
	}
	return history
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

func (s *GeminiService) ExtractOnboardingResult(response string) (*OnboardingResult, error) {
	response = strings.TrimSpace(response)

	startIdx := strings.LastIndex(response, "{")
	endIdx := strings.LastIndex(response, "}")

	if startIdx == -1 || endIdx == -1 || startIdx >= endIdx {
		return nil, nil
	}

	jsonStr := response[startIdx : endIdx+1]

	var result OnboardingResult
	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		return nil, nil
	}

	if result.Phase != "completed" {
		return nil, nil
	}

	validLevels := map[string]bool{
		"A1": true, "A2": true, "B1": true, "B2": true, "C1": true, "C2": true,
	}
	if !validLevels[result.Level] {
		return nil, nil
	}

	return &result, nil
}

func (s *GeminiService) ExtractLevelFromResponse(response string) (*LevelAssessment, error) {
	response = strings.TrimSpace(response)

	startIdx := strings.Index(response, "{")
	endIdx := strings.LastIndex(response, "}")

	if startIdx == -1 || endIdx == -1 || startIdx >= endIdx {
		return nil, nil
	}

	jsonStr := response[startIdx : endIdx+1]

	var assessment LevelAssessment
	if err := json.Unmarshal([]byte(jsonStr), &assessment); err != nil {
		return nil, nil
	}

	validLevels := map[string]bool{
		"A1": true, "A2": true, "B1": true, "B2": true, "C1": true, "C2": true,
	}
	if !validLevels[assessment.Level] {
		return nil, nil
	}

	return &assessment, nil
}
