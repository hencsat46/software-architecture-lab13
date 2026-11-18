package handlers

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"lab13/notification-service/internal/common/response"
)

// Create godoc
// @Summary Create notifications
// @Tags notifications
// @Accept json
// @Produce json
// @Param request body writeRequest true "Entity fields"
// @Success 201 {object} response.Response
// @Failure 400,404,500,503 {object} response.Response
// @Router /api/v1/notifications [post]
func (h Handler) Create(c *gin.Context) {
	l := h.logger.WithMethod("Create")
	l.Info("Entering...")
	if response.Cancelled(c, l) {
		return
	}
	ctx := c.Request.Context()
	var req writeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		l.Warn("invalid body", zap.Error(err))
		response.BadRequest(c)
		return
	}
	e, err := h.service.Create(ctx, req.toEntity())
	if err != nil {
		response.Failure(c, l, err)
		return
	}
	response.Created(c, toResponse(e))
}
