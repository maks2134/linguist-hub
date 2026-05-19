package handlers

import (
	"github.com/gofiber/fiber/v3"
	"github.com/makskozlov/linguist-hub/internal/dto"
	"github.com/makskozlov/linguist-hub/internal/services"
	"github.com/makskozlov/linguist-hub/pkg/httputil"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// Register godoc
// @Summary Register a new user
// @Description Create a new user account
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "Registration credentials"
// @Success 201 {object} httputil.Response{data=dto.AuthResponse}
// @Failure 400 {object} httputil.Response
// @Router /auth/register [post]
func (h *AuthHandler) Register(c fiber.Ctx) error {
	var req dto.RegisterRequest
	if err := httputil.ParseBody(c, &req); err != nil {
		return httputil.RespondError(c, err)
	}

	resp, err := h.authService.Register(req.Email, req.Password)
	if err != nil {
		return httputil.RespondError(c, err)
	}

	return httputil.RespondJSON(c, 201, resp)
}

// Login godoc
// @Summary Login user
// @Description Authenticate user and return JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Login credentials"
// @Success 200 {object} httputil.Response{data=dto.AuthResponse}
// @Failure 401 {object} httputil.Response
// @Router /auth/login [post]
func (h *AuthHandler) Login(c fiber.Ctx) error {
	var req dto.LoginRequest
	if err := httputil.ParseBody(c, &req); err != nil {
		return httputil.RespondError(c, err)
	}

	resp, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		return httputil.RespondError(c, err)
	}

	return httputil.RespondJSON(c, 200, resp)
}

// UpdateProfile godoc
// @Summary Update user profile
// @Description Update user nickname
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.UpdateProfileRequest true "Profile data"
// @Success 200 {object} httputil.Response{data=dto.User}
// @Failure 400 {object} httputil.Response
// @Failure 401 {object} httputil.Response
// @Router /auth/profile [put]
func (h *AuthHandler) UpdateProfile(c fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	var req dto.UpdateProfileRequest
	if err := httputil.ParseBody(c, &req); err != nil {
		return httputil.RespondError(c, err)
	}

	resp, err := h.authService.UpdateProfile(userID, req.Nickname)
	if err != nil {
		return httputil.RespondError(c, err)
	}

	return httputil.RespondJSON(c, 200, resp)
}

// ChangePassword godoc
// @Summary Change user password
// @Description Change authenticated user's password
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.ChangePasswordRequest true "Password data"
// @Success 200 {object} httputil.Response
// @Failure 400 {object} httputil.Response
// @Failure 401 {object} httputil.Response
// @Router /auth/password [put]
func (h *AuthHandler) ChangePassword(c fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	var req dto.ChangePasswordRequest
	if err := httputil.ParseBody(c, &req); err != nil {
		return httputil.RespondError(c, err)
	}

	if err := h.authService.ChangePassword(userID, req.OldPassword, req.NewPassword); err != nil {
		return httputil.RespondError(c, err)
	}

	return httputil.RespondJSON(c, 200, fiber.Map{"message": "Password updated successfully"})
}

// GetMe godoc
// @Summary Get current user info
// @Description Get authenticated user's profile information
// @Tags auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} httputil.Response{data=dto.User}
// @Failure 401 {object} httputil.Response
// @Router /auth/me [get]
func (h *AuthHandler) GetMe(c fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	resp, err := h.authService.GetMe(userID)
	if err != nil {
		return httputil.RespondError(c, err)
	}

	return httputil.RespondJSON(c, 200, resp)
}

func (h *AuthHandler) RegisterRoutes(api fiber.Router, authMiddleware fiber.Handler) {
	auth := api.Group("/auth")
	auth.Post("/register", h.Register)
	auth.Post("/login", h.Login)

	// Protected routes
	auth.Get("/me", authMiddleware, h.GetMe)
	auth.Put("/profile", authMiddleware, h.UpdateProfile)
	auth.Put("/password", authMiddleware, h.ChangePassword)
}
