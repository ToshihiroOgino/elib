package controller

import (
	"log/slog"
	"net/http"

	"github.com/ToshihiroOgino/elib/auth"
	"github.com/ToshihiroOgino/elib/usecase"
	"github.com/gin-gonic/gin"
)

type registerForm struct {
	Email    string `form:"email" binding:"required,email"`
	Password string `form:"password" binding:"required"`
}

type loginForm struct {
	Email    string `form:"email" binding:"required"`
	Password string `form:"password" binding:"required"`
}

type IUserController interface {
	// GetRegister(c *gin.Context)
	// GetLogin(c *gin.Context)
	GetProfile(c *gin.Context)
	PostRegister(c *gin.Context)
	PostLogin(c *gin.Context)
}

type userController struct {
	usecase usecase.IUserUsecase
}

const URL_ROOT = "/user"
const URL_PROFILE = URL_ROOT + ""
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

	router.Use(auth.AuthMiddleware())
	{
		router.GET(URL_PROFILE, api.GetProfile)
	}
}

/* func (u *userController) GetRegister(c *gin.Context) {
	var req registerForm

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

	token, err := auth.GenerateToken(*user.ID)
	if err != nil {
		slog.Error("failed to generate token", "error", err)
		c.HTML(http.StatusOK, "register.html", gin.H{"Error": "トークン生成に失敗しました"})
		return
	}

	c.SetCookie("auth_token", token, 3600, "/", "", false, true)
	c.Redirect(http.StatusSeeOther, "/")
}
func (u *userController) GetLogin(c *gin.Context) {
	var req loginForm
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
	c.Redirect(http.StatusSeeOther, "/")
} */

func (u *userController) GetProfile(c *gin.Context) {
	user := auth.GetUser(c)

	c.JSON(http.StatusOK, gin.H{
		"id":    user.ID,
		"email": user.Email,
	})
}

func (u *userController) PostRegister(c *gin.Context) {
	var form registerForm
	if err := c.ShouldBind(&form); err != nil {
		slog.Error("failed to bind register form", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	user, err := u.usecase.Create(form.Email, form.Password)
	if err != nil {
		slog.Error("failed to create user", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
		return
	}

	auth.SetAuthCookie(c, user.ID)
	c.Redirect(http.StatusSeeOther, URL_PROFILE)
}

func (u *userController) PostLogin(c *gin.Context) {
	var form loginForm
	if err := c.ShouldBind(&form); err != nil {
		slog.Error("failed to bind login form", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	user, err := u.usecase.FindByEmail(form.Email)
	if err != nil {
		slog.Error("failed to get user", "email", form.Email, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	isOk, err := u.usecase.Validate(user, form.Email, form.Password)
	if !isOk {
		slog.Error("validation failed", "email", form.Email, "error", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
		return
	} else if err != nil {
		slog.Error("validation error", "email", form.Email, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	auth.SetAuthCookie(c, user.ID)
	c.Redirect(http.StatusSeeOther, URL_PROFILE)
}
