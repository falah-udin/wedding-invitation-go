package staff

import (
	"bytes"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"wedding-invitation-go/internal/database"
	"wedding-invitation-go/internal/models"
	"wedding-invitation-go/internal/services"
	"wedding-invitation-go/views/layouts"
	"wedding-invitation-go/views/staff/projects"
)

// ProjectIndex — GET /staff/projects
func ProjectIndex(c *gin.Context) {
	setting := services.GetSiteSetting()
	user := services.GetCurrentUser(c.Request)
	if user == nil {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	db := database.GetDB()

	search := strings.TrimSpace(c.Query("search"))
	status := strings.TrimSpace(c.Query("status"))
	templateID := strings.TrimSpace(c.Query("template_id"))

	pageStr := c.Query("page")
	page := 1
	if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
		page = p
	}
	perPage := 20

	query := db.Model(&models.Project{})

	if search != "" {
		searchLower := "%" + strings.ToLower(search) + "%"
		query = query.Where(
			"LOWER(projects.title) LIKE ? OR projects.user_id IN (SELECT id FROM users WHERE LOWER(name) LIKE ?)",
			searchLower, searchLower,
		)
	}
	if status != "" {
		query = query.Where("projects.status = ?", status)
	}
	if templateID != "" {
		if tid, err := strconv.Atoi(templateID); err == nil && tid > 0 {
			query = query.Where("projects.template_id = ?", tid)
		}
	}

	var total int64
	query.Count(&total)

	var projectsList []models.Project
	offset := (page - 1) * perPage
	query.Preload("User").Preload("Template").
		Order("projects.created_at DESC").
		Limit(perPage).
		Offset(offset).
		Find(&projectsList)

	var templates []models.Template
	db.Where("is_active = ?", true).Order("`order` ASC, id ASC").Find(&templates)

	totalPages := int((total + int64(perPage) - 1) / int64(perPage))
	if totalPages < 1 {
		totalPages = 1
	}

	data := projects.ProjectListData{
		User:       *user,
		Setting:    setting,
		Projects:   projectsList,
		Templates:  templates,
		Search:     search,
		Status:     status,
		TemplateID: templateID,
		Page:       page,
		TotalPages: totalPages,
		Total:      int(total),
		PerPage:    perPage,
	}

	var buf bytes.Buffer
	err := layouts.StaffLayout(*user, setting, "staff.projects", projects.ProjectListPage(data)).Render(c.Request.Context(), &buf)
	if err != nil {
		c.String(http.StatusInternalServerError, "Render error: %v", err)
		return
	}
	c.Data(http.StatusOK, "text/html; charset=utf-8", buf.Bytes())
}

// ProjectShow — GET /staff/projects/:id
func ProjectShow(c *gin.Context) {
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

	allData := services.ParseJSONMap(project.DataUndangan)
	specificData := services.ParseJSONMap(project.TemplateSpecificData)
	for k, v := range specificData {
		allData[k] = v
	}

	data := projects.ProjectShowData{
		User:    *user,
		Setting: setting,
		Project: project,
		AllData: allData,
	}

	var buf bytes.Buffer
	err = layouts.StaffLayout(*user, setting, "staff.projects", projects.ProjectShowPage(data)).Render(c.Request.Context(), &buf)
	if err != nil {
		c.String(http.StatusInternalServerError, "Render error: %v", err)
		return
	}
	c.Data(http.StatusOK, "text/html; charset=utf-8", buf.Bytes())
}


// ============================================
// ProjectChangeTemplate — POST /staff/projects/:id/change-template
// Ganti template project (staff — tanpa cek ownership)
// ============================================
func ProjectChangeTemplate(c *gin.Context) {
        user := c.MustGet("user").(*models.User)

        idStr := c.Param("id")
        var projectID uint
        for _, ch := range idStr {
                if ch < '0' || ch > '9' {
                        break
                }
                projectID = projectID*10 + uint(ch-'0')
        }

        if projectID == 0 {
                c.JSON(http.StatusOK, gin.H{"success": false, "message": "ID project tidak valid"})
                return
        }

        var req struct {
                TemplateID uint `json:"template_id"`
        }
        if err := c.ShouldBindJSON(&req); err != nil {
                c.JSON(http.StatusOK, gin.H{"success": false, "message": "Data tidak valid"})
                return
        }

        if req.TemplateID == 0 {
                c.JSON(http.StatusOK, gin.H{"success": false, "message": "Template harus dipilih"})
                return
        }

        db := database.GetDB()

        var project models.Project
        if err := db.First(&project, projectID).Error; err != nil {
                c.JSON(http.StatusOK, gin.H{"success": false, "message": "Project tidak ditemukan"})
                return
        }

        var tmpl models.Template
        if err := db.Where("id = ? AND is_active = ?", req.TemplateID, true).First(&tmpl).Error; err != nil {
                c.JSON(http.StatusOK, gin.H{"success": false, "message": "Template tidak ditemukan atau tidak aktif"})
                return
        }

        if project.TemplateID == req.TemplateID {
                c.JSON(http.StatusOK, gin.H{"success": true, "message": "Template tidak berubah", "no_change": true})
                return
        }

        if err := db.Model(&project).Update("template_id", req.TemplateID).Error; err != nil {
                services.LogError("StaffProjectChangeTemplate", err, map[string]interface{}{
                        "project_id":  project.ID,
                        "template_id": req.TemplateID,
                        "staff_id":    user.ID,
                })
                c.JSON(http.StatusOK, gin.H{"success": false, "message": "Gagal mengganti template"})
                return
        }

        services.LogSuccess("StaffProjectChangeTemplate", "Staff ganti template: "+project.Slug+" -> "+tmpl.Name)
        c.JSON(http.StatusOK, gin.H{
                "success": true,
                "message": "Template berhasil diganti ke " + tmpl.Name,
        })
}
