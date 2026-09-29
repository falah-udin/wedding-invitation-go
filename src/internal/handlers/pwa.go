package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"

	"wedding-invitation-go/internal/database"
	"wedding-invitation-go/internal/models"
)

type ManifestIcon struct {
	Src     string `json:"src"`
	Sizes   string `json:"sizes"`
	Type    string `json:"type"`
	Purpose string `json:"purpose,omitempty"`
}

type ManifestJSON struct {
	Name            string         `json:"name"`
	ShortName       string         `json:"short_name"`
	Description     string         `json:"description,omitempty"`
	StartURL        string         `json:"start_url"`
	Display         string         `json:"display"`
	BackgroundColor string         `json:"background_color"`
	ThemeColor      string         `json:"theme_color"`
	Orientation     string         `json:"orientation,omitempty"`
	Icons           []ManifestIcon `json:"icons"`
}

// Manifest — GET /manifest.webmanifest
func Manifest(c *gin.Context) {
	settings := loadSettingsMap()

	name := settings["site_name"]
	if name == "" {
		name = "Wedding Invitation"
	}
	shortName := settings["pwa_short_name"]
	if shortName == "" {
		shortName = "Wedding"
	}
	themeColor := settings["pwa_theme_color"]
	if themeColor == "" {
		themeColor = "#8b5cf6"
	}
	bgColor := settings["pwa_background_color"]
	if bgColor == "" {
		bgColor = "#0f0f0f"
	}

	icons := []ManifestIcon{}
	icon192 := settings["pwa_icon_192"]
	icon512 := settings["pwa_icon_512"]
	favicon := settings["site_favicon"]

	if icon192 != "" {
		icons = append(icons, ManifestIcon{Src: "/storage/" + icon192, Sizes: "192x192", Type: "image/png", Purpose: "any maskable"})
	} else if favicon != "" {
		icons = append(icons, ManifestIcon{Src: "/storage/" + favicon, Sizes: "192x192", Type: "image/png", Purpose: "any"})
	}

	if icon512 != "" {
		icons = append(icons, ManifestIcon{Src: "/storage/" + icon512, Sizes: "512x512", Type: "image/png", Purpose: "any maskable"})
	} else if favicon != "" {
		icons = append(icons, ManifestIcon{Src: "/storage/" + favicon, Sizes: "512x512", Type: "image/png", Purpose: "any"})
	}

	manifest := ManifestJSON{
		Name:            name,
		ShortName:       shortName,
		Description:     settings["site_tagline"],
		StartURL:        "/",
		Display:         "standalone",
		BackgroundColor: bgColor,
		ThemeColor:      themeColor,
		Orientation:     "portrait",
		Icons:           icons,
	}

	c.Header("Content-Type", "application/manifest+json; charset=utf-8")
	c.Header("Cache-Control", "public, max-age=300")
	_ = json.NewEncoder(c.Writer).Encode(manifest)
}

// ServiceWorker — GET /sw.js
func ServiceWorker(c *gin.Context) {
	js := `self.addEventListener('install', (e) => { self.skipWaiting(); });
self.addEventListener('activate', (e) => { e.waitUntil(self.clients.claim()); });
self.addEventListener('fetch', (e) => { /* pass-through */ });
`
	c.Header("Content-Type", "application/javascript; charset=utf-8")
	c.Header("Cache-Control", "public, max-age=3600")
	c.String(http.StatusOK, js)
}

// loadSettingsMap — ambil semua SiteSetting dari DB
func loadSettingsMap() map[string]string {
	var list []models.SiteSetting
	db := database.GetDB()
	db.Find(&list)

	result := make(map[string]string)
	for _, s := range list {
		if s.Value != nil {
			result[s.SettingKey] = *s.Value
		} else {
			result[s.SettingKey] = ""
		}
	}
	return result
}
