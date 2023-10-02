package jwt

import (
	"fmt"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/errors"
)

// TODO from config and reset her every some time
var secret = []byte("super-secret")

var (
	ErrNilToken           = errors.NewError("get nil token")
	ErrWrongSigningMethod = errors.NewError("unexpected signing method")
	ErrInvalidToken       = errors.NewError("invalid token")
)

type UserJwtPayload struct {
	UserID uint64
	Expire int64
	Email  string
}

func NewUserJwtPayload(rawJwt string) (*UserJwtPayload, error) {
	tokenDuplicity, err := jwt.Parse(rawJwt, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("method == %v %w", token.Header["alg"], ErrWrongSigningMethod)
		}

		return secret, nil
	})
	if err != nil {
		return nil, fmt.Errorf(errors.ErrTemplate, err)
	}

	if claims, ok := tokenDuplicity.Claims.(jwt.MapClaims); ok && tokenDuplicity.Valid {
		interfaceUserID, ok1 := claims["userID"]
		interfaceExpire, ok2 := claims["expire"]
		interfaceEmail, ok3 := claims["email"]

		if !(ok1 && ok2 && ok3) {
			return nil, fmt.Errorf("error with claims: %v %w", claims, ErrInvalidToken)
		}

		userID, ok1 := interfaceUserID.(uint64)
		expire, ok2 := interfaceExpire.(int64)
		email, ok3 := interfaceEmail.(string)

		if !(ok1 && ok2 && ok3) {
			return nil, fmt.Errorf("error with casting claims: %v %w", claims, ErrInvalidToken)
		}

		return &UserJwtPayload{UserID: userID, Expire: expire, Email: email}, nil
	}

	return nil, fmt.Errorf(errors.ErrTemplate, ErrInvalidToken)
}

func (u *UserJwtPayload) getMapClaims() jwt.MapClaims {
	result := make(jwt.MapClaims)

	result["userID"] = u.UserID
	result["expire"] = u.Expire
	result["email"] = u.Email

	return result
}

func GenerateJwtToken(userToken *UserJwtPayload) (string, error) {
	if userToken == nil {
		return "", fmt.Errorf(errors.ErrTemplate, ErrNilToken)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, userToken.getMapClaims())

	tokenString, err := token.SignedString(secret)

	if err != nil {
		return "", fmt.Errorf(errors.ErrTemplate, err)
	}

	return tokenString, nil
}