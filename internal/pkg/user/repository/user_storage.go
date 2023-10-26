package repository

import (
	"fmt"
	"sync"

	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/models"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/errors"
)

var (
	ErrUserAlreadyExist = errors.NewError("user already exist")
	ErrUserNotExist     = errors.NewError("user not exist")
)

type AuthStorage interface {
	GetUser(username string) (*models.User, error)
	CreateUser(user *models.PreUser) error
	IsUserExist(username string) bool
}

type AuthStorageMap struct {
	counterUsers uint64
	users        map[string]models.User
	mu           sync.RWMutex
}

func NewAuthStorageMap() *AuthStorageMap {
	return &AuthStorageMap{
		counterUsers: 0,
		users:        make(map[string]models.User),
		mu:           sync.RWMutex{},
	}
}

func (a *AuthStorageMap) GetUsersCount() int {
	a.mu.RLock()
	defer a.mu.RUnlock()

	return len(a.users)
}

func (a *AuthStorageMap) generateIDCurUser() uint64 {
	a.counterUsers++

	return a.counterUsers
}

func (a *AuthStorageMap) GetUser(email string) (*models.User, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	if !a.IsUserExist(email) {
		return nil, fmt.Errorf("username ==%s %w", email, ErrUserNotExist)
	}

	user := a.users[email]

	return &user, nil
}

func (a *AuthStorageMap) CreateUser(user *models.PreUser) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.IsUserExist(user.Email) {
		return fmt.Errorf("email ==%s %w", user.Email, ErrUserAlreadyExist)
	}

	a.users[user.Email] = models.User{ID: a.generateIDCurUser(), Email: user.Email, Pass: user.Password}

	return nil
}

func (a *AuthStorageMap) IsUserExist(email string) bool {
	_, ok := a.users[email]

	return ok
}
