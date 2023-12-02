package usecases_test

import (
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/product/usecases"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/my_logger"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestValidatePreProduct(t *testing.T) {
	_ = my_logger.NewNop()

	data := `{"available_count": 1,
  "category_id": 1,  "city_id": 1, "saler_id": 1,
  "title": "title", "price" : 123,
  "description": "description not empty"}`

	reader := strings.NewReader(data)

	_, err := usecases.ValidatePreProduct(reader, 1000)

	assert.NoError(t, err, "Error should be nil")
}

func TestValidatePartOfPreProduct(t *testing.T) {
	_ = my_logger.NewNop()

	data := `{
        "description": "This is a test product",
        "price": 10,
        "available_count": 5
    }`

	reader := strings.NewReader(data)

	_, err := usecases.ValidatePartOfPreProduct(reader, 1000)

	assert.NoError(t, err, "Error should be nil")
}
