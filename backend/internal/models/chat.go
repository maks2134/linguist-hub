package models

import "time"

type ChatMessage struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"index;not null" json:"user_id"`
	ProfileID uint      `gorm:"index" json:"profile_id,omitempty"`
	Role      string    `gorm:"not null" json:"role"`
	Context   string    `gorm:"not null" json:"context"`
	Content   string    `gorm:"type:text;not null" json:"content"`
	CreatedAt time.Time `json:"created_at"`
}
