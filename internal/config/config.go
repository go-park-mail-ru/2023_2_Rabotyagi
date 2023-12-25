package config

import (
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/config"
)

const (
	EnvPremiumShopID     = "PREMIUM_SHOP_ID"
	EnvPremiumShopSecret = "PREMIUM_SHOP_SECRET" //nolint:gosec

	StandardPremiumShopID     = "297668"
	StandardPremiumShopSecret = "test_qlRvNM1Btl6h3upjYaWEJSxfzjqyI6CdsrbcPsFS_3M" //nolint:gosec
)

type Config struct {
	ProductionMode         bool
	MainServiceName        string
	AllowOrigin            string
	Schema                 string
	PortServer             string
	URLDataBase            string
	AddressFileServiceGrpc string
	AddressAuthServiceGrpc string
	PremiumShopID          string
	PremiumShopSecret      string
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
		MainServiceName:        config.GetEnvStr(config.EnvServiceName, config.StandardMainServiceName),
		AllowOrigin:            config.GetEnvStr(config.EnvAllowOrigin, config.StandardAllowOrigin),
		Schema:                 config.GetEnvStr(config.EnvSchema, config.StandardSchema),
		PortServer:             config.GetEnvStr(config.EnvPortBackend, config.StandardPort),
		URLDataBase:            config.GetEnvStr(config.EnvURLDataBase, config.StandardURLDataBase),
		AddressFileServiceGrpc: config.GetEnvStr(config.EnvAddressFileServiceGrpc, config.StandardAddressFileServiceGrpc),
		AddressAuthServiceGrpc: config.GetEnvStr(config.EnvAddressAuthServiceGrpc, config.StandardAddressAuthGrpc),
		PremiumShopID:          config.GetEnvStr(EnvPremiumShopID, StandardPremiumShopID),
		PremiumShopSecret:      config.GetEnvStr(EnvPremiumShopSecret, StandardPremiumShopSecret),
		OutputLogPath:          config.GetEnvStr(config.EnvOutputLogPath, config.StandardOutputLogPath),
		ErrorOutputLogPath:     config.GetEnvStr(config.EnvErrorOutputLogPath, config.StandardErrorOutputLogPath),
	}
}
