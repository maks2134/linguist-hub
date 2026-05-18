package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type JSONB []byte

func (j JSONB) Value() (driver.Value, error) {
	if len(j) == 0 {
		return nil, nil
	}
	return string(j), nil
}

func (j *JSONB) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}
	s, ok := value.([]byte)
	if !ok {
		return errors.New("invalid Scan source")
	}
	*j = s
	return nil
}

func (j JSONB) MarshalJSON() ([]byte, error) {
	if j == nil {
		return []byte("null"), nil
	}
	return j, nil
}

func (j *JSONB) UnmarshalJSON(data []byte) error {
	if j == nil {
		return errors.New("null pointer")
	}
	*j = append((*j)[0:0], data...)
	return nil
}

type StudyPlan struct {
	ID          uint            `gorm:"primaryKey" json:"id"`
	ProfileID   uint            `gorm:"index;not null" json:"profile_id"`
	Structure   JSONB           `gorm:"type:jsonb" json:"structure"`
	IsGenerated bool            `gorm:"default:false" json:"is_generated"`
	Profile     LanguageProfile `gorm:"foreignKey:ProfileID" json:"-"`
}

type Module struct {
	ID          int      `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Topics      []string `json:"topics"`
}

func (p *StudyPlan) GetModules() ([]Module, error) {
	if len(p.Structure) == 0 {
		return nil, nil
	}
	var modules []Module
	if err := json.Unmarshal(p.Structure, &modules); err != nil {
		return nil, err
	}
	return modules, nil
}

func (p *StudyPlan) SetModules(modules []Module) error {
	data, err := json.Marshal(modules)
	if err != nil {
		return err
	}
	p.Structure = data
	return nil
}
