package projects

import (
	"strconv"
	"time"
)

// formatUint — konversi uint ke string
func formatUint(n uint) string {
	return strconv.FormatUint(uint64(n), 10)
}

// formatInt — konversi int ke string
func formatInt(n int) string {
	return strconv.Itoa(n)
}

// formatDate — format tanggal
func formatDate(t time.Time) string {
	if t.IsZero() {
		return "-"
	}
	return t.Format("02/01/2006")
}

// formatDateTime — format tanggal & waktu
func formatDateTime(t time.Time) string {
	if t.IsZero() {
		return "-"
	}
	return t.Format("02/01/2006 15:04")
}

// derefStr — dereference *string dengan default
func derefStr(s *string, def string) string {
	if s == nil {
		return def
	}
	return *s
}
