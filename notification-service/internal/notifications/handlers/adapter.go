package handlers

import domain "lab13/notification-service/internal/notifications"

func (r writeRequest) toEntity() *domain.Entity {
	return &domain.Entity{UserID: r.UserID, Message: r.Message}
}

func toResponse(e *domain.Entity) entityResponse {
	return entityResponse{ID: e.ID, CreatedAt: e.CreatedAt, UserID: e.UserID, Message: e.Message}
}
