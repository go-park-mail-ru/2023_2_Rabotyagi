package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/models"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/myerrors"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/mylogger"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/pgxpool"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/repository"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/responses/statuses"
	"github.com/jackc/pgx/v5"
)

var (
	ErrProductNotFound = myerrors.NewErrorBadFormatRequest("Это объявление не найдено")
	ErrNoUpdateFields  = myerrors.NewErrorBadFormatRequest(
		"Вы пытаетесь обновить пустое количество полей объявления")
	ErrNoAffectedProductRows      = myerrors.NewErrorBadFormatRequest("Не получилось обновить данные товара")
	ErrGetUncorrectedFormatImages = myerrors.NewErrorBadFormatRequest(
		"Получили некорректный формат images внутри объявления")
	ErrUncorrectedPrice = myerrors.NewErrorInternal(
		"Получили некорректный тип price")
	ErrScanCommentID = myerrors.NewErrorInternal("Ошибка сканирования comment_id")

	NameSeqProduct = pgx.Identifier{"public", "product_id_seq"} //nolint:gochecknoglobals
)

const (
	PremiumCoefficient    = uint16(5)
	NonPremiumCoefficient = uint16(1)
	SoldByUserCoefficient = uint16(3)
	ViewsCoefficient      = uint16(2)
)

type ProductStorage struct {
	pool   pgxpool.IPgxPool
	logger *mylogger.MyLogger
}

func NewProductStorage(pool pgxpool.IPgxPool) (*ProductStorage, error) {
	logger, err := mylogger.Get()
	if err != nil {
		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return &ProductStorage{
		pool:   pool,
		logger: logger,
	}, nil
}

func (p *ProductStorage) selectImagesByProductID(ctx context.Context,
	tx pgx.Tx, productID uint64,
) ([]models.Image, error) {
	logger := p.logger.LogReqID(ctx)

	var images []models.Image

	SQLSelectImages := `SELECT url FROM public."image" WHERE product_id=$1`

	imagesRows, err := tx.Query(ctx, SQLSelectImages, productID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return images, nil
		}

		logger.Errorln(err)

		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	var curURL string

	_, err = pgx.ForEachRow(imagesRows, []any{&curURL}, func() error {
		images = append(images, models.Image{URL: curURL})

		return nil
	})
	if err != nil {
		logger.Errorln(err)

		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return images, nil
}

func (p *ProductStorage) selectProductByID(ctx context.Context,
	tx pgx.Tx, productID uint64,
) (*models.Product, error) {
	logger := p.logger.LogReqID(ctx)

	var premiumStatus uint8

	SQLSelectProduct := `SELECT saler_id, category_id, title,
       description, price, created_at, views, available_count, city_id,
       delivery, safe_deal, is_active, premium_status FROM public."product" WHERE id=$1`
	product := &models.Product{ID: productID} //nolint:exhaustruct

	productRow := tx.QueryRow(ctx, SQLSelectProduct, productID)
	if err := productRow.Scan(&product.SalerID, &product.CategoryID,
		&product.Title, &product.Description, &product.Price, &product.CreatedAt,
		&product.Views, &product.AvailableCount, &product.CityID, &product.Delivery,
		&product.SafeDeal, &product.IsActive, &premiumStatus); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf(myerrors.ErrTemplate, ErrProductNotFound)
		}

		logger.Errorf("error with productId=%d: %+v", productID, err)

		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	product.ID = productID
	product.Premium = statuses.IsIntStatusPremiumSuccessful(premiumStatus)

	return product, nil
}

func (p *ProductStorage) selectCountFavouritesByProductID(ctx context.Context,
	tx pgx.Tx,
	productID uint64,
) (uint64, error) {
	logger := p.logger.LogReqID(ctx)

	var favouritesCount uint64

	SQLCountFavourites := `SELECT COUNT(id) FROM public."favourite" WHERE product_id=$1`

	CountFavouritesRow := tx.QueryRow(ctx, SQLCountFavourites, productID)
	if err := CountFavouritesRow.Scan(&favouritesCount); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, nil
		}

		logger.Errorln(err)

		return 0, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return favouritesCount, nil
}

