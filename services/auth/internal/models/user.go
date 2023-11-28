package models

import (
	"github.com/asaskevich/govalidator"
	"github.com/microcosm-cc/bluemonday"
	"strings"
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
	ID       uint64 `json:"id"    valid:"required"`
	Email    string `json:"email" valid:"required,email~Not valid email"`
	Password string `json:"password" valid:"required,password~Password must be at least 6 symbols"`
}

func (u *User) Trim() {
	u.Email = strings.TrimSpace(u.Email)
}

func (u *User) Sanitize() {
	sanitizer := bluemonday.UGCPolicy()

	u.Email = sanitizer.Sanitize(u.Email)
}
