package services

import "wedding-invitation-go/internal/models"

// GetSiteSetting — ambil site setting
// Sementara return dummy. Nanti diganti query ke DB via GORM.
func GetSiteSetting() models.SiteSetting {
	return models.DummySiteSetting()
}

// GetTemplates — ambil semua template aktif (ordered)
// Sementara return dummy. Nanti diganti query GORM ke tabel templates.
func GetTemplates() []models.Template {
	return models.DummyTemplates()
}

// GetTemplateBySlug — ambil satu template berdasarkan slug
// Berguna untuk halaman preview template
func GetTemplateBySlug(slug string) (models.Template, bool) {
	for _, t := range DummyTemplatesInternal() {
		if t.Slug == slug {
			return t, true
		}
	}
	return models.Template{}, false
}

// DummyTemplatesInternal — helper (jangan dipakai dari luar service)
func DummyTemplatesInternal() []models.Template {
	return models.DummyTemplates()
}