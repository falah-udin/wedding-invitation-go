package admin

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"wedding-invitation-go/internal/database"
	"wedding-invitation-go/internal/models"
	"wedding-invitation-go/internal/services"
	"wedding-invitation-go/views/admin/settings"
	"wedding-invitation-go/views/layouts"
)

const brandingUploadDir = "/app/public/storage/branding"
const qrisUploadDir = "/app/public/storage/qris"

// ============================================
// SettingIndex — GET /admin/settings
// ============================================
func SettingIndex(c *gin.Context) {
	setting := services.GetSiteSetting()
	user := c.MustGet("user").(*models.User)
	errorMsg := c.Query("error")
	successMsg := c.Query("success")

	settingsMap := getSettingsMap()

	var buf bytes.Buffer
	err := layouts.AdminLayout(
		*user, setting, "admin.settings",
		settings.IndexContent(settingsMap, errorMsg, successMsg),
	).Render(c.Request.Context(), &buf)

	if err != nil {
		c.String(http.StatusInternalServerError, "Render error: %v", err)
		return
	}

	c.Data(http.StatusOK, "text/html; charset=utf-8", buf.Bytes())
}

// ============================================
// SettingUpdate — POST /admin/settings
// ============================================
func SettingUpdate(c *gin.Context) {
	textFields := []string{
		"site_name", "site_tagline",
		"contact_email", "contact_phone", "contact_whatsapp", "contact_address",
		"social_instagram", "social_facebook", "social_tiktok", "social_youtube", "social_twitter",
		"footer_description", "footer_copyright", "footer_version",
		"developer_name", "developer_email", "developer_url",
		// PWA
		"pwa_short_name", "pwa_theme_color", "pwa_background_color",
	}

	for _, key := range textFields {
		value := strings.TrimSpace(c.PostForm(key))
		if err := upsertSetting(key, value); err != nil {
			c.Redirect(http.StatusFound, "/admin/settings?error=Gagal+menyimpan+"+key)
			return
		}
	}

	// Handle favicon upload
	if file, err := c.FormFile("site_favicon"); err == nil && file.Size > 0 {
		oldFavicon := getSettingValue("site_favicon")
		if oldFavicon != "" {
			os.Remove(filepath.Join("/app/public/storage", oldFavicon))
		}

		path, uploadErr := saveBrandingFile(c, file)
		if uploadErr != nil {
			c.Redirect(http.StatusFound, "/admin/settings?error=Gagal+upload+favicon")
			return
		}
		upsertSetting("site_favicon", path)
	}

	// Handle QRIS upload
	if file, err := c.FormFile("qris_image"); err == nil && file.Size > 0 {
		oldQris := getSettingValue("qris_image")
		if oldQris != "" {
			os.Remove(filepath.Join("/app/public/storage", oldQris))
		}

		path, uploadErr := saveQrisFile(c, file)
		if uploadErr != nil {
			c.Redirect(http.StatusFound, "/admin/settings?error=Gagal+upload+QRIS")
			return
		}
		upsertSetting("qris_image", path)
	}

		// Handle PWA icon 192
	if file, err := c.FormFile("pwa_icon_192"); err == nil && file.Size > 0 {
		oldIcon := getSettingValue("pwa_icon_192")
		if oldIcon != "" {
			os.Remove(filepath.Join("/app/public/storage", oldIcon))
		}

		path, uploadErr := savePwaIconFile(c, file, "192")
		if uploadErr != nil {
			c.Redirect(http.StatusFound, "/admin/settings?error=Gagal+upload+icon+192")
			return
		}
		upsertSetting("pwa_icon_192", path)
	}

	// Handle PWA icon 512
	if file, err := c.FormFile("pwa_icon_512"); err == nil && file.Size > 0 {
		oldIcon := getSettingValue("pwa_icon_512")
		if oldIcon != "" {
			os.Remove(filepath.Join("/app/public/storage", oldIcon))
		}

		path, uploadErr := savePwaIconFile(c, file, "512")
		if uploadErr != nil {
			c.Redirect(http.StatusFound, "/admin/settings?error=Gagal+upload+icon+512")
			return
		}
		upsertSetting("pwa_icon_512", path)
	}

	c.Redirect(http.StatusFound, "/admin/settings?success=Pengaturan+berhasil+disimpan")
}

