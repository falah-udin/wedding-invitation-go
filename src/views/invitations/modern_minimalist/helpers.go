package modern_minimalist

import (
	"encoding/json"
	"strconv"
	"time"

	"wedding-invitation-go/internal/models"
)

// ============================================
// CORE HELPERS
// ============================================

// getStr — ambil string dari map dengan default "-"
func getStr(data map[string]interface{}, key string) string {
	if v, ok := data[key]; ok {
		if s, ok := v.(string); ok && s != "" {
			return s
		}
	}
	return "-"
}

// getDefaultPhoto — ambil foto, fallback ke default global
func getDefaultPhoto(data map[string]interface{}, key string) string {
	v := getStr(data, key)
	if v == "" || v == "-" {
		switch key {
		case "groom_photo":
			return "/storage/defaults/groom.jpg"
		case "bride_photo":
			return "/storage/defaults/bride.jpg"
		case "father_groom_photo", "father_bride_photo":
			return "/storage/defaults/father.jpg"
		case "mother_groom_photo", "mother_bride_photo":
			return "/storage/defaults/mother.jpg"
		case "hero_image":
			return ""
		default:
			return "/storage/defaults/photo.jpg"
		}
	}
	return v
}

// hasHeroImage — cek apakah hero image ada
func hasHeroImage(data map[string]interface{}) bool {
	v := getStr(data, "hero_image")
	return v != "" && v != "-"
}

// getMusicUrlFromProject — ambil URL musik dari project
func getMusicUrlFromProject(p models.Project) string {
	if p.CustomMusic == nil || *p.CustomMusic == "" {
		return ""
	}
	return "/storage/" + *p.CustomMusic
}

// formatUint — konversi uint ke string
func formatUint(n uint) string {
	return strconv.FormatUint(uint64(n), 10)
}

// formatInt — konversi int ke string
func formatInt(n int) string {
	return strconv.Itoa(n)
}

// currentYear — tahun sekarang
func currentYear() int {
	return time.Now().Year()
}

// jsString — escape string untuk JS context (untuk data-attribute kita tidak butuh ini,
// tapi disimpan untuk jaga-jaga kalau template lain butuh)
func jsString(s string) string {
	b, err := json.Marshal(s)
	if err != nil {
		return `""`
	}
	return string(b)
}

// ============================================
// DATA INJECTED (dari handler)
// ============================================

// getExistingRsvp — ambil existing RSVP (nil kalau tidak ada)
func getExistingRsvp(data map[string]interface{}) map[string]interface{} {
	if v, ok := data["_existing_rsvp"].(map[string]interface{}); ok {
		return v
	}
	return nil
}

// hasExistingRsvp — cek apakah tamu sudah RSVP
func hasExistingRsvp(data map[string]interface{}) bool {
	return getExistingRsvp(data) != nil
}

// getExistingRsvpField — ambil field dari existing RSVP
func getExistingRsvpField(data map[string]interface{}, key, def string) string {
	rsvp := getExistingRsvp(data)
	if rsvp == nil {
		return def
	}
	if v, ok := rsvp[key]; ok {
		if s, ok := v.(string); ok && s != "" {
			return s
		}
	}
	return def
}

// getExistingRsvpGuests — ambil total_guests dari existing RSVP
func getExistingRsvpGuests(data map[string]interface{}) int {
	rsvp := getExistingRsvp(data)
	if rsvp == nil {
		return 1
	}
	if v, ok := rsvp["total_guests"]; ok {
		switch n := v.(type) {
		case int:
			if n > 0 { return n }
		case int64:
			if n > 0 { return int(n) }
		case float64:
			if n > 0 { return int(n) }
		}
	}
	return 1
}

// getSiteConfig — ambil site config value by key
func getSiteConfig(data map[string]interface{}, key, def string) string {
	cfg, ok := data["_site_config"].(models.Config)
	if !ok {
		return def
	}
	switch key {
	case "site_name": return cfg.SiteName
	case "site_tagline": return cfg.SiteTagline
	case "site_favicon": return cfg.SiteFavicon
	case "footer_description": return cfg.FooterDescription
	case "footer_copyright": return cfg.FooterCopyright
	case "footer_version": return cfg.FooterVersion
	case "contact_email": return cfg.ContactEmail
	case "contact_phone": return cfg.ContactPhone
	case "contact_whatsapp": return cfg.ContactWhatsapp
	case "contact_address": return cfg.ContactAddress
	case "social_instagram": return cfg.SocialInstagram
	case "social_facebook": return cfg.SocialFacebook
	case "social_tiktok": return cfg.SocialTiktok
	case "social_youtube": return cfg.SocialYoutube
	case "social_twitter": return cfg.SocialTwitter
	case "developer_name": return cfg.DeveloperName
	case "developer_email": return cfg.DeveloperEmail
	case "developer_url": return cfg.DeveloperUrl
	}
	return def
}

// ============================================
// BANK / AMPLOP DIGITAL
// ============================================

