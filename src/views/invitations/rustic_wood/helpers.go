package rustic_wood

import (
	"strconv"
	"time"
	"wedding-invitation-go/internal/models"
)
// getStr — ambil string dari map dengan default
func getStr(data map[string]interface{}, key string) string {
	if v, ok := data[key]; ok {
		if s, ok := v.(string); ok && s != "" {
			return s
		}
	}
	return "-"
}

// getParentsName — bangun nama orang tua dari data
func getParentsName(data map[string]interface{}, side string) string {
	var fatherKey, motherKey string
	if side == "groom" {
		fatherKey = "father_groom"
		motherKey = "mother_groom"
	} else {
		fatherKey = "father_bride"
		motherKey = "mother_bride"
	}

	father := getStr(data, fatherKey)
	mother := getStr(data, motherKey)

	if father == "-" {
		father = ""
	}
	if mother == "-" {
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

// getMusicUrlFromProject — ambil URL musik dari project (kosong kalau tidak ada)
func getMusicUrlFromProject(p models.Project) string {
	if p.CustomMusic == nil || *p.CustomMusic == "" {
		return ""
	}
	return "/storage/" + *p.CustomMusic
}

// ============================================
// HELPER UNTUK DATA INJECTED (dari handler)
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
			if n > 0 {
				return n
			}
		case int64:
			if n > 0 {
				return int(n)
			}
		case float64:
			if n > 0 {
				return int(n)
			}
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
	case "site_name":
		return cfg.SiteName
	case "site_tagline":
		return cfg.SiteTagline
	case "site_favicon":
		return cfg.SiteFavicon
	case "footer_description":
		return cfg.FooterDescription
	case "footer_copyright":
		return cfg.FooterCopyright
	case "footer_version":
		return cfg.FooterVersion
	case "contact_email":
		return cfg.ContactEmail
	case "contact_phone":
		return cfg.ContactPhone
	case "contact_whatsapp":
		return cfg.ContactWhatsapp
	case "contact_address":
		return cfg.ContactAddress
	case "social_instagram":
		return cfg.SocialInstagram
	case "social_facebook":
		return cfg.SocialFacebook
	case "social_tiktok":
		return cfg.SocialTiktok
	case "social_youtube":
		return cfg.SocialYoutube
	case "social_twitter":
		return cfg.SocialTwitter
	case "developer_name":
		return cfg.DeveloperName
	case "developer_email":
		return cfg.DeveloperEmail
	case "developer_url":
		return cfg.DeveloperUrl
	}
	return def
}

// getBankList — ambil bank list
func getBankList(data map[string]interface{}) []map[string]string {
	if v, ok := data["_bank_list"].([]map[string]string); ok {
		return v
	}
	return nil
}

// findBank — cari bank info by code
func findBank(banks []map[string]string, code string) map[string]string {
	for _, b := range banks {
		if b["code"] == code {
			return b
		}
	}
	return nil
}

// getBankAccounts — konversi bank_accounts dari data
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

// showBankAccounts — cek apakah bank accounts ditampilkan
func showBankAccounts(data map[string]interface{}) bool {
	return getStr(data, "show_bank_accounts") == "yes"
}

// getResepsiLabel — label resepsi human-readable
func getResepsiLabel(data map[string]interface{}) string {
	label := getStr(data, "resepsi_label")
	custom := getStr(data, "resepsi_label_custom")
	switch label {
	case "walimatul_ursy":
		return "Walimatul Ursy"
	case "walimah":
		return "Walimah"
	case "custom":
		if custom != "" && custom != "-" {
			return custom
		}
		return "Resepsi"
	default:
		return "Resepsi"
	}
}

// getResepsiArabic — label arab (kalau islami)
func getResepsiArabic(data map[string]interface{}) string {
	switch getStr(data, "resepsi_label") {
	case "walimatul_ursy":
		return "وَلِيمَةُ الْعُرْسِ"
	case "walimah":
		return "وَلِيمَة"
	}
	return ""
}

// isIslamicLabel — cek label islami
func isIslamicLabel(data map[string]interface{}) bool {
	l := getStr(data, "resepsi_label")
	return l == "walimatul_ursy" || l == "walimah"
}

// showDates — normalisasi show_dates
func showDates(data map[string]interface{}) string {
	v := getStr(data, "show_dates")
	if v == "" || v == "-" {
		return "both"
	}
	return v
}

// showVenue — normalisasi show_venue
func showVenue(data map[string]interface{}) string {
	v := getStr(data, "show_venue")
	if v == "" || v == "-" {
		return "both"
	}
	return v
}

// getLoveStories — ambil love stories
func getLoveStories(data map[string]interface{}) []map[string]string {
	result := []map[string]string{}
	v, ok := data["love_stories"]
	if !ok || v == nil {
		return result
	}
	arr, ok := v.([]interface{})
	if !ok {
		return result
	}
	for _, item := range arr {
		if m, ok := item.(map[string]interface{}); ok {
			title, _ := m["title"].(string)
			desc, _ := m["desc"].(string)
			result = append(result, map[string]string{"title": title, "desc": desc})
		}
	}
	return result
}

// getGallery — ambil gallery
func getGallery(data map[string]interface{}) []string {
	result := []string{}
	v, ok := data["gallery"]
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
		} else if m, ok := item.(map[string]interface{}); ok {
			if u, ok := m["url"].(string); ok {
				result = append(result, u)
			}
		}
	}
	return result
}

// getDefaultPhoto — default foto kalau kosong
func getDefaultPhoto(data map[string]interface{}, key string) string {
	v := getStr(data, key)
	if v == "" || v == "-" {
		return "/storage/defaults/" + key + ".jpg"
	}
	return v
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

// getBankTypeBadge — label badge bank
func getBankTypeBadge(account map[string]interface{}) string {
	t, _ := account["type"].(string)
	switch t {
	case "bank":
		return "BANK"
	case "ewallet":
		return "E-WALLET"
	case "qris":
		return "QRIS"
	default:
		return "LAINNYA"
	}
}

// getAccountBankName — nama bank dari account
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

// getAccountInitial — 2 huruf pertama nama bank
func getAccountInitial(account map[string]interface{}, data map[string]interface{}) string {
	name := getAccountBankName(account, data)
	if len(name) >= 2 {
		return name[:2]
	}
	return name
}

// getAccountNumber — nomor rekening
func getAccountNumber(account map[string]interface{}) string {
	n, _ := account["account_number"].(string)
	if n == "" {
		return "-"
	}
	return n
}

// getAccountName — nama pemilik rekening
func getAccountName(account map[string]interface{}) string {
	n, _ := account["account_name"].(string)
	if n == "" {
		return "-"
	}
	return n
}

// hasCustomIcon — cek icon upload
func hasCustomIcon(account map[string]interface{}) bool {
	t, _ := account["icon_type"].(string)
	p, _ := account["icon_path"].(string)
	return t == "upload" && p != ""
}

// getAccountIconPath — path icon upload
func getAccountIconPath(account map[string]interface{}) string {
	p, _ := account["icon_path"].(string)
	return p
}