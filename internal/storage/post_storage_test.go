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
    expectedID := uint64(1)
    assert.Len(t, expectedID, storageMap.GetPostsCount())
    
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
	
	prePostWithExistingID := &storage.PrePost{
        AuthorID:        2,
        Title:           "Another Post",
        Description:     "This is another test post",
        Price:           200,
        SafeTransaction: false,
        Delivery:        false,
        City:            "New York",
    }

    storageMap.AddPost(prePostWithExistingID)
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
    
	postInMap, _ := storageMap.GetPost(uint64(1))
    // Проверяем, что добавление существующего поста заменяет его значениями нового поста
    assert.Equal(t, 1, postInMap.ID)
	assert.Equal(t, prePost1.AuthorID, postInMap.AuthorID)
	assert.Equal(t, prePost1.Description, postInMap.Description)
	assert.Equal(t, prePost1.Price, postInMap.Price)
	assert.Equal(t, prePost1.Price, postInMap.Price)
	assert.Equal(t, prePost1.SafeTransaction, postInMap.SafeTransaction)
	assert.Equal(t, prePost1.Delivery, postInMap.Delivery)
	assert.Equal(t, prePost1.City, postInMap.City)
}