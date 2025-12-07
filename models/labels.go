package models

import "github.com/google/uuid"

// Label merepresentasikan tag warna yang bisa ditempel di kartu (Card).
// Contoh: Label Merah ("Urgent"), Label Hijau ("Done").
type Label struct {
	// InternalID: Primary Key.
	// Tag `gorm:"primaryKey;autoIncrement"` -> Kunci utama, nambah sendiri.
	InternalID int64 `json:"internal_id" db:"internal_id" gorm:"primaryKey;autoIncrement"`

	// PublicID: ID unik API.
	// Tag `json:"public_id"` -> Nama field di JSON.
	PublicID uuid.UUID `json:"public_id" db:"public_id"`

	// Name: Nama label (misal: "Urgent", "Bug", "Feature").
	// GORM default: varchar(255).
	Name string `json:"name" db:"name"`

	// Color: Kode warna (Hex Code), misal: "#FF0000".
	// GORM default: varchar(255).
	Color string `json:"color" db:"color"`
}
