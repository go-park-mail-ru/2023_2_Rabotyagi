package models

import "github.com/asaskevich/govalidator"

const minLenPassword = 6

//nolint:gochecknoinits
func init() {
	govalidator.CustomTypeTagMap.Set(
		"password",
		func(i interface{}, o interface{}) bool {
			subject, ok := i.(string)
			if !ok {
				return false
			}
			if len(subject) < minLenPassword {
				return false
			}

			return true
		},
	)
}

type User struct {
	ID       uint64 `json:"id"       valid:"required"`
	Email    string `json:"email"    valid:"email"`
	Phone    string `json:"phone"    valid:"numeric,length(10|11)"`
	Name     string `json:"name"     valid:"alphanum"`
	Password string `json:"password"     valid:"password"`
	Birthday uint64 `json:"birthday" valid:"numeric"`
}

type UserWithoutPassword struct {
	ID       uint64 `json:"id"       valid:"required"`
	Email    string `json:"email"    valid:"email"`
	Phone    string `json:"phone"    valid:"numeric,length(10|11)"`
	Name     string `json:"name"     valid:"alphanum"`
	Birthday uint64 `json:"birthday" valid:"numeric"`
}

type UserWithoutID struct {
	Email    string `json:"email"    valid:"email"`
	Phone    string `json:"phone"    valid:"numeric,length(10|11)"`
	Name     string `json:"name"     valid:"alphanum"`
	Password string `json:"password"     valid:"password"`
	Birthday uint64 `json:"birthday" valid:"numeric"`
}

type PreUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
