package handlers

import (
	"bytes"
	"net/http"

	"github.com/gin-gonic/gin"

	"wedding-invitation-go/internal/services"
	"wedding-invitation-go/views/guest"
)

// ============================================
// ShowHome — GET /
// Render halaman home publik
// ============================================
func ShowHome(c *gin.Context) {
	setting := services.GetSiteSetting()
	templates := services.GetTemplates()

	// Sementara: dummy isLoggedIn & role
	// Nanti: cek session user
	isLoggedIn := false
	userRole := ""

	var buf bytes.Buffer
	err := guest.HomePage(setting, templates, isLoggedIn, userRole).Render(c.Request.Context(), &buf)

	if err != nil {
		c.String(http.StatusInternalServerError, "Render error: %v", err)
		return
	}

	c.Data(http.StatusOK, "text/html; charset=utf-8", buf.Bytes())
}

// ============================================
// ShowHealth — GET /health
// Health check endpoint
// ============================================
func ShowHealth(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"app":     "Wedding Invitation Go",
		"version": "1.0.0",
	})
}
