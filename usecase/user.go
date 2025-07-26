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
	Validate(user *domain.User, email string, password string) (bool, error)
	FindByEmail(email string) (*domain.User, error)
}

type userUsecase struct {
	db *gorm.DB
}

func NewUserUsecase() IUserUsecase {
	db := sqlite.GetDB()
	return &userUsecase{
		db: db,
	}
}

func (u *userUsecase) newQuery() (*repository.Query, repository.IUserDo) {
	q := repository.Use(u.db)
	do := q.User.WithContext(u.db.Statement.Context)
	return q, do
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

	_, repo := u.newQuery()
	if err := repo.Create(user); err != nil {
		slog.Error("failed to create user", "error", err)
		return nil, err
	}
	return user, nil
}

func (u *userUsecase) Validate(user *domain.User, email string, password string) (bool, error) {
	if !strings.EqualFold(*user.Email, email) {
		return false, errors.New("invalid email")
	}

	if user.PasswordHash == nil {
		slog.Error("user data is incomplete", "email", email)
		return false, errors.New("user data is incomplete")
	}

	err := bcrypt.CompareHashAndPassword(*user.PasswordHash, []byte(password))
	if err != nil {
		slog.Error("invalid password", "email", email, "error", err)
		return false, errors.New("invalid password")
	}

	return true, nil
}

func (u *userUsecase) FindByEmail(email string) (*domain.User, error) {
	q, repo := u.newQuery()
	user, err := repo.Where(q.User.Email.Eq(email)).First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			slog.Error("user not found", "email", email)
			return nil, errors.New("user not found")
		}
		slog.Error("failed to get user by email", "email", email, "error", err)
		return nil, err
	}
	return user, nil
}
