package handlers

import (
	"context"
	"encoding/json"
	"log/slog"
	"strings"
	"time"

	"github.com/gofiber/contrib/v3/websocket"
	"github.com/gofiber/fiber/v3"
	"github.com/makskozlov/linguist-hub/internal/services"
	"github.com/makskozlov/linguist-hub/internal/services/speech"
	"github.com/makskozlov/linguist-hub/pkg/audio"
	"github.com/makskozlov/linguist-hub/pkg/ratelimit"
)

type WebSocketMessage struct {
	Type      string `json:"type"`
	Message   string `json:"message"`
	Topic     string `json:"topic"`
	ProfileID uint   `json:"profile_id"`
	Locale    string `json:"locale,omitempty"`
	LessonID  uint   `json:"lesson_id,omitempty"`
	ModuleID  int    `json:"module_id,omitempty"`
	AudioB64  string `json:"audio_b64,omitempty"`
	AudioMIME string `json:"audio_mime,omitempty"`
	Seq       int    `json:"seq,omitempty"`
	Final     bool   `json:"final,omitempty"`
}

type WebSocketResponse struct {
	Type          string   `json:"type"`
	Content       string   `json:"content,omitempty"`
	Chunk         string   `json:"chunk,omitempty"`
	UserText      string   `json:"user_text,omitempty"`
	Level         string   `json:"level,omitempty"`
	Weaknesses    []string `json:"weaknesses,omitempty"`
	IsAssessment  bool     `json:"is_assessment,omitempty"`
	Phase         string   `json:"phase,omitempty"`
	Goals         []string `json:"goals,omitempty"`
	IsCompleted   bool     `json:"is_completed,omitempty"`
	Plan          string   `json:"plan,omitempty"`
	Error         string   `json:"error,omitempty"`
	Done          bool     `json:"done,omitempty"`
	Corrections   any      `json:"corrections,omitempty"`
	LessonID      uint     `json:"lesson_id,omitempty"`
	Stage         string   `json:"stage,omitempty"`
	StageIndex    int      `json:"stage_index,omitempty"`
	StageTotal    int      `json:"stage_total,omitempty"`
	Prompt        string   `json:"prompt,omitempty"`
	TargetPhrases []string `json:"target_phrases,omitempty"`
	HitPhrases    []string `json:"hit_phrases,omitempty"`
	Score         any      `json:"score,omitempty"`
}

type WebSocketHandler struct {
	aiService    *services.AIService
	voiceService *services.VoiceService
	aiMinute     *ratelimit.Limiter
	aiDaily      *ratelimit.DailyQuota
}

func NewWebSocketHandler(
	aiService *services.AIService,
	voiceService *services.VoiceService,
	aiMinute *ratelimit.Limiter,
	aiDaily *ratelimit.DailyQuota,
) *WebSocketHandler {
	return &WebSocketHandler{
		aiService:    aiService,
		voiceService: voiceService,
		aiMinute:     aiMinute,
		aiDaily:      aiDaily,
	}
}

func (h *WebSocketHandler) allowAI(userID uint) bool {
	key := ratelimit.UserKey(userID)
	if h.aiMinute != nil && !h.aiMinute.Allow(key) {
		return false
	}
	if h.aiDaily != nil && !h.aiDaily.Allow(key) {
		return false
	}
	return true
}

