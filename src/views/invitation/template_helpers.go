package invitation

import (
	"encoding/json"
	"strconv"
)

// formatUintV2 — konversi uint ke string (nama beda supaya tidak konflik)
func formatUintV2(n uint) string {
	return strconv.FormatUint(uint64(n), 10)
}

// getStringFromMap — ambil string dari map dengan default
func getStringFromMap(data map[string]interface{}, key, def string) string {
	if v, ok := data[key]; ok {
		if s, ok := v.(string); ok && s != "" {
			return s
		}
	}
	return def
}

// getSliceFromMap — ambil slice dari map
func getSliceFromMap(data map[string]interface{}, key string) []interface{} {
	if v, ok := data[key]; ok {
		if arr, ok := v.([]interface{}); ok {
			return arr
		}
	}
	return []interface{}{}
}

// getRepeaterItemsFromMap — ambil repeater items (untuk love stories)
func getRepeaterItemsFromMap(data map[string]interface{}, key string) []map[string]string {
	result := []map[string]string{}

	v, ok := data[key]
	if !ok || v == nil {
		return result
	}

	arr, ok := v.([]interface{})
	if !ok {
		return result
	}

	for _, item := range arr {
		m, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		title, _ := m["title"].(string)
		desc, _ := m["desc"].(string)
		result = append(result, map[string]string{
			"title": title,
			"desc":  desc,
		})
	}

	return result
}

// dataToJSON — konversi map ke JSON string
func dataToJSON(data map[string]interface{}) string {
	b, err := json.Marshal(data)
	if err != nil {
		return "{}"
	}
	return string(b)
}
