package mylogger

import (
	"context"
	"fmt"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/myerrors"
	"sync"

	"go.uber.org/zap"
)

var (
	logger *zap.SugaredLogger //nolint:gochecknoglobals
	once   sync.Once          //nolint:gochecknoglobals

	ErrNoLogger = myerrors.NewErrorInternal("my_logger.Get для отсутствующего логгера")
)

type MyLogger struct {
	*zap.SugaredLogger
}

func NewNop() *MyLogger {
	once.Do(func() {
		logger = zap.NewNop().Sugar()
	})

	return &MyLogger{logger}
}

func New(outputPaths []string, errorOutputPaths []string, options ...zap.Option) (*MyLogger, error) {
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

	return &MyLogger{logger}, nil
}

func Get() (*MyLogger, error) {
	if logger == nil {
		fmt.Println(ErrNoLogger)

		return nil, fmt.Errorf(myerrors.ErrTemplate, ErrNoLogger)
	}

	return &MyLogger{logger}, nil
}

func (m *MyLogger) LogReqID(ctx context.Context) *MyLogger {
	return &MyLogger{m.With(
		zap.String("req_id", GetRequestIDFromCtx(ctx)),
	)}
}
