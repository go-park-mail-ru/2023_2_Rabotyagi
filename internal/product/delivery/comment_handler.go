package delivery

import (
	"context"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/server/delivery"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/responses"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/utils"
	"io"
	"net/http"

	productusecases "github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/product/usecases"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/models"
)

var _ ICommentService = (*productusecases.CommentService)(nil)

type ICommentService interface {
	GetCommentList(ctx context.Context, lastCommentID uint64, count uint64, userID uint64) ([]*models.Comment, error)
	AddComment(ctx context.Context, r io.Reader, userID uint64) (commentID uint64, err error)
	DeleteComment(ctx context.Context, commentID uint64, senderID uint64) error
}

// AddCommentHandler godoc
//
//	@Summary    add comment
//	@Description  add comment by data
//	@Tags comment
//
//	@Accept      json
//	@Produce    json
//	@Param      comment  body models.PreComment true  "comment data for adding"
//	@Success    200  {object} responses.ResponseID
//	@Failure    405  {string} string
//	@Failure    500  {string} string
//	@Failure    222  {object} responses.ErrorResponse "Error". Это Http ответ 200, внутри body статус может быть badContent(4400), badFormat(4000)//nolint:lll
//	@Router      /comment/add [post]
func (p *ProductHandler) AddCommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `Method not allowed`, http.StatusMethodNotAllowed)

		return
	}

	ctx := r.Context()
	logger := p.logger.LogReqID(ctx)

	userID, err := delivery.GetUserID(ctx, r, p.sessionManagerClient)
	if err != nil {
		responses.HandleErr(w, r, logger, err)

		return
	}

	commentID, err := p.service.AddComment(ctx, r.Body, userID)
	if err != nil {
		responses.HandleErr(w, r, logger, err)

		return
	}

	responses.SendResponse(w, logger, responses.NewResponseIDRedirect(commentID))
	logger.Infof("in AddCommentHandler: added comment id= %+v", commentID)
}

// DeleteCommentHandler godoc
//
//	@Summary     delete comment
//	@Description  delete comment for sender using user id from cookies\jwt.
//	@Description  This totally removed comment. Recovery will be impossible
//	@Tags comment
//	@Accept      json
//	@Produce    json
//	@Param      comment_id  query uint64 true  "comment id"
//	@Success    200  {object} responses.ResponseSuccessful
//	@Failure    405  {string} string
//	@Failure    500  {string} string
//	@Router      /comment/delete [delete]
//
// @Failure    222  {object} responses.ErrorResponse "Error". Это Http ответ 200, внутри body статус может быть badContent(4400)//nolint:lll
func (p *ProductHandler) DeleteCommentHandler(w http.ResponseWriter, r *http.Request) { //nolint:dupl
	if r.Method != http.MethodDelete {
		http.Error(w, `Method not allowed`, http.StatusMethodNotAllowed)

		return
	}

	ctx := r.Context()
	logger := p.logger.LogReqID(ctx)

	userID, err := delivery.GetUserID(ctx, r, p.sessionManagerClient)
	if err != nil {
		responses.HandleErr(w, r, logger, err)

		return
	}

	commentID, err := utils.ParseUint64FromRequest(r, "comment_id")
	if err != nil {
		responses.HandleErr(w, r, logger, err)

		return
	}

	err = p.service.DeleteComment(ctx, commentID, userID)
	if err != nil {
		responses.HandleErr(w, r, logger, err)

		return
	}

	responses.SendResponse(w, logger,
		responses.NewResponseSuccessful(ResponseSuccessfulDeleteComment))
	logger.Infof("in DeleteCommentHandler: delete comment id=%d", commentID)
}
