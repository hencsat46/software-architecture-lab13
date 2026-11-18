package app

import (
	"github.com/jackc/pgx/v5/pgxpool"

	"lab13/user-service/internal/config"
	"lab13/user-service/internal/logger"
	"lab13/user-service/internal/users/handlers"
	"lab13/user-service/internal/users/postgres"
	"lab13/user-service/internal/users/service"
)

func initDependencies(p *pgxpool.Pool, c *config.Config, l *logger.Logger) *handlers.Handler {
	repo := postgres.NewPostgres(p, l)
	svc := service.NewService(repo, l)
	return handlers.NewHandler(svc, l)
}
