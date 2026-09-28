package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"wedding-invitation-go/internal/database"
	"wedding-invitation-go/internal/models"
	"wedding-invitation-go/internal/services"
)

// RequireAuth — middleware yang memastikan user sudah login
// Kalau belum login, redirect ke /login
func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		store := services.GetStore()
		session, err := store.Get(c.Request, "wedding_session")
		if err != nil {
			c.Redirect(http.StatusFound, "/login?error=Sesi+tidak+valid")
			c.Abort()
			return
		}

		userID, ok := session.Values["user_id"].(uint)
		if !ok || userID == 0 {
			c.Redirect(http.StatusFound, "/login?error=Silakan+login+terlebih+dahulu")
			c.Abort()
			return
		}

		// Ambil user dari DB
		var user models.User
		db := database.GetDB()
		if err := db.First(&user, userID).Error; err != nil {
			c.Redirect(http.StatusFound, "/login?error=User+tidak+ditemukan")
			c.Abort()
			return
		}

		if !user.IsActive {
			c.Redirect(http.StatusFound, "/login?error=Akun+Anda+tidak+aktif")
			c.Abort()
			return
		}

		// Simpan user di context supaya handler bisa akses
		c.Set("user", &user)
		c.Set("user_id", user.ID)
		c.Set("user_role", user.Role)

		c.Next()
	}
}

// RequireRole — middleware yang memastikan user punya role tertentu
// Contoh: RequireRole("admin") atau RequireRole("admin", "staff")
func RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("user_role")
		if !exists {
			c.Redirect(http.StatusFound, "/login?error=Silakan+login+terlebih+dahulu")
			c.Abort()
			return
		}

		roleStr, ok := userRole.(string)
		if !ok {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		// Cek apakah role user ada di daftar yang diizinkan
		allowed := false
		for _, r := range roles {
			if r == roleStr {
				allowed = true
				break
			}
		}

		if !allowed {
			c.String(http.StatusForbidden, "Akses ditolak. Halaman ini hanya untuk: %v", roles)
			c.Abort()
			return
		}

		c.Next()
	}
}
