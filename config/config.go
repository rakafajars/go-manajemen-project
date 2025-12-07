// Package config berisi konfigurasi aplikasi.
// File ini bertanggung jawab untuk:
// 1. Membaca environment variables dari file .env
// 2. Menyimpan konfigurasi secara global agar bisa diakses dari mana saja
// 3. Mengelola koneksi ke database PostgreSQL
package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv" // Library untuk membaca file .env
	"gorm.io/driver/postgres"  // Driver PostgreSQL untuk GORM
	"gorm.io/gorm"             // ORM (Object Relational Mapping) untuk Go
)

// ============================================================================
// VARIABEL GLOBAL
// ============================================================================
// Variabel-variabel ini bisa diakses dari package lain dengan cara:
// config.DB atau config.AppConfig
var (
	// DB adalah instance koneksi database GORM yang akan dipakai di seluruh aplikasi.
	// Tipe *gorm.DB adalah pointer, artinya variabel ini menyimpan "alamat" objek DB.
	DB *gorm.DB

	// AppConfig menyimpan semua konfigurasi aplikasi yang dibaca dari .env.
	// Pointer (*Config) agar kita bisa mengubah isinya dari fungsi loadEnv().
	AppConfig *Config
)

// ============================================================================
// STRUCT CONFIG
// ============================================================================
// Config adalah struct (blueprint) untuk menyimpan semua konfigurasi aplikasi.
// Ini seperti "wadah" yang mengelompokkan data yang saling berhubungan.
type Config struct {
	AppPort           string // Port server, misal "3000"
	DBHost            string // Alamat database, misal "localhost"
	DBPort            string // Port database, misal "5432"
	DBUser            string // Username database
	DBPassword        string // Password database
	DBName            string // Nama database
	JWTSecret         string // Kunci rahasia untuk menandatangani JWT token
	JWTExpiredMinutes string // Berapa menit token expired
	JWTRefreshToken   string // Durasi refresh token
	JWTExpire         string // Durasi token (format: "1h", "24h", dll)
}

// ============================================================================
// FUNGSI loadEnv
// ============================================================================
// loadEnv membaca file .env dan mengisi variabel AppConfig.
// Fungsi ini harus dipanggil di awal aplikasi (biasanya di main.go atau init()).
func loadEnv() {
	// godotenv.Load() membaca file .env di root project.
	// Isi file .env akan masuk ke environment variables sistem.
	err := godotenv.Load()
	if err != nil {
		// Jika file .env tidak ditemukan, kita lanjutkan saja (pakai default).
		// Ini berguna di production yang biasanya set env lewat Docker/Kubernetes.
		log.Println("No .env file found")
	}

	// Buat objek Config baru dan isi dengan nilai dari environment variables.
	// &Config{...} artinya kita membuat objek lalu mengambil alamatnya (pointer).
	AppConfig = &Config{
		AppPort:           getEnv("PORT", "3000"),
		DBHost:            getEnv("DB_HOST", "localhost"),
		DBPort:            getEnv("DB_PORT", "5432"),
		DBUser:            getEnv("DB_USER", "postgres"),
		DBPassword:        getEnv("DB_PASSWORD", "postgres"),
		DBName:            getEnv("DB_NAME", "project_management"),
		JWTSecret:         getEnv("JWT_SECRET", "secret"),
		JWTExpiredMinutes: getEnv("JWT_EXPIREY_MINUTES", "6000"),
		JWTRefreshToken:   getEnv("REFRESH_TOKEN_EXPIRED", "24H"),
		JWTExpire:         getEnv("JWT_EXPIRED", "1h"),
	}
}

// ============================================================================
// FUNGSI getEnv
// ============================================================================
// getEnv adalah helper function untuk mengambil nilai environment variable.
// Jika key tidak ditemukan, fungsi ini mengembalikan nilai fallback (default).
//
// Parameter:
//   - key: Nama environment variable yang dicari (misal: "DB_HOST")
//   - fallback: Nilai default jika key tidak ditemukan (misal: "localhost")
//
// Return: Nilai environment variable atau fallback.
func getEnv(key string, fallback string) string {
	// os.LookupEnv mengembalikan 2 nilai:
	// 1. value: Isi dari environment variable
	// 2. exist: Boolean, true jika key ditemukan
	value, exist := os.LookupEnv(key)
	if exist {
		return value
	} else {
		return fallback
	}
}

// ============================================================================
// FUNGSI connectDB
// ============================================================================
// connectDB membuat koneksi ke database PostgreSQL menggunakan GORM.
// Fungsi ini harus dipanggil SETELAH loadEnv() agar AppConfig sudah terisi.
func connectDB() {
	cfg := AppConfig

	// DSN (Data Source Name) adalah string koneksi ke database.
	// Format: "host=... port=... user=... password=... dbname=... sslmode=..."
	// sslmode=disable artinya tidak pakai SSL (aman untuk development lokal).
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

	// gorm.Open() membuka koneksi ke database.
	// postgres.Open(dsn) menerjemahkan DSN menjadi koneksi PostgreSQL.
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		// log.Fatal akan mencetak error dan MENGHENTIKAN aplikasi.
		// Karena tanpa database, aplikasi tidak bisa berjalan.
		log.Fatal("failed to connect to database", err)
	}

	// db.DB() mengembalikan *sql.DB (koneksi database standar Go).
	// Ini dibutuhkan untuk mengatur connection pool.
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("failed to get database instance", err)
	}

	// ========================================================================
	// CONNECTION POOL SETTINGS
	// ========================================================================
	// Ini adalah pengaturan performa untuk koneksi database.

	// MaxIdleConns: Berapa koneksi "tidur" yang tetap dibuka.
	// Koneksi idle ini siap dipakai kapan saja tanpa perlu membuat koneksi baru.
	sqlDB.SetMaxIdleConns(10)

	// MaxOpenConns: Maksimal koneksi yang boleh dibuka bersamaan.
	// Jika sudah 100 koneksi aktif, request baru harus menunggu.
	sqlDB.SetMaxOpenConns(100)

	// ConnMaxLifetime: Berapa lama koneksi boleh hidup sebelum di-recycle.
	// Ini mencegah koneksi "zombie" yang sudah tidak responsif.
	sqlDB.SetConnMaxLifetime(time.Hour)

	// Simpan koneksi ke variabel global DB agar bisa dipakai di seluruh aplikasi.
	DB = db
}
