package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/st-ember/ecommerceprocessor/internal/enum"
)

type Purchase struct {
	Id          uuid.UUID
	Product     uuid.UUID
	Customer    uuid.UUID
	PurchasedAt time.Time
	Status      enum.PurchaseStatus
}
