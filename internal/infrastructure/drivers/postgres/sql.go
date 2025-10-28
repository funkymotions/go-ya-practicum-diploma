package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/funkymotions/go-ya-practicum-diploma/internal/config"
	_ "github.com/lib/pq"
)

type SQLDriver struct {
	DB *sql.DB
}

func NewSQLDriver(conf *config.DBConfig) (*SQLDriver, error) {
	db, err := sql.Open(conf.Type, conf.DSN)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(conf.ConnTimeout)*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}
	return &SQLDriver{DB: db}, nil
}
