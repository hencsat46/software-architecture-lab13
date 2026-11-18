package main

import (
	"log"

	"go.uber.org/zap"

	"lab13/order-service/internal/app"
	"lab13/order-service/internal/config"
	"lab13/order-service/internal/logger"
)

func main() {
	cfg := config.Load()
	l, err := logger.New(cfg.LogLevel)
	if err != nil {
		log.Fatal(err)
	}
	defer l.Sync()
	if err = app.New(cfg, l).Run(); err != nil {
		l.Fatal("application stopped", zap.Error(err))
	}
}
