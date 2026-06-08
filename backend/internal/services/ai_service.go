package services

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"strings"
	"time"

	"github.com/makskozlov/linguist-hub/internal/constants"
	"github.com/makskozlov/linguist-hub/internal/dto"
	"github.com/makskozlov/linguist-hub/internal/models"
	"github.com/makskozlov/linguist-hub/internal/repository"
	"github.com/makskozlov/linguist-hub/pkg/apierror"
	"github.com/makskozlov/linguist-hub/pkg/language"
	"gorm.io/gorm"
)

type AIService struct {
	geminiService *GeminiService
	chatRepo      *repository.ChatRepository
	profileRepo   *repository.ProfileRepository
	planRepo      *repository.PlanRepository
}

func NewAIService(
	geminiService *GeminiService,
	chatRepo *repository.ChatRepository,
	profileRepo *repository.ProfileRepository,
	planRepo *repository.PlanRepository,
) *AIService {
	return &AIService{
		geminiService: geminiService,
		chatRepo:      chatRepo,
		profileRepo:   profileRepo,
		planRepo:      planRepo,
	}
}

func nativeLanguageFor(profile *models.LanguageProfile) string {
	if profile.NativeLanguage != "" {
		return profile.NativeLanguage
	}
	return string(language.DefaultNative)
}

func planToDTO(plan *models.StudyPlan) *dto.PlanResponse {
	dtoModules := make([]dto.Module, len(plan.Structure))
	for i, m := range plan.Structure {
		locked := ModuleLocked(plan.Structure, i)
		completedSet := make(map[string]bool, len(m.CompletedTopics))
		for _, ct := range m.CompletedTopics {
			completedSet[normalizeTopic(ct)] = true
		}
		topics := make([]dto.Topic, len(m.Topics))
		for j, t := range m.Topics {
			norm := normalizeTopic(t)
			topics[j] = dto.Topic{
				Name:      norm,
				Completed: completedSet[norm],
				Locked:    locked,
			}
		}
		dtoModules[i] = dto.Module{
			ID:          m.ID,
			Title:       m.Title,
			Description: m.Description,
			Topics:      topics,
			Locked:      locked,
		}
	}
	return &dto.PlanResponse{
		ID:        plan.ID,
		ProfileID: plan.ProfileID,
		Modules:   dtoModules,
	}
}

func (s *AIService) validateProfileAccess(userID, profileID uint) (*models.LanguageProfile, error) {
	profile, err := s.profileRepo.FindByID(profileID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apierror.ErrNotFound
		}
		return nil, apierror.Wrap(err, constants.StatusInternalError, "DB_ERROR")
	}

	if profile.UserID != userID {
		return nil, apierror.ErrForbidden
	}

	return profile, nil
}

func (s *AIService) DiagnosticChat(ctx context.Context, userID, profileID uint, message string) (*dto.ChatResponse, error) {
	profile, err := s.validateProfileAccess(userID, profileID)
	if err != nil {
		return nil, err
	}

	history, err := s.chatRepo.FindByUserIDAndContext(userID, "diagnostic", profileID)
	if err != nil {
		return nil, apierror.Wrap(err, constants.StatusInternalError, "DB_ERROR")
	}

	systemPrompt := GetDiagnosticSystemPrompt(profile.TargetLanguage)

	response, err := s.geminiService.DiagnosticChat(ctx, systemPrompt, history, message)
	if err != nil {
		slog.Error("gemini API error", "error", err)
		return nil, apierror.Wrap(err, constants.StatusInternalError, "AI_ERROR")
	}

	userMsg := &models.ChatMessage{
		UserID:    userID,
		ProfileID: profileID,
		Role:      "user",
		Context:   "diagnostic",
		Content:   message,
		CreatedAt: time.Now(),
	}
	if err := s.chatRepo.Create(userMsg); err != nil {
		return nil, apierror.Wrap(err, constants.StatusInternalError, "DB_ERROR")
	}

	assistantMsg := &models.ChatMessage{
		UserID:    userID,
		ProfileID: profileID,
		Role:      "assistant",
		Context:   "diagnostic",
		Content:   response,
		CreatedAt: time.Now(),
	}
	if err := s.chatRepo.Create(assistantMsg); err != nil {
		return nil, apierror.Wrap(err, constants.StatusInternalError, "DB_ERROR")
	}

	assessment, err := s.geminiService.ExtractLevelFromResponse(response)
	if err == nil && assessment != nil {
		profile.CurrentLevel = assessment.Level
		profile.Weaknesses = models.JSONStringList(assessment.Weaknesses)
		if err := s.profileRepo.Update(profile); err != nil {
			return nil, apierror.Wrap(err, constants.StatusInternalError, "DB_ERROR")
		}

		return &dto.ChatResponse{
			Message:      response,
			Level:        assessment.Level,
			Weaknesses:   assessment.Weaknesses,
			IsAssessment: true,
		}, nil
	}

	return &dto.ChatResponse{
		Message:      response,
		IsAssessment: false,
	}, nil
}

