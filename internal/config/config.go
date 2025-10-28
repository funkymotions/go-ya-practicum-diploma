package config

import (
	"flag"
	"os"
)

type Config struct {
	AppConfig *AppConfig
	DBConfig  *DBConfig
}

func NewConfig() (*Config, error) {
	flag.Parse()
	appConf, err := NewAppConfig()
	if err != nil {
		return nil, err
	}
	dbConf, err := NewDBConfig()
	if err != nil {
		return nil, err
	}
	return &Config{
		AppConfig: appConf,
		DBConfig:  dbConf,
	}, nil
}

func LookupVar(envKey string, flagKey string) (string, bool) {
	env, exists := os.LookupEnv(envKey)
	if exists {
		return env, true
	}
	f := flag.Lookup(flagKey)
	if f != nil {
		return f.Value.String(), true
	}
	return "", false
}
