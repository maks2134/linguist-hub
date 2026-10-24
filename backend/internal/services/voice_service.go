package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/makskozlov/linguist-hub/internal/constants"
	"github.com/makskozlov/linguist-hub/internal/dto"
	"github.com/makskozlov/linguist-hub/internal/models"
	"github.com/makskozlov/linguist-hub/internal/repository"
	"github.com/makskozlov/linguist-hub/internal/services/speech"
	"github.com/makskozlov/linguist-hub/pkg/apierror"
	"github.com/makskozlov/linguist-hub/pkg/language"
	"gorm.io/gorm"
)

const voiceLessonStageTotal = 4

type VoiceService struct {
	ai          *AIService
	llm      *LLMService
	lessonRepo  *repository.VoiceLessonRepository
	chatRepo    *repository.ChatRepository
	profileRepo *repository.ProfileRepository
	planRepo    *repository.PlanRepository
	recognizer  speech.Recognizer
}

func NewVoiceService(
	ai *AIService,
	llm *LLMService,
	lessonRepo *repository.VoiceLessonRepository,
	chatRepo *repository.ChatRepository,
	profileRepo *repository.ProfileRepository,
	planRepo *repository.PlanRepository,
	recognizer speech.Recognizer,
) *VoiceService {
	return &VoiceService{
		ai:          ai,
		llm:      llm,
		lessonRepo:  lessonRepo,
		chatRepo:    chatRepo,
		profileRepo: profileRepo,
		planRepo:    planRepo,
		recognizer:  recognizer,
	}
}

func (s *VoiceService) GetLocales(userID, profileID uint) (*dto.VoiceLocalesResponse, error) {
	profile, err := s.ai.validateProfileAccess(userID, profileID)
	if err != nil {
		return nil, err
	}
	code, ok := language.NormalizeLangCode(profile.TargetLanguage)
	if !ok {
		return nil, apierror.ErrInvalidLanguage
	}
	tags := code.LocaleTags()
	return &dto.VoiceLocalesResponse{
		Primary:  tags.Primary,
		Fallback: tags.Fallback,
	}, nil
}

func (s *VoiceService) TutorTurn(ctx context.Context, userID, profileID uint, text, topic, locale string, audio *speech.AudioBlob) (*dto.VoiceTutorResponse, error) {
	transcript, err := s.recognizer.Recognize(ctx, speech.AudioBlob{Text: text, Data: nil}, locale)
	if audio != nil && text == "" {
		transcript, err = s.recognizer.Recognize(ctx, *audio, locale)
	}
	if err != nil {
		if errors.Is(err, speech.ErrClientSide) && text == "" {
			return nil, apierror.New(constants.StatusBadRequest, "TRANSCRIPT_REQUIRED", "Provide message text from client speech recognition")
		}
		return nil, apierror.Wrap(err, constants.StatusInternalError, "STT_ERROR")
	}
	if strings.TrimSpace(transcript.Text) == "" {
		return nil, apierror.New(constants.StatusBadRequest, "EMPTY_TRANSCRIPT", "Empty transcript")
	}
	return s.ai.VoiceTutorChat(ctx, userID, profileID, transcript.Text, topic, locale)
}

type lessonInitJSON struct {
	Phrases []string `json:"phrases"`
	Opening string   `json:"opening"`
}

func (s *VoiceService) LessonStart(ctx context.Context, userID, profileID uint, moduleID int, topic, locale string) (*dto.LessonResponse, error) {
	profile, err := s.ai.validateProfileAccess(userID, profileID)
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(topic) == "" {
		return nil, apierror.New(constants.StatusBadRequest, "TOPIC_REQUIRED", "Topic is required")
	}

	active, err := s.lessonRepo.FindActiveByProfile(profileID)
	if err == nil && active != nil {
		return nil, apierror.New(constants.StatusConflict, "LESSON_ALREADY_ACTIVE", "Finish or abandon the current lesson first")
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, apierror.Wrap(err, constants.StatusInternalError, "DB_ERROR")
	}

	var planID uint
	plan, planErr := s.planRepo.FindByProfileID(profileID)
	if planErr == nil {
		planID = plan.ID
	}

	phrases, opening, err := s.generateLessonInit(ctx, profile, topic)
	if err != nil {
		slog.Error("lesson init generation failed", "error", err)
		opening = fmt.Sprintf("Let's practice %s. Say hello in %s.", topic, promptNameFor(profile.TargetLanguage))
		phrases = []string{}
	}

	now := time.Now()
	lesson := &models.VoiceLesson{
		UserID:        userID,
		ProfileID:     profileID,
		PlanID:        planID,
		ModuleID:      moduleID,
		Topic:         topic,
		Stage:         "warmup",
		StageIndex:    0,
		TargetPhrases: models.JSONStringList(phrases),
		HitPhrases:    models.JSONStringList{},
		Status:        models.VoiceLessonStatusInProgress,
		StartedAt:     now,
		LastTurnAt:    now,
	}
	if err := s.lessonRepo.Create(lesson); err != nil {
		return nil, apierror.Wrap(err, constants.StatusInternalError, "DB_ERROR")
	}

	resolvedLocale := s.resolveLocale(locale, profile.TargetLanguage)
	_ = s.saveLessonMessages(userID, profileID, lesson.ID, "", opening)

	return s.lessonToDTO(lesson, opening, resolvedLocale), nil
}

