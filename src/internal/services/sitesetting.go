package services

import "wedding-invitation-go/internal/models"

// GetSiteSetting — ambil site setting
// Sementara return dummy. Nanti diganti query ke DB via GORM.
func GetSiteSetting() models.SiteSetting {
	return models.DummySiteSetting()
}

// GetTemplates — ambil semua template aktif (ordered)
// Sementara return empty. Nanti diganti query GORM ke tabel templates.
func GetTemplates() []models.Template {
	return []models.Template{}
}

// GetTemplateBySlug — ambil satu template berdasarkan slug
func GetTemplateBySlug(slug string) (models.Template, bool) {
	return models.Template{}, false
}
