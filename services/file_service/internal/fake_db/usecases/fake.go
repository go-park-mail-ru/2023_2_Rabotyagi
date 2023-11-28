package usecases

import (
	"database/sql"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/models"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/utils"

	"github.com/brianvoe/gofakeit/v6"
)

func FakeUserWihtoutID(index int) (*models.UserWithoutID, error) {
	user := new(models.UserWithoutID)

	var err error

	user.Name = sql.NullString{Valid: true, String: gofakeit.Name()}
	user.Email = gofakeit.Email() + strconv.Itoa(index)
	user.Phone = sql.NullString{Valid: true, String: strconv.Itoa(gofakeit.Number(1, 8999000000) + index)}
	user.Birthday = sql.NullTime{Valid: true, Time: gofakeit.Date()}
	user.Password = gofakeit.Password(true, true, true, true, true, 16) //nolint:gomnd

	user.Password, err = utils.HashPass(user.Password)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func FakeProduct(userMaxCount uint, categoryMaxCount uint, cityMaxCount uint) *models.Product {
	wordsInDescription := 10
	maxPrice := 1000
	maxAvailableCount := 10
	maxViews := 100
	preProduct := new(models.Product)

	preProduct.SalerID = uint64(gofakeit.Number(1, int(userMaxCount)))
	preProduct.CategoryID = uint64(gofakeit.Number(1, int(categoryMaxCount)))
	preProduct.CityID = uint64(gofakeit.Number(1, int(cityMaxCount)))
	preProduct.Title = gofakeit.BookTitle()
	preProduct.Description = gofakeit.Sentence(wordsInDescription)
	preProduct.Price = uint64(gofakeit.Number(1, maxPrice))
	preProduct.AvailableCount = uint32(gofakeit.Number(1, maxAvailableCount))
	preProduct.CityID = 1
	preProduct.Delivery = gofakeit.Bool()
	preProduct.SafeDeal = gofakeit.Bool()
	preProduct.Views = uint32(gofakeit.Number(0, maxViews))

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

func FakeCity() string {
	return gofakeit.City()
}

type FakeGeneratorImg struct {
	pathToRoot string
	prefixURL  string
	imgStorage map[string][]byte
}

func NewFakeGeneratorImg(maxNameImage uint, prefixURL string, pathToRoot string) (*FakeGeneratorImg, error) {
	imgStorage := make(map[string][]byte, maxNameImage)

	for i := 1; i <= int(maxNameImage); i++ {
		file, err := os.Open(fmt.Sprintf("%s/static/images_for_fake_db/%d.png", pathToRoot, i))
		if err != nil {
			return nil, err
		}

		rawImg, err := io.ReadAll(file)
		if err != nil {
			return nil, err
		}

		URLToFile, err := utils.Hash256(rawImg)
		if err != nil {
			return nil, err
		}

		imgStorage[URLToFile] = rawImg

		err = file.Close()
		if err != nil {
			return nil, err
		}

		URLFile, err := os.OpenFile(fmt.Sprintf("%s/static/img/%s", pathToRoot, URLToFile),
			os.O_RDWR|os.O_CREATE, 0755)
		if err != nil {
			return nil, err
		}

		_, err = URLFile.Write(rawImg)
		if err != nil {
			return nil, err
		}

		err = URLFile.Close()
		if err != nil {
			return nil, err
		}
	}

	return &FakeGeneratorImg{
		imgStorage: imgStorage, pathToRoot: pathToRoot, prefixURL: prefixURL,
	}, nil
}

func (f *FakeGeneratorImg) GetURLs(countURL uint) []string {
	minCount := min(int(countURL), len(f.imgStorage))
	result := make([]string, minCount)

	i := 0
	for filename := range f.imgStorage {
		if i == minCount {
			break
		}

		result[i] = f.prefixURL + filename
		i++
	}

	return result
}
