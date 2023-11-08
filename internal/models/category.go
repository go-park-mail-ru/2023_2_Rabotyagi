package models

import "github.com/microcosm-cc/bluemonday"

type Category struct {
	ID       uint64 `json:"id"       valid:"required"`
	Name     string `json:"name"     valid:"required, regexp=^[а-яА-Яa-zA-Z0-9\s]+$~Name may contain only russian, english letter, numbers and spaces"`
	ParentID uint64 `json:"parent_id"       valid:"required"`
}

func (c *Category) Sanitize() {
	sanitizer := bluemonday.UGCPolicy()

	c.Name = sanitizer.Sanitize(c.Name)
}
