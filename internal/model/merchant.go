package model

import (
	"time"

	"github.com/google/uuid"
)

type Merchant struct {
	Id uuid.UUID
	Username string
	Country string
	Description string
	Status MerchantStatus
	JoinedAt time.Time
}

type MerchantStatus int

const (
	Pending MerchantStatus = iota
	Active
	Suspended
)