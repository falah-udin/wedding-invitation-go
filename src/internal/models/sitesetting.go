package models

import (
	"time"
)

// SiteSetting — model untuk tabel site_settings
// Kolom pakai prefix setting_ untuk hindari reserved word MySQL
type SiteSetting struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	SettingKey   string    `gorm:"column:setting_key;size:255;uniqueIndex;not null" json:"setting_key"`
	Value        *string   `gorm:"type:text" json:"value,omitempty"`
	SettingGroup string    `gorm:"column:setting_group;size:50;default:'general'" json:"setting_group"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (SiteSetting) TableName() string {
	return "site_settings"
}

type Config struct {
	SiteName          string
	SiteTagline       string
	SiteFavicon       string
	FooterDescription string
	FooterCopyright   string
	FooterVersion     string
	ContactEmail      string
	ContactPhone      string
	ContactWhatsapp   string
	ContactAddress    string
	SocialInstagram   string
	SocialFacebook    string
	SocialTiktok      string
	SocialYoutube     string
	SocialTwitter     string
	DeveloperName     string
	DeveloperEmail    string
	DeveloperUrl      string
	QrisImage         string
}

func (c Config) WithDefaults() Config {
	if c.SiteName == "" {
		c.SiteName = "Wedding SaaS"
	}
	if c.SiteTagline == "" {
		c.SiteTagline = "Platform Undangan Digital"
	}
	if c.FooterDescription == "" {
		c.FooterDescription = "Platform undangan pernikahan digital dengan template elegan, RSVP online, guest book, dan galeri foto."
	}
	if c.FooterCopyright == "" {
		c.FooterCopyright = "Wedding SaaS. Platform Undangan Digital. All rights reserved."
	}
	if c.DeveloperName == "" {
		c.DeveloperName = "Tim Developer"
	}
	return c
}
