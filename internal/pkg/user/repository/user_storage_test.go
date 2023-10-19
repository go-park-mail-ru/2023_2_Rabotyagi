package repository_test

import (
	"testing"

	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/models"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/user/repository"
	"github.com/stretchr/testify/assert"
)

func TestUserStorageMapGetUser(t *testing.T) {
	storageMap := repository.NewAuthStorageMap()
	preUser := &models.PreUser{
		Email:    "test@gmail.com",
		Password: "testpassword",
	}

	_ = storageMap.CreateUser(preUser)

	// Проверяем, что пост был добавлен успешно
	expectedID := 1
	assert.Equal(t, expectedID, storageMap.GetUsersCount())

	userInMap, _ := storageMap.GetUser(preUser.Email)

	assert.Equal(t, uint64(1), userInMap.ID)
	assert.Equal(t, preUser.Email, userInMap.Email)
	assert.Equal(t, preUser.Password, userInMap.Password)
}

func TestUserStorageMapGetUserError(t *testing.T) {
	storageMap := repository.NewAuthStorageMap()
	preUser := &models.PreUser{
		Email:    "test@gmail.com",
		Password: "testpassword",
	}

	_ = storageMap.CreateUser(preUser)

	_, err := storageMap.GetUser("non-existen-email@gmail.com")
	assert.NotNil(t, err)
}

func TestUserStorageMapIsUserExist(t *testing.T) {
	storageMap := repository.NewAuthStorageMap()
	preUser := &models.PreUser{
		Email:    "test@gmail.com",
		Password: "testpassword",
	}

	_ = storageMap.CreateUser(preUser)

	assert.Equal(t, true, storageMap.IsUserExist("test@gmail.com"))
}

func TestUserStorageMapCreateUserError(t *testing.T) {
	storageMap := repository.NewAuthStorageMap()
	preUser1 := &models.PreUser{
		Email:    "test@gmail.com",
		Password: "testpassword",
	}

	_ = storageMap.CreateUser(preUser1)

	preUser2 := &models.PreUser{
		Email:    "test@gmail.com",
		Password: "newpassword",
	}

	err := storageMap.CreateUser(preUser2)

	assert.NotNil(t, err)
}
