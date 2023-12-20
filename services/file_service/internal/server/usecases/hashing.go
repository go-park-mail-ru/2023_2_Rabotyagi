package usecases

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/myerrors"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/mylogger"
)

var ErrHashing = myerrors.NewErrorInternal("ошибка хеширования")

func HashContent(ctx context.Context, content []byte) (string, error) {
	logger, err := mylogger.Get()
	if err != nil {
		return "", fmt.Errorf(myerrors.ErrTemplate, err)
	}

	logger = logger.LogReqID(ctx)

	hash := sha256.New()

	_, err = hash.Write(content)
	if err != nil {
		logger.Errorln(err)

		return "", fmt.Errorf("%w %s", ErrHashing, err.Error())
	}

	result := hash.Sum(nil)

	return hex.EncodeToString(result), nil
}
