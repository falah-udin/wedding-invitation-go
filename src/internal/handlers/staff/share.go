package staff

import (
	"bytes"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"wedding-invitation-go/internal/database"
	"wedding-invitation-go/internal/models"
	"wedding-invitation-go/internal/services"
	"wedding-invitation-go/views/layouts"
	"wedding-invitation-go/views/staff/share"
)

// ShareIndex — GET /staff/share
// List project published + stats share
func ShareIndex(c *gin.Context) {
	setting := services.GetSiteSetting()
	user := services.GetCurrentUser(c.Request)
	if user == nil {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	db := database.GetDB()

	// Filter
	search := strings.TrimSpace(c.Query("search"))
	sortBy := c.Query("sort")
	if sortBy == "" {
		sortBy = "latest"
	}

	// Base query: project published
	query := db.Model(&models.Project{}).Where("status = ?", "published")

	// Search
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
	case "guests_desc":
		// Sort by guest count — perlu subquery, kita pakai raw
		query = query.Order("(SELECT COUNT(*) FROM invitation_guests WHERE project_id = projects.id) DESC")
	default: // latest
		query = query.Order("projects.created_at DESC")
	}

	// Ambil project
	var projects []models.Project
	query.Preload("User").Preload("Template").Find(&projects)

	// Hitung stats per project (guests, shared, with_phone)
	statsMap := make(map[uint]share.ProjectStats)

	var totalGuests, totalShared, totalWithPhone int

	for _, p := range projects {
		var total, shared, withPhone int64

		db.Model(&models.InvitationGuest{}).Where("project_id = ?", p.ID).Count(&total)
		db.Model(&models.InvitationGuest{}).Where("project_id = ? AND is_shared = ?", p.ID, true).Count(&shared)
		db.Model(&models.InvitationGuest{}).
			Where("project_id = ? AND phone IS NOT NULL AND phone != ''", p.ID).
			Count(&withPhone)

		statsMap[p.ID] = share.ProjectStats{
			GuestsCount:    int(total),
			SharedCount:    int(shared),
			WithPhoneCount: int(withPhone),
		}

		totalGuests += int(total)
		totalShared += int(shared)
		totalWithPhone += int(withPhone)
	}

	// Total published projects (global, tidak terfilter)
	var totalPublished int64
	db.Model(&models.Project{}).Where("status = ?", "published").Count(&totalPublished)

	data := share.ShareListData{
		User:           *user,
		Setting:        setting,
		Projects:       projects,
		StatsMap:       statsMap,
		TotalPublished: int(totalPublished),
		TotalGuests:    totalGuests,
		TotalShared:    totalShared,
		TotalWithPhone: totalWithPhone,
		Search:         search,
		SortBy:         sortBy,
	}

	var buf bytes.Buffer
	err := layouts.StaffLayout(*user, setting, "staff.share", share.ShareListPage(data)).Render(c.Request.Context(), &buf)
	if err != nil {
		c.String(http.StatusInternalServerError, "Render error: %v", err)
		return
	}
	c.Data(http.StatusOK, "text/html; charset=utf-8", buf.Bytes())
}
