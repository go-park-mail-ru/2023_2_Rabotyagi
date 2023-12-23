package models

import (
	"database/sql"
	"strings"
	"time"
	"unicode"

	"github.com/microcosm-cc/bluemonday"
)

type Image struct {
	URL string `json:"url" valid:"required"`
}

type Product struct {
	ID             uint64               `json:"id"              valid:"required"`
	SalerID        uint64               `json:"saler_id"        valid:"required"`
	CategoryID     uint64               `json:"category_id"     valid:"required"`
	CityID         uint64               `json:"city_id"         valid:"required"`
	Title          string               `json:"title"           valid:"required, length(1|256)~Заголовок должен быть длинной от 1 до 256 символов"`   //nolint:nolintlint
	Description    string               `json:"description"     valid:"required, length(1|4000)~Описание должно быть длинной от 1 до 4000 симвволов"` //nolint:nolintlint
	Price          sql.NullInt64        `json:"price"           swaggertype:"integer" example:"100" valid:"optional"`
	CreatedAt      time.Time            `json:"created_at"      valid:"required"`
	PremiumExpire  sql.NullTime         `json:"premium_expire"  swaggertype:"string" example:"2014-12-12T14:00:12+07:00"  valid:"optional"` //nolint:nolintlint
	Views          uint32               `json:"views"           valid:"required"`
	AvailableCount uint32               `json:"available_count" valid:"required"`
	Delivery       bool                 `json:"delivery"        valid:"optional"`
	SafeDeal       bool                 `json:"safe_deal"       valid:"optional"`
	InFavourites   bool                 `json:"in_favourites"   valid:"optional"`
	IsActive       bool                 `json:"is_active"       valid:"optional"`
	Premium        bool                 `json:"premium"         valid:"required"`
	Images         []Image              `json:"images"`
	PriceHistory   []PriceHistoryRecord `json:"price_history"   valid:"optional"`
	Favourites     uint64               `json:"favourites"      valid:"required"`
}

// PreProduct
// @Description safe_deal optional
// @Description delivery optional
type PreProduct struct {
	SalerID        uint64        `json:"saler_id"        valid:"required"`
	CategoryID     uint64        `json:"category_id"     valid:"required"`
	CityID         uint64        `json:"city_id"         valid:"required"`
	Title          string        `json:"title"           valid:"required, length(1|256)~Заголовок должен быть длинной от 1 до 256 символов"`   //nolint:nolintlint
	Description    string        `json:"description"     valid:"required, length(1|4000)~Описание должно быть длинной от 1 до 4000 симвволов"` //nolint:nolintlint
	Price          sql.NullInt64 `json:"price"           swaggertype:"integer" example:"100" valid:"optional"`
	AvailableCount uint32        `json:"available_count" valid:"required"`
	Delivery       bool          `json:"delivery"        valid:"optional"`
	SafeDeal       bool          `json:"safe_deal"       valid:"optional"`
	IsActive       bool          `json:"is_active"       valid:"optional"`
	Images         []Image       `json:"images"`
}

func (p *PreProduct) Trim() {
	p.Title = strings.TrimFunc(p.Title, unicode.IsSpace)
	p.Description = strings.TrimFunc(p.Description, unicode.IsSpace)
}

//easyjson:json
type ProductInFeed struct {
	ID             uint64  `json:"id"              valid:"required"`
	Title          string  `json:"title"           valid:"required, length(1|256)~Заголовок должен быть длинной от 1 до 256 символов"` //nolint:nolintlint
	Price          uint64  `json:"price"           valid:"required"`
	CityID         uint64  `json:"city_id"         valid:"required"`
	AvailableCount uint32  `json:"available_count" valid:"required"`
	Delivery       bool    `json:"delivery"        valid:"optional"`
	SafeDeal       bool    `json:"safe_deal"       valid:"optional"`
	InFavourites   bool    `json:"in_favourites"   valid:"optional"`
	IsActive       bool    `json:"is_active"       valid:"optional"`
	Premium        bool    `json:"premium"         valid:"required"`
	Images         []Image `json:"images"`
	Favourites     uint64  `json:"favourites"      valid:"required"`
}

//easyjson:json
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
