package secure

import (
	"log/slog"
	"time"

	"github.com/ToshihiroOgino/elib/generated/generated/domain"
	"github.com/gin-gonic/gin"
)

// SessionData holds user session information
type SessionData struct {
	UserID    string    `json:"user_id"`
	Email     string    `json:"email"`
	LoginTime time.Time `json:"login_time"`
	LastSeen  time.Time `json:"last_seen"`
}

// SessionManager provides centralized session management
type SessionManager struct {
	cookieManager *CookieManager
}

// NewSessionManager creates a new session manager
func NewSessionManager() *SessionManager {
	return &SessionManager{
		cookieManager: NewCookieManager(),
	}
}

// CreateSession creates a new user session
func (sm *SessionManager) CreateSession(c *gin.Context, user *domain.User) error {
	now := time.Now()

	sessionData := SessionData{
		UserID:    user.ID,
		Email:     user.Email,
		LoginTime: now,
		LastSeen:  now,
	}

	// Generate JWT token
	token, err := generateToken(user.ID)
	if err != nil {
		slog.Error("failed to generate JWT token", "error", err)
		return err
	}

	// Set auth cookie
	sm.cookieManager.SetAuthCookie(c, token)

	// Store session data in context for current request
	c.Set("session_data", sessionData)
	c.Set(userKey, user)

	slog.Info("session created", "user_id", user.ID, "email", user.Email)
	return nil
}

// RefreshSession updates the session's last seen time
func (sm *SessionManager) RefreshSession(c *gin.Context, user *domain.User) error {
	// Get existing session data
	sessionData, exists := c.Get("session_data")
	if !exists {
		// Create new session data if not exists
		now := time.Now()
		sessionData = SessionData{
			UserID:    user.ID,
			Email:     user.Email,
			LoginTime: now,
			LastSeen:  now,
		}
	} else {
		// Update last seen time
		data := sessionData.(SessionData)
		data.LastSeen = time.Now()
		sessionData = data
	}

	c.Set("session_data", sessionData)
	return nil
}

// GetSessionData retrieves session data from context
func (sm *SessionManager) GetSessionData(c *gin.Context) (SessionData, bool) {
	sessionData, exists := c.Get("session_data")
	if !exists {
		return SessionData{}, false
	}
	return sessionData.(SessionData), true
}

// ValidateSession checks if the current session is valid
func (sm *SessionManager) ValidateSession(c *gin.Context) (*domain.User, error) {
	// First try to get user using existing authentication
	user, err := GetLoggedInUser(c)
	if err != nil {
		return nil, err
	}

	// Refresh session data
	err = sm.RefreshSession(c, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// DestroySession removes all session data and cookies
func (sm *SessionManager) DestroySession(c *gin.Context) {
	// Clear all authentication cookies
	sm.cookieManager.ClearAllAuthCookies(c)

	// Remove session data from context
	sessionData, exists := c.Get("session_data")
	if exists {
		data := sessionData.(SessionData)
		slog.Info("session destroyed", "user_id", data.UserID, "email", data.Email)
	}

	// Clear context data
	c.Set("session_data", nil)
	c.Set(userKey, nil)
}

// IsSessionExpired checks if session has expired based on inactivity
func (sm *SessionManager) IsSessionExpired(sessionData SessionData, maxInactivity time.Duration) bool {
	return time.Since(sessionData.LastSeen) > maxInactivity
}

// Global session manager instance
var globalSessionManager = NewSessionManager()

// Package-level convenience functions for session management
func CreateUserSession(c *gin.Context, user *domain.User) error {
	return globalSessionManager.CreateSession(c, user)
}

func ValidateUserSession(c *gin.Context) (*domain.User, error) {
	return globalSessionManager.ValidateSession(c)
}

func RefreshUserSession(c *gin.Context, user *domain.User) error {
	return globalSessionManager.RefreshSession(c, user)
}

func GetUserSessionData(c *gin.Context) (SessionData, bool) {
	return globalSessionManager.GetSessionData(c)
}

func DestroyUserSession(c *gin.Context) {
	globalSessionManager.DestroySession(c)
}

func IsUserSessionExpired(sessionData SessionData, maxInactivity time.Duration) bool {
	return globalSessionManager.IsSessionExpired(sessionData, maxInactivity)
}
