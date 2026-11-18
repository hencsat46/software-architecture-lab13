package handlers

import domain "lab13/order-service/internal/orders"

func (r writeRequest) toEntity() *domain.Entity {
	return &domain.Entity{UserID: r.UserID, Product: r.Product, Quantity: r.Quantity}
}

func toResponse(e *domain.Entity) entityResponse {
	return entityResponse{ID: e.ID, CreatedAt: e.CreatedAt, UserID: e.UserID, Product: e.Product, Quantity: e.Quantity}
}
