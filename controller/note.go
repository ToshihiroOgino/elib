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
	createNewNote(c *gin.Context)
	saveNote(c *gin.Context)
	deleteNote(c *gin.Context)
	viewNote(c *gin.Context)
}

type noteController struct {
	usecase usecase.INoteUsecase
}

type saveNoteRequest struct {
	ID      string `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

const _URL_NOTE_ROOT = "/note"

func NewNoteController(router *gin.Engine) INoteController {
	instance := &noteController{
		usecase: usecase.NewNoteUsecase(),
	}
	setupNoteRoute(instance, router)
	return instance
}

func setupNoteRoute(api INoteController, router *gin.Engine) {
	noteGroup := router.Group(_URL_NOTE_ROOT)
	noteGroup.Use(auth.AuthMiddleware())
	{
		noteGroup.GET("", api.getNote)
		noteGroup.GET("/new", api.createNewNote)
		noteGroup.GET("/:id", api.getNoteById)
		noteGroup.POST("/save", api.saveNote)
		noteGroup.DELETE("/delete/:id", api.deleteNote)
	}

	// 共有用の認証なしルート
	router.GET(_URL_NOTE_ROOT+"/view/:id", api.viewNote)
}

func (n *noteController) getNote(c *gin.Context) {
	user := auth.GetUser(c)

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
		currentNote = n.usecase.NewNote(user)
		// 新規メモを保存
		savedNote, err := n.usecase.UpdateNote(currentNote)
		if err != nil {
			slog.Error("failed to save new note", "error", err)
		} else {
			currentNote = savedNote
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
	user := auth.GetUser(c)
	noteId := c.Param("id")

	// 指定されたメモを取得
	note, err := n.usecase.Find(noteId)
	if err != nil {
		slog.Error("failed to get note", "noteId", noteId, "error", err)
		c.Redirect(http.StatusSeeOther, _URL_NOTE_ROOT)
		return
	}

	// 権限チェック
	if note.AuthorID != user.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
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

func (n *noteController) createNewNote(c *gin.Context) {
	user := auth.GetUser(c)

	// 新規メモを作成
	newNote := n.usecase.NewNote(user)
	savedNote, err := n.usecase.UpdateNote(newNote)
	if err != nil {
		slog.Error("failed to create new note", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create note"})
		return
	}

	c.Redirect(http.StatusSeeOther, _URL_NOTE_ROOT+"/"+savedNote.ID)
}

func (n *noteController) saveNote(c *gin.Context) {
	user := auth.GetUser(c)

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
	user := auth.GetUser(c)
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

func (n *noteController) viewNote(c *gin.Context) {
	noteId := c.Param("id")

	// メモを取得（共有表示用）
	note, err := n.usecase.Find(noteId)
	if err != nil {
		slog.Error("failed to get note for viewing", "noteId", noteId, "error", err)
		c.HTML(http.StatusNotFound, "404.html", gin.H{
			"title": "メモが見つかりません",
		})
		return
	}

	c.HTML(http.StatusOK, "view.html", gin.H{
		"title": note.Title,
		"note":  note,
	})
}
