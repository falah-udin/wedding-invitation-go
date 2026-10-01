package invitations

import (
	"fmt"
	"io"
	"net/http"

	"wedding-invitation-go/internal/invitation"
	"wedding-invitation-go/internal/models"
	"wedding-invitation-go/views/invitations/botanical_garden"
	"wedding-invitation-go/views/invitations/dark_technology"
	"wedding-invitation-go/views/invitations/elegant_gold"
	"wedding-invitation-go/views/invitations/modern_minimalist"
	"wedding-invitation-go/views/invitations/muslim_elegan"
	"wedding-invitation-go/views/invitations/rustic_wood"
	"wedding-invitation-go/views/invitations/traditional_java"
)

func RenderTemplate(
	w io.Writer,
	r *http.Request,
	project models.Project,
	data *invitation.TemplateData,
	guestName string,
) error {
	folder := "rustic_wood"

	if project.Template != nil {
		folder = normalizeFolder(project.Template.Folder)
	}

	ctx := r.Context()

	switch folder {
	case "rustic_wood":
		return rustic_wood.RusticWood(ctx, w, project, data, guestName)
	case "muslim_elegan":
		return muslim_elegan.MuslimElegan(ctx, w, project, data, guestName)
	case "elegant_gold":
		return elegant_gold.ElegantGold(ctx, w, project, data, guestName)
	case "modern_minimalist":
		return modern_minimalist.ModernMinimalist(ctx, w, project, data, guestName)
	case "traditional_java":
		return traditional_java.TraditionalJava(ctx, w, project, data, guestName)
	case "botanical_garden":
		return botanical_garden.BotanicalGarden(ctx, w, project, data, guestName)
	case "dark_technology":
		return dark_technology.DarkTechnology(ctx, w, project, data, guestName)
	default:
		return fmt.Errorf("template %s tidak dikenal", folder)
	}
}

func normalizeFolder(folder string) string {
	result := ""
	for _, c := range folder {
		if c == '-' {
			result += "_"
		} else {
			result += string(c)
		}
	}
	return result
}
