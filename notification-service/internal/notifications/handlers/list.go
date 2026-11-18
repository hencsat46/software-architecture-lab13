package handlers

import (
	"github.com/gin-gonic/gin"

	"lab13/notification-service/internal/common/response"
)

// List godoc
// @Summary List notifications
// @Tags notifications
// @Accept json
// @Produce json
// @Success 200 {object} response.Response
// @Failure 400,404,500,503 {object} response.Response
// @Router /api/v1/notifications [get]
func (h Handler) List(c *gin.Context) {
	l := h.logger.WithMethod("List")
	l.Info("Entering...")
	if response.Cancelled(c, l) {
		return
	}
	ctx := c.Request.Context()
	es, err := h.service.List(ctx)
	if err != nil {
		response.Failure(c, l, err)
		return
	}
	out := make([]entityResponse, 0, len(es))
	for i := range es {
		out = append(out, toResponse(&es[i]))
	}
	response.OK(c, out)
}
