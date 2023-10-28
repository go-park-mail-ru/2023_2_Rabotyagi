package models

import (
	"time"

	"github.com/asaskevich/govalidator"
)

const MinLenPassword = 6

//nolint:gochecknoinits
func init() {
	govalidator.CustomTypeTagMap.Set(
		"password",
		func(i interface{}, o interface{}) bool {
			subject, ok := i.(string)
			if !ok {
				return false
			}
			if len(subject) < MinLenPassword {
				return false
			}

			return true
		},
	)
}

type User struct {
	ID       uint64    `json:"id"       valid:"required"`
	Email    string    `json:"email"    valid:"email"`
	Phone    string    `json:"phone"    valid:"numeric,length(10|11)"`
	Name     string    `json:"name"     valid:"alphanum"`
	Password string    `json:"password" valid:"password"`
	Birthday time.Time `json:"birthday" valid:"rfc3339"`
}

type UserWithoutPassword struct {
	ID       uint64    `json:"id"       valid:"required"`
	Email    string    `json:"email"    valid:"required, email"`
	Phone    string    `json:"phone"    valid:"numeric,length(10|11)"`
	Name     string    `json:"name"     valid:"required, alphanum"`
	Birthday time.Time `json:"birthday" valid:"rfc3339"`
}

type UserWithoutID struct {
	Email    string    `json:"email"    valid:"required, email"`
	Phone    string    `json:"phone"    valid:"numeric,length(10|11)"`
	Name     string    `json:"name"     valid:"required, alphanum"`
	Password string    `json:"password" valid:"required, password~Password must be at least 6 symbols"`
	Birthday time.Time `json:"birthday" valid:"rfc3339"`
}

type PreUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
