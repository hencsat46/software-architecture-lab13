package handlers

import "time"

type idRequest struct {
	ID string `uri:"id" binding:"required,uuid"`
}
type writeRequest struct {
	UserID  string `json:"user_id" binding:"required,uuid"`
	Message string `json:"message" binding:"required,max=2000"`
}
type entityResponse struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"created_at"`
}
