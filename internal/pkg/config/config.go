package config

import "os"

const (
	standardAllowOrigin = "localhost"
	standardSchema      = "http://"
	standardPort        = "8080"
	standardURLDataBase = "postgres://postgres:postgres@postgres:5432/youla?sslmode=disable"

	envAllowOrigin = "ALLOW_ORIGIN"
	envSchema      = "SCHEMA"
	envPortBackend = "PORT_BACKEND"
	envURLDataBase = "URL_DATA_BASE"
)

type Config struct {
	AllowOrigin string
	Schema      string
	PortServer  string
	URLDataBase string
}

func New() *Config {
	return &Config{
		AllowOrigin: getEnvStr(envAllowOrigin, standardAllowOrigin),
		Schema:      getEnvStr(envSchema, standardSchema),
		PortServer:  getEnvStr(envPortBackend, standardPort),
		URLDataBase: getEnvStr(envURLDataBase, standardURLDataBase),
	}
}

func getEnvStr(name string, defaultValue string) string {
	result, ok := os.LookupEnv(name)
	if !ok {
		return defaultValue
	}

	return result
}
