package config

import "github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/config"

const (
	standardOutputLogPathFST   = "stdout /var/log/backend/logs_fs.json"
	standardErrorOutputLogPath = "stderr /var/log/backend/err_logs_fs.json"
)

type Config struct {
	ProductionMode         bool
	AddressFileServiceGrpc string
	Schema                 string
	AllowOrigin            string
	Port                   string
	PathToRoot             string
	FileServiceDir         string
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
		AddressFileServiceGrpc: config.GetEnvStr(config.EnvAddressFileServiceGrpc, config.StandardAddressFileServiceGrpc),
		AllowOrigin:            config.GetEnvStr(config.EnvAllowOrigin, config.StandardAllowOrigin),
		Schema:                 config.GetEnvStr(config.EnvSchema, config.StandardSchema),
		Port:                   config.GetEnvStr(config.EnvFileServicePortHTTP, config.StandardFileServicePortHTTP),
		PathToRoot:             config.GetEnvStr(config.EnvPathToRoot, config.StandardPathToRoot),
		FileServiceDir:         config.GetEnvStr(config.EnvFileServiceDir, config.StandardFileServiceDir),
		OutputLogPath:          config.GetEnvStr(config.EnvOutputLogPath, standardOutputLogPathFST),
		ErrorOutputLogPath:     config.GetEnvStr(config.EnvErrorOutputLogPath, standardErrorOutputLogPath),
	}
}
