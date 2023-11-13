package my_logger

import (
	"fmt"
	"sync"

	"go.uber.org/zap"
)

var (
	logger *zap.SugaredLogger = nil
	once   sync.Once
)

func NewNop() *zap.SugaredLogger {
	return zap.NewNop().Sugar()
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
		return nil, fmt.Errorf("NO LOGER")
	}

	return logger, nil
}
