package main

import (
	"flag"
	"log"

	"github.com/funkymotions/go-ya-practicum-diploma/internal/config"
	appdriver "github.com/funkymotions/go-ya-practicum-diploma/internal/infrastructure/drivers/postgres"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	var direction string
	flag.StringVar(&direction, "direction", "up", "Migration direction (up/down)")
	flag.Parse()
	conf, err := config.NewConfig()
	if err != nil {
		log.Fatalf("failed to load DB config: %v", err)
	}
	appDriver, err := appdriver.NewSQLDriver(conf.DBConfig)
	if err != nil {
		log.Fatalf("failed to create PostgreSQL driver: %v", err)
	}
	driver, err := postgres.WithInstance(appDriver.DB, &postgres.Config{})
	if err != nil {
		log.Fatalf("failed to create PostgreSQL driver: %v", err)
	}
	migrator, err := migrate.NewWithDatabaseInstance("file://migrations", conf.DBConfig.Type, driver)
	if err != nil {
		log.Fatalf("failed to create new migrator: %v", err)
	}
	if direction == "up" {
		if err := migrator.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("migration up failed: %v", err)
		}
		log.Printf("migration up completed successfully")
	}
	if direction == "down" {
		if err := migrator.Down(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("migration down failed: %v", err)
		}
		log.Printf("migration down completed successfully")
	}
}
