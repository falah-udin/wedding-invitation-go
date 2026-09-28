package routes

import (
	"github.com/gin-gonic/gin"

	"wedding-invitation-go/internal/handlers"
)

// Register — daftarkan semua route ke Gin engine
func Register(r *gin.Engine) {
	// ============================================
	// PUBLIC ROUTES
	// ============================================
	r.GET("/", handlers.ShowHome)
	r.GET("/health", handlers.ShowHealth)

	// ============================================
	// AUTH ROUTES
	// ============================================
	r.GET("/login", handlers.ShowLogin)
	r.POST("/login", handlers.HandleLogin)
	r.GET("/register", handlers.ShowRegister)
	r.POST("/register", handlers.HandleRegister)
	r.POST("/logout", handlers.HandleLogout)

	// ============================================
	// NANTI: ADMIN, CLIENT, STAFF ROUTES
	// ============================================
	// r.GET("/admin/dashboard", handlers.AdminDashboard)
	// r.GET("/client/dashboard", handlers.ClientDashboard)
	// ...
}