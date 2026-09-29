package client

import (
	"bytes"
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"

	"wedding-invitation-go/internal/database"
	"wedding-invitation-go/internal/models"
	"wedding-invitation-go/internal/services"
	"wedding-invitation-go/views/client/dashboard"
	"wedding-invitation-go/views/layouts"
)

// Dashboard — GET /client/dashboard
func Dashboard(c *gin.Context) {
	setting := services.GetSiteSetting()
	user := services.GetCurrentUser(c.Request)
	if user == nil {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	db := database.GetDB()

	// Ambil semua project milik client ini
	var projects []models.Project
	db.Where("user_id = ?", user.ID).Order("created_at DESC").Find(&projects)

	// Hitung statistik
	var totalProjects = len(projects)
	var totalPublished, totalViews, totalRsvp int
	for _, p := range projects {
		if p.Status == "published" {
			totalPublished++
		}
		totalViews += p.TotalViews
		totalRsvp += p.TotalRsvp
	}

	// Recent 5 (sudah di-order)
	recent := projects
	if len(recent) > 5 {
		recent = recent[:5]
	}

	// Sort ulang by created_at desc (untuk jaga-jaga)
	sort.Slice(recent, func(i, j int) bool {
		return recent[i].CreatedAt.After(recent[j].CreatedAt)
	})

	data := dashboard.DashboardData{
		User:           *user,
		Setting:        setting,
		TotalProjects:  totalProjects,
		TotalPublished: totalPublished,
		TotalViews:     totalViews,
		TotalRsvp:      totalRsvp,
		RecentProjects: recent,
	}

	var buf bytes.Buffer
	err := layouts.ClientLayout(*user, setting, "client.dashboard", dashboard.DashboardPage(data)).Render(c.Request.Context(), &buf)
	if err != nil {
		c.String(http.StatusInternalServerError, "Render error: %v", err)
		return
	}
	c.Data(http.StatusOK, "text/html; charset=utf-8", buf.Bytes())
}
