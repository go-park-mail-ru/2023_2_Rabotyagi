package usecases

import (
	"go.uber.org/zap"
)

func NewLogger(options ...zap.Option) (*zap.SugaredLogger, error) {
	logger, err := zap.NewProduction(options...)
	if err != nil {
		return nil, err //nolint:wrapcheck
	}

	loggerSugar := logger.Sugar()

	return loggerSugar, nil
}
