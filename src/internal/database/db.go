package database

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// Connect — koneksi ke database MySQL
func Connect() error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		getEnv("DB_USERNAME", "wedding_user"),
		getEnv("DB_PASSWORD", "wedding123"),
		getEnv("DB_HOST", "mysql"),
		getEnv("DB_PORT", "3306"),
		getEnv("DB_DATABASE", "wedding_invitation_db"),
	)

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})

	if err != nil {
		return fmt.Errorf("gagal koneksi DB: %w", err)
	}

	log.Println("✅ Berhasil konek ke MySQL")
	return nil
}

// GetDB — ambil instance DB
func GetDB() *gorm.DB {
	return DB
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}