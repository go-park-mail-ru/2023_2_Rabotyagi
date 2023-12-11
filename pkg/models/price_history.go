package models

import (
	"time"
)

type PriceHistoryRecord struct {
	Price     uint64    `json:"price"           valid:"required"`
	CreatedAt time.Time `json:"created_at"      valid:"required"`
}
