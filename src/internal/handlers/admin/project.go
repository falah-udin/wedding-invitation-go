package admin

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"wedding-invitation-go/internal/database"
	"wedding-invitation-go/internal/models"
	"wedding-invitation-go/internal/services"
	"wedding-invitation-go/views/admin/projects"
	"wedding-invitation-go/views/layouts"
)

// ============================================
// ProjectIndex — GET /admin/projects
// ============================================
func ProjectIndex(c *gin.Context) {
	setting := services.GetSiteSetting()
	user := c.MustGet("user").(*models.User)

	search := strings.TrimSpace(c.Query("search"))
	statusFilter := c.Query("status")
	templateFilter := c.Query("template_id")

	db := database.GetDB()

	query := db.Model(&models.Project{}).Preload("User").Preload("Template")

	if statusFilter != "" {
		query = query.Where("status = ?", statusFilter)
	}
	if templateFilter != "" {
		query = query.Where("template_id = ?", templateFilter)
	}
	if search != "" {
		query = query.Where("title LIKE ?", "%"+search+"%")
	}

	var list []models.Project
	query.Order("created_at DESC").Find(&list)

	var totalProjects, totalDraft, totalPublished, totalArchived int64
	db.Model(&models.Project{}).Count(&totalProjects)
	db.Model(&models.Project{}).Where("status = ?", "draft").Count(&totalDraft)
	db.Model(&models.Project{}).Where("status = ?", "published").Count(&totalPublished)
	db.Model(&models.Project{}).Where("status = ?", "archived").Count(&totalArchived)

	var totalViews int64
	db.Model(&models.Project{}).Select("COALESCE(SUM(total_views), 0)").Scan(&totalViews)

	var totalRsvp int64
	db.Model(&models.Project{}).Select("COALESCE(SUM(total_rsvp), 0)").Scan(&totalRsvp)

	var templates []models.Template
	db.Order("`order` ASC, id ASC").Find(&templates)

	var buf bytes.Buffer
	err := layouts.AdminLayout(
		*user, setting, "admin.projects",
		projects.IndexContent(
			list, templates,
			search, statusFilter, templateFilter,
			totalProjects, totalDraft, totalPublished, totalArchived,
			totalViews, totalRsvp,
		),
	).Render(c.Request.Context(), &buf)

	if err != nil {
		c.String(http.StatusInternalServerError, "Render error: %v", err)
		return
	}

	c.Data(http.StatusOK, "text/html; charset=utf-8", buf.Bytes())
}

// ============================================
// ProjectShow — GET /admin/projects/:id
// ============================================
func ProjectShow(c *gin.Context) {
	setting := services.GetSiteSetting()
	user := c.MustGet("user").(*models.User)

	id := c.Param("id")
	var project models.Project
	db := database.GetDB()
	if err := db.Preload("User").Preload("Template").First(&project, id).Error; err != nil {
		c.Redirect(http.StatusFound, "/admin/projects?error=Project+tidak+ditemukan")
		return
	}

	dataUndangan := parseJSONMap(project.DataUndangan)

	var buf bytes.Buffer
	err := layouts.AdminLayout(
		*user, setting, "admin.projects",
		projects.ShowContent(project, dataUndangan),
	).Render(c.Request.Context(), &buf)

	if err != nil {
		c.String(http.StatusInternalServerError, "Render error: %v", err)
		return
	}

	c.Data(http.StatusOK, "text/html; charset=utf-8", buf.Bytes())
}

// ============================================
// ProjectDelete — POST /admin/projects/:id/delete
// ============================================
func ProjectDelete(c *gin.Context) {
	id := c.Param("id")

	db := database.GetDB()
	var project models.Project
	if err := db.First(&project, id).Error; err != nil {
		c.Redirect(http.StatusFound, "/admin/projects?error=Project+tidak+ditemukan")
		return
	}

	db.Delete(&project)
	c.Redirect(http.StatusFound, "/admin/projects")
}


// ============================================
// ProjectChangeTemplate — POST /admin/projects/:id/change-template
// Ganti template project (admin — tanpa cek ownership)
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

        // Admin bisa pilih template apapun (termasuk yang non-aktif)
        var tmpl models.Template
        if err := db.First(&tmpl, req.TemplateID).Error; err != nil {
                c.JSON(http.StatusOK, gin.H{"success": false, "message": "Template tidak ditemukan"})
                return
        }

        if project.TemplateID == req.TemplateID {
                c.JSON(http.StatusOK, gin.H{"success": true, "message": "Template tidak berubah", "no_change": true})
                return
        }

        if err := db.Model(&project).Update("template_id", req.TemplateID).Error; err != nil {
                services.LogError("ProjectChangeTemplate", err, map[string]interface{}{
                        "project_id":  project.ID,
                        "template_id": req.TemplateID,
                        "admin_id":    user.ID,
                })
                c.JSON(http.StatusOK, gin.H{"success": false, "message": "Gagal mengganti template"})
                return
        }

        services.LogSuccess("ProjectChangeTemplate", "Admin ganti template: "+project.Slug+" -> "+tmpl.Name)
        c.JSON(http.StatusOK, gin.H{
                "success": true,
                "message": "Template berhasil diganti ke " + tmpl.Name,
        })
}

// ============================================
// HELPER — parse JSON string ke map[string]string
// ============================================
func parseJSONMap(jsonStr string) map[string]string {
	result := make(map[string]string)
	if jsonStr == "" {
		return result
	}

	var raw map[string]interface{}
	if err := json.Unmarshal([]byte(jsonStr), &raw); err != nil {
		return result
	}

	for k, v := range raw {
		if v == nil {
			result[k] = "-"
			continue
		}
		switch val := v.(type) {
		case string:
			if val == "" {
				result[k] = "-"
			} else {
				result[k] = val
			}
		case float64:
			result[k] = fmt.Sprintf("%.0f", val)
		case bool:
			if val {
				result[k] = "Ya"
			} else {
				result[k] = "Tidak"
			}
		default:
			result[k] = fmt.Sprintf("%v", val)
		}
	}

	return result
}
