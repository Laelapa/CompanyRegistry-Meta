package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/Laelapa/CompanyRegistry/internal/config"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("FATAL: %v\n", err)
	}
}

func run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	cfg, err := config.Load()
	if err != nil {
		return err
	}

	return nil
}
