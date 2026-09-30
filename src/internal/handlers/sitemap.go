package handlers

import (
	"encoding/xml"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"wedding-invitation-go/internal/database"
	"wedding-invitation-go/internal/models"
)

// SitemapURL — 1 entry di sitemap
type SitemapURL struct {
	Loc        string `xml:"loc"`
	LastMod    string `xml:"lastmod,omitempty"`
	ChangeFreq string `xml:"changefreq,omitempty"`
	Priority   string `xml:"priority,omitempty"`
}

// SitemapURLSet — root element
type SitemapURLSet struct {
	XMLName xml.Name     `xml:"urlset"`
	Xmlns   string       `xml:"xmlns,attr"`
	URLs    []SitemapURL `xml:"url"`
}

// Sitemap — GET /sitemap.xml
func Sitemap(c *gin.Context) {
	baseURL := sitemapBaseURL(c)

	urls := []SitemapURL{
		{
			Loc:        baseURL + "/",
			LastMod:    time.Now().Format("2006-01-02"),
			ChangeFreq: "daily",
			Priority:   "1.0",
		},
		{
			Loc:        baseURL + "/login",
			ChangeFreq: "monthly",
			Priority:   "0.5",
		},
		{
			Loc:        baseURL + "/register",
			ChangeFreq: "monthly",
			Priority:   "0.5",
		},
	}

	// Ambil semua project published
	db := database.GetDB()
	var projects []models.Project
	db.Where("status = ?", "published").
		Order("updated_at DESC").
		Limit(1000).
		Find(&projects)

	for _, p := range projects {
		urls = append(urls, SitemapURL{
			Loc:        baseURL + "/invitation/" + p.Slug,
			LastMod:    p.UpdatedAt.Format("2006-01-02"),
			ChangeFreq: "weekly",
			Priority:   "0.8",
		})
	}

	sitemap := SitemapURLSet{
		Xmlns: "http://www.sitemaps.org/schemas/sitemap/0.9",
		URLs:  urls,
	}

	c.Header("Content-Type", "application/xml; charset=utf-8")
	c.Header("Cache-Control", "public, max-age=3600")
	_, _ = c.Writer.WriteString(xml.Header)
	_ = xml.NewEncoder(c.Writer).Encode(sitemap)
}

// Robots — GET /robots.txt
func Robots(c *gin.Context) {
	baseURL := sitemapBaseURL(c)

	content := strings.Join([]string{
		"User-agent: *",
		"Allow: /",
		"",
		"# Halaman internal — jangan diindeks",
		"Disallow: /admin/",
		"Disallow: /client/",
		"Disallow: /staff/",
		"Disallow: /share/",
		"Disallow: /invitation/create/",
		"Disallow: /preview/",
		"Disallow: /login",
		"Disallow: /register",
		"",
		"Sitemap: " + baseURL + "/sitemap.xml",
		"",
	}, "\n")

	c.Header("Content-Type", "text/plain; charset=utf-8")
	c.Header("Cache-Control", "public, max-age=3600")
	c.String(http.StatusOK, content)
}

// sitemapBaseURL — helper ambil base URL (paksa https di production)
func sitemapBaseURL(c *gin.Context) string {
	// Prioritas: APP_URL dari env
	if appURL := os.Getenv("APP_URL"); appURL != "" {
		return strings.TrimSuffix(appURL, "/")
	}

	// Fallback: dari request
	scheme := "http"
	if c.Request.Header.Get("X-Forwarded-Proto") == "https" || c.Request.TLS != nil {
		scheme = "https"
	}
	host := c.Request.Host
	return scheme + "://" + host
}
