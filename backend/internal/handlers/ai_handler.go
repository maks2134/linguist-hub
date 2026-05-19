package handlers

import (
	"github.com/gofiber/fiber/v3"
	"github.com/makskozlov/linguist-hub/internal/dto"
	"github.com/makskozlov/linguist-hub/internal/services"
	"github.com/makskozlov/linguist-hub/pkg/apierror"
	"github.com/makskozlov/linguist-hub/pkg/httputil"
)

type AIHandler struct {
	aiService *services.AIService
}

func NewAIHandler(aiService *services.AIService) *AIHandler {
	return &AIHandler{aiService: aiService}
}

// OnboardingChat godoc
// @Summary AI onboarding chat
// @Description Unified onboarding flow: goals → level check → plan generation, communicated in the user's native language
// @Tags ai
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param profileId path int true "Profile ID"
// @Param request body dto.ChatRequest true "Chat message"
// @Success 200 {object} httputil.Response{data=dto.OnboardingChatResponse}
// @Failure 400 {object} httputil.Response
// @Failure 401 {object} httputil.Response
// @Router /ai/onboarding/{profileId} [post]
func (h *AIHandler) OnboardingChat(c fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	var req dto.ChatRequest
	if err := httputil.ParseBody(c, &req); err != nil {
		return httputil.RespondError(c, err)
	}

	profileID, err := fiber.Params[int](c, "profileId"), error(nil)
	if err != nil {
		return httputil.RespondError(c, err)
	}

	resp, err := h.aiService.OnboardingChat(c.RequestCtx(), userID, uint(profileID), req.Message)
	if err != nil {
		return httputil.RespondError(c, err)
	}

	return httputil.RespondJSON(c, 200, resp)
}

// DiagnosticChat godoc
// @Summary AI diagnostic chat
// @Description Send message to AI for level assessment
// @Tags ai
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param profileId path int true "Profile ID"
// @Param request body dto.ChatRequest true "Chat message"
// @Success 200 {object} httputil.Response{data=dto.ChatResponse}
// @Failure 400 {object} httputil.Response
// @Failure 401 {object} httputil.Response
// @Router /ai/diagnostic/{profileId} [post]
func (h *AIHandler) DiagnosticChat(c fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	var req dto.ChatRequest
	if err := httputil.ParseBody(c, &req); err != nil {
		return httputil.RespondError(c, err)
	}

	profileID, err := fiber.Params[int](c, "profileId"), error(nil)
	if err != nil {
		return httputil.RespondError(c, err)
	}

	resp, err := h.aiService.DiagnosticChat(c.RequestCtx(), userID, uint(profileID), req.Message)
	if err != nil {
		return httputil.RespondError(c, err)
	}

	return httputil.RespondJSON(c, 200, resp)
}

// GeneratePlan godoc
// @Summary Generate study plan
// @Description Generate AI-powered study plan based on user level
// @Tags ai
// @Produce json
// @Security BearerAuth
// @Param profileId path int true "Profile ID"
// @Success 200 {object} httputil.Response{data=dto.PlanResponse}
// @Failure 400 {object} httputil.Response
// @Failure 401 {object} httputil.Response
// @Router /ai/generate-plan/{profileId} [post]
func (h *AIHandler) GeneratePlan(c fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	profileID, err := fiber.Params[int](c, "profileId"), error(nil)
	if err != nil {
		return httputil.RespondError(c, err)
	}

	resp, err := h.aiService.GeneratePlan(c.RequestCtx(), userID, uint(profileID))
	if err != nil {
		return httputil.RespondError(c, err)
	}

	return httputil.RespondJSON(c, 200, resp)
}

// TutorChat godoc
// @Summary AI tutor chat
// @Description Chat with AI tutor for language practice
// @Tags chat
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param profileId path int true "Profile ID"
// @Param request body dto.ChatRequest true "Chat message with topic"
// @Success 200 {object} httputil.Response{data=dto.ChatResponse}
// @Failure 400 {object} httputil.Response
// @Failure 401 {object} httputil.Response
// @Router /chat/tutor/{profileId} [post]
func (h *AIHandler) TutorChat(c fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	var req dto.ChatRequest
	if err := httputil.ParseBody(c, &req); err != nil {
		return httputil.RespondError(c, err)
	}

	profileID, err := fiber.Params[int](c, "profileId"), error(nil)
	if err != nil {
		return httputil.RespondError(c, err)
	}

	resp, err := h.aiService.TutorChat(c.RequestCtx(), userID, uint(profileID), req.Message, req.Topic)
	if err != nil {
		return httputil.RespondError(c, err)
	}

	return httputil.RespondJSON(c, 200, resp)
}

// GetPlan godoc
// @Summary Get study plan
// @Description Get existing study plan for profile
// @Tags ai
// @Produce json
// @Security BearerAuth
// @Param profileId path int true "Profile ID"
// @Success 200 {object} httputil.Response{data=dto.PlanResponse}
// @Failure 401 {object} httputil.Response
// @Failure 404 {object} httputil.Response
// @Router /ai/plan/{profileId} [get]
func (h *AIHandler) GetPlan(c fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	profileID, err := fiber.Params[int](c, "profileId"), error(nil)
	if err != nil {
		return httputil.RespondError(c, err)
	}

	resp, err := h.aiService.GetPlan(userID, uint(profileID))
	if err != nil {
		return httputil.RespondError(c, err)
	}

	return httputil.RespondJSON(c, 200, resp)
}

// GetChatHistory godoc
// @Summary Get chat history
// @Description Get chat message history for a specific context
// @Tags chat
// @Produce json
// @Security BearerAuth
// @Param profileId path int true "Profile ID"
// @Param context query string true "Chat context (diagnostic|tutor)"
// @Param limit query int false "Limit (default: 0 = no limit)"
// @Param offset query int false "Offset (default: 0)"
// @Success 200 {object} httputil.Response{data=[]dto.ChatHistoryMessage}
// @Failure 400 {object} httputil.Response
// @Failure 401 {object} httputil.Response
// @Router /chat/history/{profileId} [get]
func (h *AIHandler) GetChatHistory(c fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	profileID, err := fiber.Params[int](c, "profileId"), error(nil)
	if err != nil {
		return httputil.RespondError(c, err)
	}

	context := c.Query("context")
	if context == "" {
		return httputil.RespondError(c, apierror.New(400, "INVALID_CONTEXT", "context query parameter is required"))
	}

	limit := fiber.Query[int](c, "limit", 0)
	offset := fiber.Query[int](c, "offset", 0)

	history, err := h.aiService.GetChatHistory(userID, uint(profileID), context, limit, offset)
	if err != nil {
		return httputil.RespondError(c, err)
	}

	return httputil.RespondJSON(c, 200, history)
}

func (h *AIHandler) RegisterRoutes(api fiber.Router, authMiddleware fiber.Handler) {
	ai := api.Group("/ai", authMiddleware)
	ai.Post("/onboarding/:profileId", h.OnboardingChat)
	ai.Post("/diagnostic/:profileId", h.DiagnosticChat)
	ai.Post("/generate-plan/:profileId", h.GeneratePlan)
	ai.Get("/plan/:profileId", h.GetPlan)

	chat := api.Group("/chat", authMiddleware)
	chat.Post("/tutor/:profileId", h.TutorChat)
	chat.Get("/history/:profileId", h.GetChatHistory)
}

// fiber:context-methods migrated
