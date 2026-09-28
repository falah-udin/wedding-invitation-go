package users

import "strconv"

// formatUint — konversi uint ke string
func formatUint(n uint) string {
	return strconv.FormatUint(uint64(n), 10)
}
