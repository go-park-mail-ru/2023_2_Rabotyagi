package config

import (
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/jwt"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/config"
)

type Config struct {
	ProductionMode     bool
	AllowOrigin        string
	Schema             string
	PortServer         string
	URLDataBase        string
	PathToRoot         string
	FileServiceDir     string
	OutputLogPath      string
	ErrorOutputLogPath string
}

func New() *Config {
	secret := config.GetEnvStr(config.EnvStandardSecret, config.StandardSecret)
	if secret != config.StandardSecret {
		jwt.SetSecret([]byte(secret))
	} else {
		_ = jwt.GetSecret()
	}

	productionMode := false
	if config.GetEnvStr(config.EnvEnvironmentMode, config.StandardDevelopmentMode) == config.StandardProductionMode {
		productionMode = true
	}

	return &Config{
		AllowOrigin:        config.GetEnvStr(config.EnvAllowOrigin, config.StandardAllowOrigin),
		Schema:             config.GetEnvStr(config.EnvSchema, config.StandardSchema),
		PortServer:         config.GetEnvStr(config.EnvPortBackend, config.StandardPort),
		URLDataBase:        config.GetEnvStr(config.EnvURLDataBase, config.StandardURLDataBase),
		PathToRoot:         config.GetEnvStr(config.EnvPathToRoot, config.StandardPathToRoot),
		FileServiceDir:     config.GetEnvStr(config.EnvFileServiceDir, config.StandardFileServiceDir),
		OutputLogPath:      config.GetEnvStr(config.EnvOutputLogPath, config.StandardOutputLogPath),
		ErrorOutputLogPath: config.GetEnvStr(config.EnvErrorOutputLogPath, config.StandardErrorOutputLogPath),
		ProductionMode:     productionMode,
	}
}
