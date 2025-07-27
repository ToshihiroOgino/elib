package controller

import (
	"log/slog"
	"net/http"

	"github.com/ToshihiroOgino/elib/auth"
	"github.com/ToshihiroOgino/elib/usecase"
	"github.com/gin-gonic/gin"
)

type IShareController interface {
	getSharedNote(c *gin.Context)
	deleteShare(c *gin.Context)
}

type shareController struct {
	shareUsecase usecase.IShareUsecase
	noteUsecase  usecase.INoteUsecase
	userUsecase  usecase.IUserUsecase
}

func NewShareController(router *gin.Engine) IShareController {
	instance := &shareController{
		shareUsecase: usecase.NewShareUsecase(),
		noteUsecase:  usecase.NewNoteUsecase(),
		userUsecase:  usecase.NewUserUsecase(),
	}
	setupShareRoute(instance, router)
	return instance
}

func setupShareRoute(i IShareController, router *gin.Engine) {
	shareGroup := router.Group("/share")
	shareGroup.GET("/:id", i.getSharedNote)
	shareGroup.Use(auth.AuthMiddleware())
	{
		shareGroup.DELETE("/:id", i.deleteShare)
	}
}

func (i *shareController) getSharedNote(c *gin.Context) {
	shareId := c.Param("id")
	share, err := i.shareUsecase.Find(shareId)
	if err != nil || share == nil {
		showNotFoundPage(c)
		return
	}

	note, err := i.noteUsecase.Find(share.NoteID)
	if err != nil || note == nil {
		showNotFoundPage(c)
		return
	}

	user := auth.GetSessionUser(c)
	if user == nil {
		user = i.userUsecase.CreateGuestUser()
	}

	// 自身のメモであれば編集ページにリダイレクトする
	if note.AuthorID == user.ID {
		c.Redirect(http.StatusSeeOther, "/note/"+note.ID)
	}
	c.HTML(http.StatusOK, "shared_note.html", gin.H{
		"title":    "Shared Note",
		"note":     note,
		"editable": share.Editable,
	})
}

func (i *shareController) deleteShare(c *gin.Context) {
	shareId := c.Param("id")
	share, err := i.shareUsecase.Find(shareId)
	if err != nil || share == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Share not found."})
		return
	}

	note, err := i.noteUsecase.Find(share.NoteID)
	if err != nil || note == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Note not found."})
		return
	}

	user := auth.GetSessionUser(c)
	if user == nil || user.ID != note.AuthorID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to delete this share."})
		return
	}

	if err := i.shareUsecase.Delete(share); err != nil {
		slog.Error("failed to delete share", "shareId", shareId, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete share."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Share deleted successfully."})
}

func showNotFoundPage(c *gin.Context) {
	c.HTML(404, "not_found.html", gin.H{
		"title": "Not Found",
	})
}
