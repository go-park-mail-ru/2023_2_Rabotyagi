package storage

import (
	"fmt"
	"sync"

	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/errors"
)

var (
	ErrUserAlreadyExist = errors.NewError("user already exist")
	ErrUserNotExist     = errors.NewError("user not exist")
)

type User struct {
	ID       uint64
	Email    string `json:"email"`
	Password string `json:"password"`
}

type PreUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthStorage interface {
	GetUser(username string) (*User, error)
	CreateUser(user *PreUser) error
	IsUserExist(username string) bool
}

type AuthStorageMap struct {
	counterUsers uint64
	users        map[string]User
	mu           sync.RWMutex
}

func NewAuthStorageMap() *AuthStorageMap {
	return &AuthStorageMap{
		counterUsers: 0,
		users:        make(map[string]User),
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

func (a *AuthStorageMap) GetUser(email string) (*User, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	if !a.IsUserExist(email) {
		return nil, fmt.Errorf("username ==%s %w", email, ErrUserNotExist)
	}

	user := a.users[email]

	return &user, nil
}

func (a *AuthStorageMap) CreateUser(user *PreUser) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.IsUserExist(user.Email) {
		return fmt.Errorf("email ==%s %w", user.Email, ErrUserAlreadyExist)
	}

	a.users[user.Email] = User{ID: a.generateIDCurUser(), Email: user.Email, Password: user.Password}

	return nil
}

func (a *AuthStorageMap) IsUserExist(email string) bool {
	_, ok := a.users[email]

	return ok
}
