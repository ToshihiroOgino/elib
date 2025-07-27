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
	postShareNote(c *gin.Context)
	deleteShare(c *gin.Context)
}

type shareController struct {
	shareUsecase usecase.IShareUsecase
	noteUsecase  usecase.INoteUsecase
	userUsecase  usecase.IUserUsecase
}

type shareRequest struct {
	NoteID   string `json:"noteId"`
	Editable bool   `json:"editable"`
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
		shareGroup.POST("", i.postShareNote)
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

func (i *shareController) postShareNote(c *gin.Context) {
	var req shareRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data."})
		return
	}

	note, err := i.noteUsecase.Find(req.NoteID)
	if err != nil || note == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Note not found."})
		return
	}

	user := auth.GetSessionUser(c)
	if user == nil || user.ID != note.AuthorID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to share this note."})
		return
	}

	sharingInfo, err := i.shareUsecase.ShareNote(note, req.Editable)
	if err != nil {
		slog.Error("failed to share note", "noteId", note.ID, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to share note."})
		return
	}
	c.JSON(http.StatusOK, gin.H{"shareId": sharingInfo.ID})
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
