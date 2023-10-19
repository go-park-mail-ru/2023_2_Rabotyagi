package delivery

import (
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/models"
)

const (
	ResponseSuccessfulAddPost = "Successful add post"

	ErrPostNotExist       = "Post not exists"
	ErrNoSuchCountOfPosts = "not enough posts in storage"
)

type PostResponse struct {
	Status int          `json:"status"`
	Body   *models.Post `json:"body"`
}

func NewPostResponse(status int, body *models.Post) *PostResponse {
	return &PostResponse{
		Status: status,
		Body:   body,
	}
}

type PostsListResponse struct {
	Status int                  `json:"status"`
	Body   []*models.PostInFeed `json:"body"`
}

func NewPostsListResponse(status int, body []*models.PostInFeed) *PostsListResponse {
	return &PostsListResponse{
		Status: status,
		Body:   body,
	}
}
