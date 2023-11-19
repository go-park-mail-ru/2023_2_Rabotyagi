package models

import (
	"database/sql"
	"strings"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/microcosm-cc/bluemonday"
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
	ID       uint64         `json:"id"    valid:"required"`
	Email    string         `json:"email" valid:"required,email~Not valid email"`
	Phone    sql.NullString `json:"phone"    valid:"regexp=^(\+){0,1}[0-9\s]*$,length(0|18)~Phone may contain only one + in begin and numbers,length(1|18)~Phone length must be from 1 to 18"` //nolint
	Name     sql.NullString `json:"name"     valid:"regexp=^[а-яА-Яa-zA-Z0-9\s]+$~Name may contain only russian, english letter, numbers and spaces"`
	Password string         `json:"password" valid:"required,password~Password must be at least 6 symbols"`
	Birthday sql.NullTime   `json:"birthday"`
	Avatar   sql.NullString `json:"avatar"`
}

type UserWithoutPassword struct {
	ID        uint64         `json:"id"         valid:"required"`
	Email     string         `json:"email"      valid:"required,email~Not valid email"`
	Phone     sql.NullString `json:"phone"      swaggertype:"string"   valid:"regexp=^(\+){0,1}[0-9\s]*$,length(0|18)~Phone may contain only one + in begin and numbers,length(1|18)~Phone length must be from 1 to 18"` //nolint
	Name      sql.NullString `json:"name"       swaggertype:"string"   valid:"regexp=^[а-яА-Яa-zA-Z0-9\s]*$~Name may contain only russian, english letter, numbers and spaces"`
	Birthday  sql.NullTime   `json:"birthday"   swaggertype:"string"   example:"2014-12-12T14:00:12+07:00"`
	Avatar    sql.NullString `json:"avatar"     swaggertype:"string"`
	CreatedAt time.Time      `json:"created_at" valid:"required"`
}

func (u *UserWithoutPassword) Trim() {
	u.Email = strings.TrimSpace(u.Email)
	u.Name.String = strings.TrimSpace(u.Name.String)
	u.Phone.String = strings.TrimSpace(u.Phone.String)
}

type UserWithoutID struct {
	Email    string         `json:"email" valid:"required,email~Not valid email"`
	Phone    sql.NullString `json:"phone"    valid:"regexp=^(\+){0,1}[0-9\s]*$,length(0|18)~Phone may contain only one + in begin and numbers,length(1|18)~Phone length must be from 1 to 18"` //nolint
	Name     sql.NullString `json:"name"     valid:"regexp=^[а-яА-Яa-zA-Z0-9\s]+$~Name may contain only russian, english letter, numbers and spaces"`
	Password string         `json:"password" valid:"required,password~Password must be at least 6 symbols"`
	Birthday sql.NullTime   `json:"birthday"` //nolint
	Avatar   sql.NullString `json:"avatar"`
}

func (u *UserWithoutID) Trim() {
	u.Email = strings.TrimSpace(u.Email)
	u.Name.String = strings.TrimSpace(u.Name.String)
	u.Phone.String = strings.TrimSpace(u.Phone.String)
}

func (u *UserWithoutPassword) Sanitize() {
	sanitizer := bluemonday.UGCPolicy()

	u.Email = sanitizer.Sanitize(u.Email)
	u.Phone.String = sanitizer.Sanitize(u.Phone.String)
	u.Name.String = sanitizer.Sanitize(u.Name.String)
}
