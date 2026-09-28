package invitation

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// formatUint — konversi uint ke string
func formatUint(n uint) string {
	return strconv.FormatUint(uint64(n), 10)
}

// derefStr — dereference *string dengan default
func derefStr(s *string, def string) string {
	if s == nil {
		return def
	}
	return *s
}

// getStr — ambil string dari map
func getStr(m map[string]interface{}, key, def string) string {
	if v, ok := m[key]; ok {
		if s, ok := v.(string); ok && s != "" {
			return s
		}
	}
	return def
}

// bankAccountsJSON — encode bank accounts ke JSON string
func bankAccountsJSON(data map[string]interface{}) string {
	accounts, ok := data["bank_accounts"]
	if !ok || accounts == nil {
		return "[]"
	}
	b, err := json.Marshal(accounts)
	if err != nil {
		return "[]"
	}
	return string(b)
}

// bankListJSON — encode bank list ke JSON string
func bankListJSON(banks []map[string]string) string {
	b, err := json.Marshal(banks)
	if err != nil {
		return "[]"
	}
	return string(b)
}

// suppress unused
var _ = fmt.Sprintf
