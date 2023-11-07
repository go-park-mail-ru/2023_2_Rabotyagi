package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/models"
	myerrors "github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/errors"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/server/repository"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

var (
	ErrProductNotFound       = myerrors.NewError("Это объявление не найдено")
	ErrNoUpdateFields        = myerrors.NewError("Вы пытаетесь обновить пустое количество полей объявления")
	ErrNoAffectedProductRows = myerrors.NewError("Не получилось обновить данные товара")

	NameSeqProduct = pgx.Identifier{"public", "product_id_seq"} //nolint:gochecknoglobals
)

type ProductStorage struct {
	pool   *pgxpool.Pool
	logger *zap.SugaredLogger
}

func NewProductStorage(pool *pgxpool.Pool, logger *zap.SugaredLogger) *ProductStorage {
	return &ProductStorage{
		pool:   pool,
		logger: logger,
	}
}

func (p *ProductStorage) selectImagesByProductID(ctx context.Context,
	tx pgx.Tx, productID uint64,
) ([]models.Image, error) {
	var images []models.Image

	SQLSelectImages := `SELECT url FROM public."image" WHERE product_id=$1`

	imagesRows, err := tx.Query(ctx, SQLSelectImages, productID)
	if err != nil {
		p.logger.Errorf("in selectImagesByProductId: %+v\n", err)

		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	var curURL string

	_, err = pgx.ForEachRow(imagesRows, []any{&curURL}, func() error {
		images = append(images, models.Image{URL: curURL})

		return nil
	})
	if err != nil {
		p.logger.Errorf("in selectImagesByProductId: %+v\n", err)

		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return images, nil
}

func (p *ProductStorage) selectProductByIDAndSalerID(ctx context.Context,
	tx pgx.Tx, productID uint64, userID uint64,
) (*models.Product, error) {
	SQLSelectProduct := `SELECT saler_id, category_id, title,
       description, price, created_at, views, available_count, city,
       delivery, safe_deal FROM public."product" WHERE id=$1 AND saler_id=$2`
	product := &models.Product{ID: productID} //nolint:exhaustruct

	productRow := tx.QueryRow(ctx, SQLSelectProduct, productID, userID)
	if err := productRow.Scan(&product.SalerID, &product.CategoryID,
		&product.Title, &product.Description, &product.Price, &product.CreatedAt,
		&product.Views, &product.AvailableCount, &product.City, &product.Delivery,
		&product.SafeDeal); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf(myerrors.ErrTemplate, ErrProductNotFound)
		}

		p.logger.Errorf("error in selectProductById with productId=%d: %+v", productID, err)

		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	product.ID = productID
	product.SalerID = userID

	return product, nil
}

func (p *ProductStorage) selectCountFavouritesByProductID(ctx context.Context,
	tx pgx.Tx,
	productID uint64,
) (uint64, error) {
	var favouritesCount uint64

	SQLCountFavourites := `SELECT COUNT(id) FROM public."favourite" WHERE product_id=$1`

	CountFavouritesRow := tx.QueryRow(ctx, SQLCountFavourites, productID)
	if err := CountFavouritesRow.Scan(&favouritesCount); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, fmt.Errorf(myerrors.ErrTemplate, ErrProductNotFound)
		}

		p.logger.Errorf("in selectCountFavouritesByProductID: %+v\n", err)

		return 0, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return favouritesCount, nil
}

func (p *ProductStorage) selectIsUserFavouriteProduct(ctx context.Context,
	tx pgx.Tx, productID uint64,
	userID uint64,
) (bool, error) {
	var rawRow string

	SQLSelectIsUserFavouriteProduct := `SELECT id FROM public.favourite WHERE product_id=$1 AND owner_id=$2`

	isUserFavouriteRow := tx.QueryRow(ctx, SQLSelectIsUserFavouriteProduct, productID, userID)
	if err := isUserFavouriteRow.Scan(&rawRow); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}

		p.logger.Errorf("in selectIsUserFavouriteProduct: %+v\n", err)

		return false, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return true, nil
}

type productAddition struct {
	favourites  uint64
	images      []models.Image
	inFavourite bool
}

