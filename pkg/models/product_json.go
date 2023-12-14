package models

import (
	"encoding/json"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/utils"
	"time"
)

type productJson struct {
	ID             uint64               `json:"id"              valid:"required"`
	SalerID        uint64               `json:"saler_id"        valid:"required"`
	CategoryID     uint64               `json:"category_id"     valid:"required"`
	CityID         uint64               `json:"city_id"         valid:"required"`
	Title          string               `json:"title"           valid:"required, length(1|256)~Заголовок должен быть длинной от 1 до 256 символов"`
	Description    string               `json:"description"     valid:"required, length(1|4000)~Описание должно быть длинной от 1 до 4000 симвволов"` //nolint
	Price          uint64               `json:"price"           valid:"required"`
	CreatedAt      time.Time            `json:"created_at"      valid:"required"`
	PremiumBegin   *time.Time           `json:"closed_at"       swaggertype:"string" example:"2014-12-12T14:00:12+07:00"  valid:"required"`
	PremiumExpire  *time.Time           `json:"closed_at"       swaggertype:"string" example:"2014-12-12T14:00:12+07:00"  valid:"required"`
	Views          uint32               `json:"views"           valid:"required"`
	AvailableCount uint32               `json:"available_count" valid:"required"`
	Delivery       bool                 `json:"delivery"        valid:"optional"`
	SafeDeal       bool                 `json:"safe_deal"       valid:"optional"`
	InFavourites   bool                 `json:"in_favourites"   valid:"optional"`
	IsActive       bool                 `json:"is_active"       valid:"optional"`
	Premium        bool                 `json:"premium"         valid:"required"`
	Images         []Image              `json:"images"`
	PriceHistory   []PriceHistoryRecord `json:"price_history"`
	Favourites     uint64               `json:"favourites"      valid:"required"`
}

func (p *Product) MarshalJSON() ([]byte, error) {
	var productJs = productJson{
		ID:             p.ID,
		SalerID:        p.SalerID,
		CategoryID:     p.CategoryID,
		CityID:         p.CityID,
		Title:          p.Title,
		Description:    p.Description,
		Price:          p.Price,
		CreatedAt:      p.CreatedAt,
		Views:          p.Views,
		AvailableCount: p.AvailableCount,
		Delivery:       p.Delivery,
		SafeDeal:       p.SafeDeal,
		InFavourites:   p.InFavourites,
		IsActive:       p.IsActive,
		Premium:        p.Premium,
		Images:         p.Images,
		PriceHistory:   p.PriceHistory,
		Favourites:     p.Favourites,
		PremiumBegin:   utils.NullTimeToUnsafe(p.PremiumBegin),
		PremiumExpire:  utils.NullTimeToUnsafe(p.PremiumExpire),
	}

	return json.Marshal(productJs)
}

func (p *Product) UnmarshalJSON(bytes []byte) error {
	var productJs productJson

	if err := json.Unmarshal(bytes, &productJs); err != nil {
		return err
	}

	p.ID = productJs.ID
	p.SalerID = productJs.SalerID
	p.CategoryID = productJs.CategoryID
	p.CityID = productJs.CityID
	p.Title = productJs.Title
	p.Description = productJs.Description
	p.Price = productJs.Price
	p.CreatedAt = productJs.CreatedAt
	p.Views = productJs.Views
	p.AvailableCount = productJs.AvailableCount
	p.Delivery = productJs.Delivery
	p.SafeDeal = productJs.SafeDeal
	p.InFavourites = productJs.InFavourites
	p.IsActive = productJs.IsActive
	p.Premium = productJs.Premium
	p.Images = productJs.Images
	p.PriceHistory = productJs.PriceHistory
	p.Favourites = productJs.Favourites
	p.PremiumBegin = utils.UnsafeTimeToNull(productJs.PremiumBegin)
	p.PremiumExpire = utils.UnsafeTimeToNull(productJs.PremiumExpire)

	return nil
}
