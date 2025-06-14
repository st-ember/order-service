package model

import (
	"time"

	"github.com/google/uuid"
)

type ProductCategory struct {
	Id uuid.UUID
	Name string
	Slug string
	CreatedAt time.Time
	UpdatedAt time.Time
}
