package services

import (
	"net/http"
	"os"

	"github.com/gorilla/sessions"

	"wedding-invitation-go/internal/database"
	"wedding-invitation-go/internal/models"
)

var store *sessions.CookieStore

// InitSession — inisialisasi session store
func InitSession() {
	secret := os.Getenv("SESSION_SECRET")
	if secret == "" {
		secret = "change-this-secret-in-production-min-32-chars"
	}

	store = sessions.NewCookieStore([]byte(secret))
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7, // 7 hari
		HttpOnly: true,
		Secure:   false, // set true kalau sudah HTTPS
		SameSite: http.SameSiteLaxMode,
	}
}

// GetStore — ambil session store
func GetStore() *sessions.CookieStore {
	return store
}

// ============================================
// CURRENT USER HELPERS
// ============================================

// GetCurrentUser — ambil user yang sedang login dari session
func GetCurrentUser(r *http.Request) *models.User {
	store := GetStore()
	session, err := store.Get(r, "wedding_session")
	if err != nil {
		return nil
	}

	userID, ok := session.Values["user_id"]
	if !ok {
		return nil
	}

	var id uint
	switch v := userID.(type) {
	case uint:
		id = v
	case int:
		id = uint(v)
	case int64:
		id = uint(v)
	case float64:
		id = uint(v)
	default:
		return nil
	}

	if id == 0 {
		return nil
	}

	var user models.User
	db := database.GetDB()
	if err := db.First(&user, id).Error; err != nil {
		return nil
	}

	if !user.IsActive {
		return nil
	}

	return &user
}

// IsLoggedIn — cek apakah user login
func IsLoggedIn(r *http.Request) bool {
	return GetCurrentUser(r) != nil
}

// GetUserRole — ambil role user yang login
func GetUserRole(r *http.Request) string {
	user := GetCurrentUser(r)
	if user == nil {
		return ""
	}
	return user.Role
}

// GetUserName — ambil nama user yang login
func GetUserName(r *http.Request) string {
	user := GetCurrentUser(r)
	if user == nil {
		return ""
	}
	return user.Name
}