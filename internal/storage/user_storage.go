package storage

import (
	"fmt"

	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/errors"
)

var (
	ErrUserAlreadyExist = errors.NewError("user already exist")
	ErrUserNotExist     = errors.NewError("user not exist")
)

type User struct {
	ID       uint64
	Name     string `json:"name"`
	Password string `json:"password"`
}

type PreUser struct {
	Name     string `json:"name"`
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
}

func NewAuthStorageMap() *AuthStorageMap {
	return &AuthStorageMap{
		counterUsers: 0,
		users: make(map[string]User),
	}
}

func (a *AuthStorageMap) generateIDCurUser() uint64 {
	a.counterUsers++

	return a.counterUsers
}

func (a *AuthStorageMap) GetUser(username string) (*User, error) {
	if !a.IsUserExist(username) {
		return nil, fmt.Errorf("username ==%s %w", username, ErrUserNotExist)
	}

	user := a.users[username]

	return &user, nil
}

func (a *AuthStorageMap) CreateUser(user *PreUser) error {
	if a.IsUserExist(user.Name) {
		return fmt.Errorf("username ==%s %w", user.Name, ErrUserAlreadyExist)
	}

	a.users[user.Name] = User{ID: a.generateIDCurUser(), Name: user.Name, Password: user.Password}

	return nil
}

func (a *AuthStorageMap) IsUserExist(username string) bool {
	_, ok := a.users[username]

	return ok
}
