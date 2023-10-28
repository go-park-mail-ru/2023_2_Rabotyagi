package repository

import (
	"context"
	"fmt"
	"sync"

	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/models"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/errors"
)

var (
	ErrUserNotExist = errors.NewError("user not exist")
	ErrEmailBusy    = errors.NewError("same email already in use")
	ErrPhoneBusy    = errors.NewError("same phone already in use")
)

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

func (a *AuthStorageMap) GetUserByID(ctx context.Context, id uint64) (*models.User, error) {
	a.mu.RLock()
	for _, user := range a.users {
		if id == user.ID {
			return &user, nil
		}
	}
	a.mu.RUnlock()

	return nil, ErrUserNotExist
}

func (a *AuthStorageMap) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	if !a.isEmailBusy(ctx, email) {
		return nil, fmt.Errorf("username ==%s %w", email, ErrUserNotExist)
	}

	user := a.users[email]

	return &user, nil
}

func (a *AuthStorageMap) CreateUser(ctx context.Context, user *models.UserWithoutID) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.isEmailBusy(ctx, user.Email) {
		return fmt.Errorf("email ==%s %w", user.Email, ErrEmailBusy)
	}

	a.users[user.Email] = models.User{
		ID:       a.generateIDCurUser(),
		Name:     user.Name,
		Email:    user.Email,
		Phone:    user.Phone,
		Password: user.Password,
		Birthday: user.Birthday,
	}

	return nil
}

func (a *AuthStorageMap) isEmailBusy(ctx context.Context, email string) bool {
	_, ok := a.users[email]

	return ok
}

func (a *AuthStorageMap) IsEmailBusy(ctx context.Context, email string) bool {
	a.mu.RLock()
	ok := a.isEmailBusy(ctx, email)
	a.mu.RUnlock()

	return ok
}

func (a *AuthStorageMap) isPhoneBusy(ctx context.Context, phone string) bool {
	for _, user := range a.users {
		if phone == user.Phone {
			return true
		}
	}

	return false
}

func (a *AuthStorageMap) IsPhoneBusy(ctx context.Context, phone string) bool {
	a.mu.RLock()
	busy := a.isPhoneBusy(ctx, phone)
	a.mu.RUnlock()

	return busy
}
