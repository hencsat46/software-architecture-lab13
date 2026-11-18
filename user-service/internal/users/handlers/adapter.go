package handlers

import domain "lab13/user-service/internal/users"

func (r writeRequest) toEntity() *domain.Entity { return &domain.Entity{Name: r.Name, Email: r.Email} }

func toResponse(e *domain.Entity) entityResponse {
	return entityResponse{ID: e.ID, CreatedAt: e.CreatedAt, Name: e.Name, Email: e.Email}
}
