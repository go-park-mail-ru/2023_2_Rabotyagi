package repository

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/models"
	myerrors "github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/errors"

	"github.com/jackc/pgx/v5"
)

var (
	ErrLessStatus    = myerrors.NewError("Статус заказа должен только увеличиваться")
	ErrNotFoundOrder = myerrors.NewError("Не получилось найти такой заказ для изменения")
)

func (p *ProductStorage) selectOrdersByUserID(ctx context.Context, tx pgx.Tx, userID uint64) ([]*models.Order, error) {
	var orders []*models.Order

	SQLSelectBasketByUserID := `SELECT  id, owner_id, product_id, count, status, created_at, updated_at, closed_at 
								FROM public."order" WHERE owner_id=$1 AND status=0`

	ordersRows, err := tx.Query(ctx, SQLSelectBasketByUserID, userID)
	if err != nil {
		log.Printf("in selectBasketByUserID: %+v\n", err)

		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	curOrder := new(models.Order)

	_, err = pgx.ForEachRow(ordersRows, []any{
		&curOrder.ID, &curOrder.OwnerID,
		&curOrder.ProductID, &curOrder.Count,
		&curOrder.Status, &curOrder.CreatedAt,
		&curOrder.UpdatedAt, &curOrder.ClosedAt,
	}, func() error {
		orders = append(orders, &models.Order{
			ID:        curOrder.ID,
			OwnerID:   curOrder.OwnerID,
			ProductID: curOrder.ProductID,
			Count:     curOrder.Count,
			Status:    curOrder.Status,
			CreatedAt: curOrder.CreatedAt,
			UpdatedAt: curOrder.UpdatedAt,
			ClosedAt:  curOrder.ClosedAt,
		})

		return nil
	})
	if err != nil {
		log.Printf("in selectBasketByUserID: %+v\n", err)

		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return orders, nil
}

func (p *ProductStorage) selectOrdersInBasketByUserID(ctx context.Context,
	tx pgx.Tx, userID uint64,
) ([]*models.OrderInBasket, error) {
	var orders []*models.OrderInBasket

	SQLSelectOrdersInBasketByUserID := `SELECT  "order".id, "order".owner_id, "order".product_id,
        "product".title, "product".price, "product".city, "order".count, "product".available_count,
        "product".delivery, "product".safe_deal, "product".saler_id FROM "order"
    INNER JOIN "product" ON "order".product_id = "product".id WHERE owner_id=$1 AND status=0;`

	ordersInBasketRows, err := tx.Query(ctx, SQLSelectOrdersInBasketByUserID, userID)
	if err != nil {
		log.Printf("in selectOrdersInBasketByUserID: %+v\n", err)

		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	curOrder := new(models.OrderInBasket)

	_, err = pgx.ForEachRow(ordersInBasketRows, []any{
		&curOrder.ID, &curOrder.OwnerID, &curOrder.ProductID,
		&curOrder.Title, &curOrder.Price, &curOrder.City,
		&curOrder.Count, &curOrder.AvailableCount, &curOrder.Delivery,
		&curOrder.SafeDeal, &curOrder.SalerID,
	}, func() error {
		orders = append(orders, &models.OrderInBasket{ //nolint:exhaustruct
			ID:             curOrder.ID,
			OwnerID:        curOrder.OwnerID,
			ProductID:      curOrder.ProductID,
			Title:          curOrder.Title,
			Price:          curOrder.Price,
			City:           curOrder.City,
			Count:          curOrder.Count,
			AvailableCount: curOrder.AvailableCount,
			Delivery:       curOrder.Delivery,
			SafeDeal:       curOrder.SafeDeal,
			InFavourites:   curOrder.InFavourites,
			SalerID:        curOrder.SalerID,
		})

		return nil
	})
	if err != nil {
		log.Printf("in selectOrdersInBasketByUserID: %+v\n", err)

		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return orders, nil
}

func (p *ProductStorage) GetOrdersInBasketByUserID(ctx context.Context,
	userID uint64,
) ([]*models.OrderInBasket, error) {
	var orders []*models.OrderInBasket

	err := pgx.BeginFunc(ctx, p.pool, func(tx pgx.Tx) error {
		ordersInner, err := p.selectOrdersInBasketByUserID(ctx, tx, userID)
		if err != nil {
			return err
		}

		for _, order := range ordersInner {
			images, err := p.selectImagesByProductID(ctx, tx, order.ProductID)
			if err != nil {
				return err
			}

			inFavourites, err := p.selectIsUserFavouriteProduct(ctx, tx, order.ProductID, userID)
			if err != nil {
				return err
			}

			order.Images = images
			order.InFavourites = inFavourites
		}

		orders = ordersInner

		return nil
	})
	if err != nil {
		log.Printf("in GetOrdersInBasketByUserID: %+v\n", err)

		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return orders, nil
}

func (p *ProductStorage) updateOrderCountByOrderID(ctx context.Context,
	tx pgx.Tx, userID uint64, orderID uint64, newCount uint32,
) error {
	SQLUpdateOrderCountByOrderID := `UPDATE public."order"
		 SET count=$1
		 WHERE id=$2 AND owner_id=$3`

	_, err := tx.Exec(ctx, SQLUpdateOrderCountByOrderID, newCount, orderID, userID)
	if err != nil {
		log.Printf("in updateOrderCountByOrderID: %+v", err)

		return fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return nil
}

func (p *ProductStorage) getOrderByID(ctx context.Context, tx pgx.Tx, orderID uint64) (*models.Order, error) {
	SQLGetOrderByID := `SELECT owner_id, product_id, count, status, created_at, updated_at, closed_at 
		 FROM public."order" WHERE id=$1`

	orderRow := tx.QueryRow(ctx, SQLGetOrderByID, orderID)
	order := models.Order{ //nolint:exhaustruct
		ID: orderID,
	}

	err := orderRow.Scan(&order.OwnerID, &order.ProductID, &order.Count, &order.Status, &order.CreatedAt,
		&order.UpdatedAt, &order.CreatedAt)

	if err != nil {
		log.Printf("in getOrderByID: %+v", err)

		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return &order, nil
}

func (p *ProductStorage) UpdateOrderCount(ctx context.Context, userID uint64, orderID uint64, newCount uint32) error {
	err := pgx.BeginFunc(ctx, p.pool, func(tx pgx.Tx) error {
		err := p.updateOrderCountByOrderID(ctx, tx, userID, orderID, newCount)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		log.Printf("in UpdateOrderCount: %+v\n", err)

		return fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return nil
}

func (p *ProductStorage) updateOrderStatusByOrderID(ctx context.Context,
	tx pgx.Tx, orderID uint64, newStatus uint8,
) error {
	SQLUpdateOrderCountByOrderID := `UPDATE public."order"
		 SET status=$1
		 WHERE id=$2`

	_, err := tx.Exec(ctx, SQLUpdateOrderCountByOrderID, newStatus, orderID)
	if err != nil {
		log.Printf("in updateOrderStatusByOrderID: %+v", err)

		return fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return nil
}

func (p *ProductStorage) getStatusAndCountByOrderID(ctx context.Context,
	tx pgx.Tx, userID uint64, orderID uint64,
) (uint8, uint32, error) {
	SQLGetOrderByID := `SELECT status, count
		 FROM public."order" WHERE owner_id=$1 AND id=$2`

	orderRow := tx.QueryRow(ctx, SQLGetOrderByID, userID, orderID)

	var status uint8

	var count uint32

	err := orderRow.Scan(&status, &count)
	if err != nil {
		log.Printf("in getStatusAndCountByOrderID: %+v", err)

		return models.OrderStatusError, 0, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return status, count, nil
}

func (p *ProductStorage) decreaseAvailableCountByOrderID(ctx context.Context, tx pgx.Tx, orderID uint64, count uint32) error {
	SQLDecreaseAvailableCountByOrderID := `UPDATE public."product"
		 SET available_count = available_count - $1
		 WHERE id = (
			SELECT product_id
			FROM public."order"
			WHERE id = $2
		 )`

	_, err := tx.Exec(ctx, SQLDecreaseAvailableCountByOrderID, count, orderID)
	if err != nil {
		log.Printf("in decreaseAvailableCountByOrderID: %+v", err)

		return fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return nil
}

func (p *ProductStorage) UpdateOrderStatus(ctx context.Context,
	userID uint64, orderID uint64, newStatus uint8,
) error {
	err := pgx.BeginFunc(ctx, p.pool, func(tx pgx.Tx) error {
		curStatus, count, err := p.getStatusAndCountByOrderID(ctx, tx, userID, orderID)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return fmt.Errorf(myerrors.ErrTemplate, ErrNotFoundOrder)
			}

			return err
		}

		if newStatus <= curStatus {
			return fmt.Errorf(myerrors.ErrTemplate, ErrLessStatus)
		}

		if curStatus == 0 {
			err = p.decreaseAvailableCountByOrderID(ctx, tx, orderID, count)
			if err != nil {
				return err
			}
		}

		err = p.updateOrderStatusByOrderID(ctx, tx, orderID, newStatus)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		log.Printf("in UpdateOrderStatus: %+v\n", err)

		return fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return nil
}

func (p *ProductStorage) insertOrder(ctx context.Context, tx pgx.Tx,
	userID uint64, productID uint64, count uint32,
) error {
	SQLInsertOrder := `INSERT INTO public."order"(owner_id, product_id, count) VALUES ($1, $2, $3)`

	_, err := tx.Exec(ctx, SQLInsertOrder, userID, productID, count)
	if err != nil {
		log.Printf("in insertOrder: %+v\n", err)

		return fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return nil
}

func (p *ProductStorage) AddOrderInBasket(ctx context.Context, userID uint64, productID uint64, count uint32) error {
	err := pgx.BeginFunc(ctx, p.pool, func(tx pgx.Tx) error {
		err := p.insertOrder(ctx, tx, userID, productID, count)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		log.Printf("in AddOrderInBasket: %+v\n", err)

		return fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return nil
}

func (p *ProductStorage) updateStatusFullBasket(ctx context.Context, tx pgx.Tx, userID uint64) error {
	SQLUpdateFullBasket := `UPDATE public."order" SET status=$1 WHERE owner_id=$2`

	_, err := tx.Exec(ctx, SQLUpdateFullBasket, models.OrderStatusInProcessing, userID)
	if err != nil {
		log.Printf("in updateStatusFullBasket: %+v\n", err)

		return fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return nil
}

func (p *ProductStorage) BuyFullBasket(ctx context.Context, userID uint64) error {
	err := pgx.BeginFunc(ctx, p.pool, func(tx pgx.Tx) error {
		err := p.updateStatusFullBasket(ctx, tx, userID)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		log.Printf("in BuyFullBasketHandler: %+v\n", err)

		return fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return nil
}
