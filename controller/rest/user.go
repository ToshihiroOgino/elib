package rest

import (
	"log/slog"
	"net/http"

	"github.com/ToshihiroOgino/elib/auth"
	"github.com/ToshihiroOgino/elib/usecase"
	"github.com/gin-gonic/gin"
)

type registerRequest struct {
	Email    string `json:"email" form:"email" binding:"required,email"`
	Password string `json:"password" form:"password" binding:"required,min=8"`
}

type loginRequest struct {
	Email    string `json:"email" form:"email" binding:"required,email"`
	Password string `json:"password" form:"password" binding:"required"`
}

type authResponse struct {
	Token string `json:"token"`
}

type IUserRest interface {
	Register(c *gin.Context)
	Login(c *gin.Context)
	Profile(c *gin.Context)
}

type userRest struct {
	usecase usecase.IUserUsecase
}

func NewUserRest() IUserRest {
	return &userRest{
		usecase: usecase.NewUserUsecase(),
	}
}

func SetupRoute(api IUserRest, router *gin.Engine) {
	router.POST("/register", api.Register)
	router.POST("/login", api.Login)

	protected := router.Group("/")
	protected.Use(auth.AuthMiddleware())
	{
		protected.GET("/profile", api.Profile)
	}
}

func (u *userRest) Register(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := u.usecase.Create(req.Email, req.Password)
	if err != nil {
		slog.Error("failed to create user", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
		return
	}

	token, err := auth.GenerateToken(*user.ID, *user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	c.JSON(http.StatusCreated, authResponse{Token: token})
}

func (u *userRest) Login(c *gin.Context) {
	var req loginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if isOk, err := u.usecase.Validate(req.Email, req.Password); !isOk {
		slog.Error("validation failed", "email", req.Email, "error", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
		return
	} else if err != nil {
		slog.Error("validation error", "email", req.Email, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	token, err := auth.GenerateToken(req.Email, req.Email)

	if err != nil {
		slog.Error("login failed", "error", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
		return
	}

	c.JSON(http.StatusOK, authResponse{Token: token})
}

func (u *userRest) Profile(c *gin.Context) {
	user := auth.GetUser(c)

	c.JSON(http.StatusOK, gin.H{
		"id":    *user.ID,
		"email": *user.Email,
	})
}
