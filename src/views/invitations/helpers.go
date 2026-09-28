package invitations

import (
	"encoding/json"
	"fmt"
)

// ============================================
// BASIC HELPERS (SHARED antar template)
// ============================================

// GetStr — ambil string dari map dengan fallback default
func GetStr(data map[string]interface{}, key string) string {
	if v, ok := data[key]; ok {
		if s, ok := v.(string); ok && s != "" {
			return s
		}
	}
	return "-"
}

// GetStrDefault — ambil string dengan default custom (bukan "-")
func GetStrDefault(data map[string]interface{}, key, def string) string {
	if v, ok := data[key]; ok {
		if s, ok := v.(string); ok && s != "" {
			return s
		}
	}
	return def
}

// GetBool — ambil bool dari map
func GetBool(data map[string]interface{}, key string) bool {
	if v, ok := data[key]; ok {
		if b, ok := v.(bool); ok {
			return b
		}
	}
	return false
}

// GetSlice — ambil slice dari map
func GetSlice(data map[string]interface{}, key string) []interface{} {
	v, ok := data[key]
	if !ok || v == nil {
		return []interface{}{}
	}
	if arr, ok := v.([]interface{}); ok {
		return arr
	}
	// Kalau string JSON, decode
	if str, ok := v.(string); ok {
		var arr []interface{}
		if err := json.Unmarshal([]byte(str), &arr); err == nil {
			return arr
		}
	}
	return []interface{}{}
}

// ============================================
// PARENTS NAME
// ============================================

// GetParentsName — bangun nama orang tua dari data
// side: "groom" atau "bride"
func GetParentsName(data map[string]interface{}, side string) string {
	var fatherKey, motherKey string
	if side == "groom" {
		fatherKey = "father_groom"
		motherKey = "mother_groom"
	} else {
		fatherKey = "father_bride"
		motherKey = "mother_bride"
	}

	father := GetStrDefault(data, fatherKey, "")
	mother := GetStrDefault(data, motherKey, "")

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

// ============================================
// FOTO MEMPELAI (dengan fallback)
// ============================================

// GetGroomPhoto — ambil foto mempelai pria dengan fallback
func GetGroomPhoto(data map[string]interface{}) string {
	// Cek groom_photo
	if v := GetStrDefault(data, "groom_photo", ""); v != "" {
		return v
	}
	// Fallback ke photo_groom (legacy)
	if v := GetStrDefault(data, "photo_groom", ""); v != "" {
		return v
	}
	return "/images/defaults/groom.jpg"
}

// GetBridePhoto — ambil foto mempelai wanita dengan fallback
func GetBridePhoto(data map[string]interface{}) string {
	if v := GetStrDefault(data, "bride_photo", ""); v != "" {
		return v
	}
	if v := GetStrDefault(data, "photo_bride", ""); v != "" {
		return v
	}
	return "/images/defaults/bride.jpg"
}

// ============================================
// RESEPSI LABEL
// ============================================

// GetResepsiLabel — ambil label resepsi yang human-readable
func GetResepsiLabel(data map[string]interface{}) string {
	label := GetStrDefault(data, "resepsi_label", "resepsi")
	custom := GetStrDefault(data, "resepsi_label_custom", "")

	switch label {
	case "walimatul_ursy":
		return "Walimatul Ursy"
	case "walimah":
		return "Walimah"
	case "custom":
		if custom != "" {
			return custom
		}
		return "Resepsi"
	default:
		return "Resepsi"
	}
}

// IsIslamicResepsi — cek apakah label resepsi adalah islami
func IsIslamicResepsi(data map[string]interface{}) bool {
	label := GetStrDefault(data, "resepsi_label", "resepsi")
	return label == "walimatul_ursy" || label == "walimah"
}

// GetArabicLabel — ambil label Arab untuk resepsi islami
func GetArabicLabel(data map[string]interface{}) string {
	label := GetStrDefault(data, "resepsi_label", "resepsi")
	switch label {
	case "walimatul_ursy":
		return "وَلِيمَةُ الْعُرْسِ"
	case "walimah":
		return "وَلِيمَة"
	}
	return ""
}

// ============================================
// MUSIC
// ============================================

// GetMusicUrl — ambil URL musik dari data (kosong kalau tidak ada)
func GetMusicUrl(data map[string]interface{}) string {
	v := GetStrDefault(data, "custom_music", "")
	if v == "" {
		return ""
	}
	// Kalau path relatif, tambah /storage/
	if len(v) > 0 && v[0] != '/' && v[0] != 'h' {
		return "/storage/" + v
	}
	return v
}

// ============================================
// SHOW/HIDE SECTION
// ============================================

// ShouldShowBankAccounts — cek apakah tampilkan rekening
func ShouldShowBankAccounts(data map[string]interface{}) bool {
	return GetStrDefault(data, "show_bank_accounts", "no") == "yes"
}

// GetBankAccounts — ambil list bank accounts
func GetBankAccounts(data map[string]interface{}) []map[string]string {
	var result []map[string]string

	v, ok := data["bank_accounts"]
	if !ok || v == nil {
		return result
	}

	// Handle string JSON
	if str, ok := v.(string); ok {
		var arr []map[string]interface{}
		if err := json.Unmarshal([]byte(str), &arr); err == nil {
			for _, item := range arr {
				acc := make(map[string]string)
				for k, val := range item {
					if s, ok := val.(string); ok {
						acc[k] = s
					}
				}
				result = append(result, acc)
			}
		}
		return result
	}

	// Handle array of map
	if arr, ok := v.([]interface{}); ok {
		for _, item := range arr {
			if m, ok := item.(map[string]interface{}); ok {
				acc := make(map[string]string)
				for k, val := range m {
					if s, ok := val.(string); ok {
						acc[k] = s
					}
				}
				result = append(result, acc)
			}
		}
	}

	return result
}

// ============================================
// BANK LIST (untuk display icon inisial)
// ============================================

// GetBankName — ambil nama bank dari code
func GetBankName(code string) string {
	bankMap := map[string]string{
		"bca": "BCA", "mandiri": "Mandiri", "bni": "BNI", "bri": "BRI",
		"bsi": "BSI", "cimb": "CIMB Niaga", "permata": "Permata",
		"danamon": "Danamon", "btn": "BTN", "mega": "Mega", "bjb": "BJB",
		"jago": "Jago", "seabank": "SeaBank", "blu": "Blu BCA", "jenius": "Jenius",
		"gopay": "GoPay", "ovo": "OVO", "dana": "DANA", "shopeepay": "ShopeePay",
		"linkaja": "LinkAja", "qris": "QRIS", "other": "Lainnya",
	}
	if name, ok := bankMap[code]; ok {
		return name
	}
	return code
}

// ============================================
// HELPERS: FORMAT
// ============================================

// GetInitial — ambil huruf pertama dari nama
func GetInitial(name string) string {
	for _, c := range name {
		if c != ' ' {
			return string(c)
		}
	}
	return "?"
}

// GetGuestGreeting — greeting untuk tamu
func GetGuestGreeting(guestName string) string {
	if guestName == "" || guestName == "Tamu Undangan" {
		return "Salam Hangat untuk"
	}
	return "Kepada"
}

// StrContains — cek apakah string mengandung substring
func StrContains(s, substr string) bool {
	return len(s) >= len(substr) && contains(s, substr)
}

func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// suppress unused
var _ = fmt.Sprintf
