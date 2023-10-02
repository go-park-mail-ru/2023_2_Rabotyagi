package storage_test

import (
    "testing"

	"github.com/stretchr/testify/assert"
    //"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/errors"
    "github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/storage"
)

func TestPostStorageMapGet(t *testing.T) {
    storageMap := storage.NewPostStorageMap()
    prePost := &storage.PrePost{
        AuthorID:        1,
        Title:           "Test Post",
        Description:     "This is a test post",
        Price:           100,
        SafeTransaction: true,
        Delivery:        true,
        City:            "Moscow",
    }

    storageMap.AddPost(prePost)

    // Проверяем, что пост был добавлен успешно
    expectedID := 1
    assert.Equal(t, expectedID, storageMap.GetPostsCount())
    
	postInMap, _ := storageMap.GetPost(uint64(1))
    // Проверяем, что добавление существующего поста заменяет его значениями нового поста
    assert.Equal(t, uint64(1), postInMap.ID)
	assert.Equal(t, prePost.AuthorID, postInMap.AuthorID)
	assert.Equal(t, prePost.Description, postInMap.Description)
	assert.Equal(t, prePost.Price, postInMap.Price)
	assert.Equal(t, prePost.Price, postInMap.Price)
	assert.Equal(t, prePost.SafeTransaction, postInMap.SafeTransaction)
	assert.Equal(t, prePost.Delivery, postInMap.Delivery)
	assert.Equal(t, prePost.City, postInMap.City)
}

func TestPostStorageMapGetList(t *testing.T) {
    storageMap := storage.NewPostStorageMap()
    prePost1 := &storage.PrePost{
        AuthorID:        1,
        Title:           "Test Post",
        Description:     "This is a test post",
        Price:           100,
        SafeTransaction: true,
        Delivery:        true,
        City:            "Moscow",
    }

    storageMap.AddPost(prePost1)

	prePost2 := &storage.PrePost{
        AuthorID:        1,
        Title:           "Test Post",
        Description:     "This is a test post",
        Price:           100,
        SafeTransaction: true,
        Delivery:        true,
        City:            "Moscow",
    }

	storageMap.AddPost(prePost2)
    
	postInMap1, _ := storageMap.GetPost(uint64(1))
    // Проверяем, что добавление существующего поста заменяет его значениями нового поста
    assert.Equal(t, uint64(1), postInMap1.ID)
	assert.Equal(t, prePost1.AuthorID, postInMap1.AuthorID)
	assert.Equal(t, prePost1.Description, postInMap1.Description)
	assert.Equal(t, prePost1.Price, postInMap1.Price)
	assert.Equal(t, prePost1.Price, postInMap1.Price)
	assert.Equal(t, prePost1.SafeTransaction, postInMap1.SafeTransaction)
	assert.Equal(t, prePost1.Delivery, postInMap1.Delivery)
	assert.Equal(t, prePost1.City, postInMap1.City)

	postInMap2, _ := storageMap.GetPost(uint64(2))
    // Проверяем, что добавление существующего поста заменяет его значениями нового поста
    assert.Equal(t, uint64(2), postInMap2.ID)
	assert.Equal(t, prePost2.AuthorID, postInMap2.AuthorID)
	assert.Equal(t, prePost2.Description, postInMap2.Description)
	assert.Equal(t, prePost2.Price, postInMap2.Price)
	assert.Equal(t, prePost2.Price, postInMap2.Price)
	assert.Equal(t, prePost2.SafeTransaction, postInMap2.SafeTransaction)
	assert.Equal(t, prePost2.Delivery, postInMap2.Delivery)
	assert.Equal(t, prePost2.City, postInMap2.City)
}