package handlers

import (
	"bytes"
	"net/http"

	"github.com/gin-gonic/gin"

	"wedding-invitation-go/internal/services"
	"wedding-invitation-go/views/auth"
	"wedding-invitation-go/views/layouts"
)

// ShowLogin — GET /login
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

// HandleLogin — POST /login (placeholder)
func HandleLogin(c *gin.Context) {
	email := c.PostForm("email")
	password := c.PostForm("password")

	_ = email
	_ = password

	c.Redirect(http.StatusFound, "/login?error=Backend+login+belum+diimplementasikan&email="+email)
}

// ShowRegister — GET /register (placeholder)
func ShowRegister(c *gin.Context) {
	c.Redirect(http.StatusFound, "/login?success=Register+belum+diimplementasikan")
}

// HandleRegister — POST /register (placeholder)
func HandleRegister(c *gin.Context) {
	c.Redirect(http.StatusFound, "/login?success=Register+belum+diimplementasikan")
}

// HandleLogout — POST /logout (placeholder)
func HandleLogout(c *gin.Context) {
	c.Redirect(http.StatusFound, "/login?success=Berhasil+logout")
}
