package controller

import (
	"log/slog"
	"net/http"
	"time"

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

type IUserController interface {
	GetRegister(c *gin.Context)
	GetLogin(c *gin.Context)
	GetProfile(c *gin.Context)
	PostRegister(c *gin.Context)
	PostLogin(c *gin.Context)
}

type userController struct {
	usecase usecase.IUserUsecase
}

const URL_ROOT = "/user"
const URL_PROFILE = URL_ROOT + "/"
const URL_LOGIN = URL_ROOT + "/login"
const URL_REGISTER = URL_ROOT + "/register"

func NewUserController(router *gin.Engine) IUserController {
	instance := &userController{
		usecase: usecase.NewUserUsecase(),
	}
	setupRoute(instance, router)
	return instance
}

func setupRoute(api IUserController, router *gin.Engine) {
	router.GET(URL_LOGIN, func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", gin.H{
			"title":        "Login",
			"register_url": URL_REGISTER,
		})
	})
	router.GET(URL_REGISTER, func(c *gin.Context) {
		c.HTML(http.StatusOK, "register.html", gin.H{
			"title":     "Register",
			"login_url": URL_LOGIN,
		})
	})

	router.POST(URL_LOGIN, api.PostLogin)
	router.POST(URL_REGISTER, api.PostRegister)

	protected := router.Group("/")
	protected.Use(auth.AuthMiddleware())
	{
		protected.GET(URL_PROFILE, api.GetProfile)
	}
}

func (u *userController) GetRegister(c *gin.Context) {
	var req registerRequest

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

func (u *userController) GetLogin(c *gin.Context) {
	var req loginRequest

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

func (u *userController) GetProfile(c *gin.Context) {
	user := auth.GetUser(c)

	c.JSON(http.StatusOK, gin.H{
		"id":    *user.ID,
		"email": *user.Email,
	})
}

func (u *userController) PostRegister(c *gin.Context) {
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

	auth.SetAuthCookie(c, token)

	c.Redirect(http.StatusPermanentRedirect, URL_PROFILE)
}

func (u *userController) PostLogin(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	isOk, err := u.usecase.Validate(req.Email, req.Password)
	if !isOk {
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	auth.SetAuthCookie(c, token)

	c.Redirect(http.StatusSeeOther, URL_PROFILE)
}
