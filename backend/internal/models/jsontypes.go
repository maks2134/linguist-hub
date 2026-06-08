package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// JSONStringList stores a []string in JSONB columns.
type JSONStringList []string

func (l JSONStringList) Value() (driver.Value, error) {
	if l == nil {
		return json.Marshal([]string{})
	}
	return json.Marshal(l)
}

func (l *JSONStringList) Scan(value any) error {
	if value == nil {
		*l = JSONStringList{}
		return nil
	}
	var b []byte
	switch v := value.(type) {
	case []byte:
		b = v
	case string:
		b = []byte(v)
	default:
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, l)
}