func (p *ProductStorage) selectIsUserFavouriteProduct(ctx context.Context,
	tx pgx.Tx, productID uint64,
	userID uint64,
) (bool, error) {
	logger := p.logger.LogReqID(ctx)

	var rawRow string

	SQLSelectIsUserFavouriteProduct := `SELECT id FROM public.favourite WHERE product_id=$1 AND owner_id=$2`

	isUserFavouriteRow := tx.QueryRow(ctx, SQLSelectIsUserFavouriteProduct, productID, userID)
	if err := isUserFavouriteRow.Scan(&rawRow); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}

		logger.Errorln(err)

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
		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	favouritesCount, err := p.selectCountFavouritesByProductID(ctx, tx, productID)
	if err != nil {
		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	inFavouriteProduct, err := p.selectIsUserFavouriteProduct(ctx, tx, productID, userID)
	if err != nil {
		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	innerProductAddition.images = images
	innerProductAddition.favourites = favouritesCount
	innerProductAddition.inFavourite = inFavouriteProduct

	return innerProductAddition, nil
}

func (p *ProductStorage) selectPremiumExpireByProductID(ctx context.Context,
	tx pgx.Tx,
	productID uint64,
) (sql.NullTime, error) {
	logger := p.logger.LogReqID(ctx)

	var expire sql.NullTime

	SQLPremiumEpire := `SELECT premium_expire FROM public."product" WHERE id=$1`

	expireRow := tx.QueryRow(ctx, SQLPremiumEpire, productID)
	if err := expireRow.Scan(&expire); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return sql.NullTime{}, nil //nolint:exhaustruct
		}

		logger.Errorln(err)

		return sql.NullTime{}, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return expire, nil
}

func (p *ProductStorage) selectCommentID(ctx context.Context,
	tx pgx.Tx, userID uint64, salerID uint64,
) (sql.NullInt64, error) {
	logger := p.logger.LogReqID(ctx)

	var commentID uint64

	SQLSelectCommentID := `SELECT id FROM public."comment" WHERE sender_id=$1 AND recipient_id=$2`

	commentIDRow := tx.QueryRow(ctx, SQLSelectCommentID, userID, salerID)
	if err := commentIDRow.Scan(&commentID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return sql.NullInt64{Valid: false, Int64: 0}, nil
		}

		err = fmt.Errorf("%w %v", ErrScanCommentID, err) //nolint:errorlint
		logger.Errorln(err)

		return sql.NullInt64{Valid: false, Int64: 0}, err
	}

	return sql.NullInt64{Valid: true, Int64: int64(commentID)}, nil
}

