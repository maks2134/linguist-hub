package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v3"
	"github.com/makskozlov/linguist-hub/internal/dto"
	"github.com/makskozlov/linguist-hub/internal/services"
	"github.com/makskozlov/linguist-hub/internal/services/speech"
	"github.com/makskozlov/linguist-hub/pkg/httputil"
)

type VoiceHandler struct {
	voiceService  *services.VoiceService
	sttRecognizer speech.Recognizer
}

func NewVoiceHandler(voiceService *services.VoiceService) *VoiceHandler {
	return &VoiceHandler{voiceService: voiceService}
}

func NewVoiceHandlerWithSTT(voiceService *services.VoiceService, recognizer speech.Recognizer) *VoiceHandler {
	return &VoiceHandler{voiceService: voiceService, sttRecognizer: recognizer}
}

func (h *VoiceHandler) RegisterRoutes(
	api fiber.Router,
	authMiddleware fiber.Handler,
	aiDaily fiber.Handler,
	voiceRateLimit fiber.Handler,
) {
	voice := api.Group("/voice", authMiddleware)
	voice.Get("/locales/:profileId", h.GetLocales)
	voice.Get("/lessons/:profileId", h.ListLessons)
	voice.Post("/tutor/:profileId", aiDaily, voiceRateLimit, h.TutorTurn)
	voice.Post("/lessons/:profileId/start", aiDaily, voiceRateLimit, h.LessonStart)
	voice.Post("/lessons/:profileId/turn", aiDaily, voiceRateLimit, h.LessonTurn)
	voice.Post("/lessons/:profileId/:lessonId/abandon", voiceRateLimit, h.LessonAbandon)
	voice.Post("/stt", aiDaily, voiceRateLimit, h.STT)
}

func (h *VoiceHandler) GetLocales(c fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	profileID, err := fiber.Params[int](c, "profileId"), error(nil)
	if err != nil {
		return httputil.RespondError(c, err)
	}
	resp, err := h.voiceService.GetLocales(userID, uint(profileID))
	if err != nil {
		return httputil.RespondError(c, err)
	}
	return httputil.RespondJSON(c, 200, resp)
}

func (h *VoiceHandler) TutorTurn(c fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	profileID, err := fiber.Params[int](c, "profileId"), error(nil)
	if err != nil {
		return httputil.RespondError(c, err)
	}
	var req dto.VoiceTutorRequest
	if err := httputil.ParseBody(c, &req); err != nil {
		return httputil.RespondError(c, err)
	}
	resp, err := h.voiceService.TutorTurn(c.RequestCtx(), userID, uint(profileID), req.Message, req.Topic, req.Locale, nil)
	if err != nil {
		return httputil.RespondError(c, err)
	}
	return httputil.RespondJSON(c, 200, resp)
}

func (h *VoiceHandler) LessonStart(c fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	profileID, err := fiber.Params[int](c, "profileId"), error(nil)
	if err != nil {
		return httputil.RespondError(c, err)
	}
	var req dto.LessonStartRequest
	if err := httputil.ParseBody(c, &req); err != nil {
		return httputil.RespondError(c, err)
	}
	resp, err := h.voiceService.LessonStart(c.RequestCtx(), userID, uint(profileID), req.ModuleID, req.Topic, req.Locale)
	if err != nil {
		return httputil.RespondError(c, err)
	}
	return httputil.RespondJSON(c, 200, resp)
}

func (h *VoiceHandler) LessonTurn(c fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	profileID, err := fiber.Params[int](c, "profileId"), error(nil)
	if err != nil {
		return httputil.RespondError(c, err)
	}
	_ = profileID
	var req dto.LessonTurnRequest
	if err := httputil.ParseBody(c, &req); err != nil {
		return httputil.RespondError(c, err)
	}
	resp, err := h.voiceService.LessonTurn(c.RequestCtx(), userID, req.LessonID, req.Message, req.Locale)
	if err != nil {
		return httputil.RespondError(c, err)
	}
	return httputil.RespondJSON(c, 200, resp)
}

func (h *VoiceHandler) LessonAbandon(c fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	lessonID, err := fiber.Params[int](c, "lessonId"), error(nil)
	if err != nil {
		return httputil.RespondError(c, err)
	}
	if err := h.voiceService.LessonAbandon(c.RequestCtx(), userID, uint(lessonID)); err != nil {
		return httputil.RespondError(c, err)
	}
	return c.SendStatus(204)
}

func (h *VoiceHandler) ListLessons(c fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	profileID, err := fiber.Params[int](c, "profileId"), error(nil)
	if err != nil {
		return httputil.RespondError(c, err)
	}
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))
	resp, err := h.voiceService.ListLessons(userID, uint(profileID), limit, offset)
	if err != nil {
		return httputil.RespondError(c, err)
	}
	return httputil.RespondJSON(c, 200, resp)
}

// STT handles server-side speech-to-text when FEATURE_SERVER_STT is enabled.
func (h *VoiceHandler) STT(c fiber.Ctx) error {
	if h.sttRecognizer == nil {
		return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "STT_NOT_ENABLED",
				"message": "Server STT is not enabled; use browser Web Speech API",
			},
		})
	}
	locale := c.FormValue("locale")
	file, err := c.FormFile("audio")
	if err != nil {
		return httputil.RespondError(c, err)
	}
	f, err := file.Open()
	if err != nil {
		return httputil.RespondError(c, err)
	}
	defer f.Close()
	data := make([]byte, file.Size)
	if _, err := f.Read(data); err != nil {
		return httputil.RespondError(c, err)
	}
	mime := file.Header.Get("Content-Type")
	transcript, err := h.sttRecognizer.Recognize(c.RequestCtx(), speech.AudioBlob{
		Data: data, MIMEType: mime,
	}, locale)
	if err != nil {
		return httputil.RespondError(c, err)
	}
	return httputil.RespondJSON(c, 200, dto.STTResponse{
		Text: transcript.Text, Confidence: transcript.Confidence,
	})
}
