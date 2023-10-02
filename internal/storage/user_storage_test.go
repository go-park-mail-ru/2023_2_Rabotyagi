package storage_test

import (
    "testing"

	"github.com/stretchr/testify/assert"
    "github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/storage"
)

func TestUserStorageMapGetUser(t *testing.T) {
    storageMap := storage.NewAuthStorageMap()
    preUser := &storage.PreUser{
        Email: "test@gmail.com",
		Password: "testpassword",
    }

    _ = storageMap.CreateUser(preUser)

    // Проверяем, что пост был добавлен успешно
    expectedID := 1
    assert.Equal(t, expectedID, storageMap.GetUsersCount())
    
	userInMap, _ := storageMap.GetUser(preUser.Email)
    // Проверяем, что добавление существующего поста заменяет его значениями нового поста
    assert.Equal(t, uint64(1), userInMap.ID)
	assert.Equal(t, preUser.Email, userInMap.Email)
	assert.Equal(t, preUser.Password, userInMap.Password)
}

func TestUserStorageMapIsUserExist(t *testing.T) {
    storageMap := storage.NewAuthStorageMap()
    preUser := &storage.PreUser{
        Email: "test@gmail.com",
		Password: "testpassword",
    }

    _ = storageMap.CreateUser(preUser)

    assert.Equal(t, true, storageMap.IsUserExist("test@gmail.com"))
}
