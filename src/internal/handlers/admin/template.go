package admin

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
	"wedding-invitation-go/views/admin/templates"
	"wedding-invitation-go/views/layouts"
)

const uploadDir = "/app/public/storage/templates"

func TemplateIndex(c *gin.Context) {
	setting := services.GetSiteSetting()
	user := c.MustGet("user").(*models.User)

	var list []models.Template
	db := database.GetDB()
	db.Order("`order` ASC, id ASC").Find(&list)

	var buf bytes.Buffer
	err := layouts.AdminLayout(
		*user, setting, "admin.templates",
		templates.IndexContent(list),
	).Render(c.Request.Context(), &buf)

	if err != nil {
		c.String(http.StatusInternalServerError, "Render error: %v", err)
		return
	}

	c.Data(http.StatusOK, "text/html; charset=utf-8", buf.Bytes())
}

func TemplateCreate(c *gin.Context) {
	setting := services.GetSiteSetting()
	user := c.MustGet("user").(*models.User)
	errorMsg := c.Query("error")

	var buf bytes.Buffer
	err := layouts.AdminLayout(
		*user, setting, "admin.templates",
		templates.CreateContent(errorMsg),
	).Render(c.Request.Context(), &buf)

	if err != nil {
		c.String(http.StatusInternalServerError, "Render error: %v", err)
		return
	}

	c.Data(http.StatusOK, "text/html; charset=utf-8", buf.Bytes())
}

func TemplateStore(c *gin.Context) {
	name := strings.TrimSpace(c.PostForm("name"))
	slug := strings.TrimSpace(c.PostForm("slug"))
	folder := strings.TrimSpace(c.PostForm("folder"))
	description := strings.TrimSpace(c.PostForm("description"))
	orderStr := c.PostForm("order")

	if name == "" || slug == "" || folder == "" {
		c.Redirect(http.StatusFound, "/admin/templates/create?error=Semua+field+wajib+diisi")
		return
	}

	db := database.GetDB()

	var existing models.Template
	if err := db.Where("slug = ?", slug).First(&existing).Error; err == nil {
		c.Redirect(http.StatusFound, "/admin/templates/create?error=Slug+sudah+dipakai")
		return
	}

	var order int
	fmt.Sscanf(orderStr, "%d", &order)

	tmpl := models.Template{
		Name: name, Slug: slug, Folder: folder,
		Sections: "[]", Fields: "[]",
		IsActive: true, Order: order,
	}
	if description != "" {
		tmpl.Description = &description
	}

	if file, err := c.FormFile("thumbnail"); err == nil {
		thumbnailPath, uploadErr := saveUploadedFile(c, file, slug)
		if uploadErr != nil {
			c.Redirect(http.StatusFound, "/admin/templates/create?error=Gagal+upload+thumbnail")
			return
		}
		tmpl.Thumbnail = &thumbnailPath
	}

	if err := db.Create(&tmpl).Error; err != nil {
		c.Redirect(http.StatusFound, "/admin/templates/create?error=Gagal+menyimpan+template")
		return
	}

	c.Redirect(http.StatusFound, "/admin/templates")
}

func TemplateEdit(c *gin.Context) {
	setting := services.GetSiteSetting()
	user := c.MustGet("user").(*models.User)
	errorMsg := c.Query("error")

	id := c.Param("id")
	var tmpl models.Template
	db := database.GetDB()
	if err := db.First(&tmpl, id).Error; err != nil {
		c.Redirect(http.StatusFound, "/admin/templates?error=Template+tidak+ditemukan")
		return
	}

	var buf bytes.Buffer
	err := layouts.AdminLayout(
		*user, setting, "admin.templates",
		templates.EditContent(tmpl, errorMsg),
	).Render(c.Request.Context(), &buf)

	if err != nil {
		c.String(http.StatusInternalServerError, "Render error: %v", err)
		return
	}

	c.Data(http.StatusOK, "text/html; charset=utf-8", buf.Bytes())
}

func TemplateUpdate(c *gin.Context) {
	id := c.Param("id")
	name := strings.TrimSpace(c.PostForm("name"))
	slug := strings.TrimSpace(c.PostForm("slug"))
	folder := strings.TrimSpace(c.PostForm("folder"))
	description := strings.TrimSpace(c.PostForm("description"))
	orderStr := c.PostForm("order")

	db := database.GetDB()

	var tmpl models.Template
	if err := db.First(&tmpl, id).Error; err != nil {
		c.Redirect(http.StatusFound, "/admin/templates?error=Template+tidak+ditemukan")
		return
	}

	editURL := "/admin/templates/" + id + "/edit"

	if name == "" || slug == "" || folder == "" {
		c.Redirect(http.StatusFound, editURL+"?error=Semua+field+wajib+diisi")
		return
	}

	var existing models.Template
	if err := db.Where("slug = ? AND id != ?", slug, id).First(&existing).Error; err == nil {
		c.Redirect(http.StatusFound, editURL+"?error=Slug+sudah+dipakai+template+lain")
		return
	}

	var order int
	fmt.Sscanf(orderStr, "%d", &order)

	tmpl.Name = name
	tmpl.Slug = slug
	tmpl.Folder = folder
	tmpl.Order = order
	if description != "" {
		tmpl.Description = &description
	} else {
		tmpl.Description = nil
	}

	if file, err := c.FormFile("thumbnail"); err == nil {
		if tmpl.Thumbnail != nil && *tmpl.Thumbnail != "" {
			oldPath := strings.TrimPrefix(*tmpl.Thumbnail, "/storage/")
			os.Remove(filepath.Join("/app/public/storage", oldPath))
		}
		thumbnailPath, uploadErr := saveUploadedFile(c, file, slug)
		if uploadErr != nil {
			c.Redirect(http.StatusFound, editURL+"?error=Gagal+upload+thumbnail")
			return
		}
		tmpl.Thumbnail = &thumbnailPath
	}

	if err := db.Save(&tmpl).Error; err != nil {
		c.Redirect(http.StatusFound, editURL+"?error=Gagal+menyimpan+perubahan")
		return
	}

	c.Redirect(http.StatusFound, "/admin/templates")
}

func TemplateToggle(c *gin.Context) {
	id := c.Param("id")

	db := database.GetDB()
	var tmpl models.Template
	if err := db.First(&tmpl, id).Error; err != nil {
		c.Redirect(http.StatusFound, "/admin/templates?error=Template+tidak+ditemukan")
		return
	}

	tmpl.IsActive = !tmpl.IsActive
	db.Save(&tmpl)

	c.Redirect(http.StatusFound, "/admin/templates")
}

func TemplateDelete(c *gin.Context) {
	id := c.Param("id")

	db := database.GetDB()
	var tmpl models.Template
	if err := db.First(&tmpl, id).Error; err != nil {
		c.Redirect(http.StatusFound, "/admin/templates?error=Template+tidak+ditemukan")
		return
	}

	if tmpl.Thumbnail != nil && *tmpl.Thumbnail != "" {
		oldPath := strings.TrimPrefix(*tmpl.Thumbnail, "/storage/")
		os.Remove(filepath.Join("/app/public/storage", oldPath))
	}

	db.Delete(&tmpl)
	c.Redirect(http.StatusFound, "/admin/templates")
}

func saveUploadedFile(c *gin.Context, file *multipart.FileHeader, slug string) (string, error) {
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return "", err
	}

	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("%s_%d%s", slug, time.Now().Unix(), ext)
	fullPath := filepath.Join(uploadDir, filename)

	if err := c.SaveUploadedFile(file, fullPath); err != nil {
		return "", err
	}

	return "/storage/templates/" + filename, nil
}
