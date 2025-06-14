package model

import (
	"time"

	"github.com/google/uuid"
)

type Product struct {
	Id uuid.UUID
	Category uuid.UUID
	Name string
	Description string
	Merchant uuid.UUID
	Price float32
	CreatedAt time.Time
	UpdatedAt time.Time
}
