package models

import (
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/utils"
)

//easyjson:json
type categoryJson struct {
	ID       uint64  `json:"id"           valid:"required"`
	Name     string  `json:"name"         valid:"required, regexp=^[а-яА-Яa-zA-Z0-9\s]+$~Name may contain only russian, english letter, numbers and spaces"`
	ParentID *uint64 `json:"parent_id"    valid:"required"` //nolint
}

func (c *Category) MarshalJSON() ([]byte, error) {
	var categoryJs = categoryJson{ID: c.ID, Name: c.Name, ParentID: utils.NullInt64ToUnsafeUint(c.ParentID)}

	return categoryJs.MarshalJSON()
}

func (c *Category) SafeUnmarshalJSON(bytes []byte) error {
	var categoryJs categoryJson

	if err := categoryJs.UnmarshalJSON(bytes); err != nil {
		return err
	}

	c.ID = categoryJs.ID
	c.Name = categoryJs.Name
	c.ParentID = utils.UnsafeUint64ToNullInt(categoryJs.ParentID)

	return nil
}
