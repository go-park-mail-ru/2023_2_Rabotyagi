package authorization

import (
	"time"

	rabotyagi "github.com/go-park-mail-ru/2023_2_Rabotyagi"
	jwt "github.com/dgrijalva/jwt-go"
)


func GenerateJwtToken(user *rabotyagi.User, secret []byte) (string, error) {
    token := jwt.New(jwt.SigningMethodHS256)

    claims := token.Claims.(jwt.MapClaims)
    //claims["userID"] = user.Id
    claims["username"] = user.Name
    claims["exp"] = time.Now().Add(time.Hour * 24).Unix() // Токен будет действителен в течение 24 часов

    // Подписываем токен используя переданный секретный ключ
    tokenString, err := token.SignedString(secret)
    if err != nil {
        return "", err
    }

    return tokenString, nil
}