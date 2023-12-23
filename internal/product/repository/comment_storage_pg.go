package repository

import (
	"context"
	"fmt"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/models"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/myerrors"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/repository"
	"github.com/jackc/pgx/v5"
)

var (
	ErrNoAffectedCommentRows = myerrors.NewErrorBadFormatRequest("Не получилось обновить данные комментария")
)

func (p *ProductStorage) GetCommentList(ctx context.Context, lastCommentID uint64, count uint64, userID uint64) ([]*models.Comment, error) {
	return nil, nil
}

func (p *ProductStorage) insertComment(ctx context.Context, tx pgx.Tx, preComment *models.PreComment) error {
	logger := p.logger.LogReqID(ctx)

	SQLInsertComment := `INSERT INTO public."comment"(recipient_id, sender_id, text, rating) VALUES(
		$1, $2, $3, $4)`

	_, err := tx.Exec(ctx, SQLInsertComment, preComment.RecipientID, preComment.SenderID,
		preComment.Text, preComment.Rating)
	if err != nil {
		logger.Errorln(err)

		return fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return nil
}

func (p *ProductStorage) AddComment(ctx context.Context, preComment *models.PreComment) (uint64, error) {
	logger := p.logger.LogReqID(ctx)

	var commentID uint64

	err := pgx.BeginFunc(ctx, p.pool, func(tx pgx.Tx) error {
		err := p.insertComment(ctx, tx, preComment)
		if err != nil {
			return fmt.Errorf(myerrors.ErrTemplate, err)
		}

		lastCommentID, err := repository.GetLastValSeq(ctx, tx, logger, NameSeqProduct)
		if err != nil {
			return fmt.Errorf(myerrors.ErrTemplate, err)
		}

		commentID = lastCommentID

		return nil
	})
	if err != nil {
		logger.Errorln(err)

		return 0, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return commentID, nil
}

func (p *ProductStorage) deleteComment(ctx context.Context, tx pgx.Tx, commentID uint64, senderID uint64) error {
	logger := p.logger.LogReqID(ctx)

	SQLDeleteComment := `DELETE FROM public."comment" WHERE id=$1 AND sender_id=$2`

	result, err := tx.Exec(ctx, SQLDeleteComment, commentID, senderID)
	if err != nil {
		logger.Errorln(err)

		return fmt.Errorf(myerrors.ErrTemplate, err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf(myerrors.ErrTemplate, ErrNoAffectedCommentRows)
	}

	return nil
}

func (p *ProductStorage) DeleteComment(ctx context.Context, commentID uint64, senderID uint64) error {
	logger := p.logger.LogReqID(ctx)

	err := pgx.BeginFunc(ctx, p.pool, func(tx pgx.Tx) error {
		err := p.deleteComment(ctx, tx, commentID, senderID)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		logger.Errorln(err)

		return fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return nil
}
