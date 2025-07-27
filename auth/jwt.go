package auth

import (
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/ToshihiroOgino/elib/env"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// Claims struct for JWT claims
type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

func generateToken(userID string) (string, error) {
	expirationTime := time.Now().Add(time.Hour * 24 * 7)

	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "elib-api",
			Subject:   userID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(env.Get().JWTSecret))
	if err != nil {
		return "", err
	}

	tokenString = "Bearer " + tokenString
	return tokenString, nil
}

func ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&Claims{},
		func(token *jwt.Token) (interface{}, error) {
			// Validate the signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(env.Get().JWTSecret), nil
		},
	)

	if err != nil {
		return nil, err
	}

	// Extract and validate the claims
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

func SetAuthCookie(c *gin.Context, userID string) {
	const AGE = 24 * 7 * time.Hour
	token, err := generateToken(userID)
	if err != nil {
		slog.Error("failed to generate token", "error", err)
		c.JSON(500, gin.H{"error": "internal server error"})
		return
	}

	c.SetCookie(authTokenCookieKey, token, int(AGE.Seconds()), "/", "", true, true)
}
