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

func authFailed(c *gin.Context, err error) {
	slog.Error("authentication failed", "reason", err.Error())
	c.Redirect(http.StatusSeeOther, "/user/login")
	c.Abort()
}

const authTokenCookieKey = "auth_token"
const userKey = "user"

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := GetLoggedInUser(c)
		if err != nil {
			authFailed(c, err)
			return
		}
		c.Set(userKey, user)
		c.Next()
	}
}

func GetSessionUser(c *gin.Context) *domain.User {
	user, exists := c.Get(userKey)
	if !exists {
		return nil
	}
	return user.(*domain.User)
}

func parseBearerToken(bearerToken string) (string, error) {
	parts := strings.Split(bearerToken, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", errors.New("invalid bearer token")
	}
	return parts[1], nil
}

func GetLoggedInUser(c *gin.Context) (*domain.User, error) {
	bearerToken, err := c.Cookie(authTokenCookieKey)
	if err != nil {
		bearerToken = c.GetHeader("Authorization")
		if bearerToken == "" {
			return nil, errors.New("no auth token found")
		}
		return nil, err
	}

	tokenStr, err := parseBearerToken(bearerToken)
	if err != nil {
		return nil, err
	}

	claims, err := ValidateToken(tokenStr)
	if err != nil {
		return nil, err
	}

	db := sqlite.GetDB()
	q := repository.Use(db).User
	user, err := q.WithContext(c).Where(q.ID.Eq(claims.UserID)).First()
	if err != nil {
		return nil, err
	}
	return user, nil
}

func ClearAuthCookie(c *gin.Context) {
	c.SetCookie(authTokenCookieKey, "", -1, "/", "", false, true)
	slog.Info("cleared auth cookie")
}
