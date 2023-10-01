package authorization

import (
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/storage"
)

func GenerateJwtToken(user *storage.User, secret []byte) (string, error) {
    token := jwt.New(jwt.SigningMethodHS256)

    claims := token.Claims.(jwt.MapClaims)
    claims["userID"] = user.ID
    claims["email"] = user.Email
    claims["exp"] = time.Now().Add(time.Hour * 24).Unix() // Токен будет действителен в течение 24 часов

    // Подписываем токен используя переданный секретный ключ
    tokenString, err := token.SignedString(secret)
    if err != nil {
        return "", err
    }

    return tokenString, nil
}