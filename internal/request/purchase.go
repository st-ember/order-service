package request

import (
	"github.com/google/uuid"
)

type Purchase struct {
	Product  uuid.UUID
	Customer uuid.UUID
}
