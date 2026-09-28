package invitation

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"wedding-invitation-go/internal/database"
	"wedding-invitation-go/internal/models"
	"wedding-invitation-go/internal/services"
	"wedding-invitation-go/views/invitation"
	"wedding-invitation-go/views/layouts"
)

// ============================================
// Preview — GET /invitation/create/preview
// Step 3: Preview + Music Selection
// ============================================
func Preview(c *gin.Context) {
	setting := services.GetSiteSetting()
	user := c.MustGet("user").(*models.User)

	projectID := services.WizardGetProjectID(c.Request)
	if projectID == 0 {
		c.Redirect(http.StatusFound, "/invitation/create/select-client?error=Sesi+habis")
		return
	}

	// Ambil project
	var project models.Project
	db := database.GetDB()
	if err := db.Preload("Template").Preload("User").First(&project, projectID).Error; err != nil {
		c.Redirect(http.StatusFound, "/invitation/create/select-client?error=Project+tidak+ditemukan")
		return
	}

	// Ambil data undangan + specific data
	generalData := services.ParseJSONMap(project.DataUndangan)
	specificData := services.ParseJSONMap(project.TemplateSpecificData)

	// Merge untuk preview
	allData := make(map[string]interface{})
	for k, v := range generalData {
		allData[k] = v
	}
	for k, v := range specificData {
		allData[k] = v
	}

	// Ambil musik yang tersedia (publik + milik sendiri)
	var musics []models.Music
	if user.Role == "admin" {
		// Admin lihat semua musik
		db.Preload("User").Order("is_shared DESC, title ASC").Find(&musics)
	} else {
		// Staff/client lihat musik publik + milik sendiri
		db.Preload("User").
			Where("is_shared = ? OR user_id = ?", true, user.ID).
			Order("is_shared DESC, title ASC").
			Find(&musics)
	}

	errorMsg := c.Query("error")
	successMsg := c.Query("success")

	var buf bytes.Buffer
	err := layouts.AdminLayout(
		*user, setting, "invitation.create",
		invitation.PreviewContent(
			project,
			allData,
			musics,
			errorMsg,
			successMsg,
		),
	).Render(c.Request.Context(), &buf)

	if err != nil {
		c.String(http.StatusInternalServerError, "Render error: %v", err)
		return
	}

	c.Data(http.StatusOK, "text/html; charset=utf-8", buf.Bytes())
}

// ============================================
// Publish — POST /invitation/create/publish
// Publish undangan (set status = published)
// ============================================
func Publish(c *gin.Context) {
	projectID := services.WizardGetProjectID(c.Request)
	if projectID == 0 {
		c.Redirect(http.StatusFound, "/invitation/create/select-client?error=Sesi+habis")
		return
	}

	db := database.GetDB()
	var project models.Project
	if err := db.First(&project, projectID).Error; err != nil {
		c.Redirect(http.StatusFound, "/invitation/create/select-client?error=Project+tidak+ditemukan")
		return
	}

	// Set status published
	project.Status = "published"
	if err := db.Save(&project).Error; err != nil {
		services.LogError("Publish.UpdateStatus", err, map[string]interface{}{
			"project_id": projectID,
		})
		c.Redirect(http.StatusFound, "/invitation/create/preview?error=Gagal+publish")
		return
	}

	services.LogSuccess("Publish", fmt.Sprintf("Project %s published", project.Slug))

	// Clear session wizard
	services.WizardClear(c.Writer, c.Request)

	// Redirect berdasarkan role
	user := c.MustGet("user").(*models.User)
	if user.Role == "admin" {
		c.Redirect(http.StatusFound, "/admin/projects?success=Undangan+berhasil+dipublikasikan")
	} else if user.Role == "staff" {
		c.Redirect(http.StatusFound, "/staff/projects?success=Undangan+berhasil+dipublikasikan")
	} else {
		c.Redirect(http.StatusFound, "/client/wedding?success=Undangan+berhasil+dipublikasikan")
	}
}

// ============================================
// SaveMusic — POST /invitation/create/save-music
// Simpan musik pilihan ke project (AJAX)
// ============================================
func SaveMusic(c *gin.Context) {
	projectID := services.WizardGetProjectID(c.Request)
	if projectID == 0 {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "Sesi habis"})
		return
	}

	musicPath := strings.TrimSpace(c.PostForm("custom_music"))
	if musicPath == "" {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "Musik wajib dipilih"})
		return
	}

	db := database.GetDB()
	var project models.Project
	if err := db.First(&project, projectID).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "Project tidak ditemukan"})
		return
	}

	project.CustomMusic = &musicPath
	if err := db.Save(&project).Error; err != nil {
		services.LogError("SaveMusic", err, map[string]interface{}{
			"project_id":   projectID,
			"music_path":   musicPath,
		})
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "Gagal menyimpan musik"})
		return
	}

	services.LogSuccess("SaveMusic", fmt.Sprintf("Music saved to project %d", projectID))

	c.JSON(http.StatusOK, gin.H{
		"success":  true,
		"message":  "Musik berhasil disimpan",
		"music":    musicPath,
	})
}
