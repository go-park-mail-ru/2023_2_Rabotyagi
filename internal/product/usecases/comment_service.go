package usecases

import (
	"context"
	"fmt"
	"io"

	productrepo "github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/product/repository"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/models"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/myerrors"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/mylogger"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/utils"
)

var _ ICommentStorage = (*productrepo.ProductStorage)(nil)

type ICommentStorage interface {
	GetCommentList(ctx context.Context, offset uint64, count uint64, recipientID uint64,
		senderID uint64) ([]*models.CommentInFeed, error)
	AddComment(ctx context.Context, preComment *models.PreComment) (uint64, error)
	DeleteComment(ctx context.Context, commentID uint64, senderID uint64) error
	UpdateComment(ctx context.Context, userID uint64, commentID uint64, updateFields map[string]interface{}) error
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

func (c CommentService) GetCommentList(ctx context.Context, offset uint64, count uint64,
	recipientID uint64, senderID uint64,
) ([]*models.CommentInFeed, error) {
	comments, err := c.storage.GetCommentList(ctx, offset, count, recipientID, senderID)
	if err != nil {
		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	for _, product := range comments {
		product.Sanitize()
	}

	return comments, nil
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

func (c CommentService) DeleteComment(ctx context.Context, commentID uint64, senderID uint64) error {
	err := c.storage.DeleteComment(ctx, commentID, senderID)
	if err != nil {
		return fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return nil
}

func (c CommentService) UpdateComment(ctx context.Context, r io.Reader, userID uint64, commentID uint64) error {
	var commChanges *models.CommentChanges

	var err error

	commChanges, err = ValidateCommentChanges(r)
	if err != nil {
		return fmt.Errorf(myerrors.ErrTemplate, err)
	}

	updateFieldsMap := utils.StructToMap(commChanges)

	err = c.storage.UpdateComment(ctx, userID, commentID, updateFieldsMap)
	if err != nil {
		return fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return nil
}
