package auth

import (
	"net/http"
	"strings"

	"github.com/ToshihiroOgino/elib/generated/generated/domain"
	"github.com/ToshihiroOgino/elib/generated/repository"
	"github.com/ToshihiroOgino/elib/infra/sqlite"
	"github.com/gin-gonic/gin"
)

func authFailed(c *gin.Context) {
	c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization failed"})
	c.Abort()
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var tokenStr string
		if c.Request.Method == http.MethodGet {
			authCookie, err := c.Cookie("auth_token")
			if err != nil {
				authFailed(c)
				return
			}
			tokenStr = authCookie
		} else if c.Request.Method == http.MethodPost {
			authHeader := c.GetHeader("Authorization")
			if authHeader == "" {
				authFailed(c)
				return
			}
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				authFailed(c)
				return
			}
			tokenStr = parts[1]
		}

		claims, err := ValidateToken(tokenStr)
		if err != nil {
			authFailed(c)
			return
		}

		db := sqlite.GetDB()
		q := repository.Use(db).User
		user, err := q.WithContext(c).Where(q.ID.Eq(claims.UserID)).First()
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
			c.Abort()
			return
		}

		c.Set("user", user)
		c.Next()
	}
}

func GetUser(c *gin.Context) *domain.User {
	user, exists := c.Get("user")
	if !exists {
		return nil
	}
	return user.(*domain.User)
}

func RequireAuth() gin.HandlerFunc {
	return AuthMiddleware()
}
