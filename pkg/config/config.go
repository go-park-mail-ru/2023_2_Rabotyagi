package config

import "os"

const (
	standardAllowOrigin = "localhost:3000"
	standardPort        = "8080"

	envAllowOrigin = "ALLOW_ORIGIN"
	envPortBackend = "PORT_BACKEND"
)

type Config struct {
	PortServer  string
	AllowOrigin string
}

func New() *Config {
	return &Config{
		AllowOrigin: getEnvStr(envAllowOrigin, standardAllowOrigin),
		PortServer:  getEnvStr(envPortBackend, standardPort),
	}
}

func getEnvStr(name string, defaultValue string) string {
	result, ok := os.LookupEnv(name)
	if !ok {
		return defaultValue
	}

	return result
}
