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
	"wedding-invitation-go/views/admin/music"
	"wedding-invitation-go/views/layouts"
)

const musicUploadDir = "/app/public/storage/music/library"

// ============================================
// MusicIndex — GET /admin/music
// ============================================
func MusicIndex(c *gin.Context) {
	setting := services.GetSiteSetting()
	user := c.MustGet("user").(*models.User)

	var list []models.Music
	db := database.GetDB()
	db.Preload("User").Order("id DESC").Find(&list)

	var buf bytes.Buffer
	err := layouts.AdminLayout(
		*user, setting, "admin.music",
		music.IndexContent(list),
	).Render(c.Request.Context(), &buf)

	if err != nil {
		c.String(http.StatusInternalServerError, "Render error: %v", err)
		return
	}

	c.Data(http.StatusOK, "text/html; charset=utf-8", buf.Bytes())
}

// ============================================
// MusicCreate — GET /admin/music/create
// ============================================
func MusicCreate(c *gin.Context) {
	setting := services.GetSiteSetting()
	user := c.MustGet("user").(*models.User)
	errorMsg := c.Query("error")

	var buf bytes.Buffer
	err := layouts.AdminLayout(
		*user, setting, "admin.music",
		music.CreateContent(errorMsg),
	).Render(c.Request.Context(), &buf)

	if err != nil {
		c.String(http.StatusInternalServerError, "Render error: %v", err)
		return
	}

	c.Data(http.StatusOK, "text/html; charset=utf-8", buf.Bytes())
}

// ============================================
// MusicStore — POST /admin/music/create
// ============================================
func MusicStore(c *gin.Context) {
	title := strings.TrimSpace(c.PostForm("title"))
	artist := strings.TrimSpace(c.PostForm("artist"))
	isShared := c.PostForm("is_shared") == "1"

	if title == "" {
		c.Redirect(http.StatusFound, "/admin/music/create?error=Judul+wajib+diisi")
		return
	}

	// Cek duplikasi
	db := database.GetDB()
	var existing models.Music
	if err := db.Where("title = ? AND artist = ?", title, artist).First(&existing).Error; err == nil {
		c.Redirect(http.StatusFound, "/admin/music/create?error=Musik+dengan+judul+dan+artis+yang+sama+sudah+ada")
		return
	}

	// Upload file
	file, err := c.FormFile("music_file")
	if err != nil {
		c.Redirect(http.StatusFound, "/admin/music/create?error=File+musik+wajib+dipilih")
		return
	}

	// Validasi ukuran (max 10MB)
	if file.Size > 10*1024*1024 {
		c.Redirect(http.StatusFound, "/admin/music/create?error=Ukuran+file+melebihi+10MB")
		return
	}

	// Validasi ekstensi
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".mp3" && ext != ".wav" && ext != ".ogg" {
		c.Redirect(http.StatusFound, "/admin/music/create?error=Format+file+harus+MP3,+WAV,+atau+OGG")
		return
	}

	filePath, uploadErr := saveMusicFile(c, file)
	if uploadErr != nil {
		c.Redirect(http.StatusFound, "/admin/music/create?error=Gagal+mengupload+file")
		return
	}

	// Simpan ke DB
	userVal := c.MustGet("user").(*models.User)
	userID := userVal.ID

	m := models.Music{
		UserID:   &userID,
		Title:    title,
		FilePath: filePath,
		IsShared: isShared,
	}
	if artist != "" {
		m.Artist = &artist
	}

	if err := db.Create(&m).Error; err != nil {
		// Hapus file kalau gagal DB
		os.Remove(filepath.Join("/app/public/storage", strings.TrimPrefix(filePath, "/storage/")))
		c.Redirect(http.StatusFound, "/admin/music/create?error=Gagal+menyimpan+ke+database")
		return
	}

	c.Redirect(http.StatusFound, "/admin/music")
}

// ============================================
// MusicToggle — GET /admin/music/:id/toggle
// ============================================
func MusicToggle(c *gin.Context) {
	id := c.Param("id")

	db := database.GetDB()
	var m models.Music
	if err := db.First(&m, id).Error; err != nil {
		c.Redirect(http.StatusFound, "/admin/music?error=Musik+tidak+ditemukan")
		return
	}

	m.IsShared = !m.IsShared
	db.Save(&m)

	c.Redirect(http.StatusFound, "/admin/music")
}

// ============================================
// MusicDelete — POST /admin/music/:id/delete
// ============================================
func MusicDelete(c *gin.Context) {
	id := c.Param("id")

	db := database.GetDB()
	var m models.Music
	if err := db.First(&m, id).Error; err != nil {
		c.Redirect(http.StatusFound, "/admin/music?error=Musik+tidak+ditemukan")
		return
	}

	// Hapus file fisik
	if m.FilePath != "" {
		oldPath := strings.TrimPrefix(m.FilePath, "/storage/")
		os.Remove(filepath.Join("/app/public/storage", oldPath))
	}
	
	db.Delete(&m)
	c.Redirect(http.StatusFound, "/admin/music")
}

// ============================================
// HELPER — saveMusicFile
// ============================================
func saveMusicFile(c *gin.Context, file *multipart.FileHeader) (string, error) {
	if err := os.MkdirAll(musicUploadDir, 0755); err != nil {
		return "", err
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	baseName := strings.TrimSuffix(filepath.Base(file.Filename), filepath.Ext(file.Filename))
	baseName = strings.ReplaceAll(baseName, " ", "_")
	filename := fmt.Sprintf("%d_%s%s", time.Now().Unix(), baseName, ext)
	fullPath := filepath.Join(musicUploadDir, filename)

	if err := c.SaveUploadedFile(file, fullPath); err != nil {
		return "", err
	}

	return "music/library/" + filename, nil
}
