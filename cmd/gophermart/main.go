package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/funkymotions/go-ya-practicum-diploma/internal/app"
)

func main() {
	app := app.NewApp()
	sigChan := make(chan os.Signal, 1)
	if err := app.Init(); err != nil {
		log.Fatalf("failed to initialize app: %v", err)
	}
	signal.Notify(sigChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		if err := app.Start(); err != nil {
			log.Fatalf("failed to run app: %v", err)
		}
	}()
	<-sigChan

	// notify app to stop gracefully
	app.Shutdown()
}
