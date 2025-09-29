package cache

import (
	"github.com/st-ember/ecommerceprocessor/internal/request"
)

type Purchase struct {
	Timestamp int64
	Body      request.Purchase
}
