package config

import "os"

const (
	StandardDevelopmentMode    = "development"
	StandardProductionMode     = "production"
	StandardAllowOrigin        = "localhost:3000"
	StandardSchema             = "http://"
	StandardPort               = "8080"
	StandardURLDataBase        = "postgres://postgres:postgres@localhost:5432/youla?sslmode=disable"
	StandardPathToRoot         = "."
	StandardFileServiceDir     = "./static/img"
	StandardOutputLogPath      = "stdout /var/log/backend/logs.json"
	StandardErrorOutputLogPath = "stderr /var/log/backend/err_logs.json"
	StandardSecret             = ""

	EnvEnvironmentMode    = "ENVIRONMENT"
	EnvAllowOrigin        = "ALLOW_ORIGIN"
	EnvSchema             = "SCHEMA"
	EnvPortBackend        = "PORT_BACKEND"
	EnvURLDataBase        = "URL_DATA_BASE"
	EnvPathToRoot         = "PATH_TO_ROOT"
	EnvFileServiceDir     = "FILE_SERVICE_DIR"
	EnvOutputLogPath      = "OUTPUT_LOG_PATH"
	EnvErrorOutputLogPath = "ERROR_OUTPUT_LOG_PATH"
	EnvStandardSecret     = "JWT_SECRET"
)

func GetEnvStr(name string, defaultValue string) string {
	result, ok := os.LookupEnv(name)
	if !ok {
		return defaultValue
	}

	return result
}
