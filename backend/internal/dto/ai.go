package dto

type ChatRequest struct {
	Message string `json:"message"`
	Topic   string `json:"topic,omitempty"`
}

type ChatResponse struct {
	Message      string   `json:"message"`
	Level        string   `json:"level,omitempty"`
	Weaknesses   []string `json:"weaknesses,omitempty"`
	IsAssessment bool     `json:"is_assessment,omitempty"`
}

// ChatHistoryMessage represents a single chat message in history
type ChatHistoryMessage struct {
	ID        uint   `json:"id"`
	UserID    uint   `json:"user_id"`
	ProfileID uint   `json:"profile_id,omitempty"`
	Role      string `json:"role"`
	Context   string `json:"context"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
}

type GeneratePlanRequest struct {
	ProfileID uint `json:"profile_id"`
}

type PlanResponse struct {
	ID        uint     `json:"id"`
	ProfileID uint     `json:"profile_id"`
	Modules   []Module `json:"modules"`
}

type Module struct {
	ID          int      `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Topics      []string `json:"topics"`
}
