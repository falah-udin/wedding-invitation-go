package admin

import (
	"bytes"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"wedding-invitation-go/internal/database"
	"wedding-invitation-go/internal/models"
	"wedding-invitation-go/internal/services"
	"wedding-invitation-go/views/admin/users"
	"wedding-invitation-go/views/layouts"
)

// ============================================
// UserIndex — GET /admin/users
// ============================================
func UserIndex(c *gin.Context) {
	setting := services.GetSiteSetting()
	user := c.MustGet("user").(*models.User)

	var usersList []models.User
	db := database.GetDB()
	db.Order("id DESC").Find(&usersList)

	var buf bytes.Buffer
	err := layouts.AdminLayout(
		*user,
		setting,
		"admin.users",
		users.IndexContent(usersList),
	).Render(c.Request.Context(), &buf)

	if err != nil {
		c.String(http.StatusInternalServerError, "Render error: %v", err)
		return
	}

	c.Data(http.StatusOK, "text/html; charset=utf-8", buf.Bytes())
}

// ============================================
// UserCreate — GET /admin/users/create
// ============================================
func UserCreate(c *gin.Context) {
	setting := services.GetSiteSetting()
	user := c.MustGet("user").(*models.User)
	errorMsg := c.Query("error")

	var buf bytes.Buffer
	err := layouts.AdminLayout(
		*user,
		setting,
		"admin.users",
		users.CreateContent(errorMsg),
	).Render(c.Request.Context(), &buf)

	if err != nil {
		c.String(http.StatusInternalServerError, "Render error: %v", err)
		return
	}

	c.Data(http.StatusOK, "text/html; charset=utf-8", buf.Bytes())
}

// ============================================
// UserStore — POST /admin/users/create
// ============================================
func UserStore(c *gin.Context) {
	name := strings.TrimSpace(c.PostForm("name"))
	email := strings.ToLower(strings.TrimSpace(c.PostForm("email")))
	password := c.PostForm("password")
	passwordConfirm := c.PostForm("password_confirmation")
	role := c.PostForm("role")

	if name == "" || email == "" || password == "" || role == "" {
		c.Redirect(http.StatusFound, "/admin/users/create?error=Semua+field+wajib+diisi")
		return
	}

	if len(password) < 6 {
		c.Redirect(http.StatusFound, "/admin/users/create?error=Password+minimal+6+karakter")
		return
	}

	if password != passwordConfirm {
		c.Redirect(http.StatusFound, "/admin/users/create?error=Konfirmasi+password+tidak+cocok")
		return
	}

	if role != "admin" && role != "staff" && role != "client" {
		c.Redirect(http.StatusFound, "/admin/users/create?error=Role+tidak+valid")
		return
	}

	db := database.GetDB()

	var existing models.User
	if err := db.Where("email = ?", email).First(&existing).Error; err == nil {
		c.Redirect(http.StatusFound, "/admin/users/create?error=Email+sudah+terdaftar")
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		c.Redirect(http.StatusFound, "/admin/users/create?error=Gagal+hash+password")
		return
	}

	newUser := models.User{
		Name:     name,
		Email:    email,
		Password: string(hashedPassword),
		Role:     role,
		IsActive: true,
	}

	if err := db.Create(&newUser).Error; err != nil {
		c.Redirect(http.StatusFound, "/admin/users/create?error=Gagal+menyimpan+user")
		return
	}

	c.Redirect(http.StatusFound, "/admin/users")
}

// ============================================
// UserEdit — GET /admin/users/:id/edit
// ============================================
func UserEdit(c *gin.Context) {
	setting := services.GetSiteSetting()
	user := c.MustGet("user").(*models.User)
	errorMsg := c.Query("error")

	id := c.Param("id")
	var targetUser models.User
	db := database.GetDB()
	if err := db.First(&targetUser, id).Error; err != nil {
		c.Redirect(http.StatusFound, "/admin/users?error=User+tidak+ditemukan")
		return
	}

	var buf bytes.Buffer
	err := layouts.AdminLayout(
		*user,
		setting,
		"admin.users",
		users.EditContent(targetUser, errorMsg),
	).Render(c.Request.Context(), &buf)

	if err != nil {
		c.String(http.StatusInternalServerError, "Render error: %v", err)
		return
	}

	c.Data(http.StatusOK, "text/html; charset=utf-8", buf.Bytes())
}

// ============================================
// UserUpdate — POST /admin/users/:id/edit
// ============================================
func UserUpdate(c *gin.Context) {
	id := c.Param("id")
	name := strings.TrimSpace(c.PostForm("name"))
	email := strings.ToLower(strings.TrimSpace(c.PostForm("email")))
	password := c.PostForm("password")
	passwordConfirm := c.PostForm("password_confirmation")
	role := c.PostForm("role")

	db := database.GetDB()

	var targetUser models.User
	if err := db.First(&targetUser, id).Error; err != nil {
		c.Redirect(http.StatusFound, "/admin/users?error=User+tidak+ditemukan")
		return
	}

	editURL := "/admin/users/" + id + "/edit"

	if name == "" || email == "" || role == "" {
		c.Redirect(http.StatusFound, editURL+"?error=Semua+field+wajib+diisi")
		return
	}

	var existing models.User
	if err := db.Where("email = ? AND id != ?", email, id).First(&existing).Error; err == nil {
		c.Redirect(http.StatusFound, editURL+"?error=Email+sudah+dipakai+user+lain")
		return
	}

	targetUser.Name = name
	targetUser.Email = email
	targetUser.Role = role

	if password != "" {
		if len(password) < 6 {
			c.Redirect(http.StatusFound, editURL+"?error=Password+minimal+6+karakter")
			return
		}
		if password != passwordConfirm {
			c.Redirect(http.StatusFound, editURL+"?error=Konfirmasi+password+tidak+cocok")
			return
		}
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			c.Redirect(http.StatusFound, editURL+"?error=Gagal+hash+password")
			return
		}
		targetUser.Password = string(hashedPassword)
	}

	if err := db.Save(&targetUser).Error; err != nil {
		c.Redirect(http.StatusFound, editURL+"?error=Gagal+menyimpan+perubahan")
		return
	}

	c.Redirect(http.StatusFound, "/admin/users")
}

// ============================================
// UserDelete — POST /admin/users/:id/delete
// ============================================
func UserDelete(c *gin.Context) {
	id := c.Param("id")
	currentUser := c.MustGet("user").(*models.User)

	currentID := strconv.FormatUint(uint64(currentUser.ID), 10)
	if id == currentID {
		c.Redirect(http.StatusFound, "/admin/users?error=Tidak+bisa+hapus+akun+sendiri")
		return
	}

	db := database.GetDB()
	if err := db.Delete(&models.User{}, id).Error; err != nil {
		c.Redirect(http.StatusFound, "/admin/users?error=Gagal+menghapus+user")
		return
	}

	c.Redirect(http.StatusFound, "/admin/users")
}
