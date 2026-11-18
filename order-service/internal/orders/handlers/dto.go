package handlers

import "time"

type idRequest struct {
	ID string `uri:"id" binding:"required,uuid"`
}
type writeRequest struct {
	UserID   string `json:"user_id" binding:"required,uuid"`
	Product  string `json:"product" binding:"required,max=200"`
	Quantity int    `json:"quantity" binding:"required,min=1,max=10000"`
}
type entityResponse struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Product   string    `json:"product"`
	Quantity  int       `json:"quantity"`
	CreatedAt time.Time `json:"created_at"`
}
type createResponse struct {
	entityResponse
	NotificationSent bool `json:"notification_sent"`
}
