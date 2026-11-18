package handlers

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"lab13/order-service/internal/common/response"
)

// Update godoc
// @Summary Update orders
// @Tags orders
// @Accept json
// @Produce json
// @Param id path string true "Entity UUID"
// @Param request body writeRequest true "Entity fields"
// @Success 200 {object} response.Response
// @Failure 400,404,500,503 {object} response.Response
// @Router /api/v1/orders/{id} [put]
func (h Handler) Update(c *gin.Context) {
	l := h.logger.WithMethod("Update")
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
	var req writeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		l.Warn("invalid body", zap.Error(err))
		response.BadRequest(c)
		return
	}
	e := req.toEntity()
	e.ID = id.ID
	result, err := h.service.Update(ctx, e)
	if err != nil {
		response.Failure(c, l, err)
		return
	}
	response.OK(c, toResponse(result))
}
