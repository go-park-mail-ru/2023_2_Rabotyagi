package repository_test

//
//import (
//	"testing"
//
//	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/models"
//	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/product/repository"
//	"github.com/stretchr/testify/assert"
//)
//
//func TestPostStorageMapGet(t *testing.T) {
//	storageMap := repository.NewPostStorageMap()
//	prePost := &models.PreProduct{
//		SalerID:        1,
//		Title:           "Test Product",
//		Description:     "This is a test product",
//		Price:           100,
//		SafeDeal: true,
//		Delivery:        true,
//		City:            "Moscow",
//	}
//
//	storageMap.AddPost(prePost)
//
//	// Проверяем, что пост был добавлен успешно
//	expectedID := 1
//	assert.Equal(t, expectedID, storageMap.GetPostsCount())
//
//	postInMap, _ := storageMap.GetPost(uint64(1))
//
//	assert.Equal(t, uint64(1), postInMap.ID)
//	assert.Equal(t, prePost.SalerID, postInMap.SalerID)
//	assert.Equal(t, prePost.Description, postInMap.Description)
//	assert.Equal(t, prePost.Price, postInMap.Price)
//	assert.Equal(t, prePost.Price, postInMap.Price)
//	assert.Equal(t, prePost.SafeDeal, postInMap.SafeDeal)
//	assert.Equal(t, prePost.Delivery, postInMap.Delivery)
//	assert.Equal(t, prePost.City, postInMap.City)
//}
//
//func TestPostStorageMapGetError(t *testing.T) {
//	storageMap := repository.NewPostStorageMap()
//	prePost := &models.PreProduct{
//		SalerID:        1,
//		Title:           "Test Product",
//		Description:     "This is a test product",
//		Price:           100,
//		SafeDeal: true,
//		Delivery:        true,
//		City:            "Moscow",
//	}
//
//	storageMap.AddPost(prePost)
//
//	_, err := storageMap.GetPost(2)
//
//	assert.NotNil(t, err)
//}
//
//func TestPostStorageMapGetList(t *testing.T) {
//	storageMap := repository.NewPostStorageMap()
//	prePost1 := &models.PreProduct{
//		SalerID:        1,
//		Title:           "Test Product",
//		Description:     "This is a test product",
//		Price:           100,
//		SafeDeal: true,
//		Delivery:        true,
//		City:            "Moscow",
//	}
//
//	storageMap.AddPost(prePost1)
//
//	prePost2 := &models.PreProduct{
//		SalerID:        2,
//		Title:           "Test Post2",
//		Description:     "This is a test post2",
//		Price:           10000,
//		SafeDeal: false,
//		Delivery:        false,
//		City:            "City",
//	}
//
//	storageMap.AddPost(prePost2)
//
//	postInMap1, _ := storageMap.GetPost(uint64(1))
//
//	assert.Equal(t, uint64(1), postInMap1.ID)
//	assert.Equal(t, prePost1.SalerID, postInMap1.SalerID)
//	assert.Equal(t, prePost1.Description, postInMap1.Description)
//	assert.Equal(t, prePost1.Price, postInMap1.Price)
//	assert.Equal(t, prePost1.Price, postInMap1.Price)
//	assert.Equal(t, prePost1.SafeDeal, postInMap1.SafeDeal)
//	assert.Equal(t, prePost1.Delivery, postInMap1.Delivery)
//	assert.Equal(t, prePost1.City, postInMap1.City)
//
//	postInMap2, _ := storageMap.GetPost(uint64(2))
//
//	assert.Equal(t, uint64(2), postInMap2.ID)
//	assert.Equal(t, prePost2.SalerID, postInMap2.SalerID)
//	assert.Equal(t, prePost2.Description, postInMap2.Description)
//	assert.Equal(t, prePost2.Price, postInMap2.Price)
//	assert.Equal(t, prePost2.Price, postInMap2.Price)
//	assert.Equal(t, prePost2.SafeDeal, postInMap2.SafeDeal)
//	assert.Equal(t, prePost2.Delivery, postInMap2.Delivery)
//	assert.Equal(t, prePost2.City, postInMap2.City)
//}
//
//func TestPostStorageMapGetListError(t *testing.T) {
//	storageMap := repository.NewPostStorageMap()
//	prePost := &models.PreProduct{
//		SalerID:        1,
//		Title:           "Test Product",
//		Description:     "This is a test product",
//		Price:           100,
//		SafeDeal: true,
//		Delivery:        true,
//		City:            "Moscow",
//	}
//
//	storageMap.AddPost(prePost)
//
//	_, err := storageMap.GetNPosts(10)
//
//	assert.NotNil(t, err)
//}
