package config

import "github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/config"

const (
	standardDevelopmentMode    = "development"
	standardProductionMode     = "production"
	standardURLDataBase        = "postgres://postgres:postgres@localhost:5432/youla?sslmode=disable"
	standardOutputLogPath      = "stdout"
	standardErrorOutputLogPath = "stderr"

	envEnvironmentMode    = "ENVIRONMENT"
	envURLDataBase        = "URL_DATA_BASE"
	envOutputLogPath      = "OUTPUT_LOG_PATH"
	envErrorOutputLogPath = "ERROR_OUTPUT_LOG_PATH"
)

type Config struct {
	ProductionMode     bool
	URLDataBase        string
	OutputLogPath      string
	ErrorOutputLogPath string
}

func New() *Config {
	productionMode := false
	if config.GetEnvStr(envEnvironmentMode, standardDevelopmentMode) == standardProductionMode {
		productionMode = true
	}

	return &Config{
		URLDataBase:        config.GetEnvStr(envURLDataBase, standardURLDataBase),
		OutputLogPath:      config.GetEnvStr(envOutputLogPath, standardOutputLogPath),
		ErrorOutputLogPath: config.GetEnvStr(envErrorOutputLogPath, standardErrorOutputLogPath),
		ProductionMode:     productionMode,
	}
}
