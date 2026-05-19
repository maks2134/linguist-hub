package dto

type SelectLanguageRequest struct {
	Language       string `json:"language"`
	NativeLanguage string `json:"native_language,omitempty"`
}

type UpdateNativeLanguageRequest struct {
	NativeLanguage string `json:"native_language"`
}

type ProfileResponse struct {
	ID                 uint    `json:"id"`
	UserID             uint    `json:"user_id"`
	TargetLanguage     string  `json:"target_language"`
	NativeLanguage     string  `json:"native_language"`
	CurrentLevel       string  `json:"current_level"`
	OnboardingPhase    string  `json:"onboarding_phase"`
	Goals              string  `json:"goals,omitempty"`
	ProgressPercentage float64 `json:"progress_percentage"`
}
