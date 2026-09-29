package services

import (
	"context"
	"net/http"
	"os"
	"strings"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	googleoauth2 "google.golang.org/api/oauth2/v2"
	"google.golang.org/api/option"
)

// getBaseURL — bangun base URL dari request (scheme + host)
// Prioritas:
//  1. Header X-Forwarded-Proto & X-Forwarded-Host (kalau di belakang proxy)
//  2. Request.TLS (kalau https langsung)
//  3. Request.Host (fallback)
func getBaseURL(r *http.Request) string {
	// PRIORITAS 1: APP_URL dari env (paling pasti, terutama di balik tunnel)
	if appURL := os.Getenv("APP_URL"); appURL != "" {
		return strings.TrimSuffix(appURL, "/")
	}

	// PRIORITAS 2: header X-Forwarded-Proto (kalau di belakang proxy)
	scheme := "http"
	if proto := r.Header.Get("X-Forwarded-Proto"); proto != "" {
		scheme = proto
	} else if r.TLS != nil {
		scheme = "https"
	}

	// PRIORITAS 3: host dari header atau request
	host := r.Host
	if forwardedHost := r.Header.Get("X-Forwarded-Host"); forwardedHost != "" {
		host = forwardedHost
	}

	return scheme + "://" + host
}

// GetGoogleOAuthConfig — ambil config OAuth dengan redirect URI dinamis
//
// Cara kerja GOOGLE_REDIRECT_URI:
//   - Kalau di-set (mis. "http://192.168.0.102:8086/auth/google/callback") → pakai itu
//   - Kalau kosong atau "auto" → bangun dinamis dari request.Host
//   - Kalau hanya path (mis. "/auth/google/callback") → bangun dari request.Host
func GetGoogleOAuthConfig(r *http.Request) *oauth2.Config {
	redirectURI := os.Getenv("GOOGLE_REDIRECT_URI")

	// Kalau kosong atau "auto" → bangun dinamis
	if redirectURI == "" || redirectURI == "auto" {
		redirectURI = getBaseURL(r) + "/auth/google/callback"
	} else if strings.HasPrefix(redirectURI, "/") {
		// Kalau hanya path → tambahkan base URL
		redirectURI = getBaseURL(r) + redirectURI
	}

	return &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  redirectURI,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
}

// GoogleUserInfo — info user dari Google
type GoogleUserInfo struct {
	ID      string
	Email   string
	Name    string
	Picture string
}

// GetGoogleUserInfo — tuker code → user info
func GetGoogleUserInfo(ctx context.Context, r *http.Request, code string) (*GoogleUserInfo, error) {
	config := GetGoogleOAuthConfig(r)

	// 1. Tuker code → token
	token, err := config.Exchange(ctx, code)
	if err != nil {
		return nil, err
	}

	// 2. Buat client oauth2 dengan token
	client := config.Client(ctx, token)

	// 3. Panggil Google OAuth2 API
	svc, err := googleoauth2.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, err
	}

	// 4. Ambil user info (email, name, picture, id)
	userInfo, err := svc.Userinfo.Get().Do()
	if err != nil {
		return nil, err
	}

	return &GoogleUserInfo{
		ID:      userInfo.Id,
		Email:   userInfo.Email,
		Name:    userInfo.Name,
		Picture: userInfo.Picture,
	}, nil
}
