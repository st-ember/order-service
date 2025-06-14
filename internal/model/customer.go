package model

import (
	"time"

	"github.com/google/uuid"
)

type Customer struct {
	Id uuid.UUID
	Username string
	Country uuid.UUID
	JoinedAt time.Time
}
