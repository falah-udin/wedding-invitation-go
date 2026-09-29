package client

import (
	"bytes"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"wedding-invitation-go/internal/database"
	"wedding-invitation-go/internal/models"
	"wedding-invitation-go/internal/services"
	"wedding-invitation-go/views/client/profile"
	"wedding-invitation-go/views/layouts"
)

// ChangePasswordShow — GET /client/profile/change-password
func ChangePasswordShow(c *gin.Context) {
	setting := services.GetSiteSetting()
	user := services.GetCurrentUser(c.Request)
	if user == nil {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	data := profile.ChangePasswordData{
		User:         *user,
		Setting:      setting,
		Success:      c.Query("success"),
		Error:        c.Query("error"),
		IsGoogleUser: user.GoogleID != nil && *user.GoogleID != "",
	}

	var buf bytes.Buffer
	err := layouts.ClientLayout(*user, setting, "client.password", profile.ChangePasswordPage(data)).Render(c.Request.Context(), &buf)
	if err != nil {
		c.String(http.StatusInternalServerError, "Render error: %v", err)
		return
	}
	c.Data(http.StatusOK, "text/html; charset=utf-8", buf.Bytes())
}

// ChangePasswordUpdate — POST /client/profile/change-password
func ChangePasswordUpdate(c *gin.Context) {
	user := services.GetCurrentUser(c.Request)
	if user == nil {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	// Blokir user Google
	if user.GoogleID != nil && *user.GoogleID != "" {
		c.Redirect(http.StatusFound, "/client/profile/change-password?error=Akun+Google+tidak+bisa+ganti+password+di+sini")
		return
	}

	newPassword := c.PostForm("new_password")
	confirmPassword := c.PostForm("new_password_confirmation")

	// Validasi
	if newPassword == "" || confirmPassword == "" {
		c.Redirect(http.StatusFound, "/client/profile/change-password?error=Semua+field+wajib+diisi")
		return
	}

	if len(newPassword) < 6 {
		c.Redirect(http.StatusFound, "/client/profile/change-password?error=Password+baru+minimal+6+karakter")
		return
	}

	if newPassword != confirmPassword {
		c.Redirect(http.StatusFound, "/client/profile/change-password?error=Konfirmasi+password+tidak+cocok")
		return
	}

	// Ambil user dari DB
	db := database.GetDB()
	var dbUser models.User
	if err := db.First(&dbUser, user.ID).Error; err != nil {
		c.Redirect(http.StatusFound, "/client/profile/change-password?error=User+tidak+ditemukan")
		return
	}

	// Cek: password baru tidak boleh sama dengan yang lama
	if bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(newPassword)) == nil {
		c.Redirect(http.StatusFound, "/client/profile/change-password?error=Password+baru+harus+berbeda+dari+yang+lama")
		return
	}

	// Hash password baru
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		services.LogError("ChangePassword.Hash", err, map[string]interface{}{
			"user_id": user.ID,
		})
		c.Redirect(http.StatusFound, "/client/profile/change-password?error=Gagal+memproses+password")
		return
	}

	// Update
	dbUser.Password = string(hashedPassword)
	if err := db.Save(&dbUser).Error; err != nil {
		services.LogError("ChangePassword.Save", err, map[string]interface{}{
			"user_id": user.ID,
		})
		c.Redirect(http.StatusFound, "/client/profile/change-password?error=Gagal+menyimpan+password")
		return
	}

	services.LogSuccess("ChangePassword", "Password diubah untuk user "+dbUser.Email)
	c.Redirect(http.StatusFound, "/client/profile/change-password?success=Password+berhasil+diubah")
}