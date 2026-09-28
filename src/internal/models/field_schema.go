package models

import "encoding/json"

// FieldSchema — definisi 1 field dalam template
type FieldSchema struct {
	Name           string                 `json:"name"`
	Label          string                 `json:"label"`
	Type           string                 `json:"type"` // text, textarea, number, url, select, file, file_multiple, repeater
	Required       bool                   `json:"required"`
	Placeholder    string                 `json:"placeholder,omitempty"`
	Accept         string                 `json:"accept,omitempty"`
	Group          string                 `json:"group,omitempty"`
	Description    string                 `json:"description,omitempty"`
	Options        map[string]string      `json:"options,omitempty"`
	Fields         []FieldSchema          `json:"fields,omitempty"`
	Rows           int                    `json:"rows,omitempty"`
	DefaultCount   int                    `json:"default_count,omitempty"`
	DefaultStories []map[string]string    `json:"default_stories,omitempty"`
}

// ParseFieldSchema — parse JSON string ke []FieldSchema
func ParseFieldSchema(jsonStr string) []FieldSchema {
	if jsonStr == "" {
		return []FieldSchema{}
	}

	var schema []FieldSchema
	if err := json.Unmarshal([]byte(jsonStr), &schema); err != nil {
		return []FieldSchema{}
	}
	return schema
}

// GetFieldsSchema — method untuk Template
func (t Template) GetFieldsSchema() []FieldSchema {
	return ParseFieldSchema(t.FieldsSchema)
}

// HasSpecificFields — cek apakah template punya field dinamis
func (t Template) HasSpecificFields() bool {
	return len(t.GetFieldsSchema()) > 0
}

// ============================================
// FIELD SCHEMA HELPERS
// ============================================

// GroupFieldsByType — kelompokkan field berdasarkan tipe untuk render
func GroupFieldsByType(schema []FieldSchema) map[string][]FieldSchema {
	groups := map[string][]FieldSchema{
		"text":     {},
		"file":     {},
		"gallery":  {},
		"repeater": {},
	}

	for _, field := range schema {
		switch field.Type {
		case "file":
			groups["file"] = append(groups["file"], field)
		case "file_multiple":
			groups["gallery"] = append(groups["gallery"], field)
		case "repeater":
			groups["repeater"] = append(groups["repeater"], field)
		default:
			groups["text"] = append(groups["text"], field)
		}
	}

	return groups
}
