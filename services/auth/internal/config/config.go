package config

import (
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/config"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/services/auth/internal/jwt"
)

const (
	standardDevelopmentMode = "development"
	standardProductionMode  = "production"
	standardOutputLogPath   = "stdout"
	envEnvironmentMode      = "ENVIRONMENT"
	envStandardSecret       = "JWT_SECRET"
	standardSecret          = ""
)

type Config struct {
	ProductionMode         bool
	AddressAuthServiceGrpc string
	URLDataBase            string
	OutputLogPath          string
	ErrorOutputLogPath     string
}

func New() *Config {
	secret := config.GetEnvStr(envStandardSecret, standardSecret)
	if secret != standardSecret {
		jwt.SetSecret([]byte(secret))
	} else {
		_ = jwt.GetSecret()
	}

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