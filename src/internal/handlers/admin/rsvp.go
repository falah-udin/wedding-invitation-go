package admin

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"wedding-invitation-go/internal/database"
	"wedding-invitation-go/internal/models"
	"wedding-invitation-go/internal/services"
	"wedding-invitation-go/views/admin/projects"
	"wedding-invitation-go/views/layouts"
)

// ============================================
// RsvpIndex — GET /admin/projects/:id/rsvp
// Halaman daftar RSVP
// ============================================
func RsvpIndex(c *gin.Context) {
	setting := services.GetSiteSetting()
	user := c.MustGet("user").(*models.User)

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.Redirect(http.StatusFound, "/admin/projects?error=ID+tidak+valid")
		return
	}

	var project models.Project
	db := database.GetDB()
	if err := db.Preload("User").Preload("Template").First(&project, id).Error; err != nil {
		c.Redirect(http.StatusFound, "/admin/projects?error=Project+tidak+ditemukan")
		return
	}

	// Cek akses
	if user.Role == "client" && project.UserID != user.ID {
		c.Redirect(http.StatusFound, "/client/wedding?error=Tidak+memiliki+akses")
		return
	}

	// Ambil semua RSVP
	var rsvps []models.Rsvp
	db.Where("project_id = ?", project.ID).
		Order("created_at DESC").
		Find(&rsvps)

	errorMsg := c.Query("error")
	successMsg := c.Query("success")

	var buf bytes.Buffer
	err = layouts.AdminLayout(
		*user, setting, "admin.projects.rsvp",
		projects.RsvpContent(project, rsvps, errorMsg, successMsg),
	).Render(c.Request.Context(), &buf)

	if err != nil {
		c.String(http.StatusInternalServerError, "Render error: %v", err)
		return
	}

	c.Data(http.StatusOK, "text/html; charset=utf-8", buf.Bytes())
}

// ============================================
// RsvpDelete — DELETE /admin/projects/rsvp/:rsvpId
// Hapus RSVP
// ============================================
func RsvpDelete(c *gin.Context) {
	rsvpID, err := strconv.ParseUint(c.Param("rsvpId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "ID tidak valid"})
		return
	}

	db := database.GetDB()
	var rsvp models.Rsvp
	if err := db.First(&rsvp, rsvpID).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "RSVP tidak ditemukan"})
		return
	}

	var project models.Project
	db.First(&project, rsvp.ProjectID)

	db.Delete(&rsvp)
	if project.ID > 0 && project.TotalRsvp > 0 {
		project.TotalRsvp--
		db.Save(&project)
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "RSVP berhasil dihapus"})
}

// ============================================
// RsvpBulkDelete — POST /admin/projects/rsvp/bulk-delete
// Hapus banyak RSVP
// ============================================
func RsvpBulkDelete(c *gin.Context) {
	var req struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&req); err != nil || len(req.IDs) == 0 {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "Tidak ada data yang dipilih"})
		return
	}

	db := database.GetDB()

	// Ambil project_id dari salah satu RSVP
	var firstRsvp models.Rsvp
	if err := db.First(&firstRsvp, req.IDs[0]).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "Data tidak ditemukan"})
		return
	}

	projectID := firstRsvp.ProjectID

	// Hapus
	result := db.Where("id IN ?", req.IDs).Delete(&models.Rsvp{})
	deleted := result.RowsAffected

	// Update total_rsvp
	var project models.Project
	db.First(&project, projectID)
	var remaining int64
	db.Model(&models.Rsvp{}).Where("project_id = ?", projectID).Count(&remaining)
	project.TotalRsvp = int(remaining)
	db.Save(&project)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": fmt.Sprintf("Berhasil menghapus %d reservasi", deleted),
	})
}

// ============================================
// RsvpExport — GET /admin/projects/:id/rsvp/export
// Export CSV
// ============================================
func RsvpExport(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.Redirect(http.StatusFound, "/admin/projects?error=ID+tidak+valid")
		return
	}

	var project models.Project
	db := database.GetDB()
	if err := db.First(&project, id).Error; err != nil {
		c.Redirect(http.StatusFound, "/admin/projects?error=Project+tidak+ditemukan")
		return
	}

	var rsvps []models.Rsvp
	db.Where("project_id = ?", project.ID).Order("created_at DESC").Find(&rsvps)

	// Set header CSV
	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=reservasi-%s.csv", project.Slug))

	// Tulis CSV
	writer := csv.NewWriter(c.Writer)
	defer writer.Flush()

	writer.Write([]string{"Nama", "Kehadiran", "Jumlah Tamu", "Ucapan", "Tanggal"})
	for _, r := range rsvps {
		writer.Write([]string{
			r.GuestName,
			r.GetAttendanceLabel(),
			strconv.Itoa(r.TotalGuests),
			r.GetMessage(),
			r.CreatedAt.Format("02/01/2006 15:04"),
		})
	}
}
