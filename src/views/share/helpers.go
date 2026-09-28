package share

import "strconv"

// formatUint — konversi uint ke string
func formatUint(n uint) string {
	return strconv.FormatUint(uint64(n), 10)
}

// formatInt — konversi int ke string
func formatInt(n int) string {
	return strconv.Itoa(n)
}

// derefStr — dereference *string dengan default
func derefStr(s *string, def string) string {
	if s == nil {
		return def
	}
	return *s
}

// getInitial — ambil huruf pertama dari nama
func getInitial(name string) string {
	for _, c := range name {
		if c != ' ' {
			return string(c)
		}
	}
	return "?"
}
