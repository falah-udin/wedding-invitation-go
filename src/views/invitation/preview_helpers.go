package invitation

import (
	"encoding/json"
	"strings"
)

// getStringOrDash — ambil string, kalau kosong return "-"
func getStringOrDash(data map[string]interface{}, key string) string {
	v := getStringFromMap(data, key, "")
	if v == "" {
		return "-"
	}
	return v
}

// getSliceLen — hitung panjang slice
func getSliceLen(data map[string]interface{}, key string) int {
	return len(getSliceFromMap(data, key))
}

// isFieldFilled — cek apakah field terisi
func isFieldFilled(data map[string]interface{}, key string) bool {
	if v, ok := data[key]; ok && v != nil {
		if s, ok := v.(string); ok {
			return strings.TrimSpace(s) != ""
		}
		if arr, ok := v.([]interface{}); ok {
			return len(arr) > 0
		}
		return true
	}
	return false
}

// formatBankAccounts — format bank accounts dari data
func getBankAccounts(data map[string]interface{}) []map[string]string {
	var result []map[string]string

	v, ok := data["bank_accounts"]
	if !ok || v == nil {
		return result
	}

	// Handle kalau v adalah string JSON
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

	// Handle kalau v adalah array
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

// getResepsiLabelDisplay — tampilkan label resepsi yang human-readable
func getResepsiLabelDisplay(data map[string]interface{}) string {
	label := getStringFromMap(data, "resepsi_label", "resepsi")
	custom := getStringFromMap(data, "resepsi_label_custom", "")

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

// getShowLabel — label untuk show_dates / show_venue
func getShowLabel(data map[string]interface{}, key string) string {
	v := getStringFromMap(data, key, "both")
	switch v {
	case "both":
		return "Keduanya"
	case "akad":
		return "Akad Saja"
	case "resepsi":
		return "Resepsi Saja"
	case "none":
		return "Tidak Ada"
	}
	return v
}
