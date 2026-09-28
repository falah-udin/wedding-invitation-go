package services

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

var (
	logFile     *os.File
	logFilePath string
)

// InitLogger — setup log ke stdout + file
func InitLogger() error {
	// Buat folder log
	logDir := "/app/storage/logs"
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return fmt.Errorf("gagal membuat folder log: %w", err)
	}

	// Nama file dengan tanggal: wedding-2026-09-29.log
	today := time.Now().Format("2006-01-02")
	logFilePath = filepath.Join(logDir, fmt.Sprintf("wedding-%s.log", today))

	// Buka file (append mode)
	var err error
	logFile, err = os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("gagal membuka file log: %w", err)
	}

	// Multi-writer: stdout + file
	multiWriter := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(multiWriter)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	log.Printf("📝 Logger siap. File: %s", logFilePath)
	return nil
}

// CloseLogger — tutup file log (dipanggil saat shutdown)
func CloseLogger() {
	if logFile != nil {
		logFile.Close()
	}
}

// LogError — catat error dengan detail
func LogError(context string, err error, details map[string]interface{}) {
	log.Printf("❌ ERROR [%s]: %v", context, err)
	if len(details) > 0 {
		for k, v := range details {
			log.Printf("   ├─ %s: %v", k, v)
		}
	}
}

// LogInfo — catat info penting
func LogInfo(context string, message string) {
	log.Printf("ℹ️  [%s] %s", context, message)
}

// LogSuccess — catat operasi sukses
func LogSuccess(context string, message string) {
	log.Printf("✅ [%s] %s", context, message)
}

// LogWarn — catat warning
func LogWarn(context string, message string) {
	log.Printf("⚠️  [%s] %s", context, message)
}

// GetLogFilePath — ambil path file log
func GetLogFilePath() string {
	return logFilePath
}
