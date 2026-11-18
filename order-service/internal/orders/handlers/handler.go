package handlers

import (
	"context"

	"lab13/order-service/internal/logger"
	domain "lab13/order-service/internal/orders"
)

type IService interface {
	Create(context.Context, *domain.Entity) (*domain.Entity, bool, error)
	Fetch(context.Context, string) (*domain.Entity, error)
	List(context.Context) ([]domain.Entity, error)
	Update(context.Context, *domain.Entity) (*domain.Entity, error)
	Delete(context.Context, string) error
}
type Handler struct {
	service IService
	logger  *logger.Logger
}

func NewHandler(s IService, l *logger.Logger) *Handler {
	return &Handler{service: s, logger: l.WithContext("handler", "orders")}
}
