package models

import (
	"time"
)

//easyjson:json
type PriceHistoryRecord struct {
	Price     uint64    `json:"price"           valid:"required"`
	CreatedAt time.Time `json:"created_at"      valid:"required"`
}
