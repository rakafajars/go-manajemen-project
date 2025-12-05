package models

import (
	"time"

	"github.com/google/uuid"
)

// List merepresentasikan kolom daftar tugas (misal: "To Do", "In Progress", "Done").
type List struct {
	// InternalID: Primary Key database.
	InternalID int64 `json:"internal_id" db:"internal_id" gorm:"primaryKey;autoIncrement"`

	// PublicID: ID unik untuk API.
	PublicID uuid.UUID `json:"public_id" db:"public_id"`

	// BoardPublicID: ID Public dari Board tempat List ini berada.
	// Disimpan agar kita bisa filter List berdasarkan BoardPublicID yang dikirim dari Frontend.
	// Saya perbaiki tag gorm-nya menjadi `column:board_public_id` agar valid.
	BoardPublicID uuid.UUID `json:"board_public_id" db:"board_public_id" gorm:"column:board_public_id"`

	// Title: Judul List.
	Title string `json:"title" db:"title"`

	// CreatedAt: Waktu pembuatan.
	CreatedAt time.Time `json:"created_at" db:"created_at"`

	// BoardInternalID: Foreign Key asli ke tabel Board (menggunakan InternalID).
	// Ini yang digunakan database untuk relasi (JOIN) agar cepat.
	// Tag `json:"-"` artinya field ini RAHASIA/HIDDEN dari API. Frontend tidak perlu tahu ID internal ini.
	BoardInternalID int64 `json:"-" db:"board_internal_id"`
}
