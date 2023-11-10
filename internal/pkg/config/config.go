package config

import "os"

const (
	standardAllowOrigin    = "localhost:3000"
	standardSchema         = "http://"
	standardPort           = "8080"
	standardURLDataBase    = "postgres://postgres:postgres@postgres:5432/youla?sslmode=disable"
	standardPathToRoot     = "."
	standardFileServiceDir = "./static/img"

	envAllowOrigin    = "ALLOW_ORIGIN"
	envSchema         = "SCHEMA"
	envPortBackend    = "PORT_BACKEND"
	envURLDataBase    = "URL_DATA_BASE"
	envPathToRoot     = "PATH_TO_ROOT"
	envFileServiceDir = "FILE_SERVICE_DIR"
)

type Config struct {
	AllowOrigin    string
	Schema         string
	PortServer     string
	URLDataBase    string
	PathToRoot     string
	FileServiceDir string
}

func New() *Config {
	return &Config{
		AllowOrigin:    getEnvStr(envAllowOrigin, standardAllowOrigin),
		Schema:         getEnvStr(envSchema, standardSchema),
		PortServer:     getEnvStr(envPortBackend, standardPort),
		URLDataBase:    getEnvStr(envURLDataBase, standardURLDataBase),
		PathToRoot:     getEnvStr(envPathToRoot, standardPathToRoot),
		FileServiceDir: getEnvStr(envFileServiceDir, standardFileServiceDir),
	}
}

func getEnvStr(name string, defaultValue string) string {
	result, ok := os.LookupEnv(name)
	if !ok {
		return defaultValue
	}

	return result
}
