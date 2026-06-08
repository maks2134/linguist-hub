package models

type LanguageProfile struct {
	ID                 uint           `gorm:"primaryKey" json:"id"`
	UserID             uint           `gorm:"index;not null" json:"user_id"`
	TargetLanguage     string         `gorm:"not null" json:"target_language"`
	NativeLanguage     string         `gorm:"default:'english'" json:"native_language"`
	CurrentLevel       string         `gorm:"default:'unknown'" json:"current_level"`
	OnboardingPhase    string         `gorm:"default:'goals'" json:"onboarding_phase"`
	Goals              string         `gorm:"type:text;default:''" json:"goals"`
	Weaknesses         JSONStringList `gorm:"type:jsonb;not null;default:'[]'" json:"weaknesses"`
	ProgressPercentage float64        `gorm:"default:0" json:"progress_percentage"`
	User               User           `gorm:"foreignKey:UserID" json:"-"`
}
