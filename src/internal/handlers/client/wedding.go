package client

import (
	"bytes"
	"net/http"

	"github.com/gin-gonic/gin"

	"wedding-invitation-go/internal/database"
	"wedding-invitation-go/internal/models"
	"wedding-invitation-go/internal/services"
	"wedding-invitation-go/views/client/wedding"
	"wedding-invitation-go/views/layouts"
)

// WeddingIndex — GET /client/wedding
// List semua undangan milik client
func WeddingIndex(c *gin.Context) {
	setting := services.GetSiteSetting()
	user := services.GetCurrentUser(c.Request)
	if user == nil {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	db := database.GetDB()
	var projects []models.Project
	db.Preload("Template").
		Where("user_id = ?", user.ID).
		Order("created_at DESC").
		Find(&projects)

	data := wedding.WeddingData{
		User:     *user,
		Setting:  setting,
		Projects: projects,
		Success:  c.Query("success"),
		Error:    c.Query("error"),
	}

	var buf bytes.Buffer
	err := layouts.ClientLayout(*user, setting, "client.wedding", wedding.WeddingPage(data)).Render(c.Request.Context(), &buf)
	if err != nil {
		c.String(http.StatusInternalServerError, "Render error: %v", err)
		return
	}
	c.Data(http.StatusOK, "text/html; charset=utf-8", buf.Bytes())
}

// WeddingDestroy — DELETE /client/wedding/:id
// Hapus undangan (hard delete) — return JSON
func WeddingDestroy(c *gin.Context) {
	user := services.GetCurrentUser(c.Request)
	if user == nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "Sesi habis, silakan login ulang"})
		return
	}

	idStr := c.Param("id")
	var id uint
	for _, ch := range idStr {
		if ch < '0' || ch > '9' {
			break
		}
		id = id*10 + uint(ch-'0')
	}

	if id == 0 {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "ID tidak valid"})
		return
	}

	db := database.GetDB()
	var project models.Project
	if err := db.Where("id = ? AND user_id = ?", id, user.ID).First(&project).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "Undangan tidak ditemukan"})
		return
	}

	// Hanya draft yang bisa dihapus client
	if project.Status != "draft" {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "Hanya undangan draft yang bisa dihapus"})
		return
	}

	if err := db.Delete(&project).Error; err != nil {
		services.LogError("WeddingDestroy", err, map[string]interface{}{
			"project_id": project.ID,
			"user_id":    user.ID,
		})
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "Gagal menghapus undangan"})
		return
	}

	services.LogSuccess("WeddingDestroy", "Project dihapus: "+project.Slug)
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Undangan berhasil dihapus"})
}
