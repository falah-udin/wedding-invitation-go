package services

import (
	"strings"

	"wedding-invitation-go/internal/database"
	"wedding-invitation-go/internal/models"
)

// GetRsvpByGuest — ambil RSVP berdasarkan project + guest name (case-insensitive)
func GetRsvpByGuest(projectID uint, guestName string) *models.Rsvp {
	if guestName == "" {
		return nil
	}

	var rsvp models.Rsvp
	db := database.GetDB()
	err := db.Where("project_id = ? AND LOWER(guest_name) = ?",
		projectID, strings.ToLower(guestName)).
		First(&rsvp).Error

	if err != nil {
		return nil
	}
	return &rsvp
}