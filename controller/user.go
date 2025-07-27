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
	getProfile(c *gin.Context)
	postRegister(c *gin.Context)
	postLogin(c *gin.Context)
}

type userController struct {
	usecase usecase.IUserUsecase
}

const _URL_USER_ROOT = "/user"
const _URL_PROFILE = ""
const _URL_LOGIN = "/login"
const _URL_REGISTER = "/register"

func NewUserController(router *gin.Engine) IUserController {
	instance := &userController{
		usecase: usecase.NewUserUsecase(),
	}
	setupUserRoute(instance, router)
	return instance
}

func setupUserRoute(api IUserController, router *gin.Engine) {
	userGroup := router.Group(_URL_USER_ROOT)
	userGroup.GET(_URL_LOGIN, func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", gin.H{
			"title":        "Login",
			"register_url": _URL_REGISTER,
		})
	})
	userGroup.GET(_URL_REGISTER, func(c *gin.Context) {
		c.HTML(http.StatusOK, "register.html", gin.H{
			"title":     "Register",
			"login_url": _URL_LOGIN,
		})
	})

	userGroup.POST(_URL_LOGIN, api.postLogin)
	userGroup.POST(_URL_REGISTER, api.postRegister)

	userGroup.Use(auth.AuthMiddleware())
	{
		userGroup.GET(_URL_PROFILE, api.getProfile)
	}
}

func (u *userController) getProfile(c *gin.Context) {
	user := auth.GetUser(c)

	c.JSON(http.StatusOK, gin.H{
		"id":    user.ID,
		"email": user.Email,
	})
}

func (u *userController) postRegister(c *gin.Context) {
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
	c.Redirect(http.StatusSeeOther, "/note")
}

func (u *userController) postLogin(c *gin.Context) {
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
	c.Redirect(http.StatusSeeOther, "/note")
}
