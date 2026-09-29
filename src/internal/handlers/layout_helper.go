package handlers

import (
	"bytes"
	"net/http"

	templpkg "github.com/a-h/templ"
	"github.com/gin-gonic/gin"

	"wedding-invitation-go/internal/models"
	"wedding-invitation-go/views/layouts"
)

// renderWithRoleLayout — pilih layout sesuai role user
//   - client → ClientLayout (hijau)
//   - admin/staff → AdminLayout (ungu)
func renderWithRoleLayout(c *gin.Context, user *models.User, setting models.Config, activePage string, content templpkg.Component) {
	var buf bytes.Buffer
	var err error

	switch user.Role {
	case "client":
		err = layouts.ClientLayout(*user, setting, activePage, content).Render(c.Request.Context(), &buf)
	case "staff":
		err = layouts.StaffLayout(*user, setting, activePage, content).Render(c.Request.Context(), &buf)
	default:
		err = layouts.AdminLayout(*user, setting, activePage, content).Render(c.Request.Context(), &buf)
	}

	if err != nil {
		c.String(http.StatusInternalServerError, "Render error: %v", err)
		return
	}
	c.Data(http.StatusOK, "text/html; charset=utf-8", buf.Bytes())
}

// isClientRole — helper
func isClientRole(user *models.User) bool {
	return user.Role == "client"
}
