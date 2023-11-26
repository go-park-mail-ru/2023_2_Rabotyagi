package models

import (
	"encoding/json"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/utils"
	"time"
)

type orderJson struct {
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
	var orderJs = orderJson{
		ID:        o.ID,
		OwnerID:   o.OwnerID,
		ProductID: o.ProductID,
		Count:     o.Count,
		Status:    o.Status,
		CreatedAt: o.CreatedAt,
		UpdatedAt: o.UpdatedAt,
		ClosedAt:  utils.NullTimeToUnsafe(o.ClosedAt),
	}

	return json.Marshal(orderJs)
}

func (o *Order) UnmarshalJSON(bytes []byte) error {
	var orderJs orderJson

	if err := json.Unmarshal(bytes, &orderJs); err != nil {
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
