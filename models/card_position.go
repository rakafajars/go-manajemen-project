package models

import (
	"github.com/google/uuid"
	"github.com/rakafajars/go-manajemen-project/models/types"
)

type CardPosition struct {
	InternalID int64
	PublicID   uuid.UUID
	ListID     int64
	CardOrder  types.UUIDArray
}
