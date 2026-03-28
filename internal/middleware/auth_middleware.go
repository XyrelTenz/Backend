package middleware

import (
	"strings"

	"backend/internal/auth/service"
	"backend/pkg/response"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(jwtService auth_service.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Unauthorized(c, "Authorization header is required")
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Unauthorized(c, "Invalid authorization format")
			c.Abort()
			return
		}

		claims, err := jwtService.ValidateToken(parts[1])
		if err != nil {
			response.Unauthorized(c, "Invalid token")
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("role", claims.Role)
		c.Next()
	}
}
