package database

import (
	"log"

	"golang.org/x/crypto/bcrypt"
	"wedding-invitation-go/internal/models"
)

// Seed — isi data awal (admin, staff, client)
func Seed() error {
	log.Println("🌱 Menjalankan seeder...")

	users := []struct {
		Name     string
		Email    string
		Password string
		Role     string
	}{
		{"Admin Utama", "admin@gmail.com", "12344321", "admin"},
		{"Staff Operasional", "staff@gmail.com", "12345678", "staff"},
		{"Client Demo", "client@gmail.com", "12345678", "client"},
	}

	for _, u := range users {
		var existing models.User
		result := DB.Where("email = ?", u.Email).First(&existing)

		if result.Error == nil {
			log.Printf("⏭️  User %s sudah ada, skip", u.Email)
			continue
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}

		user := models.User{
			Name:     u.Name,
			Email:    u.Email,
			Password: string(hashedPassword),
			Role:     u.Role,
			IsActive: true,
		}

		if err := DB.Create(&user).Error; err != nil {
			return err
		}

		log.Printf("✅ User dibuat: %s (%s)", u.Email, u.Role)
	}

	log.Println("✅ Seeder selesai")
	return nil
}