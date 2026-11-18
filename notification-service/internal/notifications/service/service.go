package service

import (
	"context"

	"lab13/notification-service/internal/logger"
	domain "lab13/notification-service/internal/notifications"
	"lab13/notification-service/internal/notifications/handlers"
)

type IPostgres interface {
	Create(context.Context, *domain.Entity) (*domain.Entity, error)
	Fetch(context.Context, string) (*domain.Entity, error)
	List(context.Context) ([]domain.Entity, error)
	Update(context.Context, *domain.Entity) (*domain.Entity, error)
	Delete(context.Context, string) error
}
type service struct {
	postgres IPostgres
	logger   *logger.Logger
}

func NewService(p IPostgres, l *logger.Logger) handlers.IService {
	return &service{postgres: p, logger: l.WithContext("service", "notifications")}
}
