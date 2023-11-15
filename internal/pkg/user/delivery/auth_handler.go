package delivery

import (
	"net/http"
	"time"

	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/jwt"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/my_logger"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/server/delivery"
	userusecases "github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/user/usecases"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/utils"

	"go.uber.org/zap"
)

const (
	timeTokenLife = 24 * time.Hour
)

type UserHandler struct {
	storage userusecases.IUserStorage
	logger  *zap.SugaredLogger
}

func NewUserHandler(storage userusecases.IUserStorage) (*UserHandler, error) {
	logger, err := my_logger.Get()
	if err != nil {
		return nil, err
	}

	return &UserHandler{
		storage: storage,
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
//	@Success    200  {object} delivery.Response
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

	userWithoutID, err := userusecases.ValidateUserWithoutID(r.Body)
	if err != nil {
		delivery.HandleErr(w, u.logger, err)

		return
	}

	userWithoutID.Password, err = utils.HashPass(userWithoutID.Password)
	if err != nil {
		delivery.SendErrResponse(w, u.logger,
			delivery.NewErrResponse(delivery.StatusErrInternalServer, delivery.ErrInternalServer))

		return
	}

	user, err := u.storage.AddUser(ctx, userWithoutID)
	if err != nil {
		delivery.HandleErr(w, u.logger, err)

		return
	}

	expire := time.Now().Add(timeTokenLife)

	jwtStr, err := jwt.GenerateJwtToken(
		&jwt.UserJwtPayload{UserID: user.ID, Email: user.Email, Expire: expire.Unix()}, jwt.Secret)
	if err != nil {
		delivery.SendErrResponse(w, u.logger,
			delivery.NewErrResponse(delivery.StatusErrInternalServer, delivery.ErrInternalServer))

		return
	}

	cookie := &http.Cookie{ //nolint:exhaustruct
		Name:    delivery.CookieAuthName,
		Value:   jwtStr,
		Expires: expire,
		Path:    "/",
	}

	http.SetCookie(w, cookie)
	delivery.SendOkResponse(w, u.logger, delivery.NewResponse(delivery.StatusResponseSuccessful, ResponseSuccessfulSignUp))
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
//	@Success    200  {object} delivery.Response
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

	userWithoutID, err := userusecases.ValidateUserCredentials(
		r.URL.Query().Get("email"), r.URL.Query().Get("password"))
	if err != nil {
		u.logger.Errorf("in SignInHandler: %+v\n", err)
		delivery.HandleErr(w, u.logger, err)

		return
	}

	user, err := u.storage.GetUser(ctx, userWithoutID.Email, userWithoutID.Password)
	if err != nil {
		delivery.HandleErr(w, u.logger, err)

		return
	}

	expire := time.Now().Add(timeTokenLife)

	jwtStr, err := jwt.GenerateJwtToken(&jwt.UserJwtPayload{
		UserID: user.ID,
		Email:  user.Email,
		Expire: expire.Unix(),
	},
		jwt.Secret,
	)
	if err != nil {
		delivery.SendErrResponse(w, u.logger,
			delivery.NewErrResponse(delivery.StatusErrInternalServer, delivery.ErrInternalServer))

		return
	}

	cookie := &http.Cookie{ //nolint:exhaustruct
		Name:    delivery.CookieAuthName,
		Value:   jwtStr,
		Expires: expire,
		Path:    "/",
	}

	http.SetCookie(w, cookie)
	delivery.SendOkResponse(w, u.logger, delivery.NewResponse(delivery.StatusResponseSuccessful, ResponseSuccessfulSignIn))
	u.logger.Infof("in SignInHandler: signin user: %+v", user)
}

// LogOutHandler godoc
//
//	@Summary    logout
//	@Description  logout in app
//	@Tags auth
//	@Produce    json
//	@Success    200  {object} delivery.Response
//	@Failure    405  {string} string
//	@Failure    500  {string} string
//	@Failure    222  {object} delivery.ErrorResponse "Error"
//	@Router      /logout [post]
func (u *UserHandler) LogOutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `Method not allowed`, http.StatusMethodNotAllowed)

		return
	}

	cookie, err := r.Cookie(delivery.CookieAuthName)
	if err != nil {
		u.logger.Errorln(err)
		delivery.SendErrResponse(w, u.logger, delivery.NewErrResponse(StatusUnauthorized, ErrUnauthorized))

		return
	}

	cookie.Expires = time.Now()

	http.SetCookie(w, cookie)
	delivery.SendOkResponse(w, u.logger, delivery.NewResponse(delivery.StatusResponseSuccessful, ResponseSuccessfulLogOut))
	u.logger.Infof("in LogOutHandler: logout user with cookie: %+v", cookie)
}
