package user

import (
	"github.com/ToshihiroOgino/elib/infra/sqlite"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository interface {
	FindByID(id uuid.UUID) (*User, error)
	FindByEmail(email string) (*User, error)
	FindAll() ([]*User, error)
	Save(user *User) error
	Delete(user *User) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository() UserRepository {
	return &userRepository{db: sqlite.GetDB()}
}

func (r *userRepository) FindByID(id uuid.UUID) (*User, error) {
	var user User
	if err := r.db.First(&User{
		ID: id,
	}).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByEmail(email string) (*User, error) {
	var user User
	if err := r.db.First(&User{
		Email: email,
	}).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindAll() ([]*User, error) {
	var users []*User
	if err := r.db.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (r *userRepository) Save(user *User) error {
	if user.ID == uuid.Nil {
		user.ID = uuid.New()
	}
	if err := r.db.Save(user).Error; err != nil {
		return err
	}
	return nil
}

func (r *userRepository) Delete(user *User) error {
	if user.ID == uuid.Nil {
		return nil
	}
	if err := r.db.Delete(user).Error; err != nil {
		return err
	}
	return nil
}
