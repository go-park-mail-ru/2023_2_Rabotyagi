package config

import "github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/config"

const (
	standardDevelopmentMode = "development"
	standardProductionMode  = "production"
	standardOutputLogPath   = "stdout"
	envEnvironmentMode      = "ENVIRONMENT"
)

type Config struct {
	ProductionMode         bool
	AddressAuthServiceGrpc string
	URLDataBase            string
	OutputLogPath          string
	ErrorOutputLogPath     string
}

func New() *Config {
	productionMode := false
	if config.GetEnvStr(envEnvironmentMode, standardDevelopmentMode) == standardProductionMode {
		productionMode = true
	}

	return &Config{
		ProductionMode:         productionMode,
		AddressAuthServiceGrpc: config.GetEnvStr(config.EnvAddressAuthServiceGrpc, config.StandardAddressAuthGrpc),
		URLDataBase:            config.GetEnvStr(config.EnvURLDataBase, config.StandardURLDataBase),
		OutputLogPath:          config.GetEnvStr(config.EnvOutputLogPath, standardOutputLogPath),
		ErrorOutputLogPath:     config.GetEnvStr(config.EnvErrorOutputLogPath, config.StandardErrorOutputLogPath),
	}
}
