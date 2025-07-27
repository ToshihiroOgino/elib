package security

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/gin-gonic/gin"
)

func SecurityMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		csp := "default-src 'self'; " +
			"script-src 'self' 'unsafe-inline' https://cdn.jsdelivr.net; " +
			"style-src 'self' 'unsafe-inline' https://cdn.jsdelivr.net; " +
			"font-src 'self' https://cdn.jsdelivr.net; " +
			"img-src 'self' data: https:; " +
			"connect-src 'self'; " +
			"object-src 'none'; " +
			"base-uri 'self'; " +
			"form-action 'self'"
		c.Header("Content-Security-Policy", csp)

		c.Next()
	}
}

func CSRFMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == "GET" || c.Request.Method == "HEAD" || c.Request.Method == "OPTIONS" {
			token := generateCSRFToken()
			SetCSRFToken(c, token)
			c.Next()
			return
		}

		if c.GetHeader("Content-Type") == "application/json" {
			tokenFromHeader := c.GetHeader("X-CSRF-Token")
			if tokenFromHeader == "" {
				c.JSON(403, gin.H{"error": "CSRF token required"})
				c.Abort()
				return
			}
			// For simplicity, we'll validate that token exists and is not empty
			// In production, you would validate against a session-stored token
			// TODO: Implement proper CSRF token validation
			c.Next()
			return
		}
		tokenFromForm := c.PostForm(CSRFTokenKey)
		if tokenFromForm == "" {
			c.JSON(403, gin.H{"error": "CSRF token required"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func generateCSRFToken() string {
	bytes := make([]byte, 32)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

const CSRFTokenKey = "csrf_token"

func GetCSRFToken(c *gin.Context) string {
	token, exists := c.Get(CSRFTokenKey)
	if !exists {
		return ""
	}
	return token.(string)
}

func SetCSRFToken(c *gin.Context, token string) {
	c.Set(CSRFTokenKey, token)
	c.Header("X-CSRF-Token", token)
}
