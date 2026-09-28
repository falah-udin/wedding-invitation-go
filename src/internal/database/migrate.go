package database

import (
	"log"

	"wedding-invitation-go/internal/models"
)

// Migrate — jalankan auto migration untuk semua model
func Migrate() error {
	log.Println("🔄 Menjalankan migration...")

	err := DB.AutoMigrate(
		&models.User{},
		&models.Template{},
		&models.Project{},
		&models.Music{},
		// Nanti tambahkan model lain di sini:
		// &models.Rsvp{},
		// &models.InvitationGuest{},
		// &models.SiteSetting{},
	)

	if err != nil {
		return err
	}

	log.Println("✅ Migration selesai")
	return nil
}
