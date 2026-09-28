package models

// SiteSetting — model untuk pengaturan situs
// Nanti di-map dari tabel site_settings
type SiteSetting struct {
	// Brand
	SiteName    string
	SiteTagline string
	SiteFavicon string

	// Footer
	FooterDescription string
	FooterCopyright   string
	FooterVersion     string

	// Contact
	ContactEmail    string
	ContactPhone    string
	ContactWhatsapp string
	ContactAddress  string

	// Social
	SocialInstagram string
	SocialFacebook  string
	SocialTiktok    string
	SocialYoutube   string
	SocialTwitter   string

	// Developer
	DeveloperName  string
	DeveloperEmail string
	DeveloperUrl   string

	// Payment
	QrisImage string
}

// DummySiteSetting — data dummy untuk test UI
func DummySiteSetting() SiteSetting {
	return SiteSetting{
		SiteName:          "Wedding SaaS",
		SiteTagline:       "Platform Undangan Digital",
		SiteFavicon:       "",
		FooterDescription: "Platform undangan pernikahan digital dengan template elegan, RSVP online, guest book, dan galeri foto. Praktis, hemat, dan ramah lingkungan.",
		FooterCopyright:   "Wedding SaaS. Platform Undangan Digital. All rights reserved.",
		FooterVersion:     "v1.0.0",
		ContactEmail:      "hello@weddingsaas.com",
		ContactPhone:      "+62 812-3456-7890",
		ContactWhatsapp:   "6281234567890",
		ContactAddress:    "Jakarta, Indonesia",
		SocialInstagram:   "https://instagram.com/weddingsaas",
		SocialFacebook:    "https://facebook.com/weddingsaas",
		SocialTiktok:      "https://tiktok.com/@weddingsaas",
		SocialYoutube:     "",
		SocialTwitter:     "",
		DeveloperName:     "Tim Developer",
		DeveloperEmail:    "dev@weddingsaas.com",
		DeveloperUrl:      "https://weddingsaas.com",
		QrisImage:         "",
	}
}