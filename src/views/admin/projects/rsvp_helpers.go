package projects

import (
	"strconv"
	"wedding-invitation-go/internal/models"
)

// formatUintRSVP — konversi uint ke string
func formatUintRSVP(n uint) string {
	return strconv.FormatUint(uint64(n), 10)
}

// formatIntRSVP — konversi int ke string
func formatIntRSVP(n int) string {
	return strconv.Itoa(n)
}

// derefStrRSVP — dereference *string dengan default
func derefStrRSVP(s *string, def string) string {
	if s == nil {
		return def
	}
	return *s
}

// getAttendanceColor — warna badge berdasarkan attendance
func getAttendanceColor(attendance string) string {
	switch attendance {
	case "hadir":
		return "#10b981"
	case "tidak_hadir":
		return "#ef4444"
	case "ragu":
		return "#f59e0b"
	}
	return "#6b7280"
}

// getAttendanceIcon — ikon untuk attendance
func getAttendanceIcon(attendance string) string {
	switch attendance {
	case "hadir":
		return "✅"
	case "tidak_hadir":
		return "❌"
	case "ragu":
		return "❓"
	}
	return "•"
}

// countRsvpByAttendance — hitung jumlah rsvp berdasarkan attendance
func countRsvpByAttendance(rsvps []models.Rsvp, attendance string) int {
	count := 0
	for _, r := range rsvps {
		if r.Attendance == attendance {
			count++
		}
	}
	return count
}

// truncateRSVP — potong string panjang
func truncateRSVP(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max] + "..."
}
