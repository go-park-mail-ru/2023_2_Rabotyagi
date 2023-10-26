package models

import (
	"time"
)

type User struct {
	ID       uint64
	Email    string    `json:"email"`
	Phone    string    `json:"phone"`
	Name     string    `json:"name"`
	Pass     string    `json:"pass"`
	Birthday time.Time `json:"birthday"`
}

type UserWithoutPass struct {
	ID       uint64
	Email    string    `json:"email"`
	Phone    string    `json:"phone"`
	Name     string    `json:"name"`
	Birthday time.Time `json:"birthday"`
}

type UserWithoutID struct {
	Email    string    `json:"email"`
	Phone    string    `json:"phone"`
	Name     string    `json:"name"`
	Pass     string    `json:"pass"`
	Birthday time.Time `json:"birthday"`
}

type PreUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
