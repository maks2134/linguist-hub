package handlers

import (
	"github.com/gofiber/fiber/v3"
	"github.com/makskozlov/linguist-hub/internal/dto"
	"github.com/makskozlov/linguist-hub/internal/services"
	"github.com/makskozlov/linguist-hub/pkg/httputil"
)

type ProfileHandler struct {
	profileService *services.ProfileService
}

func NewProfileHandler(profileService *services.ProfileService) *ProfileHandler {
	return &ProfileHandler{profileService: profileService}
}

// SelectLanguage godoc
// @Summary Select learning language
// @Description Create or return existing language profile
// @Tags profile
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.SelectLanguageRequest true "Language selection"
// @Success 200 {object} httputil.Response{data=dto.ProfileResponse}
// @Failure 400 {object} httputil.Response
// @Failure 401 {object} httputil.Response
// @Router /profile/select [post]
func (h *ProfileHandler) SelectLanguage(c fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	var req dto.SelectLanguageRequest
	if err := httputil.ParseBody(c, &req); err != nil {
		return httputil.RespondError(c, err)
	}

	resp, err := h.profileService.SelectLanguage(userID, req.Language, req.NativeLanguage)
	if err != nil {
		return httputil.RespondError(c, err)
	}

	return httputil.RespondJSON(c, 200, resp)
}

// GetProfiles godoc
// @Summary Get user profiles
// @Description Get all language profiles for authenticated user
// @Tags profile
// @Produce json
// @Security BearerAuth
// @Success 200 {object} httputil.Response{data=[]dto.ProfileResponse}
// @Failure 401 {object} httputil.Response
// @Router /profile/ [get]
func (h *ProfileHandler) GetProfiles(c fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	profiles, err := h.profileService.GetProfiles(userID)
	if err != nil {
		return httputil.RespondError(c, err)
	}

	return httputil.RespondJSON(c, 200, profiles)
}

// GetProfile godoc
// @Summary Get profile by ID
// @Description Get specific language profile details
// @Tags profile
// @Produce json
// @Security BearerAuth
// @Param id path int true "Profile ID"
// @Success 200 {object} httputil.Response{data=dto.ProfileResponse}
// @Failure 401 {object} httputil.Response
// @Failure 403 {object} httputil.Response
// @Failure 404 {object} httputil.Response
// @Router /profile/{id} [get]
func (h *ProfileHandler) GetProfile(c fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	profileID, err := fiber.Params[int](c, "id"), error(nil)
	if err != nil {
		return httputil.RespondError(c, err)
	}

	profile, err := h.profileService.GetProfile(userID, uint(profileID))
	if err != nil {
		return httputil.RespondError(c, err)
	}

	return httputil.RespondJSON(c, 200, profile)
}

// UpdateNativeLanguage godoc
// @Summary Update native language of a profile
// @Description Change the native language used for AI explanations
// @Tags profile
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Profile ID"
// @Param request body dto.UpdateNativeLanguageRequest true "Native language"
// @Success 200 {object} httputil.Response{data=dto.ProfileResponse}
// @Failure 400 {object} httputil.Response
// @Failure 401 {object} httputil.Response
// @Failure 403 {object} httputil.Response
// @Router /profile/{id}/native-language [put]
func (h *ProfileHandler) UpdateNativeLanguage(c fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	profileID, err := fiber.Params[int](c, "id"), error(nil)
	if err != nil {
		return httputil.RespondError(c, err)
	}

	var req dto.UpdateNativeLanguageRequest
	if err := httputil.ParseBody(c, &req); err != nil {
		return httputil.RespondError(c, err)
	}

	resp, err := h.profileService.UpdateNativeLanguage(userID, uint(profileID), req.NativeLanguage)
	if err != nil {
		return httputil.RespondError(c, err)
	}

	return httputil.RespondJSON(c, 200, resp)
}

func (h *ProfileHandler) RegisterRoutes(api fiber.Router, authMiddleware fiber.Handler) {
	profile := api.Group("/profile", authMiddleware)
	profile.Post("/select", h.SelectLanguage)
	profile.Get("/", h.GetProfiles)
	profile.Get("/:id", h.GetProfile)
	profile.Put("/:id/native-language", h.UpdateNativeLanguage)
}
