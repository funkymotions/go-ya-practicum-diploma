package config

import "errors"

type AppConfig struct {
	AppAddress           string `env:"RUN_ADDRESS"`
	AccuralSystemAddress string `env:"ACCRUAL_SYSTEM_ADDRESS"`
}

func NewAppConfig() (*AppConfig, error) {
	appConf := &AppConfig{
		AppAddress:           ":8080",
		AccuralSystemAddress: "http://localhost:8081",
	}
	value, exists := LookupVar("RUN_ADDRESS", "a")
	if !exists {
		return nil, errors.New("no application address provided")
	}
	appConf.AppAddress = value
	value, exists = LookupVar("ACCRUAL_SYSTEM_ADDRESS", "r")
	if !exists {
		return nil, errors.New("no accrual system address provided")
	}
	appConf.AccuralSystemAddress = value

	return appConf, nil
}
