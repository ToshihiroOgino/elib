package controller

import (
	"log/slog"
	"net/http"

	"github.com/ToshihiroOgino/elib/secure"
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
	getLogin(c *gin.Context)
	getRegister(c *gin.Context)
	postRegister(c *gin.Context)
	postLogin(c *gin.Context)
	postLogout(c *gin.Context)
}

type userController struct {
	usecase usecase.IUserUsecase
}

func NewUserController(router *gin.Engine) IUserController {
	instance := &userController{
		usecase: usecase.NewUserUsecase(),
	}
	setupUserRoute(instance, router)
	return instance
}

func setupUserRoute(api IUserController, router *gin.Engine) {
	userGroup := router.Group("/user")
	userGroup.GET("/login", api.getLogin)
	userGroup.GET("/register", api.getRegister)

	userGroup.POST("/login", api.postLogin)
	userGroup.POST("/register", api.postRegister)
	userGroup.POST("/logout", api.postLogout)
}

func (u *userController) getLogin(c *gin.Context) {
	redirectIfLoggedIn(c)

	c.HTML(http.StatusOK, "login.html", gin.H{
		"title":        "Login",
		"register_url": "/user/register",
	})
}

func (u *userController) getRegister(c *gin.Context) {
	redirectIfLoggedIn(c)

	c.HTML(http.StatusOK, "register.html", gin.H{
		"title":     "Register",
		"login_url": "/user/login",
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

	secure.SetAuthCookie(c, user.ID)
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

	secure.SetAuthCookie(c, user.ID)
	c.Redirect(http.StatusSeeOther, "/note")
}

func redirectIfLoggedIn(c *gin.Context) {
	user, err := secure.GetLoggedInUser(c)
	if err == nil && user != nil {
		slog.Info("user already logged in", "user_id", user.ID, "email", user.Email)
		c.Redirect(http.StatusSeeOther, "/note")
		return
	}
}

func (u *userController) postLogout(c *gin.Context) {
	secure.ClearAuthCookie(c)
	slog.Info("user logged out")
	c.Redirect(http.StatusSeeOther, "/user/login")
}
