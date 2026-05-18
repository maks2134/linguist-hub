package handlers

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/makskozlov/linguist-hub/internal/services"
)

type WebSocketMessage struct {
	Type      string `json:"type"`
	Message   string `json:"message"`
	Topic     string `json:"topic"`
	ProfileID uint   `json:"profile_id"`
}

type WebSocketResponse struct {
	Type         string   `json:"type"`
	Content      string   `json:"content,omitempty"`
	Chunk        string   `json:"chunk,omitempty"`
	Level        string   `json:"level,omitempty"`
	Weaknesses   []string `json:"weaknesses,omitempty"`
	IsAssessment bool     `json:"is_assessment,omitempty"`
	Error        string   `json:"error,omitempty"`
	Done         bool     `json:"done,omitempty"`
}

type WebSocketHandler struct {
	aiService *services.AIService
}

func NewWebSocketHandler(aiService *services.AIService) *WebSocketHandler {
	return &WebSocketHandler{aiService: aiService}
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
			case "diagnostic":
				h.handleDiagnosticChat(c, userID, req)
			case "tutor":
				h.handleTutorChat(c, userID, req)
			case "generate_plan":
				h.handleGeneratePlan(c, userID, req)
			default:
				h.sendError(c, "unknown message type: "+req.Type, false)
			}
		}
	})
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
	var chunks []string
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
	var words []string
	var current string
	for _, r := range text {
		if r == ' ' || r == '\n' || r == '\t' {
			if current != "" {
				words = append(words, current)
				current = ""
			}
			words = append(words, string(r))
		} else {
			current += string(r)
		}
	}
	if current != "" {
		words = append(words, current)
	}
	return words
}

func joinWords(words []string) string {
	result := ""
	for _, w := range words {
		result += w
	}
	return result
}

func wsAuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := c.Query("token")
		if token != "" {
			c.Request().Header.Set("Authorization", "Bearer "+token)
		}
		return c.Next()
	}
}

func (h *WebSocketHandler) RegisterRoutes(api fiber.Router, authMiddleware fiber.Handler) {
	api.Get("/ws", wsAuthMiddleware(), authMiddleware, h.HandleWebSocket())
}
