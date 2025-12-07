package models

// CardAssignee adalah "Pivot Table" untuk menentukan siapa saja yang mengerjakan kartu ini.
// Relasi: Many-to-Many (Satu kartu bisa dikerjakan banyak user, satu user bisa kerjakan banyak kartu).
// CardAssignee adalah "Pivot Table" untuk menentukan siapa saja yang mengerjakan kartu ini.
// Relasi: Many-to-Many (Satu kartu bisa dikerjakan banyak user, satu user bisa kerjakan banyak kartu).
type CardAssignee struct {
	// CardID: ID Kartu.
	// Tag `gorm:"column:card_internal_id"`:
	// Memaksa GORM menamai kolom di database sebagai `card_internal_id`.
	// Jika tidak pakai ini, GORM mungkin akan menamainya `card_id` (sesuai nama field).
	CardID int64 `json:"card_internal_id" db:"card_internal_id" gorm:"column:card_internal_id"`

	// UserID: ID User yang ditugaskan (Assignee).
	// Tag `gorm:"column:user_internal_id"`:
	// Memaksa GORM menamai kolom di database sebagai `user_internal_id`.
	UserID int64 `json:"user_internal_id" db:"user_internal_id" gorm:"column:user_internal_id"`
}
