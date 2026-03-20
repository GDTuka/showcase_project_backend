package middleware

import (
	"net/http"
	"strings"

	"showcase_project/internal/service"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(jwtService service.JWT) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header format must be Bearer {token}"})
			return
		}

		userId, err := jwtService.ValidateToken(parts[1], "access")
		if err != nil {
			c.AbortWithStatusJSON(err.Code(), gin.H{"error": err.Error()})
			return
		}

		c.Set("user_id", userId)
		c.Next()
	}
}