func getBankList(data map[string]interface{}) []map[string]string {
	if v, ok := data["_bank_list"].([]map[string]string); ok {
		return v
	}
	return nil
}

func findBank(banks []map[string]string, code string) map[string]string {
	for _, b := range banks {
		if b["code"] == code {
			return b
		}
	}
	return nil
}

func getBankAccounts(data map[string]interface{}) []map[string]interface{} {
	v, ok := data["bank_accounts"]
	if !ok || v == nil {
		return nil
	}
	arr, ok := v.([]interface{})
	if !ok {
		return nil
	}
	result := make([]map[string]interface{}, 0, len(arr))
	for _, item := range arr {
		if m, ok := item.(map[string]interface{}); ok {
			result = append(result, m)
		}
	}
	return result
}

func showBankAccounts(data map[string]interface{}) bool {
	return getStr(data, "show_bank_accounts") == "yes"
}

func getBankTypeBadge(account map[string]interface{}) string {
	t, _ := account["type"].(string)
	switch t {
	case "bank": return "BANK"
	case "ewallet": return "E-WALLET"
	case "qris": return "QRIS"
	default: return "LAINNYA"
	}
}

func getAccountBankName(account map[string]interface{}, data map[string]interface{}) string {
	if custom, ok := account["bank_name"].(string); ok && custom != "" {
		return custom
	}
	code, _ := account["bank_code"].(string)
	banks := getBankList(data)
	if b := findBank(banks, code); b != nil {
		return b["name"]
	}
	return code
}

func getAccountInitial(account map[string]interface{}, data map[string]interface{}) string {
	name := getAccountBankName(account, data)
	if len(name) >= 2 {
		return name[:2]
	}
	return name
}

func getAccountNumber(account map[string]interface{}) string {
	n, _ := account["account_number"].(string)
	if n == "" { return "-" }
	return n
}

func getAccountName(account map[string]interface{}) string {
	n, _ := account["account_name"].(string)
	if n == "" { return "-" }
	return n
}

func hasCustomIcon(account map[string]interface{}) bool {
	t, _ := account["icon_type"].(string)
	p, _ := account["icon_path"].(string)
	return t == "upload" && p != ""
}

func getAccountIconPath(account map[string]interface{}) string {
	p, _ := account["icon_path"].(string)
	return p
}

// ============================================
// RESEPSI LABEL
// ============================================

func getResepsiLabel(data map[string]interface{}) string {
	label := getStr(data, "resepsi_label")
	custom := getStr(data, "resepsi_label_custom")
	switch label {
	case "walimatul_ursy": return "Walimatul Ursy"
	case "walimah": return "Walimah"
	case "custom":
		if custom != "" && custom != "-" { return custom }
		return "Resepsi"
	default: return "Resepsi"
	}
}

func getResepsiArabic(data map[string]interface{}) string {
	switch getStr(data, "resepsi_label") {
	case "walimatul_ursy": return "وَلِيمَةُ الْعُرْسِ"
	case "walimah": return "وَلِيمَة"
	}
	return ""
}

func isIslamicLabel(data map[string]interface{}) bool {
	l := getStr(data, "resepsi_label")
	return l == "walimatul_ursy" || l == "walimah"
}

// ============================================
// SHOW DATES / VENUE
// ============================================

func showDates(data map[string]interface{}) string {
	v := getStr(data, "show_dates")
	if v == "" || v == "-" { return "both" }
	return v
}

func showVenue(data map[string]interface{}) string {
	v := getStr(data, "show_venue")
	if v == "" || v == "-" { return "both" }
	return v
}

// ============================================
// LOVE STORIES / GALLERY
// ============================================

func getLoveStories(data map[string]interface{}) []map[string]string {
	result := []map[string]string{}
	v, ok := data["love_stories"]
	if !ok || v == nil { return result }
	arr, ok := v.([]interface{})
	if !ok { return result }
	for _, item := range arr {
		if m, ok := item.(map[string]interface{}); ok {
			title, _ := m["title"].(string)
			desc, _ := m["desc"].(string)
			result = append(result, map[string]string{"title": title, "desc": desc})
		}
	}
	return result
}

func getGallery(data map[string]interface{}) []string {
	result := []string{}
	v, ok := data["gallery"]
	if !ok || v == nil { return result }
	arr, ok := v.([]interface{})
	if !ok { return result }
	for _, item := range arr {
		if s, ok := item.(string); ok {
			result = append(result, s)
		} else if m, ok := item.(map[string]interface{}); ok {
			if u, ok := m["url"].(string); ok {
				result = append(result, u)
			}
		}
	}
	return result
}

// getParentsNameSplit — gabung nama ayah & ibu dengan " & "
func getParentsNameSplit(father, mother string) string {
	if father == "" || father == "-" {
		father = ""
	}
	if mother == "" || mother == "-" {
		mother = ""
	}
	if father != "" && mother != "" {
		return father + " & " + mother
	}
	if father != "" {
		return father
	}
	if mother != "" {
		return mother
	}
	return "Bapak & Ibu"
}