package delivery

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	myerrors "github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/errors"
)

var (
	ErrWrongCount     = myerrors.NewError("Получили некорректный count параметр. Он должен быть целым")
	ErrWrongLastID    = myerrors.NewError("Получили некорректный last_id параметр. Он должен быть целым")
	ErrWrongProductID = myerrors.NewError("Получили некорректный product_id параметр. Он должен быть целым")
)

func (p *ProductHandler) createURLToProductFromID(productID uint64) string {
	return fmt.Sprintf("%s%s:%s/api/v1/product/get/%d", p.schema, p.addrOrigin, p.portServer, productID)
}

func parseCountAndLastIDFromRequest(r *http.Request) (uint64, uint64, error) {
	countStr := r.URL.Query().Get("count")

	count, err := strconv.ParseUint(countStr, 10, 64)
	if err != nil {
		err := fmt.Errorf("%w count=%s", ErrWrongCount, countStr)

		log.Printf("in parseCountAndLastIDFromRequest: %+v\n", err)

		return 0, 0, err
	}

	lastIDStr := r.URL.Query().Get("last_id")

	lastID, err := strconv.ParseUint(lastIDStr, 10, 64)
	if err != nil {
		err := fmt.Errorf("%w last_id=%s", ErrWrongLastID, lastIDStr)
		log.Printf("in parseCountAndLastIDFromRequest: %+v\n", err)

		return 0, 0, err
	}

	return count, lastID, nil
}
