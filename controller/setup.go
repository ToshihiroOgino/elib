package controller

import (
	"github.com/gin-gonic/gin"
)

type controller struct {
	user  IUserController
	note  INoteController
	share IShareController
}

func showNotFoundPage(c *gin.Context) {
	c.HTML(404, "not_found.html", gin.H{
		"title": "Not Found",
	})
}

func setNoRoute(router *gin.Engine) {
	router.NoRoute(showNotFoundPage)
}

func NewController(router *gin.Engine) *controller {
	setNoRoute(router)
	return &controller{
		user:  NewUserController(router),
		note:  NewNoteController(router),
		share: NewShareController(router),
	}
}