func (h *WebSocketHandler) HandleWebSocket() fiber.Handler {
	return websocket.New(func(c *websocket.Conn) {
		userID, ok := c.Locals("userID").(uint)
		if !ok {
			h.sendError(c, "unauthorized", true)
			return
		}

		for {
			_, msg, err := c.ReadMessage()
			if err != nil {
				slog.Error("websocket read error", "error", err)
				return
			}

			var req WebSocketMessage
			if err := json.Unmarshal(msg, &req); err != nil {
				h.sendError(c, "invalid message format", false)
				continue
			}

			switch req.Type {
			case "onboarding", "diagnostic", "tutor", "generate_plan",
				"voice_tutor_turn", "voice_tutor_audio",
				"voice_lesson_start", "voice_lesson_turn":
				if !h.allowAI(userID) {
					h.sendError(c, "rate_limited", false)
					continue
				}
			}

			switch req.Type {
			case "onboarding":
				h.handleOnboardingChat(c, userID, req)
			case "diagnostic":
				h.handleDiagnosticChat(c, userID, req)
			case "tutor":
				h.handleTutorChat(c, userID, req)
			case "generate_plan":
				h.handleGeneratePlan(c, userID, req)
			case "voice_tutor_turn":
				h.handleVoiceTutorTurn(c, userID, req)
			case "voice_tutor_audio":
				h.handleVoiceTutorAudio(c, userID, req)
			case "voice_lesson_start":
				h.handleVoiceLessonStart(c, userID, req)
			case "voice_lesson_turn":
				h.handleVoiceLessonTurn(c, userID, req)
			case "voice_lesson_abandon":
				h.handleVoiceLessonAbandon(c, userID, req)
			default:
				h.sendError(c, "unknown message type: "+req.Type, false)
			}
		}
	})
}

func (h *WebSocketHandler) handleOnboardingChat(c *websocket.Conn, userID uint, req WebSocketMessage) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	resp, err := h.aiService.OnboardingChat(ctx, userID, req.ProfileID, req.Message)
	if err != nil {
		h.sendError(c, err.Error(), false)
		return
	}

	words := splitIntoChunks(resp.Message, 5)
	for _, chunk := range words {
		if err := c.WriteJSON(WebSocketResponse{
			Type:  "onboarding",
			Chunk: chunk,
			Done:  false,
		}); err != nil {
			slog.Error("websocket write error", "error", err)
			return
		}
		time.Sleep(50 * time.Millisecond)
	}

	final := WebSocketResponse{
		Type:        "onboarding",
		Content:     resp.Message,
		Phase:       resp.Phase,
		Level:       resp.Level,
		Goals:       resp.Goals,
		Weaknesses:  resp.Weaknesses,
		IsCompleted: resp.IsCompleted,
		Done:        true,
	}
	if resp.Plan != nil {
		planJSON, _ := json.Marshal(resp.Plan)
		final.Plan = string(planJSON)
	}
	if err := c.WriteJSON(final); err != nil {
		slog.Error("websocket write error", "error", err)
	}
}

func (h *WebSocketHandler) handleDiagnosticChat(c *websocket.Conn, userID uint, req WebSocketMessage) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	resp, err := h.aiService.DiagnosticChat(ctx, userID, req.ProfileID, req.Message)
	if err != nil {
		h.sendError(c, err.Error(), false)
		return
	}

	words := splitIntoChunks(resp.Message, 5)
	for _, chunk := range words {
		response := WebSocketResponse{
			Type:  "diagnostic",
			Chunk: chunk,
			Done:  false,
		}
		if err := c.WriteJSON(response); err != nil {
			slog.Error("websocket write error", "error", err)
			return
		}
		time.Sleep(50 * time.Millisecond)
	}

	finalResp := WebSocketResponse{
		Type:         "diagnostic",
		Content:      resp.Message,
		Level:        resp.Level,
		Weaknesses:   resp.Weaknesses,
		IsAssessment: resp.IsAssessment,
		Done:         true,
	}
	if err := c.WriteJSON(finalResp); err != nil {
		slog.Error("websocket write error", "error", err)
	}
}

func (h *WebSocketHandler) streamSentences(c *websocket.Conn, msgType, text string) {
	for _, chunk := range SplitBySentences(text) {
		if err := c.WriteJSON(WebSocketResponse{
			Type: msgType, Chunk: chunk, Done: false,
		}); err != nil {
			slog.Error("websocket write error", "error", err)
			return
		}
		time.Sleep(30 * time.Millisecond)
	}
}

func (h *WebSocketHandler) handleVoiceTutorTurn(c *websocket.Conn, userID uint, req WebSocketMessage) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	resp, err := h.voiceService.TutorTurn(ctx, userID, req.ProfileID, req.Message, req.Topic, req.Locale, nil)
	if err != nil {
		h.sendError(c, err.Error(), false)
		return
	}

	h.streamSentences(c, "voice_tutor_chunk", resp.Assistant)
	_ = c.WriteJSON(WebSocketResponse{
		Type:        "voice_tutor_final",
		Content:     resp.Assistant,
		UserText:    resp.UserText,
		Corrections: resp.Corrections,
		Done:        true,
	})
}

