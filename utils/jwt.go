// Package utils berisi fungsi-fungsi helper/utilitas yang bisa dipakai di berbagai tempat
package utils

import (
	"time" // Library standar Go untuk mengelola waktu

	"github.com/golang-jwt/jwt/v5"                      // Library untuk membuat dan memvalidasi JWT
	"github.com/google/uuid"                            // Library untuk tipe data UUID
	"github.com/rakafajars/go-manajemen-project/config" // Config aplikasi kita (berisi JWTSecret, JWTExpire, dll)
)

// GenerateToken membuat token JWT untuk autentikasi user
// Parameter:
//   - userID: ID unik user dari database (tipe int64 untuk angka besar)
//   - role: Peran user seperti "admin" atau "user" (untuk authorization)
//   - email: Email user
//   - publicID: UUID publik user (lebih aman daripada expose ID asli)
//
// Return:
//   - string: Token JWT yang sudah di-generate
//   - error: Error jika ada masalah saat generate token
func GenerateToken(userID int64, role, email string, publicID uuid.UUID) (string, error) {
	// Ambil secret key dari config (kunci rahasia untuk menandatangani token)
	// Contoh secret: "mysupersecretkey123"
	secret := config.AppConfig.JWTSecret

	// Parse durasi expired token dari config
	// Contoh: "24h" = 24 jam, "30m" = 30 menit, "7d" tidak valid (gunakan "168h" untuk 7 hari)
	// Tanda _ artinya kita mengabaikan error (dalam production sebaiknya di-handle)
	duration, _ := time.ParseDuration(config.AppConfig.JWTExpire)

	// Claims adalah data yang akan disimpan di dalam token
	// jwt.MapClaims adalah tipe map[string]interface{} untuk menyimpan data bebas
	claims := jwt.MapClaims{
		"user_id":   userID,   // ID user untuk identifikasi
		"role":      role,     // Role untuk cek permission (authorization)
		"email":     email,    // Email user
		"public_id": publicID, // Public ID untuk response ke client
		// "exp" adalah expiration time - kapan token kadaluarsa
		// time.Now() = waktu sekarang
		// .Add(duration) = tambah durasi (misal 24 jam)
		// .Unix() = convert ke Unix timestamp (angka detik sejak 1 Jan 1970)
		"exp": time.Now().Add(duration).Unix(),
	}

	// Buat token baru dengan:
	// - jwt.SigningMethodHS256: Algoritma HMAC-SHA256 (standar industri, cukup aman)
	// - claims: Data yang akan dimasukkan ke dalam token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Tandatangani token dengan secret key dan kembalikan hasilnya
	// []byte(secret) = convert string secret ke byte array (format yang dibutuhkan)
	// Hasil akhir berupa string token seperti: "eyJhbGciOiJIUzI1NiIs..."
	return token.SignedString([]byte(secret))
}

// TODO: generate refresh token
