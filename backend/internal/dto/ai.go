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
