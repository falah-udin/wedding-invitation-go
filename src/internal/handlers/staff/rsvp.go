package staff

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"wedding-invitation-go/internal/database"
	"wedding-invitation-go/internal/models"
	"wedding-invitation-go/internal/services"
	"wedding-invitation-go/views/layouts"
	"wedding-invitation-go/views/staff/rsvp"
)

// RsvpListIndex — GET /staff/rsvp
// List semua project dengan count RSVP
func RsvpListIndex(c *gin.Context) {
	setting := services.GetSiteSetting()
	user := services.GetCurrentUser(c.Request)
	if user == nil {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	db := database.GetDB()

	search := strings.TrimSpace(c.Query("search"))
	sortBy := c.Query("sort")
	if sortBy == "" {
		sortBy = "latest"
	}

	// Base query
	query := db.Model(&models.Project{})

	if search != "" {
		searchLower := "%" + strings.ToLower(search) + "%"
		query = query.Where(
			"LOWER(projects.title) LIKE ? OR projects.user_id IN (SELECT id FROM users WHERE LOWER(name) LIKE ?)",
			searchLower, searchLower,
		)
	}

	// Sort
	switch sortBy {
	case "title_asc":
		query = query.Order("projects.title ASC")
	case "rsvp_desc":
		query = query.Order("(SELECT COUNT(*) FROM rsvps WHERE project_id = projects.id) DESC")
	case "rsvp_asc":
		query = query.Order("(SELECT COUNT(*) FROM rsvps WHERE project_id = projects.id) ASC")
	default:
		query = query.Order("projects.created_at DESC")
	}

	var projects []models.Project
	query.Preload("User").Preload("Template").Find(&projects)

	// Filter hanya project yang punya RSVP (opsional) — bisa hapus kalau mau tampil semua
	// Untuk sekarang tampilkan semua project dengan count RSVP
	var filtered []models.Project
	for _, p := range projects {
		var count int64
		db.Model(&models.Rsvp{}).Where("project_id = ?", p.ID).Count(&count)
		if count > 0 {
			filtered = append(filtered, p)
		}
	}
	projects = filtered

	// Statistik per project
	type ProjectStats struct {
		RsvpCount      int
		HadirCount     int
		TidakHadirCount int
		RaguCount      int
	}
	statsMap := make(map[uint]rsvp.ProjectRsvpStats)

	// Total global
	var totalRsvp, totalHadir, totalTidakHadir, totalRagu int
	var totalGuests int

	for _, p := range projects {
		var rsvpCount, hadirCount, tidakCount, raguCount int64
		db.Model(&models.Rsvp{}).Where("project_id = ?", p.ID).Count(&rsvpCount)
		db.Model(&models.Rsvp{}).Where("project_id = ? AND attendance = ?", p.ID, "hadir").Count(&hadirCount)
		db.Model(&models.Rsvp{}).Where("project_id = ? AND attendance = ?", p.ID, "tidak_hadir").Count(&tidakCount)
		db.Model(&models.Rsvp{}).Where("project_id = ? AND attendance = ?", p.ID, "ragu").Count(&raguCount)

		// Total guests (sum total_guests)
		var guestsSum int64
		db.Model(&models.Rsvp{}).
			Where("project_id = ? AND attendance = ?", p.ID, "hadir").
			Select("COALESCE(SUM(total_guests), 0)").
			Scan(&guestsSum)

		statsMap[p.ID] = rsvp.ProjectRsvpStats{
			RsvpCount:       int(rsvpCount),
			HadirCount:      int(hadirCount),
			TidakHadirCount: int(tidakCount),
			RaguCount:       int(raguCount),
		}

		totalRsvp += int(rsvpCount)
		totalHadir += int(hadirCount)
		totalTidakHadir += int(tidakCount)
		totalRagu += int(raguCount)
		totalGuests += int(guestsSum)
	}

	data := rsvp.RsvpListData{
		User:            *user,
		Setting:         setting,
		Projects:        projects,
		StatsMap:        statsMap,
		TotalRsvp:       totalRsvp,
		TotalHadir:      totalHadir,
		TotalTidakHadir: totalTidakHadir,
		TotalRagu:       totalRagu,
		TotalGuests:     totalGuests,
		Search:          search,
		SortBy:          sortBy,
	}

	var buf bytes.Buffer
	err := layouts.StaffLayout(*user, setting, "staff.rsvp", rsvp.RsvpListPage(data)).Render(c.Request.Context(), &buf)
	if err != nil {
		c.String(http.StatusInternalServerError, "Render error: %v", err)
		return
	}
	c.Data(http.StatusOK, "text/html; charset=utf-8", buf.Bytes())
}

// RsvpProjectIndex — GET /staff/projects/:id/rsvp
// Kelola RSVP per project (tabel + bulk delete + export)
func RsvpProjectIndex(c *gin.Context) {
	setting := services.GetSiteSetting()
	user := services.GetCurrentUser(c.Request)
	if user == nil {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.Redirect(http.StatusFound, "/staff/projects?error=ID+tidak+valid")
		return
	}

	db := database.GetDB()
	var project models.Project
	if err := db.Preload("User").Preload("Template").First(&project, id).Error; err != nil {
		c.Redirect(http.StatusFound, "/staff/projects?error=Project+tidak+ditemukan")
		return
	}

	// Ambil RSVP
	var rsvps []models.Rsvp
	db.Where("project_id = ?", project.ID).Order("created_at DESC").Find(&rsvps)

	// Hitung stats
	stats := map[string]int{
		"total":       len(rsvps),
		"hadir":       0,
		"tidak_hadir": 0,
		"ragu":        0,
		"total_guests": 0,
	}
	for _, r := range rsvps {
		switch r.Attendance {
		case "hadir":
			stats["hadir"]++
			stats["total_guests"] += r.TotalGuests
		case "tidak_hadir":
			stats["tidak_hadir"]++
		case "ragu":
			stats["ragu"]++
		}
	}

	data := rsvp.RsvpProjectData{
		User:    *user,
		Setting: setting,
		Project: project,
		Rsvps:   rsvps,
		Stats:   stats,
	}

	var buf bytes.Buffer
	err = layouts.StaffLayout(*user, setting, "staff.rsvp", rsvp.RsvpProjectPage(data)).Render(c.Request.Context(), &buf)
	if err != nil {
		c.String(http.StatusInternalServerError, "Render error: %v", err)
		return
	}
	c.Data(http.StatusOK, "text/html; charset=utf-8", buf.Bytes())
}

// RsvpDelete — DELETE /staff/projects/:id/rsvp/:rsvpId
func RsvpDelete(c *gin.Context) {
	user := services.GetCurrentUser(c.Request)
	if user == nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "Sesi habis"})
		return
	}

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

