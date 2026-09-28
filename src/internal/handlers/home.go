package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ShowHome — GET /
// Sementara redirect ke /login
// Nanti: render halaman home lengkap dengan template list
func ShowHome(c *gin.Context) {
	c.Redirect(http.StatusFound, "/login")
}

// ShowHealth — GET /health
// Health check endpoint untuk monitoring
func ShowHealth(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"app":     "Wedding Invitation Go",
		"version": "1.0.0",
	})
}