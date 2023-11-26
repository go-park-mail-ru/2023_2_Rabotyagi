package config

import (
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/jwt"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/config"
)

const (
	standardDevelopmentMode    = "development"
	standardProductionMode     = "production"
	standardAllowOrigin        = "localhost:3000"
	standardSchema             = "http://"
	standardPort               = "8080"
	standardURLDataBase        = "postgres://postgres:postgres@localhost:5432/youla?sslmode=disable"
	standardPathToRoot         = "."
	standardFileServiceDir     = "./static/img"
	standardOutputLogPath      = "stdout /var/log/backend/logs.json"
	standardErrorOutputLogPath = "stderr /var/log/backend/err_logs.json"
	standardSecret             = ""

	envEnvironmentMode    = "ENVIRONMENT"
	envAllowOrigin        = "ALLOW_ORIGIN"
	envSchema             = "SCHEMA"
	envPortBackend        = "PORT_BACKEND"
	envURLDataBase        = "URL_DATA_BASE"
	envPathToRoot         = "PATH_TO_ROOT"
	envFileServiceDir     = "FILE_SERVICE_DIR"
	envOutputLogPath      = "OUTPUT_LOG_PATH"
	envErrorOutputLogPath = "ERROR_OUTPUT_LOG_PATH"
	envStandardSecret     = "JWT_SECRET"
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
		AllowOrigin:        config.GetEnvStr(envAllowOrigin, standardAllowOrigin),
		Schema:             config.GetEnvStr(envSchema, standardSchema),
		PortServer:         config.GetEnvStr(envPortBackend, standardPort),
		URLDataBase:        config.GetEnvStr(envURLDataBase, standardURLDataBase),
		PathToRoot:         config.GetEnvStr(envPathToRoot, standardPathToRoot),
		FileServiceDir:     config.GetEnvStr(envFileServiceDir, standardFileServiceDir),
		OutputLogPath:      config.GetEnvStr(envOutputLogPath, standardOutputLogPath),
		ErrorOutputLogPath: config.GetEnvStr(envErrorOutputLogPath, standardErrorOutputLogPath),
		ProductionMode:     productionMode,
	}
}
