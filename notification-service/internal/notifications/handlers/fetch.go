package handlers

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"lab13/notification-service/internal/common/response"
)

// Fetch godoc
// @Summary Fetch notifications
// @Tags notifications
// @Accept json
// @Produce json
// @Param id path string true "Entity UUID"
// @Success 200 {object} response.Response
// @Failure 400,404,500,503 {object} response.Response
// @Router /api/v1/notifications/{id} [get]
func (h Handler) Fetch(c *gin.Context) {
	l := h.logger.WithMethod("Fetch")
	l.Info("Entering...")
	if response.Cancelled(c, l) {
		return
	}
	ctx := c.Request.Context()
	var id idRequest
	if err := c.ShouldBindUri(&id); err != nil {
		l.Warn("invalid path", zap.Error(err))
		response.BadRequest(c)
		return
	}
	e, err := h.service.Fetch(ctx, id.ID)
	if err != nil {
		response.Failure(c, l, err)
		return
	}
	response.OK(c, toResponse(e))
}
