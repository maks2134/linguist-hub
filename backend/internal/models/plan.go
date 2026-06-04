package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type Modules []Module

type StudyPlan struct {
	ID          uint            `gorm:"primaryKey" json:"id"`
	ProfileID   uint            `gorm:"index;not null" json:"profile_id"`
	Structure   Modules         `gorm:"type:jsonb" json:"structure"`
	IsGenerated bool            `gorm:"default:false" json:"is_generated"`
	Profile     LanguageProfile `gorm:"foreignKey:ProfileID" json:"-"`
}

type Module struct {
	ID          int      `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Topics      []string `json:"topics"`
}

func (m Modules) Value() (driver.Value, error) {
	return json.Marshal(m)
}

func (m *Modules) Scan(value any) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, m)
}
