package delivery

import (
	"context"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/jwt"
	delivery2 "github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/server/delivery"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/server/delivery/statuses"
	userusecases "github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/user/usecases"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/my_logger"
	"io"
	"net/http"
	"time"

	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/models"
	"go.uber.org/zap"
)

const (
	timeTokenLife = 24 * time.Hour
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
	logger  *zap.SugaredLogger
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

// SignUpHandler godoc
//
//	@Summary    signup
//	@Description  signup in app
//
//	@Description Error.status can be:
//	@Description StatusErrBadRequest      = 400
//	@Description  StatusErrInternalServer  = 500
//	@Tags auth
//
//	@Accept      json
//	@Produce    json
//	@Param      preUser  body internal_models.UserWithoutID true  "user data for signup"
//	@Success    200  {object} delivery.ResponseSuccessful
//	@Failure    405  {string} string
//	@Failure    500  {string} string
//	@Failure    222  {object} delivery.ErrorResponse "Error"
//	@Router      /signup [post]
func (u *UserHandler) SignUpHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `Method not allowed`, http.StatusMethodNotAllowed)

		return
	}

	ctx := r.Context()

	user, err := u.service.AddUser(ctx, r.Body)
	if err != nil {
		delivery2.HandleErr(w, u.logger, err)

		return
	}

	expire := time.Now().Add(timeTokenLife)

	jwtStr, err := jwt.GenerateJwtToken(
		&jwt.UserJwtPayload{UserID: user.ID, Email: user.Email, Expire: expire.Unix()}, jwt.GetSecret())
	if err != nil {
		delivery2.SendResponse(w, u.logger,
			delivery2.NewErrResponse(statuses.StatusInternalServer, delivery2.ErrInternalServer))

		return
	}

	cookie := &http.Cookie{ //nolint:exhaustruct
		Name:     delivery2.CookieAuthName,
		Value:    jwtStr,
		Expires:  expire,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
	}

	http.SetCookie(w, cookie)
	delivery2.SendResponse(w, u.logger, delivery2.NewResponseSuccessful(ResponseSuccessfulSignUp))
	u.logger.Infof("in SignUpHandler: added user: %+v", user)
}

// SignInHandler godoc
//
//	@Summary    signin
//	@Description  signin in app
//	@Tags auth
//	@Produce    json
//	@Param      email  query string true  "user email for signin"
//	@Param      password  query string true  "user password for signin"
//	@Success    200  {object} delivery.ResponseSuccessful
//	@Failure    405  {string} string
//	@Failure    500  {string} string
//	@Failure    222  {object} delivery.ErrorResponse "Error"
//	@Router      /signin [get]
func (u *UserHandler) SignInHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `Method not allowed`, http.StatusMethodNotAllowed)

		return
	}

	ctx := r.Context()

	email := r.URL.Query().Get("email")
	password := r.URL.Query().Get("password")

	user, err := u.service.GetUser(ctx, email, password)
	if err != nil {
		delivery2.HandleErr(w, u.logger, err)

		return
	}

	expire := time.Now().Add(timeTokenLife)

	jwtStr, err := jwt.GenerateJwtToken(&jwt.UserJwtPayload{
		UserID: user.ID,
		Email:  user.Email,
		Expire: expire.Unix(),
	},
		jwt.GetSecret(),
	)
	if err != nil {
		delivery2.SendResponse(w, u.logger,
			delivery2.NewErrResponse(statuses.StatusInternalServer, delivery2.ErrInternalServer))

		return
	}

	cookie := &http.Cookie{ //nolint:exhaustruct
		Name:     delivery2.CookieAuthName,
		Value:    jwtStr,
		Expires:  expire,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
	}

	http.SetCookie(w, cookie)
	delivery2.SendResponse(w, u.logger, delivery2.NewResponseSuccessful(ResponseSuccessfulSignIn))
	u.logger.Infof("in SignInHandler: signin user: %+v", user)
}

// LogOutHandler godoc
//
//	@Summary    logout
//	@Description  logout in app
//	@Tags auth
//	@Produce    json
//	@Success    200  {object} delivery.ResponseSuccessful
//	@Failure    405  {string} string
//	@Failure    500  {string} string
//	@Failure    222  {object} delivery.ErrorResponse "Error"
//	@Router      /logout [post]
func (u *UserHandler) LogOutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `Method not allowed`, http.StatusMethodNotAllowed)

		return
	}

	cookie, err := r.Cookie(delivery2.CookieAuthName)
	if err != nil {
		u.logger.Errorln(err)
		delivery2.SendResponse(w, u.logger, delivery2.NewErrResponse(StatusUnauthorized, ErrUnauthorized))

		return
	}

	cookie.Expires = time.Now()

	http.SetCookie(w, cookie)
	delivery2.SendResponse(w, u.logger, delivery2.NewResponseSuccessful(ResponseSuccessfulLogOut))
	u.logger.Infof("in LogOutHandler: logout user with cookie: %+v", cookie)
}
