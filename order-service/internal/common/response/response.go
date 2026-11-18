package response

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"lab13/order-service/internal/common/exceptions"
	"lab13/order-service/internal/logger"
)

type Response struct {
	Data  any    `json:"data,omitempty"`
	Error string `json:"error,omitempty"`
}

func OK(c *gin.Context, v any) { c.JSON(http.StatusOK, Response{Data: v}) }

func Created(c *gin.Context, v any) { c.JSON(http.StatusCreated, Response{Data: v}) }

func BadRequest(c *gin.Context) { c.JSON(400, Response{Error: "invalid request"}) }

func Cancelled(c *gin.Context, l *logger.Logger) bool {
	if err := c.Request.Context().Err(); err != nil {
		l.Warn("request cancelled", zap.Error(err))
		c.JSON(408, Response{Error: "request cancelled"})
		return true
	}
	return false
}

func Failure(c *gin.Context, l *logger.Logger, err error) {
	l.Error("request failed", zap.Error(err))
	code, msg := 500, "internal server error"
	switch {
	case errors.Is(err, exceptions.ErrNotFound):
		code, msg = 404, "entity not found"
	case errors.Is(err, exceptions.ErrInvalidUser):
		code, msg = 400, "user does not exist"
	case errors.Is(err, exceptions.ErrUnavailable):
		code, msg = 503, "user service unavailable"
	}
	c.JSON(code, Response{Error: msg})
}
