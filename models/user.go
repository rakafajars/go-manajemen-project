package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User merepresentasikan data pengguna dalam aplikasi.
// Struct ini digunakan untuk mapping ke database (GORM) dan JSON (API).
type User struct {
	// InternalID: ID utama untuk database (Primary Key).
	// Menggunakan int64 agar performa indexing dan relasi database lebih cepat.
	// Tag `gorm:"primaryKey"` memberi tahu GORM bahwa ini adalah kunci utama tabel.
	InternalID int64 `json:"internal_id" db:"internal_id" gorm:"primaryKey"`

	// PublicID: ID unik yang aman untuk ditampilkan ke publik/API.
	// Menggunakan UUID agar ID tidak berurutan dan sulit ditebak orang lain.
	// Tag `gorm:"column:public_id"` memaksa nama kolom di DB jadi 'public_id'.
	PublicID uuid.UUID `json:"public_id" db:"public_id"`

	// Name: Nama lengkap pengguna.
	Name string `json:"name" db:"name"`

	// Email: Alamat email pengguna.
	// Tag `gorm:"unique"` memastikan tidak ada dua user dengan email yang sama di database.
	Email string `json:"email" db:"email" gorm:"unique"`

	// Password: Password yang sudah di-hash (bukan plain text).
	Password string `json:"password" db:"password" gorm:"column:password"`

	// Role: Peran pengguna, misal "admin" atau "user".
	Role string `json:"role" db:"role"`

	// CreatedAt: Waktu kapan data ini pertama kali dibuat.
	// GORM akan otomatis mengisi field ini saat data baru disimpan.
	CreatedAt time.Time `json:"created_at" db:"created_at"`

	// UpdatedAt: Waktu kapan data ini terakhir diubah.
	// GORM akan otomatis memperbarui field ini setiap kali data di-edit.
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`

	// DeletedAt: Digunakan untuk fitur "Soft Delete".
	// Jika data dihapus, baris di database TIDAK hilang, tapi field ini akan terisi waktu penghapusan.
	// Tag `json:"-"` berarti field ini TIDAK akan dimunculkan saat data diubah jadi JSON (disembunyikan dari API).
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// UserResponse adalah Struct khusus untuk format respon JSON ke Client (API).
// Alasan kita buat terpisah dari struct User (di atas) adalah untuk keamanan dan kerapian:
// 1. Password TIDAK BOLEH dikirim balik ke frontend/client.
// 2. InternalID (ID database asli) disembunyikan, diganti dengan PublicID (UUID) agar tidak mudah ditebak.
type UserResponse struct {
	PublicID  uuid.UUID      `json:"public_id"`
	Name      string         `json:"name"`
	Email     string         `json:"email"`
	Role      string         `json:"role"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-"`
}
