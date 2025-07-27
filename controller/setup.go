package controller

import (
	"github.com/gin-gonic/gin"
)

type controller struct {
	user  IUserController
	note  INoteController
	share IShareController
}

func NewController(router *gin.Engine) *controller {
	return &controller{
		user:  NewUserController(router),
		note:  NewNoteController(router),
		share: NewShareController(router),
	}
}
