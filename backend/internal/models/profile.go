package models

type LanguageProfile struct {
	ID                 uint    `gorm:"primaryKey" json:"id"`
	UserID             uint    `gorm:"index;not null" json:"user_id"`
	TargetLanguage     string  `gorm:"not null" json:"target_language"`
	CurrentLevel       string  `gorm:"default:'unknown'" json:"current_level"`
	ProgressPercentage float64 `gorm:"default:0" json:"progress_percentage"`
	User               User    `gorm:"foreignKey:UserID" json:"-"`
}
