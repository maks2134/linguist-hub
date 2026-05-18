package handlers

import (
	"github.com/gofiber/fiber/v2"
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
func (h *AuthHandler) Register(c *fiber.Ctx) error {
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
func (h *AuthHandler) Login(c *fiber.Ctx) error {
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

func (h *AuthHandler) RegisterRoutes(api fiber.Router) {
	auth := api.Group("/auth")
	auth.Post("/register", h.Register)
	auth.Post("/login", h.Login)
}
