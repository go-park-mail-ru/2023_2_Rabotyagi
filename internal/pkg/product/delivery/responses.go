package delivery

import (
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/models"
)

const (
	ResponseSuccessfulAddPost = "Successful add product"

	ErrPostNotExist       = "Product not exists"
	ErrNoSuchCountOfPosts = "not enough posts in storage"
)

type PostResponse struct {
	Status int             `json:"status"`
	Body   *models.Product `json:"body"`
}

func NewPostResponse(status int, body *models.Product) *PostResponse {
	return &PostResponse{
		Status: status,
		Body:   body,
	}
}

type PostsListResponse struct {
	Status int                     `json:"status"`
	Body   []*models.ProductInFeed `json:"body"`
}

func NewPostsListResponse(status int, body []*models.ProductInFeed) *PostsListResponse {
	return &PostsListResponse{
		Status: status,
		Body:   body,
	}
}
