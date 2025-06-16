package controller

import (
	"github.com/ToshihiroOgino/elib/controller/page"
	"github.com/ToshihiroOgino/elib/controller/rest"
	"github.com/gin-gonic/gin"
)

type controller struct {
	userApi  rest.IUserRest
	userPage page.IUserPage
}

func NewController() *controller {
	return &controller{
		userApi:  rest.NewUserRest(),
		userPage: page.NewUserPage(),
	}
}

func (c *controller) SetupRoutes(router *gin.Engine) {
	page.SetupRoute(c.userPage, router)
	rest.SetupRoute(c.userApi, router)
}
