package service

import (
	"context"

	"lab13/order-service/internal/logger"
	domain "lab13/order-service/internal/orders"
	"lab13/order-service/internal/orders/handlers"
)

type IPostgres interface {
	Create(context.Context, *domain.Entity) (*domain.Entity, error)
	Fetch(context.Context, string) (*domain.Entity, error)
	List(context.Context) ([]domain.Entity, error)
	Update(context.Context, *domain.Entity) (*domain.Entity, error)
	Delete(context.Context, string) error
}
type IUser interface {
	Exists(context.Context, string) error
}
type INotification interface {
	Send(context.Context, string, string) error
}
type service struct {
	postgres      IPostgres
	logger        *logger.Logger
	users         IUser
	notifications INotification
}

func NewService(p IPostgres, u IUser, n INotification, l *logger.Logger) handlers.IService {
	return &service{postgres: p, logger: l.WithContext("service", "orders"), users: u, notifications: n}
}
