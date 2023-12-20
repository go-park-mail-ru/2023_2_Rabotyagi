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
	Email    string         `json:"email" valid:"required,email~Некорректный формат email"`
	Phone    sql.NullString `json:"phone"    valid:"regexp=^(\+){0,1}[0-9\s]*$,length(0|18)~Телефон должен содержать только один + в начале и цифры после,length(1|18)~Длинна номера телефона вместе с плюсом не больше 18 символов"` //nolint
	Name     sql.NullString `json:"name"     valid:"regexp=^[а-яА-Яa-zA-Z0-9\s]+$~Имя может содержать только русские, английские буквы, цифры и пробелы"`                                                                             //nolint:lll
	Password string         `json:"password" valid:"required,password~Пароль должен быть минимум 6 символов"`
	Birthday sql.NullTime   `json:"birthday"`
	Avatar   sql.NullString `json:"avatar"`
}

type UserWithoutPassword struct {
	ID        uint64         `json:"id"         valid:"required"`
	Email     string         `json:"email"      valid:"required,email~Некорректный формат email"`
	Phone     sql.NullString `json:"phone"      swaggertype:"string"   valid:"regexp=^(\+){0,1}[0-9\s]*$,length(0|18)~PТелефон должен содержать только один + в начале и цифры после,length(1|18)~Длинна номера телефона вместе с плюсом не больше 18 символов"` //nolint:lll
	Name      sql.NullString `json:"name"       swaggertype:"string"   valid:"regexp=^[а-яА-Яa-zA-Z0-9\s]*$~Имя может содержать только русские, английские буквы, цифры и пробелы"`                                                                              //nolint:lll
	Birthday  sql.NullTime   `json:"birthday"   swaggertype:"string"   example:"2014-12-12T14:00:12+07:00"`
	Avatar    sql.NullString `json:"avatar"     swaggertype:"string"`
	CreatedAt time.Time      `json:"created_at" valid:"required"`
}

func (u *UserWithoutPassword) Trim() {
	u.Email = strings.TrimSpace(u.Email)
	u.Name.String = strings.TrimSpace(u.Name.String)
	u.Phone.String = strings.TrimSpace(u.Phone.String)
}

//easyjson:json
type UserWithoutID struct {
	Email    string         `json:"email" valid:"required,email~Некорректный формат email"`
	Phone    sql.NullString `json:"phone"    valid:"regexp=^(\+){0,1}[0-9\s]*$,length(0|18)~Телефон должен содержать только один + в начале и цифры после,length(1|18)~Длинна номера телефона вместе с плюсом не больше 18 символов"` //nolint
	Name     sql.NullString `json:"name"     valid:"regexp=^[а-яА-Яa-zA-Z0-9\s]+$~Имя может содержать только русские, английские буквы, цифры и пробелы"`                                                                             //nolint:lll
	Password string         `json:"password" valid:"required,password~Пароль должен быть минимум 6 символов"`
	Birthday sql.NullTime   `json:"birthday"`
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
