package usecases

import (
	"context"
	"fmt"
	productrepo "github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/product/repository"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/models"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/myerrors"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/mylogger"
	"io"
)

var _ ICommentStorage = (*productrepo.ProductStorage)(nil)

type ICommentStorage interface {
	GetCommentList(ctx context.Context, lastCommentID uint64, count uint64, userID uint64) ([]*models.Comment, error)
	AddComment(ctx context.Context, preComment *models.PreComment) (uint64, error)
	DeleteComment(ctx context.Context, commentID uint64, senderID uint64) error
}

type CommentService struct {
	storage ICommentStorage
	logger  *mylogger.MyLogger
}

func NewCommentService(commentStorage ICommentStorage) (*CommentService, error) {
	logger, err := mylogger.Get()
	if err != nil {
		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return &CommentService{storage: commentStorage, logger: logger}, nil
}

func (c CommentService) GetCommentList(ctx context.Context, lastCommentID uint64, count uint64,
	userID uint64) ([]*models.Comment, error) {
	return nil, nil
}

func (c CommentService) AddComment(ctx context.Context, r io.Reader, userID uint64) (uint64, error) {
	preComment, err := ValidatePreComment(r, userID)
	if err != nil {
		return 0, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	commentID, err := c.storage.AddComment(ctx, preComment)
	if err != nil {
		return 0, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return commentID, nil
}

func (c CommentService) DeleteComment(ctx context.Context, commentID uint64, senderID uint64) error { //nolint:revive
	err := c.storage.DeleteComment(ctx, commentID, senderID)
	if err != nil {
		return fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return nil
}
