package music

import "strconv"

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
