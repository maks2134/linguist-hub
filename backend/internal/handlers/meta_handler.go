package handlers

import (
	"github.com/gofiber/fiber/v3"
	"github.com/makskozlov/linguist-hub/internal/dto"
	"github.com/makskozlov/linguist-hub/pkg/httputil"
)

type MetaHandler struct {
	serverSTT bool
}

func NewMetaHandler(serverSTT bool) *MetaHandler {
	return &MetaHandler{serverSTT: serverSTT}
}

func (h *MetaHandler) RegisterRoutes(api fiber.Router) {
	api.Get("/capabilities", h.Capabilities)
}

func (h *MetaHandler) Capabilities(c fiber.Ctx) error {
	return httputil.RespondJSON(c, 200, dto.CapabilitiesResponse{
		ServerSTT: h.serverSTT,
	})
}
