package config

import "os"

const (
	standardAllowOrigin = "localhost:3000"
	standardPort        = "8080"
	standardURLDataBase = "postgres://postgres:password@localhost:5432/youla?sslmode=disable"

	envAllowOrigin = "ALLOW_ORIGIN"
	envPortBackend = "PORT_BACKEND"
	envURLDataBase = "URL_DATA_BASE"
)

type Config struct {
	PortServer  string
	AllowOrigin string
	URLDataBase string
}

func New() *Config {
	return &Config{
		AllowOrigin: getEnvStr(envAllowOrigin, standardAllowOrigin),
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
