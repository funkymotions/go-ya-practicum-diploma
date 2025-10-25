package config

import (
	"errors"
	"flag"

	"github.com/caarlos0/env/v11"
)

type DBConfig struct {
	Type        string `env:"DB_TYPE"`
	DSN         string `env:"DATABASE_URI"`
	ConnTimeout int64  `env:"DB_CONN_TIMEOUT"`
}

func NewDBConfig() (*DBConfig, error) {
	var flagName = "d"
	var DSN string
	flag.StringVar(&DSN, flagName, "", "Database connection string")
	flag.Parse()
	dbConfig := DBConfig{
		Type:        "postgres",
		ConnTimeout: 1,
	}
	if err := env.Parse(&dbConfig); err != nil {
		return nil, err
	}
	if dbConfig.DSN != "" {
		return &dbConfig, nil
	}
	if f := flag.Lookup(flagName); f != nil && f.Value.String() != "" {
		dbConfig.DSN = f.Value.String()
		return &dbConfig, nil
	}
	return nil, errors.New("no database connection string provided")
}
