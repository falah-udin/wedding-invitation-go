package admin

import (
	"fmt"
	"strconv"
)

// formatNumber — format angka dengan pemisah ribuan
// Contoh: 45678 → "45.678"
func formatNumber(n int) string {
	if n < 1000 {
		return strconv.Itoa(n)
	}

	// Konversi ke string
	s := strconv.Itoa(n)
	
	// Bagi jadi grup 3 digit dari belakang
	var result []string
	for len(s) > 3 {
		result = append([]string{s[len(s)-3:]}, result...)
		s = s[:len(s)-3]
	}
	result = append([]string{s}, result...)

	// Gabungkan dengan titik
	output := ""
	for i, part := range result {
		if i > 0 {
			output += "."
		}
		output += part
	}

	return output
}

// _ pastikan fmt dipakai (kalau nanti butuh)
var _ = fmt.Sprintf
