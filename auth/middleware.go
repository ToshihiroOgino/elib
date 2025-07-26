package auth

import (
	"net/http"
	"strings"
	"time"

	"github.com/ToshihiroOgino/elib/generated/generated/domain"
	"github.com/ToshihiroOgino/elib/generated/repository"
	"github.com/ToshihiroOgino/elib/infra/sqlite"
	"github.com/gin-gonic/gin"
)

func authFailed(c *gin.Context, reason string) {
	switch c.Request.Method {
	case http.MethodGet:
		c.Redirect(http.StatusPermanentRedirect, "/user/login")
		c.Abort()
	default:
		c.JSON(http.StatusUnauthorized, gin.H{"error": reason})
		c.Abort()
	}
}

const authTokenCookie = "auth_token"
const userCookie = "user"

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var tokenStr string
		switch c.Request.Method {
		case http.MethodGet:
			authCookie, err := c.Cookie(authTokenCookie)
			if err != nil {
				authFailed(c, "authorization failed")
				return
			}
			tokenStr = authCookie
		default:
			authHeader := c.GetHeader("Authorization")
			if authHeader == "" {
				authFailed(c, "authorization failed")
				return
			}
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				authFailed(c, "authorization failed")
				return
			}
			tokenStr = parts[1]
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
		c.Set(userCookie, user)
		c.Next()
	}
}

func GetUser(c *gin.Context) *domain.User {
	user, exists := c.Get(userCookie)
	if !exists {
		return nil
	}
	return user.(*domain.User)
}

func SetAuthCookie(c *gin.Context, token string) {
	const AGE = 24 * 7 * time.Hour
	c.SetCookie(authTokenCookie, token, int(AGE.Seconds()), "/", "", true, true)
}
