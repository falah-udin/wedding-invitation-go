package client

import (
	"bytes"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"wedding-invitation-go/internal/database"
	"wedding-invitation-go/internal/models"
	"wedding-invitation-go/internal/services"
	"wedding-invitation-go/views/client/reservations"
	"wedding-invitation-go/views/layouts"
)

// ReservationIndex — GET /client/reservations
func ReservationIndex(c *gin.Context) {
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

	rsvpCounts := make(map[uint]int)
	for _, p := range projects {
		var count int64
		db.Model(&models.Rsvp{}).Where("project_id = ?", p.ID).Count(&count)
		rsvpCounts[p.ID] = int(count)
	}

	data := reservations.IndexData{
		User:       *user,
		Setting:    setting,
		Projects:   projects,
		RsvpCounts: rsvpCounts,
	}

	var buf bytes.Buffer
	err := layouts.ClientLayout(*user, setting, "client.reservations", reservations.IndexPage(data)).Render(c.Request.Context(), &buf)
	if err != nil {
		c.String(http.StatusInternalServerError, "Render error: %v", err)
		return
	}
	c.Data(http.StatusOK, "text/html; charset=utf-8", buf.Bytes())
}

// ReservationShow — GET /client/reservations/:id
func ReservationShow(c *gin.Context) {
	setting := services.GetSiteSetting()
	user := services.GetCurrentUser(c.Request)
	if user == nil {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.Redirect(http.StatusFound, "/client/reservations?error=ID+tidak+valid")
		return
	}

	db := database.GetDB()
	var project models.Project
	if err := db.Preload("Template").First(&project, id).Error; err != nil {
		c.Redirect(http.StatusFound, "/client/reservations?error=Undangan+tidak+ditemukan")
		return
	}

	if project.UserID != user.ID {
		c.Redirect(http.StatusFound, "/client/reservations?error=Tidak+memiliki+akses")
		return
	}

	var rsvps []models.Rsvp
	db.Where("project_id = ?", project.ID).Order("created_at DESC").Find(&rsvps)

	stats := map[string]int{
		"total":       len(rsvps),
		"hadir":       0,
		"tidak_hadir": 0,
		"ragu":        0,
	}
	for _, r := range rsvps {
		switch r.Attendance {
		case "hadir":
			stats["hadir"]++
		case "tidak_hadir":
			stats["tidak_hadir"]++
		case "ragu":
			stats["ragu"]++
		}
	}

	data := reservations.ShowData{
		User:    *user,
		Setting: setting,
		Project: project,
		Rsvps:   rsvps,
		Stats:   stats,
	}

	var buf bytes.Buffer
	err = layouts.ClientLayout(*user, setting, "client.reservations", reservations.ShowPage(data)).Render(c.Request.Context(), &buf)
	if err != nil {
		c.String(http.StatusInternalServerError, "Render error: %v", err)
		return
	}
	c.Data(http.StatusOK, "text/html; charset=utf-8", buf.Bytes())
}

// ReservationDelete — DELETE /client/reservations/:id/rsvp/:rsvpId
func ReservationDelete(c *gin.Context) {
	user := services.GetCurrentUser(c.Request)
	if user == nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "Sesi habis"})
		return
	}

	projectID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "ID project tidak valid"})
		return
	}

	rsvpID, err := strconv.ParseUint(c.Param("rsvpId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "ID RSVP tidak valid"})
		return
	}

	db := database.GetDB()

	var project models.Project
	if err := db.First(&project, projectID).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "Project tidak ditemukan"})
		return
	}
	if project.UserID != user.ID {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "Tidak memiliki akses"})
		return
	}

	var rsvp models.Rsvp
	if err := db.Where("id = ? AND project_id = ?", rsvpID, projectID).First(&rsvp).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "RSVP tidak ditemukan"})
		return
	}

	if err := db.Delete(&rsvp).Error; err != nil {
		services.LogError("ReservationDelete", err, map[string]interface{}{
			"project_id": projectID,
			"rsvp_id":    rsvpID,
		})
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "Gagal menghapus RSVP"})
		return
	}

	if project.TotalRsvp > 0 {
		project.TotalRsvp--
		db.Save(&project)
	}

	services.LogSuccess("ReservationDelete", "RSVP dihapus")
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "RSVP berhasil dihapus"})
}
