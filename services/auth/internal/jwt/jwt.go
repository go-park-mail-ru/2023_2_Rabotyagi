package jwt

import (
	"fmt"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/my_logger"
	"sync"
	"time"

	"github.com/golang-jwt/jwt"
)

const (
	lenSecret     = 64
	TimeTokenLife = 24 * time.Hour
)

var (
	globalSecret []byte = nil
	once         sync.Once
)

func SetSecret(secret []byte) {
	once.Do(func() {
		globalSecret = secret
	})
}

func GetSecret() []byte {
	var result []byte

	once.Do(func() {
		if globalSecret == nil {
			globalSecret = make([]byte, lenSecret)
		}

		result = globalSecret
	})

	return result
}

var (
	ErrNilToken           = fmt.Errorf("Получили токен = nil")
	ErrWrongSigningMethod = fmt.Errorf("Неожиданный signing метод ")
	ErrInvalidToken       = fmt.Errorf("Некорректный токен")
)

type UserJwtPayload struct {
	UserID uint64
	Expire int64
	Email  string
}

func NewUserJwtPayload(rawJwt string, secret []byte) (*UserJwtPayload, error) {
	logger, err := my_logger.Get()
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	tokenDuplicity, err := jwt.Parse(rawJwt, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			logger.Errorf("method == %+v %w", token.Header["alg"], ErrWrongSigningMethod)

			return nil, fmt.Errorf("%w", ErrInvalidToken)
		}

		return secret, nil
	})
	if err != nil {
		logger.Errorf("%s", err.Error())

		return nil, fmt.Errorf("%w", ErrInvalidToken)
	}

	if claims, ok := tokenDuplicity.Claims.(jwt.MapClaims); ok && tokenDuplicity.Valid {
		interfaceUserID, ok1 := claims["userID"]
		interfaceExpire, ok2 := claims["expire"]
		interfaceEmail, ok3 := claims["email"]

		if !(ok1 && ok2 && ok3) {
			logger.Errorf("error with claims: %+v", claims)

			return nil, fmt.Errorf("%w", ErrInvalidToken)
		}

		userID, ok1 := interfaceUserID.(float64)
		expire, ok2 := interfaceExpire.(float64)
		email, ok3 := interfaceEmail.(string)

		if !(ok1 && ok2 && ok3) {
			logger.Errorf("error with casting claims: %+v", claims)

			return nil, fmt.Errorf("%w", ErrInvalidToken)
		}

		return &UserJwtPayload{UserID: uint64(userID), Expire: int64(expire), Email: email}, nil
	}

	return nil, fmt.Errorf("%w", ErrInvalidToken)
}

func (u *UserJwtPayload) getMapClaims() jwt.MapClaims {
	result := make(jwt.MapClaims)

	result["userID"] = u.UserID
	result["expire"] = u.Expire
	result["email"] = u.Email

	return result
}

func GenerateJwtToken(userToken *UserJwtPayload, secret []byte) (string, error) {
	logger, err := my_logger.Get()
	if err != nil {
		return "", err //nolint:wrapcheck
	}

	if userToken == nil {
		logger.Errorln(ErrNilToken)

		return "", fmt.Errorf("%w", ErrInvalidToken)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, userToken.getMapClaims())

	tokenString, err := token.SignedString(secret)
	if err != nil {
		logger.Errorln(err)

		return "", fmt.Errorf("%w", ErrInvalidToken)
	}

	return tokenString, nil
}
