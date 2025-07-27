package secure

import (
	"time"

	"github.com/ToshihiroOgino/elib/generated/generated/domain"
	"github.com/gin-gonic/gin"
)

// LoginUser creates a new session for user login
func LoginUser(c *gin.Context, user *domain.User) error {
	return CreateUserSession(c, user)
}

// LogoutUser destroys the current user session
func LogoutUser(c *gin.Context) {
	DestroyUserSession(c)
}

// GetCurrentUser returns the currently authenticated user
func GetCurrentUser(c *gin.Context) *domain.User {
	return GetSessionUser(c)
}

// IsUserAuthenticated checks if there's a valid authenticated user
func IsUserAuthenticated(c *gin.Context) bool {
	user := GetCurrentUser(c)
	return user != nil
}

// CheckSessionTimeout checks if the session has timed out
func CheckSessionTimeout(c *gin.Context, timeoutDuration time.Duration) bool {
	sessionData, exists := GetUserSessionData(c)
	if !exists {
		return true // No session data means timed out
	}

	return IsUserSessionExpired(sessionData, timeoutDuration)
}

// RefreshCurrentSession updates the current session's last activity
func RefreshCurrentSession(c *gin.Context) error {
	user := GetCurrentUser(c)
	if user == nil {
		return nil // No user to refresh
	}

	return RefreshUserSession(c, user)
}

// SessionTimeoutMiddleware checks for session timeout
func SessionTimeoutMiddleware(timeoutDuration time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip timeout check for login/register endpoints
		if c.Request.URL.Path == "/user/login" || c.Request.URL.Path == "/user/register" {
			c.Next()
			return
		}

		if CheckSessionTimeout(c, timeoutDuration) {
			LogoutUser(c)
			c.Redirect(302, "/user/login?timeout=1")
			c.Abort()
			return
		}

		// Refresh session on activity
		RefreshCurrentSession(c)
		c.Next()
	}
}
