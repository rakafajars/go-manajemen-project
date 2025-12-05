package models

import "time"

// BoardMember merepresentasikan hubungan "Many-to-Many" antara Board dan User.
// Satu Board bisa punya banyak Member (User), dan satu User bisa join ke banyak Board.
// Di database, ini sering disebut "Pivot Table" atau "Join Table".
type BoardMember struct {
	// BoardID: ID dari Board.
	// Tag `gorm:"primaryKey"` (digabung dengan UserID) membuat "Composite Primary Key".
	// Artinya, kombinasi BoardID + UserID harus unik. Tidak bisa ada duplikat (User yg sama join Board yg sama 2x).
	BoardID int64 `json:"board_internal_id" db:"board_internal_id" gorm:"column:board_internal_id;primaryKey"`

	// UserID: ID dari User yang menjadi member.
	// Ini juga bagian dari Primary Key.
	UserID int64 `json:"user_internal_id" db:"user_internal_id" gorm:"column:user_internal_id;primaryKey"`

	// JoinedAt: Kapan user tersebut join ke board ini.
	JoinedAt time.Time `json:"joined_at" db:"joined_at"`
}
