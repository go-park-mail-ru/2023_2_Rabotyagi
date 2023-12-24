package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/models"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/myerrors"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/repository"
	"github.com/jackc/pgx/v5"
)

var (
	NameSeqOrder = pgx.Identifier{"public", "order_id_seq"} //nolint:gochecknoglobals

	ErrLessStatus              = myerrors.NewErrorBadContentRequest("Статус заказа должен только увеличиваться")
	ErrNotFoundOrder           = myerrors.NewErrorBadContentRequest("Не получилось найти такой заказ для изменения")
	ErrNotFoundOrdersInBasket  = myerrors.NewErrorBadContentRequest("Не получилось найти заказы для покупки")
	ErrNoAffectedOrderRows     = myerrors.NewErrorBadContentRequest("Не получилось обновить данные заказа")
	ErrAvailableCountNotEnough = myerrors.NewErrorBadContentRequest(
		"Товара доступно меньше, чем вы пытаетесь довавить в корзину")
)

func (p *ProductStorage) selectOrdersInBasketByUserID(ctx context.Context,
	tx pgx.Tx, userID uint64,
) ([]*models.OrderInBasket, error) {
	logger := p.logger.LogReqID(ctx)

	var orders []*models.OrderInBasket

	SQLSelectOrdersInBasketByUserID := `SELECT  "order".id, "order".owner_id, "order".product_id,
        "product".title, "product".price, "product".city_id, "order".count, "product".available_count,
        "product".delivery, "product".safe_deal, "product".saler_id FROM public."order"
    INNER JOIN "product" ON "order".product_id = "product".id WHERE owner_id=$1 AND status=0;`

	ordersInBasketRows, err := tx.Query(ctx, SQLSelectOrdersInBasketByUserID, userID)
	if err != nil {
		logger.Errorln(err)

		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	curOrder := new(models.OrderInBasket)

	_, err = pgx.ForEachRow(ordersInBasketRows, []any{
		&curOrder.ID, &curOrder.OwnerID, &curOrder.ProductID,
		&curOrder.Title, &curOrder.Price, &curOrder.CityID,
		&curOrder.Count, &curOrder.AvailableCount, &curOrder.Delivery,
		&curOrder.SafeDeal, &curOrder.SalerID,
	}, func() error {
		orders = append(orders, &models.OrderInBasket{ //nolint:exhaustruct
			ID:             curOrder.ID,
			OwnerID:        curOrder.OwnerID,
			ProductID:      curOrder.ProductID,
			Title:          curOrder.Title,
			Price:          curOrder.Price,
			CityID:         curOrder.CityID,
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
		logger.Errorln(err)

		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return orders, nil
}

func (p *ProductStorage) GetOrdersInBasketByUserID(ctx context.Context, //nolint:dupl
	userID uint64,
) ([]*models.OrderInBasket, error) {
	logger := p.logger.LogReqID(ctx)

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
		logger.Errorln(err)

		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return orders, nil
}

func (p *ProductStorage) selectOrdersNotInBasketByUserID(ctx context.Context,
	tx pgx.Tx, userID uint64,
) ([]*models.OrderNotInBasket, error) {
	logger := p.logger.LogReqID(ctx)

	var orders []*models.OrderNotInBasket

	SQLSelectOrdersInBasketByUserID := `SELECT "order".id, "order".owner_id, "order".product_id,
        "product".title, "product".price, "product".city_id, "order".count, "order".status, "product".available_count,
        "product".delivery, "product".safe_deal, "product".saler_id FROM public."order"
    INNER JOIN "product" ON "order".product_id = "product".id WHERE owner_id=$1 AND status > 0 AND status < 255;`

	ordersInBasketRows, err := tx.Query(ctx, SQLSelectOrdersInBasketByUserID, userID)
	if err != nil {
		logger.Errorln(err)

		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	curOrder := new(models.OrderNotInBasket)

	_, err = pgx.ForEachRow(ordersInBasketRows, []any{
		&curOrder.ID, &curOrder.OwnerID, &curOrder.ProductID,
		&curOrder.Title, &curOrder.Price, &curOrder.CityID,
		&curOrder.Count, &curOrder.Status, &curOrder.AvailableCount, &curOrder.Delivery,
		&curOrder.SafeDeal, &curOrder.SalerID,
	}, func() error {
		orders = append(orders, &models.OrderNotInBasket{
			OrderInBasket: models.OrderInBasket{ //nolint:exhaustruct
				ID:             curOrder.ID,
				OwnerID:        curOrder.OwnerID,
				ProductID:      curOrder.ProductID,
				Title:          curOrder.Title,
				Price:          curOrder.Price,
				CityID:         curOrder.CityID,
				Count:          curOrder.Count,
				AvailableCount: curOrder.AvailableCount,
				Delivery:       curOrder.Delivery,
				SafeDeal:       curOrder.SafeDeal,
				InFavourites:   curOrder.InFavourites,
				SalerID:        curOrder.SalerID,
			},
			Status: curOrder.Status,
		})

		return nil
	})
	if err != nil {
		logger.Errorln(err)

		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return orders, nil
}

func (p *ProductStorage) GetOrdersNotInBasketByUserID(ctx context.Context, //nolint:dupl
	userID uint64,
) ([]*models.OrderNotInBasket, error) {
	logger := p.logger.LogReqID(ctx)

	var orders []*models.OrderNotInBasket

	err := pgx.BeginFunc(ctx, p.pool, func(tx pgx.Tx) error {
		ordersInner, err := p.selectOrdersNotInBasketByUserID(ctx, tx, userID)
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
		logger.Errorln(err)

		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return orders, nil
}

func (p *ProductStorage) selectOrdersSoldByUserID(ctx context.Context,
	tx pgx.Tx, userID uint64,
) ([]*models.OrderNotInBasket, error) {
	logger := p.logger.LogReqID(ctx)

	var orders []*models.OrderNotInBasket

	SQLSelectOrdersInBasketByUserID := `SELECT "order".id, "order".owner_id, "order".product_id,
        "product".title, "product".price, "product".city_id, "order".count, "order".status,  "product".available_count,
        "product".delivery, "product".safe_deal, "product".saler_id FROM public."order"
    INNER JOIN "product" ON "order".product_id = "product".id WHERE saler_id=$1 AND status > 0 AND status < 255;`

	ordersInBasketRows, err := tx.Query(ctx, SQLSelectOrdersInBasketByUserID, userID)
	if err != nil {
		logger.Errorln(err)

		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	curOrder := new(models.OrderNotInBasket)

	_, err = pgx.ForEachRow(ordersInBasketRows, []any{
		&curOrder.ID, &curOrder.OwnerID, &curOrder.ProductID,
		&curOrder.Title, &curOrder.Price, &curOrder.CityID,
		&curOrder.Count, &curOrder.Status, &curOrder.AvailableCount, &curOrder.Delivery,
		&curOrder.SafeDeal, &curOrder.SalerID,
	}, func() error {
		orders = append(orders, &models.OrderNotInBasket{
			OrderInBasket: models.OrderInBasket{ //nolint:exhaustruct
				ID:             curOrder.ID,
				OwnerID:        curOrder.OwnerID,
				ProductID:      curOrder.ProductID,
				Title:          curOrder.Title,
				Price:          curOrder.Price,
				CityID:         curOrder.CityID,
				Count:          curOrder.Count,
				AvailableCount: curOrder.AvailableCount,
				Delivery:       curOrder.Delivery,
				SafeDeal:       curOrder.SafeDeal,
				InFavourites:   curOrder.InFavourites,
				SalerID:        curOrder.SalerID,
			},
			Status: curOrder.Status,
		})

		return nil
	})
	if err != nil {
		logger.Errorln(err)

		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return orders, nil
}

func (p *ProductStorage) GetOrdersSoldByUserID(ctx context.Context, //nolint:dupl
	userID uint64,
) ([]*models.OrderNotInBasket, error) {
	logger := p.logger.LogReqID(ctx)

	var orders []*models.OrderNotInBasket

	err := pgx.BeginFunc(ctx, p.pool, func(tx pgx.Tx) error {
		ordersInner, err := p.selectOrdersSoldByUserID(ctx, tx, userID)
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
		logger.Errorln(err)

		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return orders, nil
}

func (p *ProductStorage) updateOrderCountByOrderID(ctx context.Context,
	tx pgx.Tx, userID uint64, orderID uint64, newCount uint32,
) error {
	logger := p.logger.LogReqID(ctx)

	SQLUpdateOrderCountByOrderID := `UPDATE public."order"
		 SET count=$1
		 WHERE id=$2 AND owner_id=$3`

	result, err := tx.Exec(ctx, SQLUpdateOrderCountByOrderID, newCount, orderID, userID)
	if err != nil {
		logger.Errorln(err)

		return fmt.Errorf(myerrors.ErrTemplate, err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf(myerrors.ErrTemplate, ErrNoAffectedOrderRows)
	}

	return nil
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
		return fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return nil
}

func (p *ProductStorage) updateOrderStatusByOrderID(ctx context.Context,
	tx pgx.Tx, orderID uint64, newStatus uint8,
) error {
	logger := p.logger.LogReqID(ctx)

	SQLUpdateOrderCountByOrderID := `UPDATE public."order"
		 SET status=$1
		 WHERE id=$2`

	result, err := tx.Exec(ctx, SQLUpdateOrderCountByOrderID, newStatus, orderID)
	if err != nil {
		logger.Errorln(err)

		return fmt.Errorf(myerrors.ErrTemplate, err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf(myerrors.ErrTemplate, ErrNoAffectedOrderRows)
	}

	return nil
}

func (p *ProductStorage) getStatusAndCountByOrderID(ctx context.Context,
	tx pgx.Tx, userID uint64, orderID uint64,
) (uint8, uint32, error) {
	logger := p.logger.LogReqID(ctx)

	SQLGetOrderByID := `SELECT status, count
		 FROM public."order" WHERE owner_id=$1 AND id=$2`

	orderRow := tx.QueryRow(ctx, SQLGetOrderByID, userID, orderID)

	var status uint8

	var count uint32

	err := orderRow.Scan(&status, &count)
	if err != nil {
		logger.Errorln(err)

		return models.OrderStatusError, 0, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return status, count, nil
}

func (p *ProductStorage) decreaseAvailableCountByOrderID(ctx context.Context,
	tx pgx.Tx, orderID uint64, count uint32,
) error {
	logger := p.logger.LogReqID(ctx)

	SQLDecreaseAvailableCountByOrderID := `UPDATE public."product"
		 SET available_count = available_count - $1
		 WHERE id = (
			SELECT product_id
			FROM public."order"
			WHERE id = $2
		 )`

	result, err := tx.Exec(ctx, SQLDecreaseAvailableCountByOrderID, count, orderID)
	if err != nil {
		logger.Errorln(err)

		return fmt.Errorf(myerrors.ErrTemplate, err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf(myerrors.ErrTemplate, ErrNoAffectedOrderRows)
	}

	return nil
}

func (p *ProductStorage) updateOrderStatus(ctx context.Context,
	tx pgx.Tx, userID uint64, orderID uint64, newStatus uint8,
) error {
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
}

func (p *ProductStorage) UpdateOrderStatus(ctx context.Context,
	userID uint64, orderID uint64, newStatus uint8,
) error {
	err := pgx.BeginFunc(ctx, p.pool, func(tx pgx.Tx) error {
		err := p.updateOrderStatus(ctx, tx, userID, orderID, newStatus)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return nil
}

func (p *ProductStorage) insertOrder(ctx context.Context, tx pgx.Tx,
	userID uint64, productID uint64, count uint32,
) error {
	logger := p.logger.LogReqID(ctx)

	SQLInsertOrder := `INSERT INTO public."order"(owner_id, product_id, count) VALUES ($1, $2, $3)`

	_, err := tx.Exec(ctx, SQLInsertOrder, userID, productID, count)
	if err != nil {
		logger.Errorln(err)

		return fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return nil
}

func (p *ProductStorage) AddOrderInBasket(ctx context.Context,
	userID uint64, productID uint64, count uint32,
) (*models.OrderInBasket, error) {
	logger := p.logger.LogReqID(ctx)

	orderInBasket := new(models.OrderInBasket)

	err := pgx.BeginFunc(ctx, p.pool, func(tx pgx.Tx) error {
		productInner, err := p.selectProductByID(ctx, tx, productID)
		if err != nil {
			return fmt.Errorf(myerrors.ErrTemplate, err)
		}

		if productInner.AvailableCount < count {
			return ErrAvailableCountNotEnough
		}

		err = p.insertOrder(ctx, tx, userID, productID, count)
		if err != nil {
			return fmt.Errorf(myerrors.ErrTemplate, err)
		}

		idOrder, err := repository.GetLastValSeq(ctx, tx, logger, NameSeqOrder)
		if err != nil {
			return fmt.Errorf(myerrors.ErrTemplate, err)
		}

		orderInBasket.ID = idOrder
		orderInBasket.OwnerID = userID
		orderInBasket.ProductID = productID
		orderInBasket.Count = count
		orderInBasket.SalerID = productInner.SalerID
		orderInBasket.Title = productInner.Title
		orderInBasket.Price = productInner.Price
		orderInBasket.CityID = productInner.CityID
		orderInBasket.AvailableCount = productInner.AvailableCount
		orderInBasket.Delivery = productInner.Delivery
		orderInBasket.SafeDeal = productInner.SafeDeal
		orderInBasket.InFavourites = productInner.InFavourites
		orderInBasket.Images = productInner.Images

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return orderInBasket, nil
}

func (p *ProductStorage) updateStatusFullBasket(ctx context.Context, tx pgx.Tx, userID uint64) error {
	logger := p.logger.LogReqID(ctx)

	SQLSelectFullBasket := `SELECT id FROM public."order" WHERE owner_id=$1 AND status=0`

	rows, err := tx.Query(ctx, SQLSelectFullBasket, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrNotFoundOrdersInBasket
		}

		logger.Errorln(err)

		return fmt.Errorf(myerrors.ErrTemplate, err)
	}

	var orderID uint64

	var slOrderID []uint64

	_, err = pgx.ForEachRow(rows, []any{&orderID}, func() error {
		slOrderID = append(slOrderID, orderID)

		return nil
	})
	if err != nil {
		logger.Errorln(err)

		return fmt.Errorf(myerrors.ErrTemplate, err)
	}

	for _, val := range slOrderID {
		err = p.updateOrderStatus(ctx, tx, userID, val, models.OrderStatusInProcessing)
		if err != nil {
			logger.Errorln(err)

			return fmt.Errorf(myerrors.ErrTemplate, err)
		}
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
		return fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return nil
}

func (p *ProductStorage) deleteOrderByOrderIDAndOwnerID(ctx context.Context,
	tx pgx.Tx, orderID uint64, ownerID uint64,
) error {
	logger := p.logger.LogReqID(ctx)

	SQLDeleteOrderByID := `DELETE FROM public."order"
		 WHERE id=$1 AND owner_id=$2`

	result, err := tx.Exec(ctx, SQLDeleteOrderByID, orderID, ownerID)
	if err != nil {
		logger.Errorln(err)

		return fmt.Errorf(myerrors.ErrTemplate, err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf(myerrors.ErrTemplate, ErrNoAffectedOrderRows)
	}

	return nil
}

func (p *ProductStorage) DeleteOrder(ctx context.Context, orderID uint64, ownerID uint64) error {
	err := pgx.BeginFunc(ctx, p.pool, func(tx pgx.Tx) error {
		err := p.deleteOrderByOrderIDAndOwnerID(ctx, tx, orderID, ownerID)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return nil
}
