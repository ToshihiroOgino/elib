package secure

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Cookie keys
const (
	AuthTokenCookieKey   = "auth_token"
	SessionDataCookieKey = "session_data"
)

// CookieConfig holds configuration for cookie settings
type CookieConfig struct {
	Name     string
	Value    string
	MaxAge   int
	Path     string
	Domain   string
	Secure   bool
	HttpOnly bool
	SameSite http.SameSite // Use http.SameSite type
}

// DefaultCookieConfig returns default secure cookie configuration
func DefaultCookieConfig() CookieConfig {
	return CookieConfig{
		MaxAge:   int((24 * 7 * time.Hour).Seconds()), // 7 days
		Path:     "/",
		Domain:   "",
		Secure:   true,                 // HTTPS only
		HttpOnly: true,                 // No JavaScript access
		SameSite: http.SameSiteLaxMode, // SameSite=Lax
	}
}

// SessionCookieConfig returns configuration for session cookies (shorter expiry)
func SessionCookieConfig() CookieConfig {
	config := DefaultCookieConfig()
	config.MaxAge = int((1 * time.Hour).Seconds()) // 1 hour for session cookies
	return config
}

// CookieManager provides centralized cookie management
type CookieManager struct {
	defaultConfig CookieConfig
}

// NewCookieManager creates a new cookie manager with default configuration
func NewCookieManager() *CookieManager {
	return &CookieManager{
		defaultConfig: DefaultCookieConfig(),
	}
}

// SetCookie sets a cookie with the provided configuration
func (cm *CookieManager) SetCookie(c *gin.Context, config CookieConfig) {
	c.SetSameSite(config.SameSite)
	c.SetCookie(
		config.Name,
		config.Value,
		config.MaxAge,
		config.Path,
		config.Domain,
		config.Secure,
		config.HttpOnly,
	)
}

// SetAuthCookie sets the authentication cookie
func (cm *CookieManager) SetAuthCookie(c *gin.Context, token string) {
	config := cm.defaultConfig
	config.Name = AuthTokenCookieKey
	config.Value = token
	cm.SetCookie(c, config)
}

// GetCookie retrieves a cookie value
func (cm *CookieManager) GetCookie(c *gin.Context, name string) (string, error) {
	return c.Cookie(name)
}

// DeleteCookie removes a cookie by setting it to expire immediately
func (cm *CookieManager) DeleteCookie(c *gin.Context, name string) {
	config := cm.defaultConfig
	config.Name = name
	config.Value = ""
	config.MaxAge = -1
	config.Secure = false // Allow deletion over HTTP as well
	cm.SetCookie(c, config)
}

// ClearAllAuthCookies removes all authentication-related cookies
func (cm *CookieManager) ClearAllAuthCookies(c *gin.Context) {
	cm.DeleteCookie(c, AuthTokenCookieKey)
	cm.DeleteCookie(c, SessionDataCookieKey)
}

// Global cookie manager instance
var globalCookieManager = NewCookieManager()

// Package-level convenience functions
func SetAuthCookieSecure(c *gin.Context, token string) {
	globalCookieManager.SetAuthCookie(c, token)
}

func GetCookieSecure(c *gin.Context, name string) (string, error) {
	return globalCookieManager.GetCookie(c, name)
}

func DeleteCookieSecure(c *gin.Context, name string) {
	globalCookieManager.DeleteCookie(c, name)
}

func ClearAllAuthCookies(c *gin.Context) {
	globalCookieManager.ClearAllAuthCookies(c)
}
