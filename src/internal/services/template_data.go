package services

import "encoding/json"

// MergeTemplateData — merge data lama + baru
// Karena field name sudah sama, cukup overwrite
func MergeTemplateData(existing map[string]interface{}, newData map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	// 1. Copy semua data lama
	for k, v := range existing {
		result[k] = v
	}

	// 2. Overwrite dengan data baru
	for k, v := range newData {
		if v == nil {
			continue
		}
		// Jangan overwrite dengan string kosong
		if s, ok := v.(string); ok && s == "" {
			continue
		}
		result[k] = v
	}

	return result
}

// ParseJSONMap — parse string JSON ke map
func ParseJSONMap(jsonStr string) map[string]interface{} {
	result := make(map[string]interface{})
	if jsonStr == "" {
		return result
	}
	json.Unmarshal([]byte(jsonStr), &result)
	return result
}

// ToJSONString — konversi map ke JSON string
func ToJSONString(data map[string]interface{}) string {
	b, err := json.Marshal(data)
	if err != nil {
		return "{}"
	}
	return string(b)
}

// GetStringFromMap — ambil nilai string dari map
func GetStringFromMap(data map[string]interface{}, key, defaultVal string) string {
	if v, ok := data[key]; ok {
		if s, ok := v.(string); ok && s != "" {
			return s
		}
	}
	return defaultVal
}

// GetSliceFromMap — ambil nilai slice dari map
func GetSliceFromMap(data map[string]interface{}, key string) []interface{} {
	if v, ok := data[key]; ok {
		if arr, ok := v.([]interface{}); ok {
			return arr
		}
	}
	return []interface{}{}
}
