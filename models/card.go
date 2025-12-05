package models

import (
	"time"

	"github.com/google/uuid"
)

// Card merepresentasikan tugas atau item dalam sebuah List (seperti kartu di Trello).
// Card merepresentasikan tugas atau item dalam sebuah List (seperti kartu di Trello).
type Card struct {
	// InternalId: Primary Key database.
	// Tag `gorm:"primaryKey"` menandakan ini kolom utama.
	InternalId int64 `json:"internal_id" db:"internal_id" gorm:"primaryKey;autoIncrement"`

	// PublicId: ID unik untuk API.
	// Tag `json:"public_id"` mengubah nama field jadi snake_case di JSON.
	PublicId uuid.UUID `json:"public_id" db:"public_id"`

	// ListID: Menandakan kartu ini ada di List mana (Foreign Key).
	// Tag `gorm:"column:list_id"` memaksa nama kolom di DB jadi 'list_id'.
	ListID int64 `json:"list_id" db:"list_id" gorm:"column:list_id"`

	// Title: Judul kartu.
	Title string `json:"title" db:"title"`

	// Description: Deskripsi detail tugas.
	Description string `json:"description" db:"description"`

	// DueDate: Tenggat waktu (Opsional).
	// Tag `json:"due_date,omitempty"` berarti field ini hilang dari JSON jika nilainya kosong (nil).
	DueDate *time.Time `json:"due_date,omitempty" db:"due_date"`

	// Position: Menentukan urutan kartu dalam List (misal: 1, 2, 3).
	// Tag `json:"position"` untuk nama field di API.
	Position int `json:"position" db:"position"`

	// CreatedAt: Waktu pembuatan.
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
