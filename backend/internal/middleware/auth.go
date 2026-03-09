package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/godongmin/open_talk/backend/internal/config"
	"github.com/godongmin/open_talk/backend/pkg/response"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		var tokenString string

		authHeader := c.GetHeader("Authorization")
		if authHeader != "" {
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || parts[0] != "Bearer" {
				response.Unauthorized(c)
				c.Abort()
				return
			}
			tokenString = parts[1]
		}

		// Fallback to query parameter for WebSocket connections
		if tokenString == "" {
			tokenString = c.Query("token")
		}

		if tokenString == "" {
			response.Unauthorized(c)
			c.Abort()
			return
		}
		claims := &jwt.RegisteredClaims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(cfg.JWTSecret), nil
		})

		if err != nil || !token.Valid {
			response.Unauthorized(c)
			c.Abort()
			return
		}

		c.Set("userID", claims.Subject)
		c.Next()
	}
}
