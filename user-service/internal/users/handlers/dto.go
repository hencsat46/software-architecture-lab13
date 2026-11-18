package handlers

import "time"

type idRequest struct {
	ID string `uri:"id" binding:"required,uuid"`
}
type writeRequest struct {
	Name  string `json:"name" binding:"required,max=100"`
	Email string `json:"email" binding:"required,email,max=254"`
}
type entityResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}
