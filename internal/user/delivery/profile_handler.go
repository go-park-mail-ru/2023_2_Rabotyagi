package delivery

import (
	"context"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/auth"
	"io"
	"net/http"

	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/server/delivery"
	userusecases "github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/user/usecases"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/models"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/my_logger"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/responses"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/utils"
)

var _ IUserService = (*userusecases.UserService)(nil)

type IUserService interface {
	GetUserWithoutPasswordByID(ctx context.Context, userID uint64) (*models.UserWithoutPassword, error)
	UpdateUser(ctx context.Context, r io.Reader, isPartialUpdate bool, userID uint64) (*models.UserWithoutPassword, error)
}

type ProfileHandler struct {
	sessionManagerClient auth.SessionMangerClient
	service              IUserService
	logger               *my_logger.MyLogger
}

func NewProfileHandler(userService IUserService,
	sessionManagerClient auth.SessionMangerClient,
) (*ProfileHandler, error) {
	logger, err := my_logger.Get()
	if err != nil {
		return nil, err
	}

	return &ProfileHandler{
		service:              userService,
		logger:               logger,
		sessionManagerClient: sessionManagerClient,
	}, nil
}

//	GetUserHandler godoc
//
// @Summary    get profile
// @Description  get profile by id
//
// @Tags profile
//
//	@Produce    json
//	@Param      id  query uint64 true  "user id"
//	@Success    200  {object} ProfileResponse
//	@Failure    405  {string} string
//	@Failure    500  {string} string
//	@Failure    222  {object} responses.ErrorResponse "Error" Внутри body статус может быть badFormat(4000)
//	@Router      /profile/get [get]
func (u *ProfileHandler) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `Method not allowed`, http.StatusMethodNotAllowed)

		return
	}

	ctx := r.Context()
	logger := u.logger.LogReqID(ctx)

	userID, err := utils.ParseUint64FromRequest(r, "id")
	if err != nil {
		responses.HandleErr(w, logger, err)

		return
	}

	user, err := u.service.GetUserWithoutPasswordByID(ctx, userID)
	if err != nil {
		responses.HandleErr(w, logger, err)

		return
	}

	responses.SendResponse(w, logger, NewProfileResponse(user))
	logger.Infof("in GetUserHandler: get product: %+v", user)
}

// PartiallyUpdateUserHandler godoc
//
//	@Summary    update profile
//	@Description  update some fields of profile
//
// @Tags profile
//
//	@Accept      json
//	@Produce    json
//	@Param      user  body models.UserWithoutPassword false  "полностью опционален"
//	@Success    200  {object} ProfileResponse
//	@Failure    405  {string} string
//	@Failure    500  {string} string
//	@Failure    222  {object} responses.ErrorResponse "Error". Внутри body статус может быть badContent(4400), badFormat(4000)
//	@Router      /profile/update [patch]
//	@Router      /profile/update [put]
func (u *ProfileHandler) PartiallyUpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch && r.Method != http.MethodPut {
		http.Error(w, `Method not allowed`, http.StatusMethodNotAllowed)

		return
	}

	ctx := r.Context()
	logger := u.logger.LogReqID(ctx)

	var err error

	userID, err := delivery.GetUserID(ctx, r, u.sessionManagerClient)
	if err != nil {
		responses.HandleErr(w, logger, err)

		return
	}

	var updatedUser *models.UserWithoutPassword
	if r.Method == http.MethodPatch {
		updatedUser, err = u.service.UpdateUser(ctx, r.Body, true, userID)
	} else {
		updatedUser, err = u.service.UpdateUser(ctx, r.Body, false, userID)
	}

	if err != nil {
		logger.Errorf("in PartiallyUpdateUserHandler: %+v\n", err)
		responses.HandleErr(w, logger, err)

		return
	}

	responses.SendResponse(w, logger, NewProfileResponse(updatedUser))
	logger.Infof("Successfully updated: %+v", userID)
}
