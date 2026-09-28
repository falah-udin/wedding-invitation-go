package models

import (
	"time"

	"gorm.io/gorm"
)

// Music — model music library, map ke tabel `music_library`
// Struktur mengikuti migration Laravel
type Music struct {
	ID       uint    `gorm:"primaryKey" json:"id"`
	UserID   *uint   `gorm:"column:user_id;index" json:"user_id,omitempty"`
	Title    string  `gorm:"size:255;not null" json:"title"`
	Artist   *string `gorm:"size:255" json:"artist,omitempty"`
	FilePath string  `gorm:"column:file_path;size:255;not null" json:"file_path"`
	Duration *string `gorm:"size:50" json:"duration,omitempty"`
	IsShared bool    `gorm:"column:is_shared;default:true" json:"is_shared"`
	PlayCount int    `gorm:"column:play_count;default:0" json:"play_count"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relasi
	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// TableName — nama tabel eksplisit
func (Music) TableName() string {
	return "music_library"
}
