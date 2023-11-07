package usecases

import (
	"database/sql"
	"strconv"

	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/models"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/utils"

	"github.com/brianvoe/gofakeit/v6"
)

func FakeUserWihtoutID(index int) (*models.UserWithoutID, error) {
	user := new(models.UserWithoutID)

	var err error

	user.Name = gofakeit.Name()
	user.Email = gofakeit.Email() + strconv.Itoa(index)
	user.Phone = strconv.Itoa(gofakeit.Number(1, 8999000000) + index) //nolint:gomnd
	user.Birthday = sql.NullTime{Valid: true, Time: gofakeit.Date()}
	user.Password = gofakeit.Password(true, true, true, true, true, 16) //nolint:gomnd

	user.Password, err = utils.HashPass(user.Password)
	if err != nil {
		return nil, err
	}

	return user, nil
}
