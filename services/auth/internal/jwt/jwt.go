package jwt

import (
	"crypto/rand"
	"fmt"
	"sync"
	"time"

	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/myerrors"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/mylogger"
	"github.com/golang-jwt/jwt"
)

const (
	lenSecret     = 64
	TimeTokenLife = 24 * time.Hour
)

var ErrGenerateSecret = myerrors.NewErrorInternal("Не получилось сгенерить секрет")

var (
	globalSecret []byte       //nolint:gochecknoglobals
	rwMu         sync.RWMutex //nolint:gochecknoglobals
)

func StartRefreshingSecret(period time.Duration, chClose <-chan struct{}) error {
	logger, err := mylogger.Get()
	if err != nil {
		return fmt.Errorf(myerrors.ErrTemplate, err)
	}

	go func() {
		for {
			select {
			case <-chClose:
				return
			default:
				time.Sleep(period)

				err = refreshSecret()
				if err != nil {
					logger.Errorln(err)
				}
			}
		}
	}()

	return nil
}

func SetSecret(secret []byte) {
	rwMu.Lock()
	globalSecret = secret
	rwMu.Unlock()
}

func refreshSecret() error {
	rwMu.Lock()
	globalSecret = make([]byte, lenSecret)

	logger, err := mylogger.Get()
	if err != nil {
		return fmt.Errorf(myerrors.ErrTemplate, err)
	}

	_, err = rand.Read(globalSecret)
	if err != nil {
		err = fmt.Errorf("%w %v", ErrGenerateSecret, err) //nolint:errorlint
		logger.Errorln(err)

		return err
	}
	rwMu.Unlock()

	return nil
}

// GetSecret return secret for jwt if it exists
// and additionally generate secret if not exist
func GetSecret() ([]byte, error) {
	var result []byte

	if globalSecret == nil {
		err := refreshSecret()
		if err != nil {
			return nil, fmt.Errorf(myerrors.ErrTemplate, err)
		}
	}

	rwMu.RLock()
	result = globalSecret
	rwMu.RUnlock()

	return result, nil
}

var (
	ErrNilToken           = myerrors.NewErrorInternal("Получили токен = nil")
	ErrWrongSigningMethod = myerrors.NewErrorBadFormatRequest("Неожиданный signing метод jwt токена")
	ErrInvalidToken       = myerrors.NewErrorBadFormatRequest("Некорректный токен")
)

type UserJwtPayload struct {
	UserID uint64
	Expire int64
}

func NewUserJwtPayload(rawJwt string, secret []byte) (*UserJwtPayload, error) {
	logger, err := mylogger.Get()
	if err != nil {
		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	tokenDuplicity, err := jwt.Parse(rawJwt, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			logger.Errorf("method == %+v %w", token.Header["alg"], ErrWrongSigningMethod)

			return nil, fmt.Errorf(myerrors.ErrTemplate, ErrInvalidToken)
		}

		return secret, nil
	})
	if err != nil {
		logger.Errorf("%s", err.Error())

		return nil, fmt.Errorf(myerrors.ErrTemplate, ErrInvalidToken)
	}

	if claims, ok := tokenDuplicity.Claims.(jwt.MapClaims); ok && tokenDuplicity.Valid {
		interfaceUserID, ok1 := claims["userID"]
		interfaceExpire, ok2 := claims["expire"]

		if !(ok1 && ok2) {
			logger.Errorf("error with claims: %+v", claims)

			return nil, fmt.Errorf(myerrors.ErrTemplate, ErrInvalidToken)
		}

		userID, ok1 := interfaceUserID.(float64)
		expire, ok2 := interfaceExpire.(float64)

		if !(ok1 && ok2) {
			logger.Errorf("error with casting claims: %+v", claims)

			return nil, fmt.Errorf(myerrors.ErrTemplate, ErrInvalidToken)
		}

		return &UserJwtPayload{UserID: uint64(userID), Expire: int64(expire)}, nil
	}

	return nil, fmt.Errorf(myerrors.ErrTemplate, ErrInvalidToken)
}

func (u *UserJwtPayload) getMapClaims() jwt.MapClaims {
	result := make(jwt.MapClaims)

	result["userID"] = u.UserID
	result["expire"] = u.Expire

	return result
}

func GenerateJwtToken(userToken *UserJwtPayload, secret []byte) (string, error) {
	logger, err := mylogger.Get()
	if err != nil {
		return "", err //nolint:wrapcheck
	}

	if userToken == nil {
		logger.Errorln(ErrNilToken)

		return "", fmt.Errorf(myerrors.ErrTemplate, ErrInvalidToken)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, userToken.getMapClaims())

	tokenString, err := token.SignedString(secret)
	if err != nil {
		logger.Errorln(err)

		return "", fmt.Errorf(myerrors.ErrTemplate, ErrInvalidToken)
	}

	return tokenString, nil
}
