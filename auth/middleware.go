package auth

import (
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/ToshihiroOgino/elib/generated/generated/domain"
	"github.com/ToshihiroOgino/elib/generated/repository"
	"github.com/ToshihiroOgino/elib/infra/sqlite"
	"github.com/gin-gonic/gin"
)

func authFailed(c *gin.Context, reason string) {
	c.JSON(http.StatusUnauthorized, gin.H{"error": reason})
	c.Abort()
}

const authTokenCookieKey = "auth_token"
const userKey = "user"

func parseBearerToken(bearerToken string) (string, error) {
	parts := strings.Split(bearerToken, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", errors.New("invalid bearer token")
	}
	return parts[1], nil
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		bearerToken, err := c.Cookie(authTokenCookieKey)
		if err != nil {
			authFailed(c, err.Error())
			return
		}

		tokenStr, err := parseBearerToken(bearerToken)
		if err != nil {
			authFailed(c, "invalid token format")
			return
		}

		claims, err := ValidateToken(tokenStr)
		if err != nil {
			authFailed(c, "authorization failed")
			return
		}

		db := sqlite.GetDB()
		q := repository.Use(db).User
		user, err := q.WithContext(c).Where(q.ID.Eq(claims.UserID)).First()
		if err != nil {
			authFailed(c, "user not found")
			return
		}
		slog.Info("user authenticated", "user_id", user.ID, "email", user.Email)
		c.Set(userKey, user)
		c.Next()
	}
}

func GetUser(c *gin.Context) *domain.User {
	user, exists := c.Get(userKey)
	if !exists {
		return nil
	}
	return user.(*domain.User)
}
