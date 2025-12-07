package models

import (
	"time"

	"github.com/google/uuid"
)

// Comment merepresentasikan komentar user pada sebuah kartu (Card).
type Comment struct {
	// InternalID: Primary Key database.
	InternalID int64 `json:"internal_id" db:"internal_id" gorm:"primaryKey;autoIncrement"`

	// PublicID: ID unik API.
	PublicID uuid.UUID `json:"public_id" db:"public_id"`

	// CardID: ID Internal Kartu (Foreign Key).
	// Digunakan untuk relasi database yang efisien (JOIN).
	// Tag `gorm:"column:card_internal_id"` memaksa nama kolom di DB.
	CardID int64 `json:"card_internal_id" db:"card_internal_id" gorm:"column:card_internal_id"`

	// CardPubID: ID Public Kartu.
	// Disimpan agar saat API minta data comment, kita bisa langsung kasih ID Kartu-nya (UUID)
	// tanpa harus JOIN ke tabel Card dulu.
	CardPubID uuid.UUID `json:"card_id" db:"card_id"`

	// UserID: ID Internal User yang membuat komentar (Foreign Key).
	UserID int64 `json:"user_internal_id" db:"user_internal_id" gorm:"column:user_internal_id"`

	// UserPubID: ID Public User.
	// Sama alasannya, biar frontend langsung dapat UUID user tanpa join tabel User.
	UserPubID uuid.UUID `json:"user_id" db:"user_id"`

	// Message: Isi komentar.
	Message string `json:"message" db:"message"`

	// CreatedAt: Waktu komentar dibuat.
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
