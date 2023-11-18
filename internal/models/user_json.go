package models

import "time"

type userJson struct {
	ID        uint64     `json:"id"`
	Email     string     `json:"email"`
	Phone     *string    `json:"phone"`
	Name      *string    `json:"name"`
	Birthday  *time.Time `json:"birthday"`
	Avatar    *string    `json:"avatar"`
	CreatedAt time.Time  `json:"created_at"`
}