func (s *VoiceService) LessonTurn(ctx context.Context, userID uint, lessonID uint, text, locale string) (*dto.LessonResponse, error) {
	lesson, err := s.lessonRepo.FindByID(lessonID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apierror.ErrNotFound
		}
		return nil, apierror.Wrap(err, constants.StatusInternalError, "DB_ERROR")
	}
	if lesson.UserID != userID {
		return nil, apierror.ErrForbidden
	}
	if lesson.Status != models.VoiceLessonStatusInProgress {
		return nil, apierror.New(constants.StatusBadRequest, "LESSON_NOT_ACTIVE", "Lesson is not in progress")
	}

	profile, err := s.profileRepo.FindByID(lesson.ProfileID)
	if err != nil {
		return nil, apierror.Wrap(err, constants.StatusInternalError, "DB_ERROR")
	}

	transcript, err := s.recognizer.Recognize(ctx, speech.AudioBlob{Text: text}, locale)
	if err != nil && !errors.Is(err, speech.ErrClientSide) {
		return nil, apierror.Wrap(err, constants.StatusInternalError, "STT_ERROR")
	}
	userText := strings.TrimSpace(transcript.Text)
	if userText == "" {
		userText = strings.TrimSpace(text)
	}
	if userText == "" {
		return nil, apierror.New(constants.StatusBadRequest, "EMPTY_TRANSCRIPT", "Empty transcript")
	}

	s.updateHitPhrases(lesson, userText)
	lesson.StageIndex++
	lesson.LastTurnAt = time.Now()

	prompt, completed, score, err := s.advanceLesson(ctx, profile, lesson, userText)
	if err != nil {
		return nil, err
	}

	_ = s.saveLessonMessages(userID, lesson.ProfileID, lesson.ID, userText, prompt)

	if completed && score != nil {
		lesson.Score = score
		lesson.Status = models.VoiceLessonStatusCompleted
		now := time.Now()
		lesson.CompletedAt = &now
		lesson.Stage = "completed"
	}
	if err := s.lessonRepo.Update(lesson); err != nil {
		return nil, apierror.Wrap(err, constants.StatusInternalError, "DB_ERROR")
	}

	resp := s.lessonToDTO(lesson, prompt, s.resolveLocale(locale, profile.TargetLanguage))
	return resp, nil
}

func (s *VoiceService) LessonAbandon(ctx context.Context, userID, lessonID uint) error {
	lesson, err := s.lessonRepo.FindByID(lessonID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apierror.ErrNotFound
		}
		return apierror.Wrap(err, constants.StatusInternalError, "DB_ERROR")
	}
	if lesson.UserID != userID {
		return apierror.ErrForbidden
	}
	lesson.Status = models.VoiceLessonStatusAbandoned
	return s.lessonRepo.Update(lesson)
}

func (s *VoiceService) ListLessons(userID, profileID uint, limit, offset int) ([]dto.LessonResponse, error) {
	if _, err := s.ai.validateProfileAccess(userID, profileID); err != nil {
		return nil, err
	}
	lessons, err := s.lessonRepo.ListByProfile(profileID, limit, offset)
	if err != nil {
		return nil, apierror.Wrap(err, constants.StatusInternalError, "DB_ERROR")
	}
	out := make([]dto.LessonResponse, len(lessons))
	for i, l := range lessons {
		out[i] = *s.lessonToDTO(&l, "", "")
	}
	return out, nil
}

func (s *VoiceService) generateLessonInit(ctx context.Context, profile *models.LanguageProfile, topic string) ([]string, string, error) {
	prompt := GetVoiceLessonInitPrompt(profile.TargetLanguage, nativeLanguageFor(profile), topic, profile.CurrentLevel)
	raw, err := s.llm.DiagnosticChat(ctx, prompt, nil, "Generate lesson phrases.")
	if err != nil {
		return nil, "", err
	}
	start := strings.Index(raw, "{")
	end := strings.LastIndex(raw, "}")
	if start == -1 || end <= start {
		return nil, "", fmt.Errorf("no JSON in lesson init response")
	}
	var init lessonInitJSON
	if err := json.Unmarshal([]byte(raw[start:end+1]), &init); err != nil {
		return nil, "", err
	}
	return init.Phrases, init.Opening, nil
}

