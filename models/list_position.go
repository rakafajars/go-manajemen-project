package models

import (
	"github.com/google/uuid"
	"github.com/rakafajars/go-manajemen-project/models/types"
)

// ListPosition menyimpan urutan List di dalam sebuah Board.
// Kenapa dipisah? Agar saat kita geser-geser posisi List (Drag & Drop), kita tidak perlu mengupdate semua baris di tabel List.
// Cukup update satu baris di tabel ini yang berisi urutan ID-nya.
type ListPosition struct {
	// InternalId: Primary Key standar.
	InternalId int64 `json:"internal_id" db:"internal_id" gorm:"primaryKey;autoIncrement"`

	// PublicId: ID unik untuk API.
	PublicId uuid.UUID `json:"public_id" db:"public_id"`

	// BoardID: Menandakan posisi ini milik Board mana.
	BoardID int64 `json:"board_id" db:"board_id" gorm:"column:board_id"`

	// ListOrder: Berisi urutan ID List (UUID) dalam bentuk Array.
	// Contoh isi: ["uuid-list-A", "uuid-list-B", "uuid-list-C"]
	//
	// Tipe `types.UUIDArray` ini spesial karena kita buat sendiri (Custom Type).
	// Berkat method `GormDataType()`, `Scan()`, dan `Value()` yang sudah kita buat,
	// GORM jadi tahu cara simpan array ini ke kolom PostgreSQL bertipe `uuid[]`.
	ListOrder types.UUIDArray `json:"list_order" db:"list_order"`
}
