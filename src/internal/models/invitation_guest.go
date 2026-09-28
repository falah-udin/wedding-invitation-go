package models

import (
	"time"

	"gorm.io/gorm"
)

// InvitationGuest — daftar tamu yang diundang
// Map ke tabel `invitation_guests`
type InvitationGuest struct {
	ID           uint       `gorm:"primaryKey" json:"id"`
	ProjectID    uint       `gorm:"not null;index" json:"project_id"`
	RsvpID       *uint      `gorm:"column:rsvp_id;index" json:"rsvp_id,omitempty"`
	Name         string     `gorm:"size:255;not null" json:"name"`
	Phone        *string    `gorm:"size:30" json:"phone,omitempty"`
	IsShared     bool       `gorm:"column:is_shared;default:false" json:"is_shared"`
	LastSharedAt *time.Time `gorm:"column:last_shared_at" json:"last_shared_at,omitempty"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relasi
	Project *Project `gorm:"foreignKey:ProjectID" json:"project,omitempty"`
	Rsvp    *Rsvp    `gorm:"foreignKey:RsvpID" json:"rsvp,omitempty"`
}

// TableName — nama tabel
func (InvitationGuest) TableName() string {
	return "invitation_guests"
}

// GetPhone — ambil nomor telepon (string kosong kalau nil)
func (g InvitationGuest) GetPhone() string {
	if g.Phone == nil {
		return ""
	}
	return *g.Phone
}

// IsHasPhone — cek apakah punya nomor telepon
func (g InvitationGuest) IsHasPhone() bool {
	return g.Phone != nil && *g.Phone != ""
}
