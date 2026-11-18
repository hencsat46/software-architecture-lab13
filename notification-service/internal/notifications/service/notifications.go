package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	domain "lab13/notification-service/internal/notifications"
)

func (s service) Create(ctx context.Context, e *domain.Entity) (*domain.Entity, error) {
	l := s.logger.WithMethod("Create")
	l.Info("Entering...")
	e.ID = uuid.NewString()
	e.CreatedAt = time.Now().UTC()
	result, err := s.postgres.Create(ctx, e)
	if err != nil {
		l.Error("create failed", zap.Error(err))
		return nil, err
	}
	l.Info("notification sent (log stub)", zap.String("user_id", e.UserID), zap.String("message", e.Message))
	return result, nil
}

func (s service) Fetch(ctx context.Context, id string) (*domain.Entity, error) {
	l := s.logger.WithMethod("Fetch")
	l.Info("Entering...")
	result, err := s.postgres.Fetch(ctx, id)
	if err != nil {
		l.Error("fetch failed", zap.Error(err))
	}
	return result, err
}

func (s service) List(ctx context.Context) ([]domain.Entity, error) {
	l := s.logger.WithMethod("List")
	l.Info("Entering...")
	result, err := s.postgres.List(ctx)
	if err != nil {
		l.Error("list failed", zap.Error(err))
	}
	return result, err
}

func (s service) Update(ctx context.Context, e *domain.Entity) (*domain.Entity, error) {
	l := s.logger.WithMethod("Update")
	l.Info("Entering...")
	result, err := s.postgres.Update(ctx, e)
	if err != nil {
		l.Error("update failed", zap.Error(err))
	}
	return result, err
}

func (s service) Delete(ctx context.Context, id string) error {
	l := s.logger.WithMethod("Delete")
	l.Info("Entering...")
	err := s.postgres.Delete(ctx, id)
	if err != nil {
		l.Error("delete failed", zap.Error(err))
	}
	return err
}
