package staff

import (
	"bytes"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"wedding-invitation-go/internal/database"
	"wedding-invitation-go/internal/models"
	"wedding-invitation-go/internal/services"
	"wedding-invitation-go/views/layouts"
	"wedding-invitation-go/views/staff/dashboard"
)

// Dashboard — GET /staff/dashboard
func Dashboard(c *gin.Context) {
	setting := services.GetSiteSetting()
	user := services.GetCurrentUser(c.Request)
	if user == nil {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	db := database.GetDB()

	// Stats utama
	var totalProjects, totalDraft, totalPublished, totalArchived int64
	db.Model(&models.Project{}).Count(&totalProjects)
	db.Model(&models.Project{}).Where("status = ?", "draft").Count(&totalDraft)
	db.Model(&models.Project{}).Where("status = ?", "published").Count(&totalPublished)
	db.Model(&models.Project{}).Where("status = ?", "archived").Count(&totalArchived)

	// Sum views & rsvp
	var totalViews, totalRsvp int64
	db.Model(&models.Project{}).Select("COALESCE(SUM(total_views), 0)").Scan(&totalViews)
	db.Model(&models.Project{}).Select("COALESCE(SUM(total_rsvp), 0)").Scan(&totalRsvp)

	// Total client
	var totalClients int64
	db.Model(&models.User{}).Where("role = ?", "client").Count(&totalClients)

	// Stale drafts (> 7 hari belum publish)
	sevenDaysAgo := time.Now().AddDate(0, 0, -7)
	var staleDrafts []models.Project
	db.Preload("User").Preload("Template").
		Where("status = ? AND created_at < ?", "draft", sevenDaysAgo).
		Order("created_at ASC").
		Limit(10).
		Find(&staleDrafts)

	// Recent projects (10 terbaru)
	var recentProjects []models.Project
	db.Preload("User").Preload("Template").
		Order("created_at DESC").
		Limit(10).
		Find(&recentProjects)

	// Top viewed (5 teratas)
	var topViewed []models.Project
	db.Preload("User").
		Where("total_views > ?", 0).
		Order("total_views DESC").
		Limit(5).
		Find(&topViewed)

	data := dashboard.DashboardData{
		User:           *user,
		Setting:        setting,
		TotalProjects:  int(totalProjects),
		TotalDraft:     int(totalDraft),
		TotalPublished: int(totalPublished),
		TotalArchived:  int(totalArchived),
		TotalViews:     int(totalViews),
		TotalRsvp:      int(totalRsvp),
		TotalClients:   int(totalClients),
		StaleDrafts:    staleDrafts,
		RecentProjects: recentProjects,
		TopViewed:      topViewed,
	}

	var buf bytes.Buffer
	err := layouts.StaffLayout(*user, setting, "staff.dashboard", dashboard.DashboardPage(data)).Render(c.Request.Context(), &buf)
	if err != nil {
		c.String(http.StatusInternalServerError, "Render error: %v", err)
		return
	}
	c.Data(http.StatusOK, "text/html; charset=utf-8", buf.Bytes())
}
