package middleware

import (
	"net/http"
	"strings"

	"github.com/Wrehat/ewallet-wallet/internal/domain"
	"github.com/Wrehat/ewallet-wallet/pkg/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func AuthMiddleware(ums domain.UMSGateway, log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			log.Warn("authorization header is empty")
			response.JSON(c, http.StatusUnauthorized, "authorization header is empty", nil)
			c.Abort()
			return
		}

		token := authHeader
		if strings.HasPrefix(authHeader, "Bearer ") {
			token = strings.TrimPrefix(authHeader, "Bearer ")
		}

		tokenData, err := ums.ValidateToken(c.Request.Context(), token)
		if err != nil {
			log.Warn("token validation failed", zap.Error(err))
			response.JSON(c, http.StatusUnauthorized, "token invalid", nil)
			c.Abort()
			return
		}

		c.Set("user", tokenData)
		c.Next()
	}
}