func (s *VoiceService) advanceLesson(ctx context.Context, profile *models.LanguageProfile, lesson *models.VoiceLesson, userText string) (prompt string, completed bool, score *models.VoiceScore, err error) {
	stage := lesson.Stage
	switch stage {
	case "warmup":
		if lesson.StageIndex >= 2 {
			lesson.Stage = "drills"
			lesson.StageIndex = 0
		}
	case "drills":
		if lesson.StageIndex >= 3 || s.allPhrasesHit(lesson) {
			lesson.Stage = "free_talk"
			lesson.StageIndex = 0
		}
	case "free_talk":
		if lesson.StageIndex >= 2 {
			lesson.Stage = "recap"
			lesson.StageIndex = 0
		}
	case "recap":
		// score returned by AI on this stage
	}

	ctxKey := fmt.Sprintf("voice_lesson:%d", lesson.ID)
	history, err := s.chatRepo.FindByUserIDAndContext(lesson.UserID, ctxKey, lesson.ProfileID)
	if err != nil {
		return "", false, nil, apierror.Wrap(err, constants.StatusInternalError, "DB_ERROR")
	}

	systemPrompt := GetVoiceLessonStagePrompt(
		profile.TargetLanguage, nativeLanguageFor(profile), profile.CurrentLevel,
		lesson.Stage, lesson.Topic, []string(lesson.TargetPhrases), []string(lesson.HitPhrases),
	)

	raw, err := s.llm.TutorChat(ctx, systemPrompt, history, userText)
	if err != nil {
		return "", false, nil, apierror.Wrap(err, constants.StatusInternalError, "AI_ERROR")
	}

	text, scorePayload := splitAssistantAndLessonScore(raw)
	if scorePayload != nil {
		return text, true, &models.VoiceScore{
			Fluency:       scorePayload.Score.Fluency,
			Grammar:       scorePayload.Score.Grammar,
			Vocabulary:    scorePayload.Score.Vocabulary,
			Pronunciation: scorePayload.Score.Pronunciation,
			Notes:         scorePayload.Score.Notes,
		}, nil
	}
	return text, false, nil, nil
}

func (s *VoiceService) updateHitPhrases(lesson *models.VoiceLesson, userText string) {
	lower := strings.ToLower(userText)
	for _, phrase := range lesson.TargetPhrases {
		if strings.Contains(lower, strings.ToLower(phrase)) && !containsString(lesson.HitPhrases, phrase) {
			lesson.HitPhrases = append(lesson.HitPhrases, phrase)
		}
	}
}

func (s *VoiceService) allPhrasesHit(lesson *models.VoiceLesson) bool {
	if len(lesson.TargetPhrases) == 0 {
		return true
	}
	return len(lesson.HitPhrases) >= len(lesson.TargetPhrases)
}

func containsString(list models.JSONStringList, item string) bool {
	for _, v := range list {
		if v == item {
			return true
		}
	}
	return false
}

func (s *VoiceService) saveLessonMessages(userID, profileID, lessonID uint, userText, assistantText string) error {
	ctx := fmt.Sprintf("voice_lesson:%d", lessonID)
	now := time.Now()
	if userText != "" {
		if err := s.chatRepo.Create(&models.ChatMessage{
			UserID: userID, ProfileID: profileID, Role: "user",
			Context: ctx, Content: userText, CreatedAt: now,
		}); err != nil {
			return err
		}
	}
	return s.chatRepo.Create(&models.ChatMessage{
		UserID: userID, ProfileID: profileID, Role: "assistant",
		Context: ctx, Content: assistantText, CreatedAt: now,
	})
}

func (s *VoiceService) lessonToDTO(lesson *models.VoiceLesson, prompt, locale string) *dto.LessonResponse {
	resp := &dto.LessonResponse{
		LessonID:      lesson.ID,
		Stage:         lesson.Stage,
		StageIndex:    lesson.StageIndex,
		StageTotal:    voiceLessonStageTotal,
		Prompt:        prompt,
		TargetPhrases: []string(lesson.TargetPhrases),
		HitPhrases:    []string(lesson.HitPhrases),
		Locale:        locale,
	}
	if lesson.Score != nil {
		resp.Score = &dto.LessonScoreDTO{
			Fluency:       lesson.Score.Fluency,
			Grammar:       lesson.Score.Grammar,
			Vocabulary:    lesson.Score.Vocabulary,
			Pronunciation: lesson.Score.Pronunciation,
			Notes:         lesson.Score.Notes,
		}
	}
	return resp
}

func (s *VoiceService) resolveLocale(locale, targetLang string) string {
	if locale != "" {
		return locale
	}
	if code, ok := language.NormalizeLangCode(targetLang); ok {
		return code.LocaleTags().Primary
	}
	return "en-US"
}