func (s *AIService) GeneratePlan(ctx context.Context, userID, profileID uint) (*dto.PlanResponse, error) {
	profile, err := s.validateProfileAccess(userID, profileID)
	if err != nil {
		return nil, err
	}

	if profile.CurrentLevel == "unknown" {
		return nil, apierror.New(constants.StatusBadRequest, "LEVEL_REQUIRED", "Please complete diagnostic first")
	}

	existingPlan, err := s.planRepo.FindByProfileID(profileID)
	if err == nil {
		return planToDTO(existingPlan), nil
	}

	planJSON, err := s.geminiService.GeneratePlan(ctx, profile.TargetLanguage, nativeLanguageFor(profile), profile.CurrentLevel)
	if err != nil {
		return nil, apierror.Wrap(err, constants.StatusInternalError, "AI_ERROR")
	}

	var modules models.Modules
	if err := json.Unmarshal(planJSON, &modules); err != nil {
		return nil, apierror.Wrap(err, constants.StatusInternalError, "PARSE_ERROR")
	}

	plan := &models.StudyPlan{
		ProfileID:   profileID,
		Structure:   modules,
		IsGenerated: true,
	}

	if err := s.planRepo.Create(plan); err != nil {
		return nil, apierror.Wrap(err, constants.StatusInternalError, "DB_ERROR")
	}

	return planToDTO(plan), nil
}

func (s *AIService) VoiceTutorChat(ctx context.Context, userID, profileID uint, message, topic, locale string) (*dto.VoiceTutorResponse, error) {
	profile, err := s.validateProfileAccess(userID, profileID)
	if err != nil {
		return nil, err
	}

	targetCode, ok := language.NormalizeLangCode(profile.TargetLanguage)
	if !ok {
		return nil, apierror.ErrInvalidLanguage
	}
	if locale != "" && !language.ValidateLocale(locale, targetCode) {
		return nil, apierror.New(constants.StatusBadRequest, "INVALID_LOCALE", "Unsupported locale for this profile")
	}
	resolvedLocale := locale
	if resolvedLocale == "" {
		resolvedLocale = targetCode.LocaleTags().Primary
	}

	history, err := s.chatRepo.FindByContexts(userID, []string{"tutor", "voice_tutor"}, profileID)
	if err != nil {
		return nil, apierror.Wrap(err, constants.StatusInternalError, "DB_ERROR")
	}

	if topic == "" {
		topic = "general conversation"
	}
	systemPrompt := GetVoiceTutorSystemPrompt(profile.TargetLanguage, nativeLanguageFor(profile), topic, profile.CurrentLevel, []string(profile.Weaknesses))

	rawResponse, err := s.geminiService.TutorChat(ctx, systemPrompt, history, message)
	if err != nil {
		return nil, apierror.Wrap(err, constants.StatusInternalError, "AI_ERROR")
	}

	assistantText, corrections := splitAssistantAndCorrections(rawResponse)

	now := time.Now()
	userMsg := &models.ChatMessage{
		UserID: userID, ProfileID: profileID, Role: "user",
		Context: "voice_tutor", Content: message, CreatedAt: now,
	}
	if err := s.chatRepo.Create(userMsg); err != nil {
		return nil, apierror.Wrap(err, constants.StatusInternalError, "DB_ERROR")
	}
	assistantMsg := &models.ChatMessage{
		UserID: userID, ProfileID: profileID, Role: "assistant",
		Context: "voice_tutor", Content: assistantText, CreatedAt: now,
	}
	if err := s.chatRepo.Create(assistantMsg); err != nil {
		return nil, apierror.Wrap(err, constants.StatusInternalError, "DB_ERROR")
	}

	return &dto.VoiceTutorResponse{
		UserText:    message,
		Assistant:   assistantText,
		Corrections: corrections,
		Locale:      resolvedLocale,
	}, nil
}

func (s *AIService) TutorChat(ctx context.Context, userID, profileID uint, message, topic string) (*dto.ChatResponse, error) {
	profile, err := s.validateProfileAccess(userID, profileID)
	if err != nil {
		return nil, err
	}

	history, err := s.chatRepo.FindByUserIDAndContext(userID, "tutor", profileID)
	if err != nil {
		return nil, apierror.Wrap(err, constants.StatusInternalError, "DB_ERROR")
	}
	systemPrompt := GetTutorSystemPrompt(profile.TargetLanguage, nativeLanguageFor(profile), topic, profile.CurrentLevel, []string(profile.Weaknesses))

	response, err := s.geminiService.TutorChat(ctx, systemPrompt, history, message)
	if err != nil {
		return nil, apierror.Wrap(err, constants.StatusInternalError, "AI_ERROR")
	}

	userMsg := &models.ChatMessage{
		UserID:    userID,
		ProfileID: profileID,
		Role:      "user",
		Context:   "tutor",
		Content:   message,
		CreatedAt: time.Now(),
	}
	if err := s.chatRepo.Create(userMsg); err != nil {
		return nil, apierror.Wrap(err, constants.StatusInternalError, "DB_ERROR")
	}

	assistantMsg := &models.ChatMessage{
		UserID:    userID,
		ProfileID: profileID,
		Role:      "assistant",
		Context:   "tutor",
		Content:   response,
		CreatedAt: time.Now(),
	}
	if err := s.chatRepo.Create(assistantMsg); err != nil {
		return nil, apierror.Wrap(err, constants.StatusInternalError, "DB_ERROR")
	}

	return &dto.ChatResponse{
		Message: response,
	}, nil
}

