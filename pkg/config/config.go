package config

import "os"

const (
	// Main env

	StandardDevelopmentMode    = "development"
	StandardProductionMode     = "production"
	StandardAllowOrigin        = "localhost:3000"
	StandardSchema             = "http://"
	StandardPort               = "8080"
	StandardOutputLogPath      = "stdout /var/log/backend/logs.json"
	StandardErrorOutputLogPath = "stderr /var/log/backend/err_logs.json"
	StandardURLDataBase        = "postgres://postgres:postgres@localhost:5432/youla?sslmode=disable"
	StandardSecret             = ""

	EnvEnvironmentMode    = "ENVIRONMENT"
	EnvAllowOrigin        = "ALLOW_ORIGIN"
	EnvSchema             = "SCHEMA"
	EnvPortBackend        = "PORT_BACKEND"
	EnvOutputLogPath      = "OUTPUT_LOG_PATH"
	EnvErrorOutputLogPath = "ERROR_OUTPUT_LOG_PATH"
	EnvURLDataBase        = "URL_DATA_BASE"
	EnvStandardSecret     = "JWT_SECRET"

	// File service env

	StandardPathToRoot             = "."
	StandardFileServiceDir         = "./static/img"
	StandardAddressFileServiceGrpc = "127.0.0.1:8081"
	StandardFileServicePortHTTP    = "8018"

	EnvPathToRoot             = "PATH_TO_ROOT"
	EnvFileServiceDir         = "FILE_SERVICE_DIR"
	EnvAddressFileServiceGrpc = "ADDRESS_FS_GRPC"
	EnvFileServicePortHTTP    = "PORT_FS"
)

func GetEnvStr(name string, defaultValue string) string {
	result, ok := os.LookupEnv(name)
	if !ok {
		return defaultValue
	}

	return result
}
