package controller

import (
	"github.com/gin-gonic/gin"
)

type controller struct {
	user IUserController
	note INoteController
}

func NewController(router *gin.Engine) *controller {
	return &controller{
		user: NewUserController(router),
		note: NewNoteController(router),
	}
}
