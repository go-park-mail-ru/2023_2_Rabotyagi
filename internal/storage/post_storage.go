package storage

import (
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/errors"
)

var (
	ErrPostNotExist     = errors.NewError("post not exist")
)

type Post struct {
	ID       uint64
	AuthorID     uint64 `json:"author"`
	Title string `json:"title"`
	Description string `json:"description"`
	Price int `json:"price"`
	SafeTransaction bool `json:"safe"`
	Delivery bool `json:"delivery"`
	City string `json:"city"`
}

type PrePost struct {
	AuthorID     uint64 `json:"author"`
	Title string `json:"title"`
	Description string `json:"description"`
	Price int `json:"price"`
	SafeTransaction bool `json:"safe"`
	Delivery bool `json:"delivery"`
	City string `json:"city"`
}

type PostStorage interface {
	GetPost(postID uint64) (*Post, error)
	AddPost(user *PreUser) 
}

type PostStorageMap struct {
	counterPosts uint64
	posts        map[uint64]Post
}

func NewPostStorageSlice() *PostStorageMap {
	return &PostStorageMap{
		counterPosts: 0,
		posts: make(map[uint64]Post),
	}
}

func (a *PostStorageMap) generatePostID() uint64 {
	a.counterPosts++

	return a.counterPosts
}

func (a *PostStorageMap) GetPost(postID uint64) (*Post, error) {
	post, exists := a.posts[postID]
	
	if exists{
		return &post, nil
	}

	return nil, ErrPostNotExist
}

func (a *PostStorageMap) AddPost(post *PrePost) {
	id := a.generatePostID()

	a.posts[id] = Post{ID: id, AuthorID: post.AuthorID, 
		Title: post.Title, Description: post.Description, 
		Price: post.Price, SafeTransaction: post.SafeTransaction, 
		Delivery: post.Delivery, City: post.City,
	}
}
