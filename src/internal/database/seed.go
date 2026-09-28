package database

import (
	"log"

	"golang.org/x/crypto/bcrypt"
	"wedding-invitation-go/internal/models"
)

func Seed() error {
	log.Println("🌱 Menjalankan seeder...")

	if err := seedUsers(); err != nil {
		return err
	}

	if err := seedSiteSettings(); err != nil {
		return err
	}

	log.Println("✅ Seeder selesai")
	return nil
}

func seedUsers() error {
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

	return nil
}

func seedSiteSettings() error {
	defaults := []struct {
		Key   string
		Value string
		Group string
	}{
		{"site_name", "Wedding SaaS", "brand"},
		{"site_tagline", "Platform Undangan Digital", "brand"},
		{"site_favicon", "", "brand"},
		{"footer_description", "Platform undangan pernikahan digital dengan template elegan, RSVP online, guest book, dan galeri foto. Praktis, hemat, dan ramah lingkungan.", "footer"},
		{"footer_copyright", "Wedding SaaS. Platform Undangan Digital. All rights reserved.", "footer"},
		{"footer_version", "v1.0.0", "footer"},
		{"contact_email", "hello@weddingsaas.com", "contact"},
		{"contact_phone", "+62 812-3456-7890", "contact"},
		{"contact_whatsapp", "6281234567890", "contact"},
		{"contact_address", "Jakarta, Indonesia", "contact"},
		{"social_instagram", "https://instagram.com/weddingsaas", "social"},
		{"social_facebook", "https://facebook.com/weddingsaas", "social"},
		{"social_tiktok", "https://tiktok.com/@weddingsaas", "social"},
		{"social_youtube", "", "social"},
		{"social_twitter", "", "social"},
		{"developer_name", "Tim Developer", "developer"},
		{"developer_email", "dev@weddingsaas.com", "developer"},
		{"developer_url", "https://weddingsaas.com", "developer"},
		{"qris_image", "", "payment"},
	}

	for _, item := range defaults {
		var existing models.SiteSetting
		result := DB.Where("setting_key = ?", item.Key).First(&existing)
		if result.Error == nil {
			continue
		}

		value := item.Value
		setting := models.SiteSetting{
			SettingKey:   item.Key,
			Value:        &value,
			SettingGroup: item.Group,
		}
		if err := DB.Create(&setting).Error; err != nil {
			return err
		}
	}

	log.Printf("✅ Site settings di-seed (%d items)", len(defaults))
	return nil
}
