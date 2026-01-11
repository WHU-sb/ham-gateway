package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware creates a middleware that validates the X-Gateway-Token header
func AuthMiddleware(token string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// If no token is configured, skip authentication
		if token == "" {
			c.Next()
			return
		}

		// Get token from header
		clientToken := c.GetHeader("X-Gateway-Token")
		if clientToken != token {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    "401",
				"message": "Unauthorized: invalid or missing token",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// CORSMiddleware adds CORS headers for Cloudflare Worker requests
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-Gateway-Token, Authorization")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
