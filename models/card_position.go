package models

import (
	"github.com/google/uuid"
	"github.com/rakafajars/go-manajemen-project/models/types"
)

// CardPosition menyimpan urutan Kartu (Card) di dalam sebuah List.
// Fungsinya sama persis dengan `ListPosition`, tapi ini untuk level kartu.
// Dengan ini, Drag & Drop kartu antar posisi (atau bahkan antar List) jadi sangat cepat.
type CardPosition struct {
	// InternalID: Primary Key.
	// Tag `gorm:"primarykey"` (huruf kecil semua juga boleh, tapi standar GORM biasanya camelCase `primaryKey`).
	InternalID int64 `json:"internal_id" gorm:"primarykey;autoIncrement"`

	// PublicID: ID unik API.
	// Tag `not null` memastikan kolom ini wajib diisi di database.
	PublicID uuid.UUID `json:"public_id" gorm:"type:uuid;not null"`

	// ListID: Menandakan urutan kartu ini milik List yang mana.
	// Tag `column:list_internal_id` memaksa nama kolom di DB jadi `list_internal_id`.
	ListID int64 `json:"list_internal_id" gorm:"column:list_internal_id;not null"`

	// CardOrder: Array UUID yang menyimpan urutan kartu.
	// Contoh: ["uuid-kartu-atas", "uuid-kartu-tengah", "uuid-kartu-bawah"]
	// Tag `type:uuid[]` memastikan kolom di PostgreSQL bertipe Array.
	CardOrder types.UUIDArray `json:"card_order" gorm:"type:uuid[]"`
}
