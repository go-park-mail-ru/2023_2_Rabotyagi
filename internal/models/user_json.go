package models

import (
	"encoding/json"
	"time"

	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/utils"
)

type userJson struct {
	ID        uint64     `json:"id"`
	Email     string     `json:"email"`
	Phone     *string    `json:"phone"`
	Name      *string    `json:"name"`
	Birthday  *time.Time `json:"birthday"`
	Avatar    *string    `json:"avatar"`
	CreatedAt time.Time  `json:"created_at"`
}

func (u *UserWithoutPassword) MarshalJSON() ([]byte, error) {
	var userJs = userJson{
		ID:        u.ID,
		Email:     u.Email,
		Phone:     utils.NullStringToUnsafe(u.Phone),
		Name:      utils.NullStringToUnsafe(u.Name),
		Birthday:  utils.NullTimeToUnsafe(u.Birthday),
		Avatar:    utils.NullStringToUnsafe(u.Avatar),
		CreatedAt: u.CreatedAt,
	}

	return json.Marshal(userJs)
}

func (u *UserWithoutPassword) UnmarshalJSON(bytes []byte) error {
	var userJs userJson

	if err := json.Unmarshal(bytes, &userJs); err != nil {
		return err
	}

	u.ID = userJs.ID
	u.Email = userJs.Email
	u.Phone = utils.UnsafeStringToNull(userJs.Phone)
	u.Name = utils.UnsafeStringToNull(userJs.Name)
	u.Birthday = utils.UnsafeTimeToNull(userJs.Birthday)
	u.Avatar = utils.UnsafeStringToNull(userJs.Avatar)
	u.CreatedAt = userJs.CreatedAt

	return nil
}
