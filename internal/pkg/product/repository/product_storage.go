package repository

import (
	"fmt"
	"sync"

	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/models"
	myerrors "github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/errors"
)

var (
	ErrProductNotExist       = myerrors.NewError("product not exist")
	ErrNoSuchCountOfProducts = myerrors.NewError("n > products count")
)

type ProductStorageMap struct {
	counterProducts uint64
	products        map[uint64]*models.Product
	mu              sync.RWMutex
}

func GenerateProducts(productStorageMap *ProductStorageMap) *ProductStorageMap {
	for i := 1; i <= 40; i++ {
		productID := productStorageMap.generateProductID()
		productStorageMap.products[productID] = &models.Product{ //nolint:exhaustruct
			ID:      productID,
			SalerID: 1,
			Title:   fmt.Sprintf("product %d", productID),
			Images: []models.Image{{
				URL: "http://84.23.53.28:8080/api/v1/img/" +
					"�%7D�̙�%7F�w���f%7C.WebP",
			}},
			Description: fmt.Sprintf("description of product %d", productID),
			Price:       100 * productID,
			SafeDeal:    true,
			Delivery:    true,
			City:        "Moscow",
		}
	}

	return productStorageMap
}

func NewProductStorageMap() *ProductStorageMap {
	return &ProductStorageMap{
		counterProducts: 0,
		products:        make(map[uint64]*models.Product),
		mu:              sync.RWMutex{},
	}
}

func (a *ProductStorageMap) GetProductsCount() int {
	a.mu.RLock()
	defer a.mu.RUnlock()

	return len(a.products)
}

func (a *ProductStorageMap) generateProductID() uint64 {
	a.counterProducts++

	return a.counterProducts
}

func (a *ProductStorageMap) GetProduct(productID uint64) (*models.Product, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	product, exists := a.products[productID]

	if exists {
		return product, nil
	}

	return nil, ErrProductNotExist
}

func (a *ProductStorageMap) AddProduct(product *models.PreProduct) {
	a.mu.Lock()
	defer a.mu.Unlock()

	id := a.generateProductID()

	a.products[id] = &models.Product{
		ID:          id,
		SalerID:     product.SalerID,
		Title:       product.Title,
		Images:      product.Images,
		Description: product.Description,
		Price:       product.Price,
		SafeDeal:    product.SafeDeal,
		Delivery:    product.Delivery,
		City:        product.City,
	}
}

func (a *ProductStorageMap) GetNProducts(n int) ([]*models.ProductInFeed, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	if n > int(a.counterProducts) {
		return nil, ErrNoSuchCountOfProducts
	}

	productsInFeedSlice := make([]*models.ProductInFeed, 0, n)

	for _, product := range a.products {
		n--

		productsInFeedSlice = append(productsInFeedSlice, &models.ProductInFeed{
			ID:       product.ID,
			Title:    product.Title,
			Images:   product.Images,
			Price:    product.Price,
			SafeDeal: product.SafeDeal,
			Delivery: product.Delivery,
			City:     product.City,
		})

		if n == 0 {
			return productsInFeedSlice, nil
		}
	}

	return productsInFeedSlice, nil
}
