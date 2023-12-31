package usecases

import (
	"context"
	"fmt"
	"io"

	productrepo "github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/product/repository"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/models"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/myerrors"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/mylogger"
)

var _ IBasketStorage = (*productrepo.ProductStorage)(nil)

type IBasketStorage interface {
	AddOrderInBasket(ctx context.Context, userID uint64, productID uint64, count uint32) (*models.OrderInBasket, error)
	GetOrdersInBasketByUserID(ctx context.Context, userID uint64) ([]*models.OrderInBasket, error)
	GetOrdersNotInBasketByUserID(ctx context.Context, userID uint64) ([]*models.OrderNotInBasket, error)
	GetOrdersSoldByUserID(ctx context.Context, userID uint64) ([]*models.OrderNotInBasket, error)
	UpdateOrderCount(ctx context.Context, userID uint64, orderID uint64, newCount uint32) error
	UpdateOrderStatus(ctx context.Context, userID uint64, orderID uint64, newStatus uint8) error
	BuyFullBasket(ctx context.Context, userID uint64) error
	DeleteOrder(ctx context.Context, orderID uint64, ownerID uint64) error
}

type BasketService struct {
	storage IBasketStorage
	logger  *mylogger.MyLogger
}

func NewBasketService(basketStorage IBasketStorage) (*BasketService, error) {
	logger, err := mylogger.Get()
	if err != nil {
		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return &BasketService{storage: basketStorage, logger: logger}, nil
}

func (b BasketService) AddOrder(ctx context.Context, r io.Reader, userID uint64) (*models.OrderInBasket, error) {
	preOrder, err := ValidatePreOrder(r)
	if err != nil {
		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	orderInBasket, err := b.storage.AddOrderInBasket(ctx, userID, preOrder.ProductID, preOrder.Count)
	if err != nil {
		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return orderInBasket, nil
}

func (b BasketService) GetOrdersByUserID(ctx context.Context, userID uint64) ([]*models.OrderInBasket, error) {
	orders, err := b.storage.GetOrdersInBasketByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	for _, order := range orders {
		order.Sanitize()
	}

	return orders, nil
}

func (b BasketService) GetOrdersNotInBasketByUserID(ctx context.Context, userID uint64,
) ([]*models.OrderNotInBasket, error) {
	orders, err := b.storage.GetOrdersNotInBasketByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	for _, order := range orders {
		order.Sanitize()
	}

	return orders, nil
}

func (b BasketService) GetOrdersSoldByUserID(ctx context.Context, userID uint64) ([]*models.OrderNotInBasket, error) {
	orders, err := b.storage.GetOrdersSoldByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	for _, order := range orders {
		order.Sanitize()
	}

	return orders, nil
}

func (b BasketService) UpdateOrderCount(ctx context.Context, r io.Reader, userID uint64) error {
	orderChanges, err := ValidateOrderChangesCount(r)
	if err != nil {
		return fmt.Errorf(myerrors.ErrTemplate, err)
	}

	err = b.storage.UpdateOrderCount(ctx, userID, orderChanges.ID, orderChanges.Count)
	if err != nil {
		return fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return nil
}

func (b BasketService) UpdateOrderStatus(ctx context.Context, r io.Reader, userID uint64) error {
	orderChanges, err := ValidateOrderChangesStatus(r)
	if err != nil {
		return fmt.Errorf(myerrors.ErrTemplate, err)
	}

	err = b.storage.UpdateOrderStatus(ctx, userID, orderChanges.ID, orderChanges.Status)
	if err != nil {
		return fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return nil
}

func (b BasketService) BuyFullBasket(ctx context.Context, userID uint64) error {
	err := b.storage.BuyFullBasket(ctx, userID)
	if err != nil {
		return fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return nil
}

func (b BasketService) DeleteOrder(ctx context.Context, orderID uint64, ownerID uint64) error {
	err := b.storage.DeleteOrder(ctx, orderID, ownerID)
	if err != nil {
		return fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return nil
}
