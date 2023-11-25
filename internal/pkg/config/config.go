package config

import (
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/jwt"
	"os"
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
	envStandardSecret     = "STANDARD_SECRET"
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
	secret := getEnvStr(envStandardSecret, standardSecret)
	if secret != standardSecret {
		jwt.SetSecret([]byte(secret))
	} else {
		_ = jwt.GetSecret()
	}

	productionMode := false
	if getEnvStr(envEnvironmentMode, standardDevelopmentMode) == standardProductionMode {
		productionMode = true
	}

	return &Config{
		AllowOrigin:        getEnvStr(envAllowOrigin, standardAllowOrigin),
		Schema:             getEnvStr(envSchema, standardSchema),
		PortServer:         getEnvStr(envPortBackend, standardPort),
		URLDataBase:        getEnvStr(envURLDataBase, standardURLDataBase),
		PathToRoot:         getEnvStr(envPathToRoot, standardPathToRoot),
		FileServiceDir:     getEnvStr(envFileServiceDir, standardFileServiceDir),
		OutputLogPath:      getEnvStr(envOutputLogPath, standardOutputLogPath),
		ErrorOutputLogPath: getEnvStr(envErrorOutputLogPath, standardErrorOutputLogPath),
		ProductionMode:     productionMode,
	}
}

func getEnvStr(name string, defaultValue string) string {
	result, ok := os.LookupEnv(name)
	if !ok {
		return defaultValue
	}

	return result
}
