package repository

import (
	"context"
	"fmt"
	"log"

	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/models"
	myerrors "github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductStorage struct {
	pool *pgxpool.Pool
}

func NewProductStorage(pool *pgxpool.Pool) *ProductStorage {
	return &ProductStorage{
		pool: pool,
	}
}

func (p *ProductStorage) GetProduct(ctx context.Context, productID uint64) (*models.Product, error) {
	return nil, nil
}

func (p *ProductStorage) GetNProducts(ctx context.Context) ([]*models.Product, error) {
	return nil, nil
}

func (p *ProductStorage) insertProduct(ctx context.Context, tx pgx.Tx, preProduct *models.PreProduct) error {
	SQLInsertProduct := `INSERT INTO public."product"(saler_id,
		category_id, title, description, price,available_count,
		city, delivery, safe_deal) VALUES(
		$1, $2, $3, $4, $5, $6, $7, $8, $9)`
	_, err := tx.Exec(ctx, SQLInsertProduct, preProduct.SalerID, preProduct.CategoryID,
		preProduct.Title, preProduct.Description, preProduct.Price, preProduct.AvailableCount,
		preProduct.City, preProduct.Delivery, preProduct.SafeDeal)

	if err != nil {
		log.Printf("in insertProduct: %+v\n", err)

		return fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return nil
}

func (p *ProductStorage) AddProduct(ctx context.Context, preProduct *models.PreProduct) error {
	err := pgx.BeginFunc(ctx, p.pool, func(tx pgx.Tx) error {
		err := p.insertProduct(ctx, tx, preProduct)

		return err
	})
	if err != nil {
		log.Printf("in AddProduct: %+v\n", err)

		return fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return nil
}
