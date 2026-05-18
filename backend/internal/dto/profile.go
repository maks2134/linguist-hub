package dto

type SelectLanguageRequest struct {
	Language string `json:"language"`
}

type ProfileResponse struct {
	ID                 uint    `json:"id"`
	UserID             uint    `json:"user_id"`
	TargetLanguage     string  `json:"target_language"`
	CurrentLevel       string  `json:"current_level"`
	ProgressPercentage float64 `json:"progress_percentage"`
}
