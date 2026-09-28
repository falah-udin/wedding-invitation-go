package models

import (
	"time"

	"gorm.io/gorm"
)

// Template — model template undangan, map ke tabel `templates`
// Struktur mengikuti migration Laravel
type Template struct {
	ID          uint    `gorm:"primaryKey" json:"id"`
	Name        string  `gorm:"size:255;not null" json:"name"`
	Slug        string  `gorm:"size:255;uniqueIndex;not null" json:"slug"`
	Folder      string  `gorm:"size:255;not null" json:"folder"`
	Sections    string  `gorm:"type:json" json:"sections,omitempty"`
	Thumbnail   *string `gorm:"size:255" json:"thumbnail,omitempty"`
	Description *string `gorm:"type:text" json:"description,omitempty"`
	Fields      string  `gorm:"type:json" json:"fields,omitempty"`
	FieldsSchema string `gorm:"column:fields_schema;type:json" json:"fields_schema,omitempty"`
	IsActive    bool    `gorm:"column:is_active;default:true" json:"is_active"`
	Order       int     `gorm:"default:0" json:"order"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName — nama tabel eksplisit
func (Template) TableName() string {
	return "templates"
}
