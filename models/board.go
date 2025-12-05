package models

import (
	"time"

	"github.com/google/uuid"
)

// Board merepresentasikan papan kerja (seperti Trello/Jira board).
// Board merepresentasikan papan kerja (seperti Trello/Jira board).
type Board struct {
	// InternalID: Primary Key untuk database.
	// Tag `gorm:"primaryKey;autoIncrement"` artinya kolom ini adalah kunci utama dan nilainya nambah sendiri (1, 2, 3...).
	InternalID int64 `json:"internal_id" db:"internal_id" gorm:"primaryKey;autoIncrement"`

	// PublicID: ID unik untuk API.
	// Tag `json:"public_id"` berarti di response API field ini bernama "public_id".
	PublicID uuid.UUID `json:"public_id" db:"public_id"`

	// Title: Judul board.
	Title string `json:"title" db:"title"`

	// Description: Deskripsi board.
	Description string `json:"description" db:"description"`

	// OwnerID: ID User pemilik board ini (Foreign Key).
	// Tag `gorm:"column:owner_internal_id"` memaksa nama kolom di database jadi 'owner_internal_id'.
	// Tanpa tag ini, GORM mungkin akan menamainya 'owner_id' secara default.
	OwnerID int64 `json:"owner_internal_id" db:"owner_internal_id" gorm:"column:owner_internal_id"`

	// OwnerPublicID: ID Public pemilik.
	// Disimpan agar frontend bisa tahu siapa pemiliknya tanpa kita harus join ke tabel User dulu.
	OwnerPublicID uuid.UUID `json:"owner_public_id" db:"owner_public_id"`

	// CreatedAt: Waktu pembuatan.
	CreatedAt time.Time `json:"created_at" db:"created_at"`

	// DueDate: Tenggat waktu board (Opsional).
	// Tag `json:"due_date,omitempty"`:
	// - `due_date`: Nama field di JSON.
	// - `omitempty`: Jika nilainya kosong (nil), field ini HILANG dari JSON (hemat bandwidth).
	DueDate *time.Time `json:"due_date,omitempty" db:"due_date"`
}
