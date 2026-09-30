package invitation

import (
	"encoding/json"
	"strings"

	"wedding-invitation-go/internal/models"
	"wedding-invitation-go/internal/services"
)

func LoadData(project models.Project, guestName string) *TemplateData {
	generalData := services.ParseJSONMap(project.DataUndangan)
	specificData := services.ParseJSONMap(project.TemplateSpecificData)

	merged := make(map[string]interface{}, len(generalData)+len(specificData))
	for k, v := range generalData {
		merged[k] = v
	}
	for k, v := range specificData {
		merged[k] = v
	}

	d := &TemplateData{
		GroomName: getStr(merged, "groom_name"),
		BrideName: getStr(merged, "bride_name"),

		GroomPhoto: getStr(merged, "groom_photo"),
		BridePhoto: getStr(merged, "bride_photo"),
		HeroImage:  getStr(merged, "hero_image"),

		FatherGroomPhoto: getStr(merged, "father_groom_photo"),
		MotherGroomPhoto: getStr(merged, "mother_groom_photo"),
		FatherBridePhoto: getStr(merged, "father_bride_photo"),
		MotherBridePhoto: getStr(merged, "mother_bride_photo"),

		FatherGroom: getStr(merged, "father_groom"),
		MotherGroom: getStr(merged, "mother_groom"),
		FatherBride: getStr(merged, "father_bride"),
		MotherBride: getStr(merged, "mother_bride"),

		GroomFamilyOrigin: getStr(merged, "groom_family_origin"),
		BrideFamilyOrigin: getStr(merged, "bride_family_origin"),
		KembarMayang:      getStr(merged, "kembar_mayang"),

		GroomInstagram: getStr(merged, "groom_instagram"),
		BrideInstagram: getStr(merged, "bride_instagram"),
		EventInstagram: getStr(merged, "event_instagram"),
		Whatsapp:       getStr(merged, "whatsapp"),
		Youtube:        getStr(merged, "youtube"),

		AkadDate:  getStr(merged, "akad_date"),
		AkadTime:  getStr(merged, "akad_time"),
		AkadVenue: getStr(merged, "akad_venue"),

		ResepsiDate:        getStr(merged, "resepsi_date"),
		ResepsiTime:        getStr(merged, "resepsi_time"),
		ResepsiVenue:       getStr(merged, "resepsi_venue"),
		ResepsiLabel:       getStr(merged, "resepsi_label"),
		ResepsiLabelCustom: getStr(merged, "resepsi_label_custom"),

		MapsURLAkad:    getStr(merged, "maps_url_akad"),
		MapsURLResepsi: getStr(merged, "maps_url_resepsi"),

		ShowDates:        getStrDefault(merged, "show_dates", "both"),
		ShowVenue:        getStrDefault(merged, "show_venue", "both"),
		ShowBankAccounts: getStrDefault(merged, "show_bank_accounts", "no"),

		BankAccounts: parseBankAccounts(merged),
		Gallery:      parseGallery(merged),
		LoveStories:  parseLoveStories(merged),

		Specific: specificData,

		Project:      project,
		GuestName:    guestName,
		IsNamedGuest: guestName != "" && guestName != "Tamu Undangan",
	}

	d.SiteConfig = services.GetSiteSetting()
	d.BankList = models.GetBankList()
	d.ExistingRsvp = loadExistingRsvp(project, guestName)
	d.MusicURL = getMusicURL(project)

	return d
}

func parseBankAccounts(m map[string]interface{}) []BankAccount {
	var result []BankAccount
	v, ok := m["bank_accounts"]
	if !ok || v == nil {
		return result
	}
	if str, ok := v.(string); ok {
		var arr []map[string]interface{}
		if err := json.Unmarshal([]byte(str), &arr); err == nil {
			for _, item := range arr {
				result = append(result, mapToBankAccount(item))
			}
		}
		return result
	}
	if arr, ok := v.([]interface{}); ok {
		for _, item := range arr {
			if mm, ok := item.(map[string]interface{}); ok {
				result = append(result, mapToBankAccount(mm))
			}
		}
	}
	return result
}

func mapToBankAccount(m map[string]interface{}) BankAccount {
	return BankAccount{
		ID:            getStr(m, "id"),
		Type:          getStr(m, "type"),
		BankCode:      getStr(m, "bank_code"),
		BankName:      getStr(m, "bank_name"),
		AccountNumber: getStr(m, "account_number"),
		AccountName:   getStr(m, "account_name"),
		IconType:      getStr(m, "icon_type"),
		IconPath:      getStr(m, "icon_path"),
	}
}

func parseGallery(m map[string]interface{}) []string {
	result := []string{}
	v, ok := m["gallery"]
	if !ok || v == nil {
		return result
	}
	arr, ok := v.([]interface{})
	if !ok {
		return result
	}
	for _, item := range arr {
		if s, ok := item.(string); ok {
			result = append(result, s)
			continue
		}
		if mm, ok := item.(map[string]interface{}); ok {
			if u := getStr(mm, "url"); u != "" {
				result = append(result, u)
			}
		}
	}
	return result
}

func parseLoveStories(m map[string]interface{}) []LoveStory {
	result := []LoveStory{}
	v, ok := m["love_stories"]
	if !ok || v == nil {
		return result
	}
	arr, ok := v.([]interface{})
	if !ok {
		return result
	}
	for _, item := range arr {
		if mm, ok := item.(map[string]interface{}); ok {
			result = append(result, LoveStory{
				Title: getStr(mm, "title"),
				Desc:  getStr(mm, "desc"),
			})
		}
	}
	return result
}

func loadExistingRsvp(project models.Project, guestName string) *ExistingRsvp {
	if guestName == "" || guestName == "Tamu Undangan" {
		return nil
	}
	rsvp := services.GetRsvpByGuest(project.ID, guestName)
	if rsvp == nil {
		return nil
	}
	return &ExistingRsvp{
		Attendance:  rsvp.Attendance,
		TotalGuests: rsvp.TotalGuests,
		Message:     rsvp.GetMessage(),
	}
}

func getMusicURL(project models.Project) string {
	if project.CustomMusic == nil || *project.CustomMusic == "" {
		return ""
	}
	v := *project.CustomMusic
	if len(v) > 0 && v[0] != '/' && v[0] != 'h' {
		return "/storage/" + v
	}
	return v
}

func getStr(m map[string]interface{}, key string) string {
	if v, ok := m[key]; ok {
		if s, ok := v.(string); ok {
			return strings.TrimSpace(s)
		}
	}
	return ""
}

func getStrDefault(m map[string]interface{}, key, def string) string {
	if v := getStr(m, key); v != "" {
		return v
	}
	return def
}
