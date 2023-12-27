package config

import (
	"fmt"

	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/config"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/myerrors"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/services/auth/internal/jwt"
)

const (
	standardDevelopmentMode        = "development"
	standardProductionMode         = "production"
	envEnvironmentMode             = "ENVIRONMENT"
	envStandardSecret              = "JWT_SECRET"
	standardSecret                 = ""
	standardOutputLogPathAuth      = "stdout /var/log/backend/logs_auth.json"
	standardErrorOutputLogPathAuth = "stderr /var/log/backend/err_logs_auth.json"
)

type Config struct {
	ProductionMode         bool
	AuthServiceName        string
	AuthServicePort        string
	AddressAuthServiceGrpc string
	URLDataBase            string
	PathCertFile           string
	PathKeyFile            string
}

func New() (*Config, error) {
	secret := config.GetEnvStr(envStandardSecret, standardSecret)
	if secret != standardSecret {
		jwt.SetSecret([]byte(secret))
	} else {
		_, err := jwt.GetSecret()
		if err != nil {
			return nil, fmt.Errorf(myerrors.ErrTemplate, err)
		}
	}

	productionMode := false
	if config.GetEnvStr(envEnvironmentMode, standardDevelopmentMode) == standardProductionMode {
		productionMode = true
	}

	return &Config{
		ProductionMode:         productionMode,
		AuthServiceName:        config.GetEnvStr(config.EnvServiceName, config.StandardAuthServiceName),
		AuthServicePort:        config.GetEnvStr(config.EnvAuthServicePortHTTP, config.StandardAuthServicePortHTTP),
		AddressAuthServiceGrpc: config.GetEnvStr(config.EnvAddressAuthServiceGrpc, config.StandardAddressAuthGrpc),
		URLDataBase:            config.GetEnvStr(config.EnvURLDataBase, config.StandardURLDataBase),
		PathCertFile:           config.GetEnvStr(config.EnvPathCertFile, config.StandardPathCertFile),
		PathKeyFile:            config.GetEnvStr(config.EnvPathKeyFile, config.StandardPathKeyFile),
	}, nil
}

type Logger struct {
	OutputLogPath      string
	ErrorOutputLogPath string
}

func NewConfigLogger() (*Logger, error) {
	return &Logger{
		OutputLogPath:      config.GetEnvStr(config.EnvOutputLogPath, standardOutputLogPathAuth),
		ErrorOutputLogPath: config.GetEnvStr(config.EnvErrorOutputLogPath, standardErrorOutputLogPathAuth),
	}, nil
}
