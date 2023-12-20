package models

import (
	"database/sql"

	"github.com/microcosm-cc/bluemonday"
)

type Category struct {
	ID       uint64        `json:"id"         valid:"required"`
	Name     string        `json:"name"       valid:"required, regexp=^[а-яА-Яa-zA-Z0-9\s]+$~Name may contain only russian, english letter, numbers and spaces"` //nolint
	ParentID sql.NullInt64 `json:"parent_id"  swaggertype:"integer" valid:"required"`
}

func (c *Category) Sanitize() {
	sanitizer := bluemonday.UGCPolicy()

	c.Name = sanitizer.Sanitize(c.Name)
}
