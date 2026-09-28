package routes

import (
	"github.com/gin-gonic/gin"

	"wedding-invitation-go/internal/handlers"
	adminHandlers "wedding-invitation-go/internal/handlers/admin"
	"wedding-invitation-go/internal/middleware"
)

// Register — daftarkan semua route
func Register(r *gin.Engine) {
	// PUBLIC ROUTES
	r.GET("/", handlers.ShowHome)
	r.GET("/health", handlers.ShowHealth)

	// AUTH ROUTES
	r.GET("/login", handlers.ShowLogin)
	r.POST("/login", handlers.HandleLogin)
	r.GET("/register", handlers.ShowRegister)
	r.POST("/register", handlers.HandleRegister)
	r.POST("/logout", handlers.HandleLogout)

	// ADMIN ROUTES
	admin := r.Group("/admin")
	admin.Use(middleware.RequireAuth())
	admin.Use(middleware.RequireRole("admin"))
	{
		admin.GET("/dashboard", adminHandlers.DashboardHandler)

		// USER CRUD
		admin.GET("/users", adminHandlers.UserIndex)
		admin.GET("/users/create", adminHandlers.UserCreate)
		admin.POST("/users/create", adminHandlers.UserStore)
		admin.GET("/users/:id/edit", adminHandlers.UserEdit)
		admin.POST("/users/:id/edit", adminHandlers.UserUpdate)
		admin.POST("/users/:id/delete", adminHandlers.UserDelete)
	}
}
