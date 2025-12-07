package models

// CardLabel adalah "Pivot Table" untuk hubungan Many-to-Many antara Card dan Label.
// Artinya: Satu Card bisa punya banyak Label, dan satu Label bisa dipakai di banyak Card.
type CardLabel struct {
	// CardID: ID Kartu.
	// Tag `gorm:"column:card_internal_id"` memaksa nama kolom di database.
	CardID int64 `json:"card_internal_id" gorm:"column:card_internal_id"`

	// LabelID: ID Label.
	// Tag `gorm:"column:label_internal_id"` memaksa nama kolom di database.
	LabelID int64 `json:"label_internal_id" gorm:"column:label_internal_id"`
}
