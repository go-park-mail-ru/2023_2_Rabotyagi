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

func FakePreProduct(userMaxCount uint, categoryMaxCount uint) *models.PreProduct {
	wordsInDescription := 10
	maxPrice := 1000
	maxAvailableCount := 10
	preProduct := new(models.PreProduct)

	preProduct.SalerID = uint64(gofakeit.Number(1, int(userMaxCount)))
	preProduct.CategoryID = uint64(gofakeit.Number(1, int(categoryMaxCount)))
	preProduct.Title = gofakeit.BookTitle()
	preProduct.Description = gofakeit.Sentence(wordsInDescription)
	preProduct.Price = uint64(gofakeit.Number(1, maxPrice))
	preProduct.AvailableCount = uint32(gofakeit.Number(1, maxAvailableCount))
	preProduct.City = gofakeit.City()
	preProduct.Delivery = gofakeit.Bool()
	preProduct.SafeDeal = gofakeit.Bool()

	return preProduct
}

func FakePreOrder(userMaxCount uint, productMaxCount uint) *models.Order {
	preOrder := new(models.Order)

	preOrder.Count = 1
	preOrder.ProductID = uint64(gofakeit.Number(1, int(productMaxCount)))
	preOrder.OwnerID = uint64(gofakeit.Number(1, int(userMaxCount)))

	return preOrder
}

func FakeFavourite(userMaxCount uint, productMaxCount uint) (uint64, uint64) {
	ownerID := uint64(gofakeit.Number(1, int(userMaxCount)))
	productID := uint64(gofakeit.Number(1, int(productMaxCount)))

	return ownerID, productID
}