// RsvpBulkDelete — POST /staff/projects/:id/rsvp/bulk-delete
func RsvpBulkDelete(c *gin.Context) {
	user := services.GetCurrentUser(c.Request)
	if user == nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "Sesi habis"})
		return
	}

	var body struct {
		IDs []uint `json:"ids"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "Data tidak valid"})
		return
	}

	if len(body.IDs) == 0 {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "Tidak ada RSVP yang dipilih"})
		return
	}

	db := database.GetDB()
	result := db.Where("id IN ?", body.IDs).Delete(&models.Rsvp{})

	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "Gagal menghapus"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": fmt.Sprintf("%d RSVP berhasil dihapus", result.RowsAffected),
	})
}

// RsvpExport — GET /staff/projects/:id/rsvp/export
// Export CSV
func RsvpExport(c *gin.Context) {
	user := services.GetCurrentUser(c.Request)
	if user == nil {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.Redirect(http.StatusFound, "/staff/projects?error=ID+tidak+valid")
		return
	}

	db := database.GetDB()
	var project models.Project
	if err := db.First(&project, id).Error; err != nil {
		c.Redirect(http.StatusFound, "/staff/projects?error=Project+tidak+ditemukan")
		return
	}

	var rsvps []models.Rsvp
	db.Where("project_id = ?", project.ID).Order("created_at DESC").Find(&rsvps)

	// Set header
	c.Header("Content-Type", "text/csv; charset=utf-8")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=rsvp-%s.csv", project.Slug))

	writer := csv.NewWriter(c.Writer)
	defer writer.Flush()

	// Header CSV
	_ = writer.Write([]string{"No", "Nama", "Kehadiran", "Jumlah", "Ucapan", "Tanggal"})

	for i, r := range rsvps {
		attendance := r.GetAttendanceLabel()
		message := r.GetMessage()
		_ = writer.Write([]string{
			fmt.Sprintf("%d", i+1),
			r.GuestName,
			attendance,
			fmt.Sprintf("%d", r.TotalGuests),
			message,
			r.CreatedAt.Format("02/01/2006 15:04"),
		})
	}
}
