package app

import (
	"github.com/jackc/pgx/v5/pgxpool"

	"lab13/notification-service/internal/config"
	"lab13/notification-service/internal/logger"
	"lab13/notification-service/internal/notifications/handlers"
	"lab13/notification-service/internal/notifications/postgres"
	"lab13/notification-service/internal/notifications/service"
)

func initDependencies(p *pgxpool.Pool, c *config.Config, l *logger.Logger) *handlers.Handler {
	repo := postgres.NewPostgres(p, l)
	svc := service.NewService(repo, l)
	return handlers.NewHandler(svc, l)
}
