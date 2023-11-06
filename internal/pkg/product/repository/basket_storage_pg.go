package repository

import (
	"context"
	"fmt"
	"log"

	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/models"
	myerrors "github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/errors"

	"github.com/jackc/pgx/v5"
)

var ErrLessStatus = myerrors.NewError("статус заказа должен только увеличиваться")

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
         "product".title, "product".price, "product".city, "order".count, 
         "product".delivery, "product".safe_deal, "product".in_favourites, 
		 FROM public."order" WHERE owner_id=$1 AND status=0
		 LEFT JOIN public."product" ON "order".product_id = "product".id`

	ordersInBasketRows, err := tx.Query(ctx, SQLSelectOrdersInBasketByUserID, userID)
	if err != nil {
		log.Printf("in selectOrdersInBasketByUserID: %+v\n", err)

		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	curOrder := new(models.OrderInBasket)

	_, err = pgx.ForEachRow(ordersInBasketRows, []any{
		&curOrder.ID, &curOrder.OwnerID, &curOrder.ProductID,
		&curOrder.Title, &curOrder.Price, &curOrder.City,
		&curOrder.Count, &curOrder.Delivery, &curOrder.SafeDeal,
		&curOrder.InFavourites,
	}, func() error {
		orders = append(orders, &models.OrderInBasket{ //nolint:exhaustruct
			ID:           curOrder.ID,
			OwnerID:      curOrder.OwnerID,
			ProductID:    curOrder.ProductID,
			Title:        curOrder.Title,
			Price:        curOrder.Price,
			City:         curOrder.City,
			Count:        curOrder.Count,
			Delivery:     curOrder.Delivery,
			SafeDeal:     curOrder.SafeDeal,
			InFavourites: curOrder.InFavourites,
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

			order.Images = images
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
	tx pgx.Tx, orderID uint64, newCount uint32,
) error {
	SQLUpdateOrderCountByOrderID := `UPDATE public."order"
		 SET count=$1
		 WHERE id=$2`

	_, err := tx.Exec(ctx, SQLUpdateOrderCountByOrderID, newCount, orderID)
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

func (p *ProductStorage) UpdateOrderCount(ctx context.Context, orderID uint64, newCount uint32) (*models.Order, error) {
	updatedOrder := &models.Order{} //nolint:exhaustruct

	err := pgx.BeginFunc(ctx, p.pool, func(tx pgx.Tx) error {
		err := p.updateOrderCountByOrderID(ctx, tx, orderID, newCount)
		if err != nil {
			return err
		}

		updatedOrder, err = p.getOrderByID(ctx, tx, orderID)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		log.Printf("in UpdateOrderCount: %+v\n", err)

		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return updatedOrder, nil
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

func (p *ProductStorage) getStatusByOrderID(ctx context.Context, tx pgx.Tx, orderID uint64) (uint8, error) {
	SQLGetOrderByID := `SELECT status 
		 FROM public."order" WHERE id=$1`

	orderRow := tx.QueryRow(ctx, SQLGetOrderByID, orderID)

	var status uint8
	err := orderRow.Scan(status)

	if err != nil {
		log.Printf("in getStatusByOrderID: %+v", err)

		return 255, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return status, nil
}

func (p *ProductStorage) decreaseAvailableCountByOrderID(ctx context.Context, tx pgx.Tx, orderID uint64) error {
	SQLDecreaseAvailableCountByOrderID := `UPDATE public."product"
		 SET available_count = available_count - 1
		 WHERE id = (
			SELECT product_id
			FROM public."order"
			WHERE id = $1
		 )`

	_, err := tx.Exec(ctx, SQLDecreaseAvailableCountByOrderID, orderID)
	if err != nil {
		log.Printf("in decreaseAvailableCountByOrderID: %+v", err)

		return fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return nil
}

func (p *ProductStorage) UpdateOrderStatus(ctx context.Context,
	orderID uint64, newStatus uint8,
) (*models.Order, error) {
	updatedOrder := &models.Order{} //nolint:exhaustruct

	err := pgx.BeginFunc(ctx, p.pool, func(tx pgx.Tx) error {
		curStatus, err := p.getStatusByOrderID(ctx, tx, orderID)
		if err != nil {
			return err
		}

		if newStatus <= curStatus {
			return fmt.Errorf(myerrors.ErrTemplate, ErrLessStatus)
		}

		if curStatus == 0 {
			err = p.decreaseAvailableCountByOrderID(ctx, tx, orderID)
			if err != nil {
				return err
			}
		}

		err = p.updateOrderStatusByOrderID(ctx, tx, orderID, newStatus)
		if err != nil {
			return err
		}

		updatedOrder, err = p.getOrderByID(ctx, tx, orderID)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		log.Printf("in UpdateOrderStatus: %+v\n", err)

		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return updatedOrder, nil
}