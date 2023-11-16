package my_logger

import (
	"fmt"
	"sync"

	myerrors "github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/my_errors"

	"go.uber.org/zap"
)

var (
	logger *zap.SugaredLogger //nolint:gochecknoglobals
	once   sync.Once          //nolint:gochecknoglobals

	ErrNoLogger = fmt.Errorf("my_logger.Get для отсутствующего логгера")
)

func NewNop() *zap.SugaredLogger {
	once.Do(func() {
		logger = zap.NewNop().Sugar()
	})

	return logger
}

func New(outputPaths []string, errorOutputPaths []string, options ...zap.Option) (*zap.SugaredLogger, error) {
	var err error

	once.Do(func() {
		cfg := zap.NewProductionConfig()
		cfg.OutputPaths = outputPaths
		cfg.ErrorOutputPaths = errorOutputPaths
		zapLogger, innerErr := cfg.Build(options...)
		if innerErr != nil {
			err = innerErr

			return
		}
		logger = zapLogger.Sugar()
	})

	if err != nil {
		return nil, err
	}

	return logger, nil
}

func Get() (*zap.SugaredLogger, error) {
	if logger == nil {
		fmt.Println(ErrNoLogger)

		return nil, fmt.Errorf(myerrors.ErrTemplate, ErrNoLogger)
	}

	return logger, nil
}