func (p *ProductStorage) getProductAddition(ctx context.Context,
	tx pgx.Tx, productID uint64, userID uint64,
) (*productAddition, error) {
	innerProductAddition := new(productAddition)

	images, err := p.selectImagesByProductID(ctx, tx, productID)
	if err != nil {
		p.logger.Errorf("in getProductAddition: %+v\n", err)

		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	favouritesCount, err := p.selectCountFavouritesByProductID(ctx, tx, productID)
	if err != nil {
		p.logger.Errorf("in getProductAddition: %+v\n", err)

		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	inFavouriteProduct, err := p.selectIsUserFavouriteProduct(ctx, tx, productID, userID)
	if err != nil {
		p.logger.Errorf("in getProductAddition: %+v\n", err)

		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	innerProductAddition.images = images
	innerProductAddition.favourites = favouritesCount
	innerProductAddition.inFavourite = inFavouriteProduct

	return innerProductAddition, nil
}

func (p *ProductStorage) GetProduct(ctx context.Context, productID uint64, userID uint64) (*models.Product, error) {
	var product *models.Product

	err := pgx.BeginFunc(ctx, p.pool, func(tx pgx.Tx) error {
		productInner, err := p.selectProductByIDAndSalerID(ctx, tx, productID, userID)
		if err != nil {
			return err
		}

		productAdditionInner, err := p.getProductAddition(ctx, tx, productID, userID)
		if err != nil {
			return err
		}

		product = productInner
		product.Images = productAdditionInner.images
		product.Favourites = productAdditionInner.favourites
		product.InFavourites = productAdditionInner.inFavourite

		return nil
	})
	if err != nil {
		p.logger.Errorf("in GetProduct: %+v\n", err)

		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return product, nil
}

// selectProductsInFeedWithWhereOrderLimit accepts arguments in the appropriate format:
//
// whereClause can be:
// nil - ignored.
//
// map[string]interface{} OR squirrel.Eq - map of SQL expressions to values. Each key is transformed into
// an expression like "<key> = ?", with the corresponding value bound to the placeholder. If the value is
// nil, the expression will be "<key> IS NULL". If the value is an array or slice, the expression will be
// "<key> IN (?,?,...)", with one placeholder for each item in the value.
//
// orderByClause:
// nil - ignored.
//
// another cases add ORDER BY expressions to the query
//
// limit sets a LIMIT clause on the query.
func (p *ProductStorage) selectProductsInFeedWithWhereOrderLimit(ctx context.Context, tx pgx.Tx,
	limit uint64, whereClause any, orderByClause []string,
) ([]*models.ProductInFeed, error) {
	query := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).Select("id, title," +
		"price, city, delivery, safe_deal").From(`public."product"`).
		Where(whereClause).OrderBy(orderByClause...).Limit(limit)

	SQLQuery, args, err := query.ToSql()
	if err != nil {
		p.logger.Errorf("in selectProductsInFeedWithWhereOrderLimit: %+v\n", err)

		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	rowsProducts, err := tx.Query(ctx, SQLQuery, args...)
	if err != nil {
		p.logger.Errorf("in selectProductsInFeedWithWhereOrderLimit: %+v\n", err)

		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	curProduct := new(models.ProductInFeed)

	var slProduct []*models.ProductInFeed

	_, err = pgx.ForEachRow(rowsProducts, []any{
		&curProduct.ID, &curProduct.Title,
		&curProduct.Price, &curProduct.City,
		&curProduct.Delivery, &curProduct.SafeDeal,
	}, func() error {
		slProduct = append(slProduct, &models.ProductInFeed{ //nolint:exhaustruct
			ID:       curProduct.ID,
			Title:    curProduct.Title,
			Price:    curProduct.Price,
			City:     curProduct.City,
			Delivery: curProduct.Delivery,
			SafeDeal: curProduct.SafeDeal,
		})

		return nil
	})
	if err != nil {
		p.logger.Errorf("in selectProductsInFeedWithWhereOrderLimit: %+v\n", err)

		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return slProduct, nil
}

func (p *ProductStorage) GetNewProducts(ctx context.Context,
	lastProductID uint64, count uint64, userID uint64,
) ([]*models.ProductInFeed, error) {
	var slProduct []*models.ProductInFeed

	err := pgx.BeginFunc(ctx, p.pool, func(tx pgx.Tx) error {
		slProductInner, err := p.selectProductsInFeedWithWhereOrderLimit(ctx,
			tx, count, squirrel.Gt{"id": lastProductID}, []string{"created_at DESC"})
		if err != nil {
			return err
		}

		for _, product := range slProductInner {
			productAdditionInner, err := p.getProductAddition(ctx, tx, product.ID, userID)
			if err != nil {
				return err
			}

			product.Images = productAdditionInner.images
			product.Favourites = productAdditionInner.favourites
			product.InFavourites = productAdditionInner.inFavourite

			slProduct = append(slProduct, product)
		}

		return nil
	})
	if err != nil {
		p.logger.Errorf("in GetNewProducts: %+v\n", err)

		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return slProduct, nil
}

func (p *ProductStorage) GetProductsOfSaler(ctx context.Context,
	lastProductID uint64, count uint64, userID uint64,
) ([]*models.ProductInFeed, error) {
	var slProduct []*models.ProductInFeed

	err := pgx.BeginFunc(ctx, p.pool, func(tx pgx.Tx) error {
		whereClause := fmt.Sprintf("id > %d AND saler_id = %d", lastProductID, userID)

		slProductInner, err := p.selectProductsInFeedWithWhereOrderLimit(ctx,
			tx, count, whereClause, []string{"created_at DESC"})
		if err != nil {
			return err
		}

		for _, product := range slProductInner {
			productAdditionInner, err := p.getProductAddition(ctx, tx, product.ID, userID)
			if err != nil {
				return err
			}

			product.Images = productAdditionInner.images
			product.Favourites = productAdditionInner.favourites
			product.InFavourites = productAdditionInner.inFavourite

			slProduct = append(slProduct, product)
		}

		return nil
	})
	if err != nil {
		p.logger.Errorf("in GetProductsOfSaler: %+v\n", err)

		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return slProduct, nil
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
		p.logger.Errorf("in insertProduct: %+v\n", err)

		return fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return nil
}

func (p *ProductStorage) AddProduct(ctx context.Context, preProduct *models.PreProduct) (uint64, error) {
	var productID uint64

	err := pgx.BeginFunc(ctx, p.pool, func(tx pgx.Tx) error {
		err := p.insertProduct(ctx, tx, preProduct)
		if err != nil {
			return err
		}

		LastProductID, err := repository.GetLastValSeq(ctx, tx, p.logger, NameSeqProduct)
		if err != nil {
			return err
		}

		productID = LastProductID

		return err
	})
	if err != nil {
		p.logger.Errorf("in AddProduct: %+v\n", err)

		return 0, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return productID, nil
}

func (p *ProductStorage) updateProduct(ctx context.Context, tx pgx.Tx,
	productID uint64, updateFields map[string]interface{},
) error {
	if len(updateFields) == 0 {
		return ErrNoUpdateFields
	}

	query := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).Update(`public."product"`).
		Where(squirrel.Eq{"id": productID}).SetMap(updateFields)

	queryString, args, err := query.ToSql()
	if err != nil {
		p.logger.Errorf("in updateProduct: %+v", err)

		return fmt.Errorf(myerrors.ErrTemplate, err)
	}

	result, err := tx.Exec(ctx, queryString, args...)
	if err != nil {
		p.logger.Errorf("updateProduct: %+v", err)

		return fmt.Errorf(myerrors.ErrTemplate, err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf(myerrors.ErrTemplate, ErrNoAffectedProductRows)
	}

	return nil
}

func (p *ProductStorage) UpdateProduct(ctx context.Context, productID uint64,
	updateFields map[string]interface{},
) error {
	err := pgx.BeginFunc(ctx, p.pool, func(tx pgx.Tx) error {
		err := p.updateProduct(ctx, tx, productID, updateFields)

		return err
	})
	if err != nil {
		p.logger.Errorf("in UpdateProduct: %+v\n", err)

		return fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return nil
}

func (p *ProductStorage) closeProduct(ctx context.Context, tx pgx.Tx, productID uint64, userID uint64) error {
	SQLCloseProduct := `UPDATE public."product" SET available_count=0, is_active=false WHERE id=$1 AND saler_id=$2`

	result, err := tx.Exec(ctx, SQLCloseProduct, productID, userID)
	if err != nil {
		p.logger.Errorf("in closeProduct: %+v\n", err)

		return fmt.Errorf(myerrors.ErrTemplate, err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf(myerrors.ErrTemplate, ErrNoAffectedProductRows)
	}

	return nil
}

func (p *ProductStorage) CloseProduct(ctx context.Context, productID uint64, userID uint64) error {
	err := pgx.BeginFunc(ctx, p.pool, func(tx pgx.Tx) error {
		err := p.closeProduct(ctx, tx, productID, userID)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		p.logger.Errorf("in CloseProduct: %+v\n", err)

		return fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return nil
}

func (p *ProductStorage) deleteProduct(ctx context.Context, tx pgx.Tx, productID uint64, userID uint64) error {
	SQLCloseProduct := `DELETE FROM public."product" WHERE id=$1 AND saler_id=$2`

	result, err := tx.Exec(ctx, SQLCloseProduct, productID, userID)
	if err != nil {
		p.logger.Errorf("in deleteProduct: %+v\n", err)

		return fmt.Errorf(myerrors.ErrTemplate, err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf(myerrors.ErrTemplate, ErrNoAffectedProductRows)
	}

	return nil
}

func (p *ProductStorage) DeleteProduct(ctx context.Context, productID uint64, userID uint64) error {
	err := pgx.BeginFunc(ctx, p.pool, func(tx pgx.Tx) error {
		err := p.deleteProduct(ctx, tx, productID, userID)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		p.logger.Errorf("in DeleteProduct: %+v\n", err)

		return fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return nil
}