func (p *ProductStorage) getProduct(ctx context.Context,
	tx pgx.Tx, productID uint64, userID uint64,
) (*models.Product, error) {
	product, err := p.selectProductByID(ctx, tx, productID)
	if err != nil {
		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	if product.SalerID == userID && product.Premium {
		product.PremiumExpire, err = p.selectPremiumExpireByProductID(ctx, tx, productID)
		if err != nil {
			return nil, fmt.Errorf(myerrors.ErrTemplate, err)
		}
	}

	productAdditionInner, err := p.getProductAddition(ctx, tx, productID, userID)
	if err != nil {
		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	productPriceHistory, err := p.selectPriceHistory(ctx, tx, productID)
	if err != nil {
		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	productCommentID, err := p.selectCommentID(ctx, tx, userID, product.SalerID)
	if err != nil {
		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	product.PriceHistory = productPriceHistory
	product.Images = productAdditionInner.images
	product.Favourites = productAdditionInner.favourites
	product.InFavourites = productAdditionInner.inFavourite
	product.CommentID = productCommentID

	return product, nil
}

func (p *ProductStorage) GetProduct(ctx context.Context, productID uint64, userID uint64) (*models.Product, error) {
	var product *models.Product

	err := pgx.BeginFunc(ctx, p.pool, func(tx pgx.Tx) error {
		productInner, err := p.getProduct(ctx, tx, productID, userID)
		if err != nil {
			return err
		}

		viewExist, err := p.viewExist(ctx, tx, userID, productID)
		if err != nil {
			return err
		}

		if !viewExist && userID != 0 {
			err = p.addView(ctx, tx, userID, productID)
			if err != nil {
				return err
			}

			err = p.incViews(ctx, tx, productID)
			if err != nil {
				return err
			}
		}

		product = productInner

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return product, nil
}

// selectProductsInFeedWithWhereOrderLimitOffset accepts arguments in the appropriate format:
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
func (p *ProductStorage) selectProductsInFeedWithWhereOrderLimitOffset(ctx context.Context, tx pgx.Tx,
	limit uint64, whereClause any, orderByClause []string, offset uint64,
) ([]*models.ProductInFeed, error) {
	logger := p.logger.LogReqID(ctx)

	query := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).Select("id, title," +
		"price, city_id, delivery, safe_deal, is_active, available_count, premium_status").From(`public."product"`).
		Where(whereClause).OrderBy(orderByClause...).Limit(limit).Offset(offset)

	SQLQuery, args, err := query.ToSql()
	if err != nil {
		logger.Errorln(err)

		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	rowsProducts, err := tx.Query(ctx, SQLQuery, args...)
	if err != nil {
		logger.Errorln(err)

		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	curProduct := new(models.ProductInFeed)

	var slProduct []*models.ProductInFeed

	var premiumStatus uint8

	_, err = pgx.ForEachRow(rowsProducts, []any{
		&curProduct.ID, &curProduct.Title,
		&curProduct.Price, &curProduct.CityID,
		&curProduct.Delivery, &curProduct.SafeDeal, &curProduct.IsActive, &curProduct.AvailableCount, &premiumStatus,
	}, func() error {
		slProduct = append(slProduct, &models.ProductInFeed{ //nolint:exhaustruct
			ID:             curProduct.ID,
			Title:          curProduct.Title,
			Price:          curProduct.Price,
			CityID:         curProduct.CityID,
			Delivery:       curProduct.Delivery,
			SafeDeal:       curProduct.SafeDeal,
			IsActive:       curProduct.IsActive,
			AvailableCount: curProduct.AvailableCount,
			Premium:        statuses.IsIntStatusPremiumSuccessful(premiumStatus),
		})

		return nil
	})
	if err != nil {
		logger.Errorln(err)

		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return slProduct, nil
}

func (p *ProductStorage) GetPopularProducts(ctx context.Context,
	offset uint64, count uint64, userID uint64,
) ([]*models.ProductInFeed, error) {
	logger := p.logger.LogReqID(ctx)

	var slProduct []*models.ProductInFeed

	err := pgx.BeginFunc(ctx, p.pool, func(tx pgx.Tx) error {
		whereClause := "is_active = true"

		slProductInner, err := p.selectProductsInFeedWithWhereOrderLimitOffset(ctx,
			tx, count, whereClause, []string{OrderByClauseForProductList(PremiumCoefficient,
				NonPremiumCoefficient, SoldByUserCoefficient, ViewsCoefficient)}, offset)
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
		logger.Errorln(err)

		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return slProduct, nil
}

func (p *ProductStorage) GetProductsOfSaler(ctx context.Context,
	offset uint64, count uint64, userID uint64, isMy bool,
) ([]*models.ProductInFeed, error) {
	logger := p.logger.LogReqID(ctx)

	var slProduct []*models.ProductInFeed

	err := pgx.BeginFunc(ctx, p.pool, func(tx pgx.Tx) error {
		var whereClause string
		if isMy {
			whereClause = fmt.Sprintf("saler_id = %d", userID)
		} else {
			whereClause = fmt.Sprintf("saler_id = %d AND is_active = true",
				userID)
		}

		slProductInner, err := p.selectProductsInFeedWithWhereOrderLimitOffset(ctx,
			tx, count, whereClause, []string{OrderByClauseForProductList(PremiumCoefficient,
				NonPremiumCoefficient, SoldByUserCoefficient, ViewsCoefficient)}, offset)
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
		logger.Errorln(err)

		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return slProduct, nil
}

func (p *ProductStorage) deleteAllImagesOfProduct(ctx context.Context, tx pgx.Tx, productID uint64) error {
	logger := p.logger.LogReqID(ctx)

	SQLDeleteImage := `DELETE FROM public."image" WHERE product_id=$1;`

	_, err := tx.Exec(ctx, SQLDeleteImage, productID)
	if err != nil {
		logger.Errorln(err)

		return fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return nil
}

func (p *ProductStorage) insertImages(ctx context.Context, tx pgx.Tx, productID uint64, slImg []models.Image) error {
	logger := p.logger.LogReqID(ctx)

	SQLInsertImage := `INSERT INTO public."image" (url, product_id) VALUES(@imgURL, @productID)`
	batch := &pgx.Batch{}

	for _, image := range slImg {
		args := pgx.NamedArgs{
			"imgURL":    image.URL,
			"productID": productID,
		}
		batch.Queue(SQLInsertImage, args)
	}

	results := tx.SendBatch(ctx, batch)
	defer results.Close()

	for _, image := range slImg {
		_, err := results.Exec()
		if err != nil {
			logger.Errorf("with product_id=%d with URL=%s %+v", productID, image.URL, err)

			return fmt.Errorf(myerrors.ErrTemplate, err)
		}
	}

	return results.Close() //nolint:wrapcheck
}

func (p *ProductStorage) insertProduct(ctx context.Context, tx pgx.Tx, preProduct *models.PreProduct) error {
	logger := p.logger.LogReqID(ctx)

	SQLInsertProduct := `INSERT INTO public."product"(saler_id,
		category_id, title, description, price,available_count,
		city_id, delivery, safe_deal) VALUES(
		$1, $2, $3, $4, $5, $6, $7, $8, $9)`

	_, err := tx.Exec(ctx, SQLInsertProduct, preProduct.SalerID, preProduct.CategoryID,
		preProduct.Title, preProduct.Description, preProduct.Price, preProduct.AvailableCount,
		preProduct.CityID, preProduct.Delivery, preProduct.SafeDeal)
	if err != nil {
		logger.Errorln(err)

		return fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return nil
}

func (p *ProductStorage) AddProduct(ctx context.Context, preProduct *models.PreProduct) (uint64, error) {
	logger := p.logger.LogReqID(ctx)

	var productID uint64

	err := pgx.BeginFunc(ctx, p.pool, func(tx pgx.Tx) error {
		err := p.insertProduct(ctx, tx, preProduct)
		if err != nil {
			return fmt.Errorf(myerrors.ErrTemplate, err)
		}

		LastProductID, err := repository.GetLastValSeq(ctx, tx, logger, NameSeqProduct)
		if err != nil {
			return fmt.Errorf(myerrors.ErrTemplate, err)
		}

		err = p.insertImages(ctx, tx, LastProductID, preProduct.Images)
		if err != nil {
			return fmt.Errorf(myerrors.ErrTemplate, err)
		}

		productID = LastProductID

		err = p.addPriceHistoryRecord(ctx, tx, productID, preProduct.Price)
		if err != nil {
			return fmt.Errorf(myerrors.ErrTemplate, err)
		}

		return nil
	})
	if err != nil {
		logger.Errorln(err)

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

	logger := p.logger.LogReqID(ctx)

	query := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).Update(`public."product"`).
		Where(squirrel.Eq{"id": productID}).SetMap(updateFields)

	queryString, args, err := query.ToSql()
	if err != nil {
		logger.Errorln(err)

		return fmt.Errorf(myerrors.ErrTemplate, err)
	}

	result, err := tx.Exec(ctx, queryString, args...)
	if err != nil {
		logger.Errorln(err)

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
	logger := p.logger.LogReqID(ctx)

	err := pgx.BeginFunc(ctx, p.pool, func(tx pgx.Tx) error {
		var err error

		updateImages, imagesExist := updateFields["images"]
		if imagesExist {
			delete(updateFields, "images")

			err = p.deleteAllImagesOfProduct(ctx, tx, productID)
			if err != nil {
				return err
			}

			slImages, ok := updateImages.([]models.Image)
			if !ok {
				errMessage := fmt.Errorf("%w product_id=%d", ErrGetUncorrectedFormatImages, productID)
				logger.Errorln(errMessage)

				return errMessage
			}

			err = p.insertImages(ctx, tx, productID, slImages)
			if err != nil {
				return err
			}
		}

		err = p.updateProduct(ctx, tx, productID, updateFields)
		if err != nil {
			return err
		}

		price, ok := updateFields["price"]
		if ok {
			priceUint64, ok := price.(uint64)
			if !ok {
				return ErrUncorrectedPrice
			}

			err = p.addPriceHistoryRecord(ctx, tx, productID, priceUint64)
			if err != nil {
				return err
			}
		}

		return err
	})
	if err != nil {
		logger.Errorln(err)

		return fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return nil
}

func (p *ProductStorage) closeProduct(ctx context.Context, tx pgx.Tx, productID uint64, userID uint64) error {
	logger := p.logger.LogReqID(ctx)

	SQLCloseProduct := `UPDATE public."product" SET is_active=false WHERE id=$1 AND saler_id=$2`

	result, err := tx.Exec(ctx, SQLCloseProduct, productID, userID)
	if err != nil {
		logger.Errorln(err)

		return fmt.Errorf(myerrors.ErrTemplate, err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf(myerrors.ErrTemplate, ErrNoAffectedProductRows)
	}

	return nil
}

func (p *ProductStorage) CloseProduct(ctx context.Context, productID uint64, userID uint64) error {
	logger := p.logger.LogReqID(ctx)

	err := pgx.BeginFunc(ctx, p.pool, func(tx pgx.Tx) error {
		err := p.closeProduct(ctx, tx, productID, userID)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		logger.Errorln(err)

		return fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return nil
}

func (p *ProductStorage) activateProduct(ctx context.Context, tx pgx.Tx, productID uint64, userID uint64) error {
	logger := p.logger.LogReqID(ctx)

	SQLActivateProduct := `UPDATE public."product" SET is_active=true WHERE id=$1 AND saler_id=$2`

	result, err := tx.Exec(ctx, SQLActivateProduct, productID, userID)
	if err != nil {
		logger.Errorln(err)

		return fmt.Errorf(myerrors.ErrTemplate, err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf(myerrors.ErrTemplate, ErrNoAffectedProductRows)
	}

	return nil
}

func (p *ProductStorage) ActivateProduct(ctx context.Context, productID uint64, userID uint64) error {
	logger := p.logger.LogReqID(ctx)

	err := pgx.BeginFunc(ctx, p.pool, func(tx pgx.Tx) error {
		err := p.activateProduct(ctx, tx, productID, userID)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		logger.Errorln(err)

		return fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return nil
}

func (p *ProductStorage) deleteProduct(ctx context.Context, tx pgx.Tx, productID uint64, userID uint64) error {
	logger := p.logger.LogReqID(ctx)

	SQLCloseProduct := `DELETE FROM public."product" WHERE id=$1 AND saler_id=$2`

	result, err := tx.Exec(ctx, SQLCloseProduct, productID, userID)
	if err != nil {
		logger.Errorln(err)

		return fmt.Errorf(myerrors.ErrTemplate, err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf(myerrors.ErrTemplate, ErrNoAffectedProductRows)
	}

	return nil
}

func (p *ProductStorage) DeleteProduct(ctx context.Context, productID uint64, userID uint64) error {
	logger := p.logger.LogReqID(ctx)

	err := pgx.BeginFunc(ctx, p.pool, func(tx pgx.Tx) error {
		err := p.deleteProduct(ctx, tx, productID, userID)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		logger.Errorln(err)

		return fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return nil
}

func (p *ProductStorage) viewExist(ctx context.Context, tx pgx.Tx, userID uint64, productID uint64) (bool, error) {
	logger := p.logger.LogReqID(ctx)

	SQLViewExist := `SELECT EXISTS(SELECT * FROM public."view" WHERE user_id = $1 AND product_id = $2);`

	exist := false

	existRow := tx.QueryRow(ctx, SQLViewExist, userID, productID)
	if err := existRow.Scan(&exist); err != nil {
		logger.Errorln(err)

		return false, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return exist, nil
}

func (p *ProductStorage) addView(ctx context.Context, tx pgx.Tx, userID uint64, productID uint64) error {
	logger := p.logger.LogReqID(ctx)

	SQLAddView := `INSERT INTO public."view" (user_id, product_id)
				   VALUES ($1, $2)`

	_, err := tx.Exec(ctx, SQLAddView, userID, productID)
	if err != nil {
		logger.Errorln(err)

		return fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return nil
}

func (p *ProductStorage) incViews(ctx context.Context, tx pgx.Tx, productID uint64) error {
	logger := p.logger.LogReqID(ctx)

	SQLAddView := `UPDATE public."product" 
				   SET views = views + 1 
				   WHERE id=$1`

	_, err := tx.Exec(ctx, SQLAddView, productID)
	if err != nil {
		logger.Errorln(err)

		return fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return nil
}

func (p *ProductStorage) searchProduct(ctx context.Context, tx pgx.Tx, searchInput string) ([]string, error) {
	logger := p.logger.LogReqID(ctx)

	SQLSearchProduct := `SELECT title
FROM (
  SELECT DISTINCT ON (title) title
  FROM product
  WHERE to_tsvector(title) @@ to_tsquery(replace($1 || ':*', ' ', ' | '))
    AND ts_rank(to_tsvector(title), to_tsquery(replace($1 || ':*', ' ', ' | '))) > 0
    AND is_active = true
) AS t
ORDER BY ts_rank(to_tsvector(title), to_tsquery(replace($1 || ':*', ' ', ' | '))) DESC
LIMIT 5;`

	var products []string

	productsRows, err := tx.Query(ctx, SQLSearchProduct, searchInput)
	if err != nil {
		logger.Errorln(err)

		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	var curProduct string

	_, err = pgx.ForEachRow(productsRows, []any{
		&curProduct,
	}, func() error {
		products = append(products, curProduct)

		return nil
	})
	if err != nil {
		logger.Errorln(err)

		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return products, nil
}

func (p *ProductStorage) SearchProduct(ctx context.Context, searchInput string) ([]string, error) {
	var products []string

	logger := p.logger.LogReqID(ctx)

	err := pgx.BeginFunc(ctx, p.pool, func(tx pgx.Tx) error {
		productsInner, err := p.searchProduct(ctx, tx, searchInput)
		if err != nil {
			return err
		}

		products = productsInner

		return nil
	})
	if err != nil {
		logger.Errorln(err)

		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return products, nil
}

func (p *ProductStorage) searchProductFeed(ctx context.Context, tx pgx.Tx,
	searchInput string, lastNumber uint64, limit uint64,
) ([]*models.ProductInFeed, error) {
	logger := p.logger.LogReqID(ctx)

	var premiumStatus uint8

	SQLSearchProduct := `SELECT id, title, price, city_id, delivery, safe_deal, is_active, available_count, premium_status
	FROM product
	WHERE (to_tsvector(title) @@ to_tsquery(replace($1 || ':*', ' ', ' | '))
	   OR to_tsvector(description) @@ to_tsquery(replace($1 || ':*', ' ', ' | ')))
	   AND is_active = true
	ORDER BY ts_rank(to_tsvector(title), to_tsquery(replace($1 || ':*', ' ', ' | '))) DESC,
			 ts_rank(to_tsvector(description), to_tsquery(replace($1 || ':*', ' ', ' | '))) DESC
	OFFSET $2
	LIMIT $3;`

	rowsProducts, err := tx.Query(ctx, SQLSearchProduct, searchInput, lastNumber, limit)
	if err != nil {
		logger.Errorln(err)

		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	curProduct := new(models.ProductInFeed)

	var slProduct []*models.ProductInFeed

	_, err = pgx.ForEachRow(rowsProducts, []any{
		&curProduct.ID, &curProduct.Title,
		&curProduct.Price, &curProduct.CityID,
		&curProduct.Delivery, &curProduct.SafeDeal, &curProduct.IsActive, &curProduct.AvailableCount, &premiumStatus,
	}, func() error {
		slProduct = append(slProduct, &models.ProductInFeed{ //nolint:exhaustruct
			ID:             curProduct.ID,
			Title:          curProduct.Title,
			Price:          curProduct.Price,
			CityID:         curProduct.CityID,
			Delivery:       curProduct.Delivery,
			SafeDeal:       curProduct.SafeDeal,
			IsActive:       curProduct.IsActive,
			AvailableCount: curProduct.AvailableCount,
			Premium:        statuses.IsIntStatusPremiumSuccessful(premiumStatus),
		})

		return nil
	})
	if err != nil {
		logger.Errorln(err)

		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return slProduct, nil
}

func (p *ProductStorage) GetSearchProductFeed(ctx context.Context,
	searchInput string, lastNumber uint64, limit uint64, userID uint64,
) ([]*models.ProductInFeed, error) {
	var slProduct []*models.ProductInFeed

	logger := p.logger.LogReqID(ctx)

	err := pgx.BeginFunc(ctx, p.pool, func(tx pgx.Tx) error {
		slProductInner, err := p.searchProductFeed(ctx,
			tx, searchInput, lastNumber, limit)
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
		logger.Errorln(err)

		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return slProduct, nil
}
