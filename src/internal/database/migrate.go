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
		// Nanti tambahkan model lain di sini:
		// &models.Template{},
		// &models.Project{},
		// &models.Rsvp{},
		// dst
	)

	if err != nil {
		return err
	}

	log.Println("✅ Migration selesai")
	return nil
}