func (h *WebSocketHandler) handleVoiceTutorAudio(c *websocket.Conn, userID uint, req WebSocketMessage) {
	if !req.Final {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()

	var blob *speech.AudioBlob
	if req.AudioB64 != "" {
		data, err := audio.DecodeBase64(req.AudioB64)
		if err != nil {
			h.sendError(c, "invalid audio_b64", false)
			return
		}
		mime := req.AudioMIME
		if mime == "" {
			mime = "audio/webm;codecs=opus"
		}
		if err := audio.Validate(mime, len(data)); err != nil {
			h.sendError(c, err.Error(), false)
			return
		}
		blob = &speech.AudioBlob{Data: data, MIMEType: mime}
	}

	resp, err := h.voiceService.TutorTurn(ctx, userID, req.ProfileID, req.Message, req.Topic, req.Locale, blob)
	if err != nil {
		h.sendError(c, err.Error(), false)
		return
	}
	h.streamSentences(c, "voice_tutor_chunk", resp.Assistant)
	_ = c.WriteJSON(WebSocketResponse{
		Type:        "voice_tutor_final",
		Content:     resp.Assistant,
		UserText:    resp.UserText,
		Corrections: resp.Corrections,
		Done:        true,
	})
}

func (h *WebSocketHandler) handleVoiceLessonStart(c *websocket.Conn, userID uint, req WebSocketMessage) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	resp, err := h.voiceService.LessonStart(ctx, userID, req.ProfileID, req.ModuleID, req.Topic, req.Locale)
	if err != nil {
		h.sendError(c, err.Error(), false)
		return
	}
	h.streamSentences(c, "voice_lesson_stage", resp.Prompt)
	_ = c.WriteJSON(WebSocketResponse{
		Type:          "voice_lesson_stage",
		LessonID:      resp.LessonID,
		Stage:         resp.Stage,
		StageIndex:    resp.StageIndex,
		StageTotal:    resp.StageTotal,
		Prompt:        resp.Prompt,
		TargetPhrases: resp.TargetPhrases,
		HitPhrases:    resp.HitPhrases,
		Done:          true,
	})
}

func (h *WebSocketHandler) handleVoiceLessonTurn(c *websocket.Conn, userID uint, req WebSocketMessage) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	lessonID := req.LessonID
	if lessonID == 0 {
		h.sendError(c, "lesson_id required", false)
		return
	}
	resp, err := h.voiceService.LessonTurn(ctx, userID, lessonID, req.Message, req.Locale)
	if err != nil {
		h.sendError(c, err.Error(), false)
		return
	}

	if resp.Stage == "completed" && resp.Score != nil {
		h.streamSentences(c, "voice_lesson_stage", resp.Prompt)
		_ = c.WriteJSON(WebSocketResponse{
			Type:     "voice_lesson_complete",
			LessonID: resp.LessonID,
			Stage:    resp.Stage,
			Prompt:   resp.Prompt,
			Score:    resp.Score,
			Done:     true,
		})
		return
	}

	h.streamSentences(c, "voice_lesson_stage", resp.Prompt)
	_ = c.WriteJSON(WebSocketResponse{
		Type:          "voice_lesson_stage",
		LessonID:      resp.LessonID,
		Stage:         resp.Stage,
		StageIndex:    resp.StageIndex,
		StageTotal:    resp.StageTotal,
		Prompt:        resp.Prompt,
		TargetPhrases: resp.TargetPhrases,
		HitPhrases:    resp.HitPhrases,
		Done:          true,
	})
}

func (h *WebSocketHandler) handleVoiceLessonAbandon(c *websocket.Conn, userID uint, req WebSocketMessage) {
	if req.LessonID == 0 {
		h.sendError(c, "lesson_id required", false)
		return
	}
	if err := h.voiceService.LessonAbandon(context.Background(), userID, req.LessonID); err != nil {
		h.sendError(c, err.Error(), false)
		return
	}
	_ = c.WriteJSON(WebSocketResponse{Type: "voice_lesson_abandon", Done: true})
}

