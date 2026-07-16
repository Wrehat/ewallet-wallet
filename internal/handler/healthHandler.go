package handler

import (
	"net/http"

	"github.com/Wrehat/ewallet-wallet/internal/domain"
	"github.com/gin-gonic/gin"
)

type HealthHandler struct {
	uc domain.HealthCheckUsecase
}

func NewHealthHandler(uc domain.HealthCheckUsecase) *HealthHandler {
	return &HealthHandler{
		uc: uc,
	}
}

func (h *HealthHandler) HealthCheck(c *gin.Context) {

	ctx := c.Request.Context()
	if err := h.uc.Check(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "unhealthy",
			"error":  "connection failed",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "healthy",
	})
}
