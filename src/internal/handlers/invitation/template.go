package invitation

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"wedding-invitation-go/internal/database"
	"wedding-invitation-go/internal/models"
	"wedding-invitation-go/internal/services"
	"wedding-invitation-go/views/invitation"
	"wedding-invitation-go/views/layouts"
)

// ============================================
// TemplateList — GET /invitation/create/template
// ============================================
func TemplateList(c *gin.Context) {
	setting := services.GetSiteSetting()
	user := c.MustGet("user").(*models.User)

	projectID := services.WizardGetProjectID(c.Request)
	if projectID == 0 {
		c.Redirect(http.StatusFound, "/invitation/create/select-client?error=Sesi+habis")
		return
	}

	var project models.Project
	db := database.GetDB()
	if err := db.First(&project, projectID).Error; err != nil {
		c.Redirect(http.StatusFound, "/invitation/create/select-client?error=Project+tidak+ditemukan")
		return
	}

	// Ambil semua template aktif
	var templates []models.Template
	db.Where("is_active = ?", true).Order("`order` ASC, id ASC").Find(&templates)

	// Handle ganti template via query ?template_id=X
	selectedTemplateID := c.Query("template_id")
	if selectedTemplateID != "" {
		var newTID uint
		for _, ch := range selectedTemplateID {
			if ch < '0' || ch > '9' {
				break
			}
			newTID = newTID*10 + uint(ch-'0')
		}
		if newTID > 0 {
			var newTmpl models.Template
			if err := db.First(&newTmpl, newTID).Error; err == nil {
				project.TemplateID = newTmpl.ID
				db.Save(&project)
			}
		}
	}

	// Ambil template aktif
	var selectedTemplate models.Template
	db.First(&selectedTemplate, project.TemplateID)
	fieldsSchema := selectedTemplate.GetFieldsSchema()

	// Ambil data existing (dengan field name yang sama)
	existingData := services.ParseJSONMap(project.TemplateSpecificData)

	errorMsg := c.Query("error")
	successMsg := c.Query("success")

	var buf bytes.Buffer
	err := layouts.AdminLayout(
		*user, setting, "invitation.create",
		invitation.TemplateContent(
			project,
			templates,
			selectedTemplate,
			fieldsSchema,
			existingData,
			errorMsg,
			successMsg,
		),
	).Render(c.Request.Context(), &buf)

	if err != nil {
		c.String(http.StatusInternalServerError, "Render error: %v", err)
		return
	}

	c.Data(http.StatusOK, "text/html; charset=utf-8", buf.Bytes())
}