func (s *AIService) OnboardingChat(ctx context.Context, userID, profileID uint, message string) (*dto.OnboardingChatResponse, error) {
	profile, err := s.validateProfileAccess(userID, profileID)
	if err != nil {
		return nil, err
	}

	if profile.OnboardingPhase == "completed" {
		return nil, apierror.New(constants.StatusBadRequest, "ONBOARDING_COMPLETE", "Onboarding already completed, use tutor chat")
	}

	history, err := s.chatRepo.FindByUserIDAndContext(userID, "onboarding", profileID)
	if err != nil {
		return nil, apierror.Wrap(err, constants.StatusInternalError, "DB_ERROR")
	}

	systemPrompt := GetOnboardingSystemPrompt(profile.TargetLanguage, nativeLanguageFor(profile))

	response, err := s.geminiService.DiagnosticChat(ctx, systemPrompt, history, message)
	if err != nil {
		slog.Error("gemini API error", "error", err)
		return nil, apierror.Wrap(err, constants.StatusInternalError, "AI_ERROR")
	}

	userMsg := &models.ChatMessage{
		UserID:    userID,
		ProfileID: profileID,
		Role:      "user",
		Context:   "onboarding",
		Content:   message,
		CreatedAt: time.Now(),
	}
	if err := s.chatRepo.Create(userMsg); err != nil {
		return nil, apierror.Wrap(err, constants.StatusInternalError, "DB_ERROR")
	}

	assistantMsg := &models.ChatMessage{
		UserID:    userID,
		ProfileID: profileID,
		Role:      "assistant",
		Context:   "onboarding",
		Content:   response,
		CreatedAt: time.Now(),
	}
	if err := s.chatRepo.Create(assistantMsg); err != nil {
		return nil, apierror.Wrap(err, constants.StatusInternalError, "DB_ERROR")
	}

	result, _ := s.geminiService.ExtractOnboardingResult(response)
	if result != nil {
		profile.CurrentLevel = result.Level
		profile.OnboardingPhase = "completed"
		profile.Goals = strings.Join(result.Goals, ", ")
		profile.Weaknesses = models.JSONStringList(result.Weaknesses)
		if err := s.profileRepo.Update(profile); err != nil {
			return nil, apierror.Wrap(err, constants.StatusInternalError, "DB_ERROR")
		}

		var planResp *dto.PlanResponse
		planJSON, planErr := s.geminiService.GeneratePlan(ctx, profile.TargetLanguage, nativeLanguageFor(profile), profile.CurrentLevel)
		if planErr == nil {
			var modules models.Modules
			if parseErr := json.Unmarshal(planJSON, &modules); parseErr == nil {
				plan := &models.StudyPlan{
					ProfileID:   profileID,
					Structure:   modules,
					IsGenerated: true,
				}
				if createErr := s.planRepo.Create(plan); createErr == nil {
					planResp = planToDTO(plan)
				}
			}
		}

		return &dto.OnboardingChatResponse{
			Message:     response,
			Phase:       "completed",
			Level:       result.Level,
			Goals:       result.Goals,
			Weaknesses:  result.Weaknesses,
			Plan:        planResp,
			IsCompleted: true,
		}, nil
	}

	return &dto.OnboardingChatResponse{
		Message:     response,
		Phase:       profile.OnboardingPhase,
		IsCompleted: false,
	}, nil
}

func (s *AIService) GetPlan(userID, profileID uint) (*dto.PlanResponse, error) {
	_, err := s.validateProfileAccess(userID, profileID)
	if err != nil {
		return nil, err
	}

	plan, err := s.planRepo.FindByProfileID(profileID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apierror.New(constants.StatusNotFound, "PLAN_NOT_FOUND", "Study plan not found")
		}
		return nil, apierror.Wrap(err, constants.StatusInternalError, "DB_ERROR")
	}

	return planToDTO(plan), nil
}

func (s *AIService) GetChatHistory(userID, profileID uint, context string, limit, offset int) ([]dto.ChatHistoryMessage, error) {
	_, err := s.validateProfileAccess(userID, profileID)
	if err != nil {
		return nil, err
	}

	history, err := s.chatRepo.FindByUserIDAndContextPaginated(userID, context, profileID, limit, offset)
	if err != nil {
		return nil, apierror.Wrap(err, constants.StatusInternalError, "DB_ERROR")
	}

	result := make([]dto.ChatHistoryMessage, len(history))
	for i, msg := range history {
		result[i] = dto.ChatHistoryMessage{
			ID:        msg.ID,
			UserID:    msg.UserID,
			ProfileID: msg.ProfileID,
			Role:      msg.Role,
			Context:   msg.Context,
			Content:   msg.Content,
			CreatedAt: msg.CreatedAt.Format("2006-01-02T15:04:05Z"),
		}
	}

	return result, nil
}
