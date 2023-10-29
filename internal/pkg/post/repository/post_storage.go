package repository

import (
	"fmt"
	"sync"

	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/models"
	myerrors "github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/errors"
)

var (
	ErrPostNotExist       = myerrors.NewError("post not exist")
	ErrNoSuchCountOfPosts = myerrors.NewError("n > posts count")
)

type PostStorage interface {
	GetPost(postID uint64) (*models.Post, error)
	GetNPosts() []*models.Post
	AddPost(user *models.PreUser)
}

type PostStorageMap struct {
	counterPosts uint64
	posts        map[uint64]*models.Post
	mu           sync.RWMutex
}

func GeneratePosts(postStorageMap *PostStorageMap) *PostStorageMap {
	for i := 1; i <= 40; i++ {
		postID := postStorageMap.generatePostID()
		postStorageMap.posts[postID] = &models.Post{
			ID:       postID,
			AuthorID: 1,
			Title:    fmt.Sprintf("post %d", postID),
			Image: models.Image{
				URL: "http://84.23.53.28:8080/api/v1/img/" +
					"�%7D�̙�%7F�w���f%7C.WebP",
				Alt: "http://84.23.53.28:8080/api/v1/img/" +
					"�%7D�̙�%7F�w���f%7C.WebP",
			},
			Description:     fmt.Sprintf("description of post %d", postID),
			Price:           uint(100 * postID),
			SafeTransaction: true,
			Delivery:        true,
			City:            "Moscow",
		}
	}

	return postStorageMap
}

func NewPostStorageMap() *PostStorageMap {
	return &PostStorageMap{
		counterPosts: 0,
		posts:        make(map[uint64]*models.Post),
		mu:           sync.RWMutex{},
	}
}

func (a *PostStorageMap) GetPostsCount() int {
	a.mu.RLock()
	defer a.mu.RUnlock()

	return len(a.posts)
}

func (a *PostStorageMap) generatePostID() uint64 {
	a.counterPosts++

	return a.counterPosts
}

func (a *PostStorageMap) GetPost(postID uint64) (*models.Post, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	post, exists := a.posts[postID]

	if exists {
		return post, nil
	}

	return nil, ErrPostNotExist
}

func (a *PostStorageMap) AddPost(post *models.PrePost) {
	a.mu.Lock()
	defer a.mu.Unlock()

	id := a.generatePostID()

	a.posts[id] = &models.Post{
		ID:       id,
		AuthorID: post.AuthorID,
		Title:    post.Title,
		Image: models.Image{
			URL: post.Image.URL,
			Alt: post.Image.Alt,
		},
		Description:     post.Description,
		Price:           post.Price,
		SafeTransaction: post.SafeTransaction,
		Delivery:        post.Delivery,
		City:            post.City,
	}
}

func (a *PostStorageMap) GetNPosts(n int) ([]*models.PostInFeed, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	if n > int(a.counterPosts) {
		return nil, ErrNoSuchCountOfPosts
	}

	postsInFeedSlice := make([]*models.PostInFeed, 0, n)

	for _, post := range a.posts {
		n--

		postsInFeedSlice = append(postsInFeedSlice, &models.PostInFeed{
			ID:    post.ID,
			Title: post.Title,
			Image: models.Image{
				URL: post.Image.URL,
				Alt: post.Image.Alt,
			},
			Price:           post.Price,
			SafeTransaction: post.SafeTransaction,
			Delivery:        post.Delivery,
			City:            post.City,
		})

		if n == 0 {
			return postsInFeedSlice, nil
		}
	}

	return postsInFeedSlice, nil
}
