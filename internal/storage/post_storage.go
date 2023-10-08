package storage

import (
	"fmt"
	"sync"

	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/errors"
)

var (
	ErrPostNotExist       = errors.NewError("post not exist")
	ErrNoSuchCountOfPosts = errors.NewError("n > posts count")
)

type Image struct {
	Url string
	Alt string
}

type Post struct {
	ID              uint64 `json:"id"`
	AuthorID        uint64 `json:"author"`
	Title           string `json:"title"`
	Image           Image  `jsom:"image"`
	Description     string `json:"description"`
	Price           int    `json:"price"`
	SafeTransaction bool   `json:"safe"`
	Delivery        bool   `json:"delivery"`
	City            string `json:"city"`
}

type PrePost struct {
	AuthorID        uint64 `json:"author"`
	Title           string `json:"title"`
	Image           Image  `jsom:"image"`
	Description     string `json:"description"`
	Price           int    `json:"price"`
	SafeTransaction bool   `json:"safe"`
	Delivery        bool   `json:"delivery"`
	City            string `json:"city"`
}

type PostInFeed struct {
	ID              uint64 `json:"id"`
	Title           string `json:"title"`
	Image           Image  `json:"image"`
	Price           int    `json:"price"`
	SafeTransaction bool   `json:"safe"`
	Delivery        bool   `json:"delivery"`
	City            string `json:"city"`
}

type PostStorage interface {
	GetPost(postID uint64) (*Post, error)
	GetNPosts() []*Post
	AddPost(user *PreUser)
}

type PostStorageMap struct {
	counterPosts uint64
	posts        map[uint64]*Post
	mu           sync.RWMutex
}

func GeneratePosts(postStorageMap *PostStorageMap) *PostStorageMap {
	for i := 1; i <= 40; i++ {
		postID := postStorageMap.generatePostID()
		postStorageMap.posts[postID] = &Post{
			ID:       postID,
			AuthorID: 1,
			Title:    fmt.Sprintf("post %d", postID),
			Image: Image{
				fmt.Sprintf("img_url%d", postID),
				fmt.Sprintf("img_alt%d", postID),
			},
			Description:     fmt.Sprintf("description of post %d", postID),
			Price:           int(100 * postID),
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
		posts:        make(map[uint64]*Post),
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

func (a *PostStorageMap) GetPost(postID uint64) (*Post, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	post, exists := a.posts[postID]

	if exists {
		return post, nil
	}

	return nil, ErrPostNotExist
}

func (a *PostStorageMap) AddPost(post *PrePost) {
	a.mu.Lock()
	defer a.mu.Unlock()

	id := a.generatePostID()

	a.posts[id] = &Post{
		ID:       id,
		AuthorID: post.AuthorID,
		Title:    post.Title,
		Image: Image{
			post.Image.Url,
			post.Image.Alt,
		},
		Description:     post.Description,
		Price:           post.Price,
		SafeTransaction: post.SafeTransaction,
		Delivery:        post.Delivery,
		City:            post.City,
	}
}

func (a *PostStorageMap) GetNPosts(n int) ([]*PostInFeed, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	if n > int(a.counterPosts) {
		return nil, ErrNoSuchCountOfPosts
	}

	postsInFeedSlice := make([]*PostInFeed, 0, n)

	for _, post := range a.posts {
		n--

		postsInFeedSlice = append(postsInFeedSlice, &PostInFeed{
			ID:    post.ID,
			Title: post.Title,
			Image: Image{
				post.Image.Url,
				post.Image.Alt,
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
