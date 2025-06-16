package page

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/ToshihiroOgino/elib/auth"
	"github.com/ToshihiroOgino/elib/usecase"
	"github.com/gin-gonic/gin"
)

type RegisterRequest struct {
	Email    string `json:"email" form:"email" binding:"required,email"`
	Password string `json:"password" form:"password" binding:"required,min=8"`
}

type LoginRequest struct {
	Email    string `json:"email" form:"email" binding:"required,email"`
	Password string `json:"password" form:"password" binding:"required"`
}

type IUserPage interface {
	Register(c *gin.Context)
	Login(c *gin.Context)
	Profile(c *gin.Context)
}

type userApi struct {
	usecase usecase.IUserUsecase
}

func NewUserPage() IUserPage {
	return &userApi{
		usecase: usecase.NewUserUsecase(),
	}
}

func SetupRoute(api IUserPage, router *gin.Engine) {
	router.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", gin.H{
			"title": "Login",
		})
	})
	router.GET("/register", func(c *gin.Context) {
		c.HTML(http.StatusOK, "register.html", gin.H{
			"title": "Register",
		})
	})

	router.POST("/login", api.Login)
	router.POST("/register", api.Register)

	protected := router.Group("/")
	protected.Use(auth.AuthMiddleware())
	{
		protected.GET("/profile", api.Profile)
	}
}

func (u *userApi) Register(c *gin.Context) {
	var req RegisterRequest

	if err := c.ShouldBind(&req); err != nil {
		c.HTML(http.StatusOK, "register.html", gin.H{"Error": err.Error()})
		return
	}

	user, err := u.usecase.Create(req.Email, req.Password)
	if err != nil {
		slog.Error("failed to create user", "error", err)
		c.HTML(http.StatusOK, "register.html", gin.H{"Error": "登録に失敗しました"})
		return
	}

	token, err := auth.GenerateToken(*user.ID, *user.Email)
	if err != nil {
		slog.Error("failed to generate token", "error", err)
		c.HTML(http.StatusOK, "register.html", gin.H{"Error": "トークン生成に失敗しました"})
		return
	}

	c.SetCookie("auth_token", token, 3600, "/", "", false, true)
	c.Redirect(http.StatusSeeOther, "/")
}

func (u *userApi) Login(c *gin.Context) {
	var req LoginRequest

	if err := c.ShouldBind(&req); err != nil {
		c.HTML(http.StatusOK, "login.html", gin.H{"Error": err.Error()})
		return
	}

	isOk, err := u.usecase.Validate(req.Email, req.Password)
	if !isOk {
		slog.Error("validation failed", "email", req.Email, "error", err)
		c.HTML(http.StatusOK, "login.html", gin.H{"Error": "メールアドレスまたはパスワードが正しくありません"})
		return
	} else if err != nil {
		slog.Error("validation error", "email", req.Email, "error", err)
		c.HTML(http.StatusOK, "login.html", gin.H{"Error": "内部サーバーエラーが発生しました"})
		return
	}

	token, err := auth.GenerateToken(req.Email, req.Email)
	if err != nil {
		slog.Error("login failed", "error", err)
		c.HTML(http.StatusOK, "login.html", gin.H{"Error": "メールアドレスまたはパスワードが正しくありません"})
		return
	}

	const AGE = time.Hour * 24 * 7
	c.SetCookie("auth_token", token, int(AGE.Seconds()), "/", "", false, true)
	c.Redirect(http.StatusSeeOther, "/")
}

func (u *userApi) Profile(c *gin.Context) {
	user := auth.GetUser(c)

	c.JSON(http.StatusOK, gin.H{
		"id":    *user.ID,
		"email": *user.Email,
	})
}
