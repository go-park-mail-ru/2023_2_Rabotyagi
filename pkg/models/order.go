package models

import (
	"database/sql"
	"github.com/microcosm-cc/bluemonday"
	"time"
)

type Order struct {
	ID        uint64       `json:"id"          valid:"required"`
	OwnerID   uint64       `json:"owner_id"    valid:"required"`
	ProductID uint64       `json:"product_id"  valid:"required"`
	Count     uint32       `json:"count"       valid:"required"`
	Status    uint8        `json:"status"      valid:"required"`
	CreatedAt time.Time    `json:"created_at"  valid:"required"`
	UpdatedAt time.Time    `json:"updated_at"  valid:"required"`
	ClosedAt  sql.NullTime `json:"closed_at"   swaggertype:"string" example:"2014-12-12T14:00:12+07:00"  valid:"required"`
}

type OrderChanges struct {
	ID     uint64 `json:"id"      valid:"required"`
	Count  uint32 `json:"count"   valid:"required"`
	Status uint8  `json:"status"  valid:"required"`
}

type PreOrder struct {
	ProductID uint64 `json:"product_id" valid:"required"`
	Count     uint32 `json:"count"      valid:"required"`
}

type OrderInBasket struct {
	ID             uint64  `json:"id"              valid:"required"`
	OwnerID        uint64  `json:"owner_id"        valid:"required"`
	SalerID        uint64  `json:"saler_id"        valid:"required"`
	ProductID      uint64  `json:"product_id"      valid:"required"`
	CityID         uint64  `json:"city_id"         valid:"required"`
	Title          string  `json:"title"           valid:"required, length(1|256)~Title length must be from 1 to 256"`
	Price          uint64  `json:"price"           valid:"required"`
	Count          uint32  `json:"count"           valid:"required"`
	AvailableCount uint32  `json:"available_count" valid:"required"`
	Delivery       bool    `json:"delivery"        valid:"required"`
	SafeDeal       bool    `json:"safe_deal"       valid:"required"`
	InFavourites   bool    `json:"in_favourites"   valid:"required"`
	Images         []Image `json:"images"`
}

const (
	OrderStatusInBasket = iota
	OrderStatusInProcessing
	OrderStatusPaid
	OrderStatusClosed
	OrderStatusError = 255
)

func (o *OrderInBasket) Sanitize() {
	sanitizer := bluemonday.UGCPolicy()

	o.Title = sanitizer.Sanitize(o.Title)
}
