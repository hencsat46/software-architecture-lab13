package handlers

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"lab13/user-service/internal/common/response"
)

// Delete godoc
// @Summary Delete users
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "Entity UUID"
// @Success 200 {object} response.Response
// @Failure 400,404,500,503 {object} response.Response
// @Router /api/v1/users/{id} [delete]
func (h Handler) Delete(c *gin.Context) {
	l := h.logger.WithMethod("Delete")
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
	if err := h.service.Delete(ctx, id.ID); err != nil {
		response.Failure(c, l, err)
		return
	}
	response.OK(c, gin.H{"deleted": true})
}
