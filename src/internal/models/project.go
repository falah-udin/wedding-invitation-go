package models

import (
	"time"

	"gorm.io/gorm"
)

// Project — model project undangan, map ke tabel `projects`
// Struktur mengikuti migration Laravel
type Project struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	UserID    uint           `gorm:"not null;index" json:"user_id"`
	TemplateID uint          `gorm:"not null;index" json:"template_id"`
	Title     string         `gorm:"size:255;not null" json:"title"`
	Slug      string         `gorm:"size:255;uniqueIndex;not null" json:"slug"`

	// JSON columns
	DataUndangan         string `gorm:"column:data_undangan;type:json" json:"data_undangan"`
	TemplateSpecificData string `gorm:"column:template_specific_data;type:json" json:"template_specific_data,omitempty"`

	// Status
	Status string `gorm:"type:enum('draft','published','archived');default:'draft'" json:"status"`

	// Music & counters
	CustomMusic *string `gorm:"column:custom_music;size:255" json:"custom_music,omitempty"`
	TotalViews  int     `gorm:"column:total_views;default:0" json:"total_views"`
	TotalRsvp   int     `gorm:"column:total_rsvp;default:0" json:"total_rsvp"`

	// Timestamps
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relasi (loaded via Preload)
	User     *User     `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Template *Template `gorm:"foreignKey:TemplateID" json:"template,omitempty"`
}

// TableName — nama tabel eksplisit
func (Project) TableName() string {
	return "projects"
}
