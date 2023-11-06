package models

import (
	"database/sql"
	"strings"

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
	ID       uint64       `json:"id"       valid:"required"`
	Email    string       `json:"email"    valid:"required,email~Not valid email"`
	Phone    string       `json:"phone"    valid:"required, numeric~Phone may contain only numbers,length(1|18)~Phone length must be from 1 to 18"` //nolint
	Name     string       `json:"name"     valid:"required, regexp=^[а-яА-Яa-zA-Z0-9\s]+$~Name may contain only russian, english letter, numbers and spaces"`
	Password string       `json:"password" valid:"required,password~Password must be at least 6 symbols"`
	Birthday sql.NullTime `json:"birthday"` //nolint
}

type UserWithoutPassword struct {
	ID       uint64       `json:"id"       valid:"required"`
	Email    string       `json:"email"    valid:"required,email~Not valid email"`
	Phone    string       `json:"phone"    valid:"required, numeric~Phone may contain only numbers,length(1|18)~Phone length must be from 1 to 18"` //nolint
	Name     string       `json:"name"     valid:"required, regexp=^[а-яА-Яa-zA-Z0-9\s]+$~Name may contain only russian, english letter, numbers and spaces"`
	Birthday sql.NullTime `json:"birthday"` //nolint
}

func (u *UserWithoutPassword) Trim() {
	u.Email = strings.TrimSpace(u.Email)
	u.Name = strings.TrimSpace(u.Name)
	u.Phone = strings.TrimSpace(u.Phone)
}

type UserWithoutID struct {
	Email    string       `json:"email"    valid:"required,email~Not valid email"`
	Phone    string       `json:"phone"    valid:"required, numeric~Phone may contain only numbers,length(1|18)~Phone length must be from 1 to 18"` //nolint
	Name     string       `json:"name"     valid:"required, regexp=^[а-яА-Яa-zA-Z0-9\s]+$~Name may contain only russian, english letter, numbers and spaces"`
	Password string       `json:"password" valid:"required,password~Password must be at least 6 symbols"`
	Birthday sql.NullTime `json:"birthday"` //nolint
}

func (u *UserWithoutID) Trim() {
	u.Email = strings.TrimSpace(u.Email)
	u.Name = strings.TrimSpace(u.Name)
	u.Phone = strings.TrimSpace(u.Phone)
}
