package delivery

import (
	"context"
	userusecases "github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/user/usecases"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/my_logger"
	"io"
	"net/http"

	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/models"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/server/delivery"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/utils"
)

var _ IUserService = (*userusecases.UserService)(nil)

type IUserService interface {
	AddUser(ctx context.Context, r io.Reader) (*models.User, error)
	GetUser(ctx context.Context, email string, password string) (*models.UserWithoutPassword, error)
	GetUserWithoutPasswordByID(ctx context.Context, userID uint64) (*models.UserWithoutPassword, error)
	UpdateUser(ctx context.Context, r io.Reader, isPartialUpdate bool, userID uint64) (*models.UserWithoutPassword, error)
}

type UserHandler struct {
	service IUserService
	logger  *my_logger.MyLogger
}

func NewUserHandler(userService IUserService) (*UserHandler, error) {
	logger, err := my_logger.Get()
	if err != nil {
		return nil, err
	}

	return &UserHandler{
		service: userService,
		logger:  logger,
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
//	@Failure    222  {object} delivery.ErrorResponse "Error"
//	@Router      /profile/get [get]
func (u *UserHandler) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `Method not allowed`, http.StatusMethodNotAllowed)

		return
	}

	ctx := utils.AddRequestIDToCtx(r.Context())

	userID, err := utils.ParseUint64FromRequest(r, "id")
	if err != nil {
		delivery.HandleErr(w, u.logger, err)

		return
	}

	user, err := u.service.GetUserWithoutPasswordByID(ctx, userID)
	if err != nil {
		delivery.HandleErr(w, u.logger, err)

		return
	}

	delivery.SendResponse(w, u.logger, NewProfileResponse(user))
	u.logger.LogReqID(ctx).Infof("in GetUserHandler: get product: %+v", user)
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
//	@Failure    222  {object} delivery.ErrorResponse "Error". Внутри body статус может быть badContent(4400), badFormat(4000)
//	@Router      /profile/update [patch]
//	@Router      /profile/update [put]
func (u *UserHandler) PartiallyUpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch && r.Method != http.MethodPut {
		http.Error(w, `Method not allowed`, http.StatusMethodNotAllowed)

		return
	}

	ctx := r.Context()

	var err error

	userID, err := delivery.GetUserIDFromCookie(r)
	if err != nil {
		delivery.HandleErr(w, u.logger, err)

		return
	}

	var updatedUser *models.UserWithoutPassword
	if r.Method == http.MethodPatch {
		updatedUser, err = u.service.UpdateUser(ctx, r.Body, true, userID)
	} else {
		updatedUser, err = u.service.UpdateUser(ctx, r.Body, false, userID)
	}

	if err != nil {
		u.logger.Errorf("in PartiallyUpdateUserHandler: %+v\n", err)
		delivery.HandleErr(w, u.logger, err)

		return
	}

	delivery.SendResponse(w, u.logger, NewProfileResponse(updatedUser))
	u.logger.Infof("Successfully updated: %+v", userID)
}
