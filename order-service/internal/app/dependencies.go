package app

import (
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"lab13/order-service/internal/config"
	"lab13/order-service/internal/logger"
	"lab13/order-service/internal/orders/handlers"
	"lab13/order-service/internal/orders/httpclient"
	"lab13/order-service/internal/orders/postgres"
	"lab13/order-service/internal/orders/service"
)

func initDependencies(p *pgxpool.Pool, c *config.Config, l *logger.Logger) *handlers.Handler {
	repo := postgres.NewPostgres(p, l)
	client := httpclient.New(&http.Client{Timeout: 3 * time.Second}, c.UserURL, c.NotificationURL)
	svc := service.NewService(repo, client, client, l)
	return handlers.NewHandler(svc, l)
}
