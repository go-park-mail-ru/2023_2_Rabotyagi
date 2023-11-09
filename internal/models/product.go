package models

import (
	"github.com/microcosm-cc/bluemonday"
	"strings"
	"time"
	"unicode"
)

type Image struct {
	URL string `json:"url"`
}

type Product struct {
	ID             uint64    `json:"id"              valid:"required"`
	SalerID        uint64    `json:"saler_id"        valid:"required"`
	CategoryID     uint64    `json:"category_id"     valid:"required"`
	Title          string    `json:"title"           valid:"required, length(1|256)~Title length must be from 1 to 256"`
	Description    string    `json:"description"     valid:"required, length(1|4000)~Description length must be from 1 to 4000"` //nolint
	Price          uint64    `json:"price"           valid:"required"`
	CreatedAt      time.Time `json:"created_at"      valid:"required"`
	Views          uint32    `json:"views"           valid:"required"`
	AvailableCount uint32    `json:"available_count" valid:"required"`
	City           string    `json:"city"            valid:"required, length(1|256)~City length must be from 1 to 256"` //nolint
	Delivery       bool      `json:"delivery"        valid:"required"`
	SafeDeal       bool      `json:"safe_deal"       valid:"required"`
	InFavourites   bool      `json:"in_favourites"   valid:"required"`
	Images         []Image   `json:"images"`
	Favourites     uint64    `json:"favourites"      valid:"required"`
}

// PreProduct
// @Description safe_deal optional
// @Description delivery optional
type PreProduct struct {
	SalerID        uint64  `json:"saler_id"        valid:"required"`
	CategoryID     uint64  `json:"category_id"     valid:"required"`
	Title          string  `json:"title"           valid:"required, length(1|256)~Title length must be from 1 to 256"`
	Description    string  `json:"description"     valid:"required, length(1|4000)~Description length must be from 1 to 4000"` //nolint
	Price          uint64  `json:"price"           valid:"required"`
	City           string  `json:"city"            valid:"required, length(1|256)~City length must be from 1 to 256"` //nolint
	AvailableCount uint32  `json:"available_count" valid:"required"`
	Delivery       bool    `json:"delivery"        valid:"optional"`
	SafeDeal       bool    `json:"safe_deal"       valid:"optional"`
	Images         []Image `json:"images"`
}

func (p *PreProduct) Trim() {
	p.Title = strings.TrimFunc(p.Title, unicode.IsSpace)
	p.Description = strings.TrimFunc(p.Description, unicode.IsSpace)
	p.City = strings.TrimFunc(p.City, unicode.IsSpace)
}

type ProductInFeed struct {
	ID           uint64  `json:"id"            valid:"required"`
	Title        string  `json:"title"         valid:"required, length(1|256)~Title length must be from 1 to 256"`
	Price        uint64  `json:"price"         valid:"required"`
	City         string  `json:"city"          valid:"required, length(1|256)~City length must be from 1 to 256"`
	Delivery     bool    `json:"delivery"      valid:"required"`
	SafeDeal     bool    `json:"safe_deal"     valid:"required"`
	InFavourites bool    `json:"in_favourites" valid:"required"`
	Images       []Image `json:"images"`
	Favourites   uint64  `json:"favourites"    valid:"required"`
}

func (p *Product) Sanitize() {
	sanitizer := bluemonday.UGCPolicy()

	p.City = sanitizer.Sanitize(p.City)
	p.Title = sanitizer.Sanitize(p.Title)
	p.Description = sanitizer.Sanitize(p.Description)
}

func (p *ProductInFeed) Sanitize() {
	sanitizer := bluemonday.UGCPolicy()

	p.City = sanitizer.Sanitize(p.City)
	p.Title = sanitizer.Sanitize(p.Title)
}
