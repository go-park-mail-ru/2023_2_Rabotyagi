package usecases

import (
	"go.uber.org/zap"
)

func NewLogger(outputPaths []string, errorOutputPaths []string, options ...zap.Option) (*zap.SugaredLogger, error) {
	config := zap.NewProductionConfig()
	config.DisableStacktrace = true
	config.OutputPaths = outputPaths
	config.ErrorOutputPaths = errorOutputPaths

	logger, err := config.Build(options...)
	if err != nil {
		return nil, err //nolint:wrapcheck
	}

	loggerSugar := logger.Sugar()

	return loggerSugar, nil
}
