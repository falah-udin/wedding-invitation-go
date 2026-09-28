package services

import (
	"net/http"
	"os"

	"github.com/gorilla/sessions"
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