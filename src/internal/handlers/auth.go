package handlers

import (
	"bytes"
	"crypto/rand"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"

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

// ============================================
// GoogleLogin — GET /auth/google
// Redirect ke Google untuk login
// ============================================
func GoogleLogin(c *gin.Context) {
	config := services.GetGoogleOAuthConfig(c.Request)

	if config.ClientID == "" || config.ClientSecret == "" {
		c.Redirect(http.StatusFound, "/login?error=Google+OAuth+belum+dikonfigurasi")
		return
	}

	// State untuk CSRF protection
	state := "wedding-oauth-state"

	url := config.AuthCodeURL(state, oauth2.AccessTypeOffline)
	c.Redirect(http.StatusFound, url)
}

// ============================================
// GoogleCallback — GET /auth/google/callback
// Handle callback dari Google setelah user login
// ============================================
func GoogleCallback(c *gin.Context) {
	// 1. Cek error dari Google
	if errMsg := c.Query("error"); errMsg != "" {
		c.Redirect(http.StatusFound, "/login?error=Google+login+dibatalkan")
		return
	}

	// 2. Ambil code dari query
	code := c.Query("code")
	if code == "" {
		c.Redirect(http.StatusFound, "/login?error=Code+tidak+valid")
		return
	}

	// 3. Tuker code → user info
	userInfo, err := services.GetGoogleUserInfo(c.Request.Context(), c.Request, code)
	if err != nil {
		services.LogError("GoogleCallback", err, nil)
		c.Redirect(http.StatusFound, "/login?error=Gagal+login+Google")
		return
	}

	// 4. Cari user berdasarkan email atau google_id
	db := database.GetDB()
	var user models.User

	err = db.Where("email = ? OR google_id = ?", userInfo.Email, userInfo.ID).
		First(&user).Error

	if err != nil {
		// 5. User belum ada → daftar baru sebagai client
		// Generate password random (karena kolom password NOT NULL)
		randomBytes := make([]byte, 32)
		rand.Read(randomBytes)
		randomPassword, _ := bcrypt.GenerateFromPassword(randomBytes, bcrypt.DefaultCost)

		user = models.User{
			Name:     userInfo.Name,
			Email:    userInfo.Email,
			Password: string(randomPassword),
			Role:     "client",
			IsActive: true,
			GoogleID: &userInfo.ID,
		}
		if userInfo.Picture != "" {
			user.Avatar = &userInfo.Picture
		}

		if err := db.Create(&user).Error; err != nil {
			services.LogError("GoogleCallback.CreateUser", err, map[string]interface{}{
				"email": userInfo.Email,
			})
			c.Redirect(http.StatusFound, "/login?error=Gagal+membuat+akun")
			return
		}

		services.LogSuccess("GoogleCallback", "User baru via Google: "+userInfo.Email)
	} else {
		// 6. User sudah ada → update google_id & avatar (kalau belum ada)
		if user.GoogleID == nil {
			user.GoogleID = &userInfo.ID
		}
		if userInfo.Picture != "" && (user.Avatar == nil || *user.Avatar == "") {
			user.Avatar = &userInfo.Picture
		}
		db.Save(&user)
	}

	// 7. Cek user aktif
	if !user.IsActive {
		c.Redirect(http.StatusFound, "/login?error=Akun+Anda+tidak+aktif")
		return
	}

	// 8. Set session (sama seperti login biasa)
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

	services.LogSuccess("GoogleCallback", "Login sukses via Google: "+user.Email)

	// 9. Redirect berdasarkan role
	switch user.Role {
	case "admin":
		c.Redirect(http.StatusFound, "/admin/dashboard")
	case "staff":
		c.Redirect(http.StatusFound, "/staff/dashboard")
	default:
		c.Redirect(http.StatusFound, "/client/dashboard")
	}
}
