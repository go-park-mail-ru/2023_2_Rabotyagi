package usecases

import (
	"log"
	"net/http"

	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/jwt"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/server/delivery"
)

// GetUserIDFromCookie return 0 if error happen and return userID if success
func GetUserIDFromCookie(r *http.Request) uint64 {
	cookie, err := r.Cookie(delivery.CookieAuthName)
	if err != nil {
		log.Printf("in getUserIDFromCookie: %+v\n", err)

		return 0
	}

	rawJwt := cookie.Value

	userPayload, err := jwt.NewUserJwtPayload(rawJwt, jwt.Secret)
	if err != nil {
		log.Printf("in getUserIDFromCookie: %+v\n", err)

		return 0
	}

	return userPayload.UserID
}
