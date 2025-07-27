package controller

import (
	"log/slog"
	"net/http"

	"github.com/ToshihiroOgino/elib/auth"
	"github.com/ToshihiroOgino/elib/generated/generated/domain"
	"github.com/ToshihiroOgino/elib/usecase"
	"github.com/gin-gonic/gin"
)

type INoteController interface {
	getNote(c *gin.Context)
	getNoteById(c *gin.Context)
	getCreateNewNote(c *gin.Context)
	postSaveNote(c *gin.Context)
	deleteNote(c *gin.Context)
}

type noteController struct {
	usecase usecase.INoteUsecase
}

type saveNoteRequest struct {
	ID      string `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

func NewNoteController(router *gin.Engine) INoteController {
	instance := &noteController{
		usecase: usecase.NewNoteUsecase(),
	}
	setupNoteRoute(instance, router)
	return instance
}

func setupNoteRoute(api INoteController, router *gin.Engine) {
	noteGroup := router.Group("/note")
	noteGroup.Use(auth.AuthMiddleware())
	{
		noteGroup.GET("", api.getNote)
		noteGroup.GET("/:id", api.getNoteById)
		noteGroup.GET("/new", api.getCreateNewNote)
		noteGroup.POST("/save", api.postSaveNote)
		noteGroup.DELETE("/delete/:id", api.deleteNote)
	}
}

func (n *noteController) getNote(c *gin.Context) {
	user := auth.GetSessionUser(c)

	// ユーザーの全メモを取得
	notes, err := n.usecase.FindNotesByUserID(user.ID)
	if err != nil {
		slog.Error("failed to get notes", "error", err)
		notes = []*domain.Note{}
	}

	// 最初のメモを選択、なければ新規作成
	var currentNote *domain.Note
	if len(notes) > 0 {
		currentNote = notes[0]
	} else {
		newNote, err := n.usecase.CreateNote(user)
		if err != nil {
			slog.Error("failed to save new note", "error", err)
		} else {
			currentNote = newNote
			notes = []*domain.Note{currentNote}
		}
	}

	c.HTML(http.StatusOK, "editor.html", gin.H{
		"title": "メモエディター",
		"note":  currentNote,
		"notes": notes,
	})
}

func (n *noteController) getNoteById(c *gin.Context) {
	user := auth.GetSessionUser(c)
	noteId := c.Param("id")

	note, err := n.usecase.Find(noteId)
	if err != nil {
		slog.Error("failed to get note", "noteId", noteId, "error", err)
		c.Redirect(http.StatusSeeOther, "/note")
		return
	}

	if note.AuthorID != user.ID {
		c.Redirect(http.StatusSeeOther, "/note")
		return
	}

	// ユーザーの全メモを取得
	notes, err := n.usecase.FindNotesByUserID(user.ID)
	if err != nil {
		slog.Error("failed to get notes", "error", err)
		notes = []*domain.Note{}
	}

	c.HTML(http.StatusOK, "editor.html", gin.H{
		"title": "メモエディター",
		"note":  note,
		"notes": notes,
	})
}

func (n *noteController) getCreateNewNote(c *gin.Context) {
	user := auth.GetSessionUser(c)

	// 新規メモを作成
	newNote, err := n.usecase.CreateNote(user)
	if err != nil {
		slog.Error("failed to create new note", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create note"})
		return
	}

	c.Redirect(http.StatusSeeOther, "/note/"+newNote.ID)
}

func (n *noteController) postSaveNote(c *gin.Context) {
	user := auth.GetSessionUser(c)

	var req saveNoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Error("failed to bind save note request", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	// メモを取得して権限チェック
	note, err := n.usecase.Find(req.ID)
	if err != nil {
		slog.Error("failed to get note for saving", "noteId", req.ID, "error", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "note not found"})
		return
	}

	if note.AuthorID != user.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	// メモを更新
	note.Title = req.Title
	note.Content = req.Content

	_, err = n.usecase.UpdateNote(note)
	if err != nil {
		slog.Error("failed to update note", "noteId", req.ID, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save note"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

func (n *noteController) deleteNote(c *gin.Context) {
	user := auth.GetSessionUser(c)
	noteId := c.Param("id")

	// メモを取得して権限チェック
	note, err := n.usecase.Find(noteId)
	if err != nil {
		slog.Error("failed to get note for deletion", "noteId", noteId, "error", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "note not found"})
		return
	}

	if note.AuthorID != user.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	// メモを削除
	err = n.usecase.Delete(note)
	if err != nil {
		slog.Error("failed to delete note", "noteId", noteId, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete note"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}
