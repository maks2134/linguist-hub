package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/makskozlov/linguist-hub/internal/dto"
	"github.com/makskozlov/linguist-hub/internal/services"
	"github.com/makskozlov/linguist-hub/pkg/httputil"
)

type AIHandler struct {
	aiService *services.AIService
}

func NewAIHandler(aiService *services.AIService) *AIHandler {
	return &AIHandler{aiService: aiService}
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
func (h *AIHandler) DiagnosticChat(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	var req dto.ChatRequest
	if err := httputil.ParseBody(c, &req); err != nil {
		return httputil.RespondError(c, err)
	}

	profileID, err := c.ParamsInt("profileId")
	if err != nil {
		return httputil.RespondError(c, err)
	}

	resp, err := h.aiService.DiagnosticChat(c.Context(), userID, uint(profileID), req.Message)
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
func (h *AIHandler) GeneratePlan(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	profileID, err := c.ParamsInt("profileId")
	if err != nil {
		return httputil.RespondError(c, err)
	}

	resp, err := h.aiService.GeneratePlan(c.Context(), userID, uint(profileID))
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
func (h *AIHandler) TutorChat(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	var req dto.ChatRequest
	if err := httputil.ParseBody(c, &req); err != nil {
		return httputil.RespondError(c, err)
	}

	profileID, err := c.ParamsInt("profileId")
	if err != nil {
		return httputil.RespondError(c, err)
	}

	resp, err := h.aiService.TutorChat(c.Context(), userID, uint(profileID), req.Message, req.Topic)
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
func (h *AIHandler) GetPlan(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	profileID, err := c.ParamsInt("profileId")
	if err != nil {
		return httputil.RespondError(c, err)
	}

	resp, err := h.aiService.GetPlan(userID, uint(profileID))
	if err != nil {
		return httputil.RespondError(c, err)
	}

	return httputil.RespondJSON(c, 200, resp)
}

func (h *AIHandler) RegisterRoutes(api fiber.Router, authMiddleware fiber.Handler) {
	ai := api.Group("/ai", authMiddleware)
	ai.Post("/diagnostic/:profileId", h.DiagnosticChat)
	ai.Post("/generate-plan/:profileId", h.GeneratePlan)
	ai.Get("/plan/:profileId", h.GetPlan)

	chat := api.Group("/chat", authMiddleware)
	chat.Post("/tutor/:profileId", h.TutorChat)
}
