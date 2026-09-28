package handlers

import (
	"bytes"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"wedding-invitation-go/internal/database"
	"wedding-invitation-go/internal/models"
	"wedding-invitation-go/internal/services"
	"wedding-invitation-go/views/auth"
	"wedding-invitation-go/views/layouts"
)

// ============================================
// ShowLogin — GET /login
// ============================================
func ShowLogin(c *gin.Context) {
	setting := services.GetSiteSetting()

	errorMsg := c.Query("error")
	successMsg := c.Query("success")
	emailValue := c.Query("email")

	var buf bytes.Buffer
	err := layouts.AuthLayout(
		setting,
		"Login",
		auth.LoginContent(setting, errorMsg, successMsg, emailValue),
	).Render(c.Request.Context(), &buf)

	if err != nil {
		c.String(http.StatusInternalServerError, "Render error: %v", err)
		return
	}

	c.Data(http.StatusOK, "text/html; charset=utf-8", buf.Bytes())
}

// ============================================
// HandleLogin — POST /login
// Validasi email & password ke database
// ============================================
func HandleLogin(c *gin.Context) {
	email := strings.ToLower(strings.TrimSpace(c.PostForm("email")))
	password := c.PostForm("password")

	// Validasi input
	if email == "" || password == "" {
		c.Redirect(http.StatusFound, "/login?error=Email+dan+password+wajib+diisi")
		return
	}

	// Cari user di database
	var user models.User
	db := database.GetDB()
	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		c.Redirect(http.StatusFound, "/login?error=Email+tidak+terdaftar&email="+email)
		return
	}

	// Cek apakah user aktif
	if !user.IsActive {
		c.Redirect(http.StatusFound, "/login?error=Akun+Anda+tidak+aktif.+Hubungi+administrator&email="+email)
		return
	}

	// Verifikasi password dengan bcrypt
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		c.Redirect(http.StatusFound, "/login?error=Password+salah.+Silakan+coba+lagi&email="+email)
		return
	}

	// Password benar → simpan user di session
	store := services.GetStore()
	session, err := store.Get(c.Request, "wedding_session")
	if err != nil {
		c.Redirect(http.StatusFound, "/login?error=Gagal+membuat+sesi")
		return
	}

	session.Values["user_id"] = user.ID
	session.Values["user_role"] = user.Role
	session.Values["user_name"] = user.Name
	session.Options.MaxAge = 86400 * 7 // 7 hari

	if err := session.Save(c.Request, c.Writer); err != nil {
		c.Redirect(http.StatusFound, "/login?error=Gagal+menyimpan+sesi")
		return
	}

	// Redirect berdasarkan role
	switch user.Role {
	case "admin":
		c.Redirect(http.StatusFound, "/admin/dashboard")
	case "staff":
		c.Redirect(http.StatusFound, "/staff/dashboard")
	default:
		c.Redirect(http.StatusFound, "/client/dashboard")
	}
}

// ============================================
// ShowRegister — GET /register
// ============================================
func ShowRegister(c *gin.Context) {
	setting := services.GetSiteSetting()

	errorMsg := c.Query("error")
	successMsg := c.Query("success")

	var buf bytes.Buffer
	err := layouts.AuthLayout(
		setting,
		"Daftar",
		auth.RegisterContent(setting, errorMsg, successMsg),
	).Render(c.Request.Context(), &buf)

	if err != nil {
		c.String(http.StatusInternalServerError, "Render error: %v", err)
		return
	}

	c.Data(http.StatusOK, "text/html; charset=utf-8", buf.Bytes())
}

// ============================================
// HandleRegister — POST /register (placeholder)
// ============================================
func HandleRegister(c *gin.Context) {
	c.Redirect(http.StatusFound, "/login?error=Register+belum+diimplementasikan")
}

// ============================================
// HandleLogout — POST /logout
// Hapus session & redirect ke login
// ============================================
func HandleLogout(c *gin.Context) {
	store := services.GetStore()
	session, err := store.Get(c.Request, "wedding_session")
	if err == nil {
		session.Options.MaxAge = -1 // hapus cookie
		session.Save(c.Request, c.Writer)
	}

	c.Redirect(http.StatusFound, "/login?success=Berhasil+logout")
}
