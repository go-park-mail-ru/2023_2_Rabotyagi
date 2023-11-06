package models

import (
	"time"
)

type Order struct {
	ID        uint64    `json:"id"          valid:"required"`
	OwnerID   uint64    `json:"owner_id"    valid:"required"`
	ProductID uint64    `json:"product_id"  valid:"required"`
	Count     uint32    `json:"count"       valid:"required"`
	Status    uint8     `json:"status"      valid:"required"`
	CreatedAt time.Time `json:"created_at"  valid:"required"`
	UpdatedAt time.Time `json:"updated_at"  valid:"required"`
	ClosedAt  time.Time `json:"closed_at"   valid:"required"`
}

type OrderInBasket struct {
	ID           uint64  `json:"id"            valid:"required"`
	OwnerID      uint64  `json:"owner_id"      valid:"required"`
	ProductID    uint64  `json:"product_id"    valid:"required"`
	Title        string  `json:"title"         valid:"required, length(1|256)~Title length must be from 1 to 256"`
	Price        uint64  `json:"price"         valid:"required"`
	City         string  `json:"city"          valid:"required, length(1|256)~City length must be from 1 to 256"`
	Count        uint32  `json:"count"         valid:"required"`
	Delivery     bool    `json:"delivery"      valid:"required"`
	SafeDeal     bool    `json:"safe_deal"     valid:"required"`
	InFavourites bool    `json:"in_favourites" valid:"required"`
	Images       []Image `json:"images"`
}
