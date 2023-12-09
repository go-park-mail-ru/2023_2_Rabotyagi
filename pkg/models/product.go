package models

import (
	"github.com/microcosm-cc/bluemonday"
	"strings"
	"time"
	"unicode"
)

type Image struct {
	URL string `json:"url" valid:"required"`
}

type Product struct {
	ID             uint64               `json:"id"              valid:"required"`
	SalerID        uint64               `json:"saler_id"        valid:"required"`
	CategoryID     uint64               `json:"category_id"     valid:"required"`
	CityID         uint64               `json:"city_id"         valid:"required"`
	Title          string               `json:"title"           valid:"required, length(1|256)~Заголовок должен быть длинной от 1 до 256 символов"`
	Description    string               `json:"description"     valid:"required, length(1|4000)~Описание должно быть длинной от 1 до 4000 симвволов"` //nolint
	Price          uint64               `json:"price"           valid:"required"`
	CreatedAt      time.Time            `json:"created_at"      valid:"required"`
	Views          uint32               `json:"views"           valid:"required"`
	AvailableCount uint32               `json:"available_count" valid:"required"`
	Delivery       bool                 `json:"delivery"        valid:"optional"`
	SafeDeal       bool                 `json:"safe_deal"       valid:"optional"`
	InFavourites   bool                 `json:"in_favourites"   valid:"optional"`
	IsActive       bool                 `json:"is_active"       valid:"optional"`
	Images         []Image              `json:"images"`
	PriceHistory   []PriceHistoryRecord `json:"price_history"`
	Favourites     uint64               `json:"favourites"      valid:"required"`
}

// PreProduct
// @Description safe_deal optional
// @Description delivery optional
type PreProduct struct {
	SalerID        uint64  `json:"saler_id"        valid:"required"`
	CategoryID     uint64  `json:"category_id"     valid:"required"`
	CityID         uint64  `json:"city_id"         valid:"required"`
	Title          string  `json:"title"           valid:"required, length(1|256)~Заголовок должен быть длинной от 1 до 256 символов"`
	Description    string  `json:"description"     valid:"required, length(1|4000)~Описание должно быть длинной от 1 до 4000 симвволов"` //nolint
	Price          uint64  `json:"price"           valid:"required"`
	AvailableCount uint32  `json:"available_count" valid:"required"`
	Delivery       bool    `json:"delivery"        valid:"optional"`
	SafeDeal       bool    `json:"safe_deal"       valid:"optional"`
	IsActive       bool    `json:"is_active"       valid:"optional"`
	Images         []Image `json:"images"`
}

func (p *PreProduct) Trim() {
	p.Title = strings.TrimFunc(p.Title, unicode.IsSpace)
	p.Description = strings.TrimFunc(p.Description, unicode.IsSpace)
}

type ProductInFeed struct {
	ID             uint64  `json:"id"              valid:"required"`
	Title          string  `json:"title"           valid:"required, length(1|256)~Заголовок должен быть длинной от 1 до 256 символов"`
	Price          uint64  `json:"price"           valid:"required"`
	CityID         uint64  `json:"city_id"         valid:"required"`
	AvailableCount uint32  `json:"available_count" valid:"required"`
	Delivery       bool    `json:"delivery"        valid:"optional"`
	SafeDeal       bool    `json:"safe_deal"       valid:"optional"`
	InFavourites   bool    `json:"in_favourites"   valid:"optional"`
	IsActive       bool    `json:"is_active"       valid:"optional"`
	Images         []Image `json:"images"`
	Favourites     uint64  `json:"favourites"      valid:"required"`
}

type ProductID struct {
	ProductID uint64 `json:"product_id"`
}

func (p *Product) Sanitize() {
	sanitizer := bluemonday.UGCPolicy()

	p.Title = sanitizer.Sanitize(p.Title)
	p.Description = sanitizer.Sanitize(p.Description)
}

func (p *ProductInFeed) Sanitize() {
	sanitizer := bluemonday.UGCPolicy()

	p.Title = sanitizer.Sanitize(p.Title)
}
