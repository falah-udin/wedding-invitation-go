package routes

import (
	"github.com/gin-gonic/gin"

	"wedding-invitation-go/internal/handlers"
	adminHandlers "wedding-invitation-go/internal/handlers/admin"
	"wedding-invitation-go/internal/middleware"
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
	// ADMIN ROUTES (Require Auth + Role Admin)
	// ============================================
	admin := r.Group("/admin")
	admin.Use(middleware.RequireAuth())
	admin.Use(middleware.RequireRole("admin"))
	{
		admin.GET("/dashboard", adminHandlers.DashboardHandler)
	}
}
