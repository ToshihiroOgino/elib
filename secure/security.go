package secure

import (
	"crypto/rand"
	"crypto/subtle"
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
			setCSRFToken(c, token)
			c.Next()
			return
		}

		if isAuthenticationEndpoint(c.Request.URL.Path) {
			if c.Request.URL.Path == "/user/login" && c.Request.Method == "POST" {
				token := generateCSRFToken()
				setCSRFToken(c, token)
			}
			c.Next()
			return
		}

		expectedToken := GetCSRFToken(c)
		if expectedToken == "" {
			var err error
			expectedToken, err = GetCookieSecure(c, CSRFTokenCookieKey)
			if err != nil || expectedToken == "" {
				c.JSON(403, gin.H{"error": "CSRF token required"})
				c.Abort()
				return
			}
		}

		var clientToken string
		if c.GetHeader("Content-Type") == "application/json" {
			clientToken = c.GetHeader("X-CSRF-Token")
		} else {
			clientToken = c.PostForm(CSRFTokenKey)
		}

		if clientToken == "" {
			c.JSON(403, gin.H{"error": "CSRF token required"})
			c.Abort()
			return
		}

		if !validateCSRFToken(expectedToken, clientToken) {
			c.JSON(403, gin.H{"error": "CSRF token validation failed", "expected": expectedToken, "received": clientToken})
			c.Abort()
			return
		}

		c.Next()
	}
}

func isAuthenticationEndpoint(path string) bool {
	authEndpoints := []string{
		"/user/login",
		"/user/register",
		"/user/logout",
	}

	for _, endpoint := range authEndpoints {
		if path == endpoint {
			return true
		}
	}
	return false
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

func setCSRFToken(c *gin.Context, token string) {
	c.Set(CSRFTokenKey, token)
	c.Header("X-CSRF-Token", token)
	SetCSRFCookieSecure(c, token)
}

func validateCSRFToken(expected, actual string) bool {
	if len(expected) != len(actual) {
		return false
	}
	return subtle.ConstantTimeCompare([]byte(expected), []byte(actual)) == 1
}

func GetCSRFTokenFromAuth(c *gin.Context) (string, error) {
	// Verify user is authenticated
	_, err := GetLoggedInUser(c)
	if err != nil {
		return "", err
	}

	// Get CSRF token from secure cookie using centralized cookie manager
	csrfToken, err := GetCookieSecure(c, CSRFTokenCookieKey)
	if err != nil || csrfToken == "" {
		return "", err
	}

	return csrfToken, nil
}
