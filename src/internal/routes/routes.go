package routes

import (
	"github.com/gin-gonic/gin"

	"wedding-invitation-go/internal/handlers"
	adminHandlers "wedding-invitation-go/internal/handlers/admin"
	clientHandlers "wedding-invitation-go/internal/handlers/client"
	invitationHandlers "wedding-invitation-go/internal/handlers/invitation"
	staffHandlers "wedding-invitation-go/internal/handlers/staff"
	"wedding-invitation-go/internal/middleware"
)

func Register(r *gin.Engine) {
	// PUBLIC
	r.GET("/", handlers.ShowHome)
	r.GET("/health", handlers.ShowHealth)
	r.GET("/preview/template/:slug", handlers.PreviewTemplate)

	// PWA
	r.GET("/manifest.webmanifest", handlers.Manifest)
	r.GET("/sw.js", handlers.ServiceWorker)

	// INVITATION PUBLIC (tamu akses)
	r.GET("/invitation/:slug", invitationHandlers.ShowInvitation)
	r.POST("/invitation/:slug/rsvp", invitationHandlers.RsvpSubmit)
	r.GET("/invitation/:slug/rsvp-list", invitationHandlers.RsvpList)

	// AUTH
	r.GET("/login", handlers.ShowLogin)
	r.POST("/login", handlers.HandleLogin)
	r.GET("/register", handlers.ShowRegister)
	r.POST("/register", handlers.HandleRegister)
	r.POST("/logout", handlers.HandleLogout)

	// GOOGLE OAUTH
	r.GET("/auth/google", handlers.GoogleLogin)
	r.GET("/auth/google/callback", handlers.GoogleCallback)

	// STAFF AREA
	staff := r.Group("/staff")
	staff.Use(middleware.RequireAuth())
	staff.Use(middleware.RequireRole("staff"))
	{
		staff.GET("/dashboard", staffHandlers.Dashboard)
		staff.GET("/share", staffHandlers.ShareIndex)
		staff.GET("/projects", staffHandlers.ProjectIndex)
		staff.GET("/projects/:id", staffHandlers.ProjectShow)
		staff.GET("/rsvp", staffHandlers.RsvpListIndex)
		staff.GET("/projects/:id/rsvp", staffHandlers.RsvpProjectIndex)
		staff.DELETE("/projects/:id/rsvp/:rsvpId", staffHandlers.RsvpDelete)
		staff.POST("/projects/:id/rsvp/bulk-delete", staffHandlers.RsvpBulkDelete)
		staff.GET("/projects/:id/rsvp/export", staffHandlers.RsvpExport)
		staff.GET("/profile/change-password", staffHandlers.ChangePasswordShow)
		staff.POST("/profile/change-password", staffHandlers.ChangePasswordUpdate)
	}

	// CLIENT AREA
	client := r.Group("/client")
	client.Use(middleware.RequireAuth())
	client.Use(middleware.RequireRole("client"))
	{
		client.GET("/dashboard", clientHandlers.Dashboard)
		client.GET("/wedding", clientHandlers.WeddingIndex)
		client.DELETE("/wedding/:id", clientHandlers.WeddingDestroy)

		// Reservasi
		client.GET("/reservations", clientHandlers.ReservationIndex)
		client.GET("/reservations/:id", clientHandlers.ReservationShow)
		client.DELETE("/reservations/:id/rsvp/:rsvpId", clientHandlers.ReservationDelete)

		// Profile
		client.GET("/profile/change-password", clientHandlers.ChangePasswordShow)
		client.POST("/profile/change-password", clientHandlers.ChangePasswordUpdate)
	}
		
	// ADMIN
	admin := r.Group("/admin")
	admin.Use(middleware.RequireAuth())
	admin.Use(middleware.RequireRole("admin"))
	{
		admin.GET("/dashboard", adminHandlers.DashboardHandler)

		// User
		admin.GET("/users", adminHandlers.UserIndex)
		admin.GET("/users/create", adminHandlers.UserCreate)
		admin.POST("/users/create", adminHandlers.UserStore)
		admin.GET("/users/:id/edit", adminHandlers.UserEdit)
		admin.POST("/users/:id/edit", adminHandlers.UserUpdate)
		admin.POST("/users/:id/delete", adminHandlers.UserDelete)

		// Template
		admin.GET("/templates", adminHandlers.TemplateIndex)
		admin.GET("/templates/create", adminHandlers.TemplateCreate)
		admin.POST("/templates/create", adminHandlers.TemplateStore)
		admin.GET("/templates/:id/edit", adminHandlers.TemplateEdit)
		admin.POST("/templates/:id/edit", adminHandlers.TemplateUpdate)
		admin.GET("/templates/:id/toggle", adminHandlers.TemplateToggle)
		admin.POST("/templates/:id/delete", adminHandlers.TemplateDelete)

		// Music
		admin.GET("/music", adminHandlers.MusicIndex)
		admin.GET("/music/create", adminHandlers.MusicCreate)
		admin.POST("/music/create", adminHandlers.MusicStore)
		admin.GET("/music/:id/toggle", adminHandlers.MusicToggle)
		admin.POST("/music/:id/delete", adminHandlers.MusicDelete)

		// Projects
		admin.GET("/projects", adminHandlers.ProjectIndex)
		admin.GET("/projects/:id", adminHandlers.ProjectShow)
		admin.POST("/projects/:id/delete", adminHandlers.ProjectDelete)

		// RSVP
		admin.GET("/projects/:id/rsvp", adminHandlers.RsvpIndex)
		admin.GET("/projects/:id/rsvp/export", adminHandlers.RsvpExport)
		admin.DELETE("/projects/rsvp/:rsvpId", adminHandlers.RsvpDelete)
		admin.POST("/projects/rsvp/bulk-delete", adminHandlers.RsvpBulkDelete)
		
		// Settings
		admin.GET("/settings", adminHandlers.SettingIndex)
		admin.POST("/settings", adminHandlers.SettingUpdate)
		admin.POST("/settings/delete-favicon", adminHandlers.SettingDeleteFavicon)
		admin.POST("/settings/delete-qris", adminHandlers.SettingDeleteQris)
		admin.POST("/settings/delete-pwa-icon", adminHandlers.SettingDeletePwaIcon)
	}

	// INVITATION WIZARD (auth required — admin, staff, client semua boleh)
	wizard := r.Group("/invitation/create")
	wizard.Use(middleware.RequireAuth())
	{
		wizard.GET("/select-client", invitationHandlers.SelectClient)
		wizard.POST("/select-client", invitationHandlers.StoreClient)
		wizard.GET("/edit/:id", invitationHandlers.EditStart)      
		wizard.GET("/general", invitationHandlers.General)
		wizard.POST("/general", invitationHandlers.GeneralSave)
		wizard.GET("/template", invitationHandlers.TemplateList)
		wizard.POST("/template", invitationHandlers.TemplateSave)
		wizard.GET("/preview", invitationHandlers.Preview)
		wizard.POST("/publish", invitationHandlers.Publish)
		wizard.POST("/save-music", invitationHandlers.SaveMusic)
	}

	// ============================================
	// SHARE INVITATION (auth required)
	// ============================================
	share := r.Group("/share")
	share.Use(middleware.RequireAuth())
	{
		share.GET("/invitation/:id", handlers.ShareInvitation)
		share.POST("/invitation/:id/guest", handlers.StoreGuest)
		share.POST("/invitation/:id/guests/bulk-update", handlers.BulkUpdateGuests)
		share.PUT("/guest/:guestId", handlers.UpdateGuest)
		share.DELETE("/guest/:guestId", handlers.DestroyGuest)
	}

}