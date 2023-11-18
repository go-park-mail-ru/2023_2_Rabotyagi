package models

import (
	"github.com/microcosm-cc/bluemonday"
)

type City struct {
	ID   uint64 `json:"id"     valid:"required"`
	Name string `json:"name"   valid:"required, length(1|256)~City name length must be from 1 to 256"` //nolint
}

func (c *City) Sanitize() {
	sanitizer := bluemonday.UGCPolicy()

	c.Name = sanitizer.Sanitize(c.Name)
}
