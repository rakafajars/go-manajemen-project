package models

import (
	"time"

	"github.com/google/uuid"
)

// CardAttachment merepresentasikan file lampiran yang di-upload ke sebuah kartu.
type CardAttachment struct {
	// InternalID: Primary Key.
	InternalID int64 `json:"internal_id" db:"internal_id" gorm:"primarykey;autoIncrement"`

	// PublicID: ID unik API.
	PublicID uuid.UUID `json:"public_id" db:"public_id"`

	// File: Menyimpan path atau URL file yang di-upload.
	// Contoh: "/uploads/images/foto.jpg" atau "https://s3.aws.com/bucket/file.pdf".
	// Database tidak menyimpan file fisiknya (blob), tapi hanya lokasinya (string).
	File string `json:"file" db:"file"`

	// UserID: Siapa yang upload file ini.
	UserID int64 `json:"user_internal_id" db:"user_internal_id"`

	// CardID: File ini milik kartu yang mana.
	CardID int64 `json:"card_internal_id" db:"card_internal_id"`

	// CreatedAt: Kapan file di-upload.
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
