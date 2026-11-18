package app

import (
	"github.com/gin-gonic/gin"

	"lab13/order-service/internal/orders/handlers"
)

func setupRoutes(h *handlers.Handler) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	g := r.Group("/api/v1/orders")
	g.POST("", h.Create)
	g.GET("", h.List)
	g.GET("/:id", h.Fetch)
	g.PUT("/:id", h.Update)
	g.DELETE("/:id", h.Delete)
	return r
}
