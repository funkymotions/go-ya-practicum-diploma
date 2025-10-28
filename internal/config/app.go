package config

import (
	"errors"
	"flag"
)

type AppConfig struct {
	AppAddress           string `env:"RUN_ADDRESS"`
	AccrualSystemAddress string `env:"ACCRUAL_SYSTEM_ADDRESS"`
}

func init() {
	flag.String("a", "", "Application run address")
	flag.String("r", "", "Accrual system address")
}

func NewAppConfig() (*AppConfig, error) {
	appConf := &AppConfig{
		AppAddress:           ":8080",
		AccrualSystemAddress: "http://localhost:8081",
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
	appConf.AccrualSystemAddress = value
	return appConf, nil
}
