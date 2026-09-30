package client

import (
	"bytes"
	"net/http"

	"github.com/gin-gonic/gin"

	"wedding-invitation-go/internal/services"
	"wedding-invitation-go/views/client/support"
	"wedding-invitation-go/views/layouts"
)

// SupportIndex — GET /client/support
// Halaman "Dukung Developer" — semua info dari setting
func SupportIndex(c *gin.Context) {
	setting := services.GetSiteSetting()
	user := services.GetCurrentUser(c.Request)
	if user == nil {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	data := support.SupportData{
		User:              *user,
		Setting:           setting,
		// QRIS & Developer
		QrisImage:         setting.QrisImage,
		SupportName:       setting.DeveloperName,
		SupportEmail:      setting.DeveloperEmail,
		SupportUrl:        setting.DeveloperUrl,
		SupportWhatsapp:   setting.ContactWhatsapp,
		// Site info
		SiteName:          setting.SiteName,
		SiteTagline:       setting.SiteTagline,
		FooterDescription: setting.FooterDescription,
		// Contact
		ContactEmail:      setting.ContactEmail,
		ContactPhone:      setting.ContactPhone,
		ContactAddress:    setting.ContactAddress,
		// Social
		SocialInstagram:   setting.SocialInstagram,
		SocialFacebook:    setting.SocialFacebook,
		SocialTiktok:      setting.SocialTiktok,
		SocialYoutube:     setting.SocialYoutube,
		SocialTwitter:     setting.SocialTwitter,
	}

	var buf bytes.Buffer
	err := layouts.ClientLayout(*user, setting, "client.support", support.SupportPage(data)).Render(c.Request.Context(), &buf)
	if err != nil {
		c.String(http.StatusInternalServerError, "Render error: %v", err)
		return
	}
	c.Data(http.StatusOK, "text/html; charset=utf-8", buf.Bytes())
}
