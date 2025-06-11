package usecase

import (
	"errors"
	"log/slog"
	"strings"

	"github.com/ToshihiroOgino/elib/generated/generated/domain"
	"github.com/ToshihiroOgino/elib/generated/repository"
	"github.com/ToshihiroOgino/elib/infra/sqlite"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type IUserUsecase interface {
	Create(email string, password string) (*domain.User, error)
}

type userUsecase struct {
	db *gorm.DB
}

func isValidEmail(email string) bool {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false
	}
	localPart := parts[0]
	domainPart := parts[1]
	if len(localPart) == 0 || len(domainPart) == 0 {
		return false
	}
	if strings.ContainsAny(localPart, "!#$%&'*+/=?^`{|}~ ") || strings.ContainsAny(domainPart, "!#$%&'*+/=?^`{|}~ ") {
		return false
	}
	if strings.Contains(localPart, "..") || strings.Contains(domainPart, "..") {
		return false
	}
	return true
}

func hashPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func Create(
	email string,
	password string,
) (*domain.User, error) {
	id := newUUID()
	if !isValidEmail(email) {
		slog.Error("invalid email format", "email", email)
		return nil, errors.New("invalid email format")
	}
	passwordHash, err := hashPassword(password)
	if err != nil {
		slog.Error("failed to hash password", "error", err)
		return nil, err
	}
	user := &domain.User{
		ID:           &id,
		Email:        &email,
		PasswordHash: &passwordHash,
	}
	db := sqlite.GetDB()
	q := repository.Use(db).User
	if err := q.WithContext(db.Statement.Context).Create(user); err != nil {
		slog.Error("failed to create user", "error", err)
		return nil, err
	}
	return user, nil
}

func (u *userUsecase) Create(email string, password string) (*domain.User, error) {
	id := newUUID()
	if !isValidEmail(email) {
		slog.Error("invalid email format", "email", email)
		return nil, errors.New("invalid email format")
	}
	passwordHash, err := hashPassword(password)
	if err != nil {
		slog.Error("failed to hash password", "error", err)
		return nil, err
	}
	user := &domain.User{
		ID:           &id,
		Email:        &email,
		PasswordHash: &passwordHash,
	}
	q := repository.Use(u.db).User
	if err := q.WithContext(u.db.Statement.Context).Create(user); err != nil {
		slog.Error("failed to create user", "error", err)
		return nil, err
	}
	return user, nil
}
