package usecase

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/ToshihiroOgino/elib/generated/generated/domain"
	"github.com/ToshihiroOgino/elib/generated/repository"
	"github.com/ToshihiroOgino/elib/infra/sqlite"
	"gorm.io/gorm"
)

type INoteUsecase interface {
	CreateNote(user *domain.User) (*domain.Note, error)
	UpdateNote(note *domain.Note) (*domain.Note, error)
	Find(noteId string) (*domain.Note, error)
	FindNotesByUserID(userID string) ([]*domain.Note, error)
	Delete(note *domain.Note) error
}

type noteUsecase struct {
	db *gorm.DB
}

func NewNoteUsecase() INoteUsecase {
	db := sqlite.GetDB()
	return &noteUsecase{
		db: db,
	}
}

func (n *noteUsecase) newQuery() (*repository.Query, repository.INoteDo) {
	q := repository.Use(n.db)
	do := q.Note.WithContext(n.db.Statement.Context)
	return q, do
}

func defaultTitle() string {
	timestamp := time.Now().Format("20060102150405")
	title := fmt.Sprintf("Untitled_%s", timestamp)
	return title
}

func (n *noteUsecase) CreateNote(user *domain.User) (*domain.Note, error) {
	note := &domain.Note{
		ID:       newUUID(),
		Title:    defaultTitle(),
		Content:  "",
		AuthorID: user.ID,
	}
	_, do := n.newQuery()
	if err := do.Create(note); err != nil {
		slog.Error("failed to create note", "error", err)
		return nil, err
	}
	return note, nil
}

func (n *noteUsecase) UpdateNote(note *domain.Note) (*domain.Note, error) {
	_, do := n.newQuery()
	if err := do.Save(note); err != nil {
		return nil, err
	}
	return note, nil
}

func (n *noteUsecase) Find(noteId string) (*domain.Note, error) {
	q, do := n.newQuery()
	note, err := do.Where(q.Note.ID.Eq(noteId)).First()
	if err != nil {
		return nil, err
	}
	return note, nil
}

func (n *noteUsecase) FindNotesByUserID(userID string) ([]*domain.Note, error) {
	q, do := n.newQuery()
	notes, err := do.Where(q.Note.AuthorID.Eq(userID)).Find()
	if err != nil {
		return nil, err
	}
	return notes, nil
}

func (n *noteUsecase) Delete(note *domain.Note) error {
	q, do := n.newQuery()
	res, err := do.Where(q.Note.ID.Eq(note.ID)).Delete()
	if err != nil {
		return err
	}
	if res.Error != nil {
		return res.Error
	}
	slog.Info("Note deleted successfully", "noteID", note.ID, "rowsAffected", res.RowsAffected)
	return nil
}
