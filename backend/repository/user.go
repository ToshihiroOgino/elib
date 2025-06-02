package repository

import (
	"github.com/ToshihiroOgino/elib/backend/domain"
	"github.com/ToshihiroOgino/elib/backend/infra"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository interface {
	FindByID(id uuid.UUID) (*domain.User, error)
	FindByEmail(email string) (*domain.User, error)
	FindAll() ([]*domain.User, error)
	Save(user *domain.User) error
	Delete(user *domain.User) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository() UserRepository {
	return &userRepository{db: infra.GetDB()}
}

func (r *userRepository) FindByID(id uuid.UUID) (*domain.User, error) {
	var user domain.User
	if err := r.db.First(&domain.User{
		ID: id,
	}).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByEmail(email string) (*domain.User, error) {
	var user domain.User
	if err := r.db.First(&domain.User{
		Email: email,
	}).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindAll() ([]*domain.User, error) {
	var users []*domain.User
	if err := r.db.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (r *userRepository) Save(user *domain.User) error {
	if user.ID == uuid.Nil {
		user.ID = uuid.New()
	}
	if err := r.db.Save(user).Error; err != nil {
		return err
	}
	return nil
}

func (r *userRepository) Delete(user *domain.User) error {
	if user.ID == uuid.Nil {
		return nil // No action needed if user ID is not set
	}
	if err := r.db.Delete(user).Error; err != nil {
		return err
	}
	return nil
}
