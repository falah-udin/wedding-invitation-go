package models

import (
	"time"

	"gorm.io/gorm"
)

// Rsvp — konfirmasi kehadiran tamu
// Map ke tabel `rsvps`
type Rsvp struct {
	ID          uint    `gorm:"primaryKey" json:"id"`
	ProjectID   uint    `gorm:"not null;index" json:"project_id"`
	GuestName   string  `gorm:"column:guest_name;size:255;not null" json:"guest_name"`
	Attendance  string  `gorm:"type:enum('hadir','tidak_hadir','ragu');default:'ragu'" json:"attendance"`
	TotalGuests int     `gorm:"column:total_guests;default:1" json:"total_guests"`
	Message     *string `gorm:"type:text" json:"message,omitempty"`
	EditToken   *string `gorm:"column:edit_token;size:64" json:"edit_token,omitempty"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relasi
	Project *Project `gorm:"foreignKey:ProjectID" json:"project,omitempty"`
}

// TableName — nama tabel
func (Rsvp) TableName() string {
	return "rsvps"
}

// GetMessage — ambil pesan (string kosong kalau nil)
func (r Rsvp) GetMessage() string {
	if r.Message == nil {
		return ""
	}
	return *r.Message
}

// GetAttendanceLabel — label untuk tampilan
func (r Rsvp) GetAttendanceLabel() string {
	switch r.Attendance {
	case "hadir":
		return "Hadir"
	case "tidak_hadir":
		return "Tidak Hadir"
	case "ragu":
		return "Ragu-ragu"
	}
	return r.Attendance
}
