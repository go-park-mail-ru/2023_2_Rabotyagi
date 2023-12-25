package config

import "github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/config"

const (
	standardOutputLogPathFS      = "stdout /var/log/backend/logs_fs.json"
	standardErrorOutputLogPathFS = "stderr /var/log/backend/err_logs_fs.json"
)

type Config struct {
	ProductionMode         bool
	ServiceName            string
	AddressFileServiceGrpc string
	Schema                 string
	AllowOrigin            string
	Port                   string
	PathToRoot             string
	FileServiceDir         string
	PathCertFile           string
	PathKeyFile            string
	OutputLogPath          string
	ErrorOutputLogPath     string
}

func New() *Config {
	productionMode := false
	if config.GetEnvStr(config.EnvEnvironmentMode, config.StandardDevelopmentMode) == config.StandardProductionMode {
		productionMode = true
	}

	return &Config{
		ProductionMode:         productionMode,
		ServiceName:            config.GetEnvStr(config.EnvServiceName, config.StandardFileServiceName),
		AddressFileServiceGrpc: config.GetEnvStr(config.EnvAddressFileServiceGrpc, config.StandardAddressFileServiceGrpc),
		AllowOrigin:            config.GetEnvStr(config.EnvAllowOrigin, config.StandardAllowOrigin),
		Schema:                 config.GetEnvStr(config.EnvSchema, config.StandardSchema),
		Port:                   config.GetEnvStr(config.EnvFileServicePortHTTP, config.StandardFileServicePortHTTP),
		PathToRoot:             config.GetEnvStr(config.EnvPathToRoot, config.StandardPathToRoot),
		FileServiceDir:         config.GetEnvStr(config.EnvFileServiceDir, config.StandardFileServiceDir),
		PathCertFile:           config.GetEnvStr(config.EnvPathCertFile, config.StandardPathCertFile),
		PathKeyFile:            config.GetEnvStr(config.EnvPathKeyFile, config.StandardPathKeyFile),
		OutputLogPath:          config.GetEnvStr(config.EnvOutputLogPath, standardOutputLogPathFS),
		ErrorOutputLogPath:     config.GetEnvStr(config.EnvErrorOutputLogPath, standardErrorOutputLogPathFS),
	}
}