func (h *WebSocketHandler) handleTutorChat(c *websocket.Conn, userID uint, req WebSocketMessage) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	resp, err := h.aiService.TutorChat(ctx, userID, req.ProfileID, req.Message, req.Topic)
	if err != nil {
		h.sendError(c, err.Error(), false)
		return
	}

	words := splitIntoChunks(resp.Message, 5)
	for _, chunk := range words {
		response := WebSocketResponse{
			Type:  "tutor",
			Chunk: chunk,
			Done:  false,
		}
		if err := c.WriteJSON(response); err != nil {
			slog.Error("websocket write error", "error", err)
			return
		}
		time.Sleep(50 * time.Millisecond)
	}

	finalResp := WebSocketResponse{
		Type:    "tutor",
		Content: resp.Message,
		Done:    true,
	}
	if err := c.WriteJSON(finalResp); err != nil {
		slog.Error("websocket write error", "error", err)
	}
}

func (h *WebSocketHandler) handleGeneratePlan(c *websocket.Conn, userID uint, req WebSocketMessage) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	resp, err := h.aiService.GeneratePlan(ctx, userID, req.ProfileID)
	if err != nil {
		h.sendError(c, err.Error(), false)
		return
	}

	response := WebSocketResponse{
		Type: "plan",
		Done: true,
	}

	modulesJSON, _ := json.Marshal(resp.Modules)
	response.Content = string(modulesJSON)

	if err := c.WriteJSON(response); err != nil {
		slog.Error("websocket write error", "error", err)
	}
}

func (h *WebSocketHandler) sendError(c *websocket.Conn, errMsg string, closeConn bool) {
	response := WebSocketResponse{
		Type:  "error",
		Error: errMsg,
		Done:  true,
	}
	if err := c.WriteJSON(response); err != nil {
		slog.Error("websocket write error", "error", err)
	}
	if closeConn {
		c.Close()
	}
}

func splitIntoChunks(text string, chunkSize int) []string {
	words := splitWords(text)
	chunks := make([]string, 0, len(words)/chunkSize+1)
	for i := 0; i < len(words); i += chunkSize {
		end := i + chunkSize
		if end > len(words) {
			end = len(words)
		}
		chunk := joinWords(words[i:end])
		chunks = append(chunks, chunk)
	}
	return chunks
}

func splitWords(text string) []string {
	words := []string{}
	var current strings.Builder
	for _, r := range text {
		if r == ' ' || r == '\n' || r == '\t' {
			if current.Len() > 0 {
				words = append(words, current.String())
				current.Reset()
			}
			words = append(words, string(r))
		} else {
			current.WriteRune(r)
		}
	}
	if current.Len() > 0 {
		words = append(words, current.String())
	}
	return words
}

func joinWords(words []string) string {
	var result strings.Builder
	for _, w := range words {
		result.WriteString(w)
	}
	return result.String()
}

func wsAuthMiddleware() fiber.Handler {
	return func(c fiber.Ctx) error {
		token := extractWSToken(c)
		if token != "" {
			c.Request().Header.Set("Authorization", "Bearer "+token)
		}
		return c.Next()
	}
}

// extractWSToken prefers Sec-WebSocket-Protocol (bearer,<jwt>) over ?token=
// so long-lived JWTs are less likely to appear in access logs / Referer.
func extractWSToken(c fiber.Ctx) string {
	if proto := c.Get("Sec-WebSocket-Protocol"); proto != "" {
		parts := strings.Split(proto, ",")
		for i := range parts {
			parts[i] = strings.TrimSpace(parts[i])
		}
		if len(parts) >= 2 && strings.EqualFold(parts[0], "bearer") && parts[1] != "" {
			return parts[1]
		}
		if len(parts) == 1 && parts[0] != "" && !strings.EqualFold(parts[0], "bearer") {
			return parts[0]
		}
	}
	return strings.TrimSpace(c.Query("token"))
}

func (h *WebSocketHandler) RegisterRoutes(api fiber.Router, authMiddleware fiber.Handler) {
	api.Get("/ws", wsAuthMiddleware(), authMiddleware, h.HandleWebSocket())
}
