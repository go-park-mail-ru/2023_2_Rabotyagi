package config

import "os"

const (
	// Main .env.

	StandardDevelopmentMode    = "development"
	StandardProductionMode     = "production"
	StandardPathCertFile       = "/etc/ssl/goods-galaxy.ru.crt"
	StandardPathKeyFile        = "/etc/ssl/goods-galaxy.ru.key"
	StandardAllowOrigin        = "localhost:3000"
	StandardSchema             = "http://"
	StandardPort               = "8080"
	StandardOutputLogPath      = "stdout /var/log/backend/logs.json"
	StandardErrorOutputLogPath = "stderr /var/log/backend/err_logs.json"
	StandardURLDataBase        = "postgres://postgres:postgres@localhost:5432/youla?sslmode=disable"
	StandardMainServiceName    = "backend"

	EnvEnvironmentMode    = "ENVIRONMENT"
	EnvPathCertFile       = "PATH_CERT_FILE"
	EnvPathKeyFile        = "PATH_KEY_FILE"
	EnvAllowOrigin        = "ALLOW_ORIGIN"
	EnvSchema             = "SCHEMA"
	EnvPortBackend        = "PORT_BACKEND"
	EnvOutputLogPath      = "OUTPUT_LOG_PATH"
	EnvErrorOutputLogPath = "ERROR_OUTPUT_LOG_PATH"
	EnvURLDataBase        = "URL_DATA_BASE"
	EnvServiceName        = "SERVICE_NAME"

	// File service .env.

	StandardPathToRoot             = "."
	StandardFileServiceDir         = "./static/img"
	StandardAddressFileServiceGrpc = ":8011"
	StandardFileServicePortHTTP    = "8081"
	StandardFileServiceName        = "backend_fs"

	EnvPathToRoot             = "PATH_TO_ROOT"
	EnvFileServiceDir         = "FILE_SERVICE_DIR"
	EnvAddressFileServiceGrpc = "ADDRESS_FS_GRPC"
	EnvFileServicePortHTTP    = "PORT_FS"

	// Auth service .env

	StandardAddressAuthGrpc     = ":8012"
	StandardAuthServicePortHTTP = "8082"
	StandardAuthServiceName     = "backend_auth"

	EnvAddressAuthServiceGrpc = "ADDRESS_AUTH_GRPC"
	EnvAuthServicePortHTTP    = "PORT_AUTH"
)

func GetEnvStr(name string, defaultValue string) string {
	result, ok := os.LookupEnv(name)
	if !ok {
		return defaultValue
	}

	return result
}
