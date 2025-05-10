package auth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// Middleware handles authentication middleware functions
type Middleware struct {
	service *Service
}

// NewMiddleware creates a new auth middleware
func NewMiddleware(service *Service) *Middleware {
	return &Middleware{
		service: service,
	}
}

// RequireAuth is a middleware that enforces authentication
func (m *Middleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		// Check if the auth header has the right format
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header format must be Bearer {token}"})
			c.Abort()
			return
		}

		// Validate the token
		token := parts[1]
		claims, err := m.service.ValidateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Set user ID in the context for later use
		c.Set("userID", claims.UserID)
		c.Next()
	}
}
