package models

import (
	"time"

	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/utils"
)

//easyjson:json
type orderJSON struct {
	ID        uint64     `json:"id"          valid:"required"`
	OwnerID   uint64     `json:"owner_id"    valid:"required"`
	ProductID uint64     `json:"product_id"  valid:"required"`
	Count     uint32     `json:"count"       valid:"required"`
	Status    uint8      `json:"status"      valid:"required"`
	CreatedAt time.Time  `json:"created_at"  valid:"required"`
	UpdatedAt time.Time  `json:"updated_at"  valid:"required"`
	ClosedAt  *time.Time `json:"closed_at"   valid:"required"`
}

func (o *Order) MarshalJSON() ([]byte, error) {
	orderJs := orderJSON{
		ID:        o.ID,
		OwnerID:   o.OwnerID,
		ProductID: o.ProductID,
		Count:     o.Count,
		Status:    o.Status,
		CreatedAt: o.CreatedAt,
		UpdatedAt: o.UpdatedAt,
		ClosedAt:  utils.NullTimeToUnsafe(o.ClosedAt),
	}

	return orderJs.MarshalJSON()
}

func (o *Order) UnmarshalJSON(bytes []byte) error {
	var orderJs orderJSON

	if err := orderJs.UnmarshalJSON(bytes); err != nil {
		return err
	}

	o.ID = orderJs.ID
	o.OwnerID = orderJs.OwnerID
	o.ProductID = orderJs.ProductID
	o.Count = orderJs.Count
	o.Status = orderJs.Status
	o.CreatedAt = orderJs.CreatedAt
	o.UpdatedAt = orderJs.UpdatedAt
	o.ClosedAt = utils.UnsafeTimeToNull(orderJs.ClosedAt)

	return nil
}
