package admin

import (
	"bytes"
	"net/http"

	"github.com/gin-gonic/gin"

	"wedding-invitation-go/internal/database"
	"wedding-invitation-go/internal/models"
	"wedding-invitation-go/internal/services"
	"wedding-invitation-go/views/admin"
	"wedding-invitation-go/views/layouts"
)

// DashboardHandler — GET /admin/dashboard
// Render halaman dashboard admin
func DashboardHandler(c *gin.Context) {
	setting := services.GetSiteSetting()

	// Ambil user dari context (sudah di-set middleware RequireAuth)
	userVal, exists := c.Get("user")
	if !exists {
		c.Redirect(http.StatusFound, "/login")
		return
	}
	user := userVal.(*models.User)

	// Stats dari database
	db := database.GetDB()

	var totalClients int64
	db.Model(&models.User{}).Where("role = ?", "client").Count(&totalClients)

	var totalStaff int64
	db.Model(&models.User{}).Where("role = ?", "staff").Count(&totalStaff)

	var totalProjects int64
	db.Model(&models.Project{}).Count(&totalProjects)

	var totalTemplates int64
	db.Model(&models.Template{}).Count(&totalTemplates)

	var totalDraft int64
	db.Model(&models.Project{}).Where("status = ?", "draft").Count(&totalDraft)

	var totalPublished int64
	db.Model(&models.Project{}).Where("status = ?", "published").Count(&totalPublished)

	var totalRsvp int64
	// Nanti: db.Model(&models.Rsvp{}).Count(&totalRsvp)
	totalRsvp = 0

	var totalViews int64
	db.Model(&models.Project{}).Select("COALESCE(SUM(total_views), 0)").Scan(&totalViews)

	// Recent projects
	var recentProjects []models.Project
	db.Preload("User").Preload("Template").
		Order("created_at DESC").
		Limit(5).
		Find(&recentProjects)

	var buf bytes.Buffer
	err := layouts.AdminLayout(
		*user,
		setting,
		"admin.dashboard",
		admin.DashboardContent(
			int(totalClients),
			int(totalStaff),
			int(totalProjects),
			int(totalTemplates),
			int(totalDraft),
			int(totalPublished),
			int(totalRsvp),
			int(totalViews),
			recentProjects,
			*user,
		),
	).Render(c.Request.Context(), &buf)

	if err != nil {
		c.String(http.StatusInternalServerError, "Render error: %v", err)
		return
	}

	c.Data(http.StatusOK, "text/html; charset=utf-8", buf.Bytes())
}