// ============================================
// TemplateSave — POST /invitation/create/template
// ============================================
func TemplateSave(c *gin.Context) {
	projectID := services.WizardGetProjectID(c.Request)
	if projectID == 0 {
		c.Redirect(http.StatusFound, "/invitation/create/select-client?error=Sesi+habis")
		return
	}

	var project models.Project
	db := database.GetDB()
	if err := db.First(&project, projectID).Error; err != nil {
		c.Redirect(http.StatusFound, "/invitation/create/select-client?error=Project+tidak+ditemukan")
		return
	}

	// Update template ID (kalau user ganti)
	newTemplateIDStr := c.PostForm("template_id")
	if newTemplateIDStr != "" {
		var newTID uint
		for _, ch := range newTemplateIDStr {
			if ch < '0' || ch > '9' {
				break
			}
			newTID = newTID*10 + uint(ch-'0')
		}
		if newTID > 0 && newTID != project.TemplateID {
			var newTmpl models.Template
			if err := db.First(&newTmpl, newTID).Error; err == nil {
				project.TemplateID = newTmpl.ID
			}
		}
	}

	// Ambil template aktif
	var currentTemplate models.Template
	db.First(&currentTemplate, project.TemplateID)
	fieldsSchema := currentTemplate.GetFieldsSchema()

	// Data lama
	existingData := services.ParseJSONMap(project.TemplateSpecificData)

	// Data baru dari form
	newData := make(map[string]interface{})

	for _, field := range fieldsSchema {
		fieldName := field.Name

		switch field.Type {
		case "file":
			// Cek upload baru dulu
			if file, err := c.FormFile("specific." + fieldName); err == nil && file.Size > 0 {
				path, uploadErr := saveProjectFile(c, file, project.ID, "specific")
				if uploadErr == nil {
					newData[fieldName] = "/storage/" + path
				}
			} else {
				// Cek flag hapus
				if c.PostForm("specific."+fieldName+"_delete") == "1" {
					// Hapus file lama
					if oldVal, ok := existingData[fieldName]; ok {
						if oldStr, ok := oldVal.(string); ok && oldStr != "" {
							oldPath := strings.TrimPrefix(oldStr, "/storage/")
							os.Remove(filepath.Join("/app/public/storage", oldPath))
						}
					}
					delete(newData, fieldName)
				}
			}

		case "file_multiple":
			// Galeri — file existing + file baru
			var files []string

			// Ambil file existing
			if val, ok := existingData[fieldName]; ok {
				if arr, ok := val.([]interface{}); ok {
					for _, item := range arr {
						if s, ok := item.(string); ok {
							files = append(files, s)
						}
					}
				}
			}

			// Upload file baru
			if form, err := c.MultipartForm(); err == nil {
				if fileHeaders, ok := form.File["specific."+fieldName]; ok {
					for _, fileHeader := range fileHeaders {
						if fileHeader.Size > 0 {
							path, uploadErr := saveProjectFile(c, fileHeader, project.ID, "gallery")
							if uploadErr == nil {
								files = append(files, "/storage/"+path)
							}
						}
					}
				}
			}

			if len(files) > 0 {
				newData[fieldName] = files
			}

		case "repeater":
			items := parseRepeaterFromForm(c, fieldName)
			if len(items) > 0 {
				newData[fieldName] = items
			}

		default:
			// Text, textarea, url, select
			value := strings.TrimSpace(c.PostForm("specific." + fieldName))
			if value != "" {
				newData[fieldName] = value
			}
		}
	}

	// Merge data lama + baru (data lama tetap tersimpan)
	mergedData := services.MergeTemplateData(existingData, newData)

	// Simpan
	project.TemplateSpecificData = services.ToJSONString(mergedData)
	if err := db.Save(&project).Error; err != nil {
		c.Redirect(http.StatusFound, "/invitation/create/template?error=Gagal+menyimpan")
		return
	}

	c.Redirect(http.StatusFound, "/invitation/create/preview")
}

// ============================================
// HELPER — save uploaded file
// ============================================
func saveProjectFile(c *gin.Context, file *multipart.FileHeader, projectID uint, subdir string) (string, error) {
	uploadDir := fmt.Sprintf("/app/public/storage/projects/%d/%s", projectID, subdir)
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return "", err
	}

	ext := filepath.Ext(file.Filename)
	cleanName := strings.ReplaceAll(filepath.Base(file.Filename), " ", "_")
	filename := fmt.Sprintf("%d_%s", time.Now().UnixNano()/1e6, cleanName)
	if len(filename) > 100 {
		filename = filename[:100] + ext
	}

	fullPath := filepath.Join(uploadDir, filename)
	if err := c.SaveUploadedFile(file, fullPath); err != nil {
		return "", err
	}

	return fmt.Sprintf("projects/%d/%s/%s", projectID, subdir, filename), nil
}

// ============================================
// HELPER — parse repeater dari form
// ============================================
func parseRepeaterFromForm(c *gin.Context, fieldName string) []map[string]string {
	var items []map[string]string

	for i := 0; i < 20; i++ {
		titleKey := fmt.Sprintf("specific.%s[%d][title]", fieldName, i)
		descKey := fmt.Sprintf("specific.%s[%d][desc]", fieldName, i)

		title := strings.TrimSpace(c.PostForm(titleKey))
		desc := strings.TrimSpace(c.PostForm(descKey))

		if title == "" && desc == "" {
			continue
		}

		items = append(items, map[string]string{
			"title": title,
			"desc":  desc,
		})
	}

	return items
}
