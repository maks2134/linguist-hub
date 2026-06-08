package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

const (
	VoiceLessonStatusInProgress = "in_progress"
	VoiceLessonStatusCompleted  = "completed"
	VoiceLessonStatusAbandoned  = "abandoned"
)

type VoiceScore struct {
	Fluency       float32 `json:"fluency"`
	Grammar       float32 `json:"grammar"`
	Vocabulary    float32 `json:"vocabulary"`
	Pronunciation float32 `json:"pronunciation"`
	Notes         string  `json:"notes"`
}

func (s VoiceScore) Value() (driver.Value, error) {
	return json.Marshal(s)
}

func (s *VoiceScore) Scan(value any) error {
	if value == nil {
		return nil
	}
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, s)
}

type VoiceLesson struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	UserID        uint           `gorm:"index;not null" json:"user_id"`
	ProfileID     uint           `gorm:"index;not null" json:"profile_id"`
	PlanID        uint           `gorm:"index" json:"plan_id,omitempty"`
	ModuleID      int            `json:"module_id"`
	Topic         string         `gorm:"not null" json:"topic"`
	Stage         string         `gorm:"not null;default:warmup" json:"stage"`
	StageIndex    int            `gorm:"default:0" json:"stage_index"`
	TargetPhrases JSONStringList `gorm:"type:jsonb;not null;default:'[]'" json:"target_phrases"`
	HitPhrases    JSONStringList `gorm:"type:jsonb;not null;default:'[]'" json:"hit_phrases"`
	Score         *VoiceScore    `gorm:"type:jsonb" json:"score,omitempty"`
	Status        string         `gorm:"not null;default:in_progress" json:"status"`
	StartedAt     time.Time      `json:"started_at"`
	CompletedAt   *time.Time     `json:"completed_at,omitempty"`
	LastTurnAt    time.Time      `gorm:"index" json:"last_turn_at"`
}