// ============================================
// HELPER
// ============================================
func getSettingsMap() map[string]string {
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

func getSettingValue(key string) string {
	var s models.SiteSetting
	db := database.GetDB()
	if err := db.Where("setting_key = ?", key).First(&s).Error; err != nil {
		return ""
	}
	if s.Value != nil {
		return *s.Value
	}
	return ""
}

func upsertSetting(key, value string) error {
	db := database.GetDB()
	var s models.SiteSetting

	result := db.Where("setting_key = ?", key).First(&s)
	if result.Error != nil {
		v := value
		s = models.SiteSetting{
			SettingKey:   key,
			Value:        &v,
			SettingGroup: "general",
		}
	}

	v := value
	s.Value = &v
	return db.Save(&s).Error
}

func saveBrandingFile(c *gin.Context, file *multipart.FileHeader) (string, error) {
	if err := os.MkdirAll(brandingUploadDir, 0755); err != nil {
		return "", err
	}

	// Validasi ukuran (max 512KB)
	if file.Size > 512*1024 {
		return "", fmt.Errorf("file terlalu besar (max 512KB)")
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	filename := fmt.Sprintf("favicon_%d%s", time.Now().Unix(), ext)
	fullPath := filepath.Join(brandingUploadDir, filename)

	if err := c.SaveUploadedFile(file, fullPath); err != nil {
		return "", err
	}

	return "branding/" + filename, nil
}

// savePwaIconFile — simpan icon PWA ke /app/public/storage/settings/pwa/
// size: "192" atau "512" — buat nama file unik
func savePwaIconFile(c *gin.Context, file *multipart.FileHeader, size string) (string, error) {
	// Validasi ekstensi
	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowed := map[string]bool{".png": true, ".jpg": true, ".jpeg": true, ".webp": true}
	if !allowed[ext] {
		return "", fmt.Errorf("format tidak didukung: %s", ext)
	}

	// Validasi ukuran file (max 1 MB untuk icon)
	if file.Size > 1024*1024 {
		return "", fmt.Errorf("file terlalu besar (max 1 MB)")
	}

	// Buat nama unik
	filename := fmt.Sprintf("pwa-%s-%d%s", size, time.Now().Unix(), ext)

	// Path di container
	dstDir := "/app/public/storage/settings/pwa"
	if err := os.MkdirAll(dstDir, 0755); err != nil {
		return "", err
	}
	dstPath := filepath.Join(dstDir, filename)

	// Save
	if err := c.SaveUploadedFile(file, dstPath); err != nil {
		return "", err
	}

	// Return path relatif (untuk DB)
	return "settings/pwa/" + filename, nil
}

func saveQrisFile(c *gin.Context, file *multipart.FileHeader) (string, error) {
	if err := os.MkdirAll(qrisUploadDir, 0755); err != nil {
		return "", err
	}

	// Validasi ukuran (max 2MB)
	if file.Size > 2*1024*1024 {
		return "", fmt.Errorf("file terlalu besar (max 2MB)")
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	filename := fmt.Sprintf("qris_%d%s", time.Now().Unix(), ext)
	fullPath := filepath.Join(qrisUploadDir, filename)

	if err := c.SaveUploadedFile(file, fullPath); err != nil {
		return "", err
	}

	return "qris/" + filename, nil
}

// ============================================
// SettingDeleteFavicon — POST /admin/settings/delete-favicon
// ============================================
func SettingDeleteFavicon(c *gin.Context) {
	oldFavicon := getSettingValue("site_favicon")
	if oldFavicon != "" {
		os.Remove(filepath.Join("/app/public/storage", oldFavicon))
	}
	upsertSetting("site_favicon", "")
	c.Redirect(http.StatusFound, "/admin/settings?success=Favicon+berhasil+dihapus")
}

// ============================================
// SettingDeleteQris — POST /admin/settings/delete-qris
// ============================================
func SettingDeleteQris(c *gin.Context) {
	oldQris := getSettingValue("qris_image")
	if oldQris != "" {
		os.Remove(filepath.Join("/app/public/storage", oldQris))
	}
	upsertSetting("qris_image", "")
	c.Redirect(http.StatusFound, "/admin/settings?success=QRIS+berhasil+dihapus")
}

// SettingDeletePwaIcon — POST /admin/settings/delete-pwa-icon
func SettingDeletePwaIcon(c *gin.Context) {
	size := c.PostForm("size")
	key := "pwa_icon_" + size
	oldIcon := getSettingValue(key)
	if oldIcon != "" {
		os.Remove(filepath.Join("/app/public/storage", oldIcon))
	}
	upsertSetting(key, "")
	c.Redirect(http.StatusFound, "/admin/settings?success=Icon+PWA+berhasil+dihapus")
}