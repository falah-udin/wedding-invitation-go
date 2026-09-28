package models

import (
	"time"

	"gorm.io/gorm"
)

// User — model user, map ke tabel `users`
// Struktur mengikuti migration Laravel
type User struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	Name          string         `gorm:"size:255;not null" json:"name"`
	Email         string         `gorm:"size:255;uniqueIndex;not null" json:"email"`
	EmailVerified *time.Time     `gorm:"column:email_verified_at" json:"email_verified_at,omitempty"`
	Password      string         `gorm:"size:255;not null" json:"-"`
	GoogleID      *string        `gorm:"column:google_id;size:255" json:"google_id,omitempty"`
	Avatar        *string        `gorm:"size:255" json:"avatar,omitempty"`
	Role          string         `gorm:"type:enum('admin','staff','client');default:'client'" json:"role"`
	IsActive      bool           `gorm:"column:is_active;default:true" json:"is_active"`
	RememberToken *string        `gorm:"column:remember_token;size:100" json:"-"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName — nama tabel eksplisit
func (User) TableName() string {
	return "users"
}

// IsAdmin — cek apakah user admin
func (u User) IsAdmin() bool {
	return u.Role == "admin"
}

// IsStaff — cek apakah user staff
func (u User) IsStaff() bool {
	return u.Role == "staff"
}

// IsClient — cek apakah user client
func (u User) IsClient() bool {
	return u.Role == "client"
}