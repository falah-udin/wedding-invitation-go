package invitation

import (
	"bytes"
	"net/http"

	templpkg "github.com/a-h/templ"
	"github.com/gin-gonic/gin"

	"wedding-invitation-go/internal/models"
	"wedding-invitation-go/views/layouts"
)

// renderWithLayout — pilih layout sesuai role user
//   - client → ClientLayout (hijau)
//   - admin/staff → AdminLayout (ungu)
func renderWithLayout(c *gin.Context, user *models.User, setting models.Config, content templpkg.Component) {
	var buf bytes.Buffer
	var err error

	switch user.Role {
	case "client":
		err = layouts.ClientLayout(*user, setting, "client.create", content).Render(c.Request.Context(), &buf)
	default: // admin, staff
		err = layouts.AdminLayout(*user, setting, "invitation.create", content).Render(c.Request.Context(), &buf)
	}

	if err != nil {
		c.String(http.StatusInternalServerError, "Render error: %v", err)
		return
	}
	c.Data(http.StatusOK, "text/html; charset=utf-8", buf.Bytes())
}

// isClient — cek apakah user client
func isClient(user *models.User) bool {
	return user.Role == "client"
}
