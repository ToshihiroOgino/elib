package usecase

import (
	"errors"
	"log/slog"

	"github.com/ToshihiroOgino/elib/generated/generated/domain"
	"github.com/ToshihiroOgino/elib/generated/repository"
	"github.com/ToshihiroOgino/elib/infra/sqlite"
	"gorm.io/gorm"
)

type IShareUsecase interface {
	ShareNote(note *domain.Note, editable bool) (*domain.SharingInfo, error)
	FindByNote(note *domain.Note) ([]*domain.SharingInfo, error)
	Delete(share *domain.SharingInfo) error
	Find(shareId string) (*domain.SharingInfo, error)
}

type shareUsecase struct {
	db *gorm.DB
}

func NewShareUsecase() IShareUsecase {
	db := sqlite.GetDB()
	return &shareUsecase{
		db: db,
	}
}

func (s *shareUsecase) newQuery() (*repository.Query, repository.ISharingInfoDo) {
	q := repository.Use(s.db)
	do := q.SharingInfo.WithContext(s.db.Statement.Context)
	return q, do
}

func (s *shareUsecase) ShareNote(note *domain.Note, editable bool) (*domain.SharingInfo, error) {
	if note == nil {
		return nil, errors.New("note cannot be nil")
	}

	sharingInfo := &domain.SharingInfo{
		ID:       newUUID(),
		NoteID:   note.ID,
		Editable: editable,
	}

	_, do := s.newQuery()
	if err := do.Create(sharingInfo); err != nil {
		return nil, err
	}

	return sharingInfo, nil
}

func (s *shareUsecase) FindByNote(note *domain.Note) ([]*domain.SharingInfo, error) {
	if note == nil {
		return nil, errors.New("note cannot be nil")
	}

	q, do := s.newQuery()
	res, err := do.Where(q.SharingInfo.NoteID.Eq(note.ID)).Find()
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *shareUsecase) Delete(share *domain.SharingInfo) error {
	if share == nil {
		return errors.New("share cannot be nil")
	}
	q, do := s.newQuery()
	res, err := do.Where(q.SharingInfo.ID.Eq(share.ID)).Delete()
	if err != nil {
		return err
	}
	slog.Info("deleted sharing info", "id", share.ID, "rowsAffected", res.RowsAffected)
	return nil
}

func (s *shareUsecase) Find(shareId string) (*domain.SharingInfo, error) {
	if shareId == "" {
		return nil, errors.New("shareId cannot be empty")
	}

	q, do := s.newQuery()
	share, err := do.Where(q.SharingInfo.ID.Eq(shareId)).First()
	if err != nil {
		return nil, err
	}

	return share, nil
}
