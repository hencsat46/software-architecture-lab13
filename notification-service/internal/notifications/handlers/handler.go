package handlers

import (
	"context"

	"lab13/notification-service/internal/logger"
	domain "lab13/notification-service/internal/notifications"
)

type IService interface {
	Create(context.Context, *domain.Entity) (*domain.Entity, error)
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
	return &Handler{service: s, logger: l.WithContext("handler", "notifications")}
}
