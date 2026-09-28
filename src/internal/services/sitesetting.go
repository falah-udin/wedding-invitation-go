package services

import (
	"wedding-invitation-go/internal/database"
	"wedding-invitation-go/internal/models"
)

func GetSiteSetting() models.Config {
	var settings []models.SiteSetting
	db := database.GetDB()
	db.Find(&settings)

	m := make(map[string]string)
	for _, s := range settings {
		if s.Value != nil {
			m[s.SettingKey] = *s.Value
		}
	}

	cfg := models.Config{
		SiteName:          m["site_name"],
		SiteTagline:       m["site_tagline"],
		SiteFavicon:       m["site_favicon"],
		FooterDescription: m["footer_description"],
		FooterCopyright:   m["footer_copyright"],
		FooterVersion:     m["footer_version"],
		ContactEmail:      m["contact_email"],
		ContactPhone:      m["contact_phone"],
		ContactWhatsapp:   m["contact_whatsapp"],
		ContactAddress:    m["contact_address"],
		SocialInstagram:   m["social_instagram"],
		SocialFacebook:    m["social_facebook"],
		SocialTiktok:      m["social_tiktok"],
		SocialYoutube:     m["social_youtube"],
		SocialTwitter:     m["social_twitter"],
		DeveloperName:     m["developer_name"],
		DeveloperEmail:    m["developer_email"],
		DeveloperUrl:      m["developer_url"],
		QrisImage:         m["qris_image"],
	}

	return cfg.WithDefaults()
}

func GetTemplates() []models.Template {
	var list []models.Template
	db := database.GetDB()
	db.Where("is_active = ?", true).Order("`order` ASC, id ASC").Find(&list)
	return list
}

func GetTemplateBySlug(slug string) (models.Template, bool) {
	var tmpl models.Template
	db := database.GetDB()
	err := db.Where("slug = ?", slug).First(&tmpl).Error
	return tmpl, err == nil
}
