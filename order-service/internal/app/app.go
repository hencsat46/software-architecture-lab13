package app

import (
	"context"
	"errors"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"lab13/order-service/internal/config"
	"lab13/order-service/internal/database"
	"lab13/order-service/internal/logger"
)

type App struct {
	cfg    *config.Config
	logger *logger.Logger
}

func New(c *config.Config, l *logger.Logger) *App { return &App{cfg: c, logger: l} }

func (a *App) Run() error {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	dbctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	pool, err := database.NewPostgres(dbctx, a.cfg.Postgres)
	if err != nil {
		return err
	}
	defer pool.Close()
	h := initDependencies(pool, a.cfg, a.logger)
	r := setupRoutes(h)
	server := &http.Server{Addr: a.cfg.Port, Handler: r, ReadHeaderTimeout: 5 * time.Second, ReadTimeout: 10 * time.Second, WriteTimeout: 15 * time.Second, IdleTimeout: 60 * time.Second}
	done := make(chan error, 1)
	go func() { done <- server.ListenAndServe() }()
	select {
	case err := <-done:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	case <-ctx.Done():
		shutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return server.Shutdown(shutdown)
	}
}
