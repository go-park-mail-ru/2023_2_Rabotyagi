package delivery

import (
	"log"
	"net/http"
	"time"

	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/jwt"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/server/delivery"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/user/usecases"
)

const (
	timeTokenLife = 24 * time.Hour
)

type UserHandler struct {
	storage    usecases.IUserStorage
	addrOrigin string
	schema     string
}

func NewUserHandler(storage usecases.IUserStorage, addrOrigin string, schema string) *UserHandler {
	return &UserHandler{
		storage:    storage,
		addrOrigin: addrOrigin,
		schema:     schema,
	}
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
	defer r.Body.Close()
	delivery.SetupCORS(w, u.addrOrigin, u.schema)

	if r.Method == http.MethodOptions {
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, `Method not allowed`, http.StatusMethodNotAllowed)

		return
	}

	ctx := r.Context()

	userWithoutID, err := usecases.ValidateUserWithoutID(r.Body)
	if err != nil {
		delivery.HandleErr(w, "in SignUpHandler:", err)

		return
	}

	user, err := u.storage.AddUser(ctx, userWithoutID)
	if err != nil {
		delivery.HandleErr(w, "error in SignUpHandler:", err)

		return
	}

	expire := time.Now().Add(timeTokenLife)

	jwtStr, err := jwt.GenerateJwtToken(
		&jwt.UserJwtPayload{UserID: user.ID, Email: user.Email, Expire: expire.Unix()}, jwt.Secret)
	if err != nil {
		log.Printf("in SignUpHandler: %+v", err)
		delivery.SendErrResponse(w, delivery.NewErrResponse(delivery.StatusErrInternalServer, delivery.ErrInternalServer))

		return
	}

	cookie := &http.Cookie{ //nolint:exhaustruct
		Name:    delivery.CookieAuthName,
		Value:   jwtStr,
		Expires: expire,
		Path:    "/",
	}

	http.SetCookie(w, cookie)
	delivery.SendOkResponse(w, delivery.NewResponse(delivery.StatusResponseSuccessful, ResponseSuccessfulSignUp))
	log.Printf("in SignUpHandler: added user: %+v", user)
}

// SignInHandler godoc
//
//	@Summary    signin
//	@Description  signin in app
//	@Tags auth
//	@Accept      json
//	@Produce    json
//	@Param      preUser  body internal_models.UserWithoutID true  "user data for signin"
//	@Success    200  {object} delivery.Response
//	@Failure    405  {string} string
//	@Failure    500  {string} string
//	@Failure    222  {object} delivery.ErrorResponse "Error"
//	@Router      /signin [post]
func (u *UserHandler) SignInHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	delivery.SetupCORS(w, u.addrOrigin, u.schema)

	if r.Method == http.MethodOptions {
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, `Method not allowed`, http.StatusMethodNotAllowed)

		return
	}

	ctx := r.Context()

	userWithoutID, err := usecases.ValidateUserCredentials(r.Body)
	if err != nil {
		delivery.HandleErr(w, "in SignUpHandler:", err)

		return
	}

	user, err := u.storage.GetUser(ctx, userWithoutID.Email, userWithoutID.Password)
	if err != nil {
		delivery.HandleErr(w, "in SignInHandler:", err)

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
		log.Printf("in SignInHandler: %+v\n", err)
		delivery.SendErrResponse(w, delivery.NewErrResponse(delivery.StatusErrInternalServer, delivery.ErrInternalServer))

		return
	}

	cookie := &http.Cookie{ //nolint:exhaustruct
		Name:    delivery.CookieAuthName,
		Value:   jwtStr,
		Expires: expire,
		Path:    "/",
	}

	http.SetCookie(w, cookie)
	delivery.SendOkResponse(w, delivery.NewResponse(delivery.StatusResponseSuccessful, ResponseSuccessfulSignIn))
	log.Printf("in SignInHandler: signin user: %+v", user)
}

// LogOutHandler godoc
//
//	@Summary    logout
//	@Description  logout in app
//	@Tags auth
//	@Accept      json
//	@Produce    json
//	@Success    200  {object} delivery.Response
//	@Failure    405  {string} string
//	@Failure    500  {string} string
//	@Failure    222  {object} delivery.ErrorResponse "Error"
//	@Router      /logout [post]
func (u *UserHandler) LogOutHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	delivery.SetupCORS(w, u.addrOrigin, u.schema)

	if r.Method == http.MethodOptions {
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, `Method not allowed`, http.StatusMethodNotAllowed)

		return
	}

	cookie, err := r.Cookie(delivery.CookieAuthName)
	if err != nil {
		log.Printf("in LogOutHandler: %+v\n", err)
		delivery.SendErrResponse(w, delivery.NewErrResponse(StatusUnauthorized, ErrUnauthorized))

		return
	}

	cookie.Expires = time.Now()

	http.SetCookie(w, cookie)
	delivery.SendOkResponse(w, delivery.NewResponse(delivery.StatusResponseSuccessful, ResponseSuccessfulLogOut))
	log.Printf("in LogOutHandler: logout user with cookie: %+v", cookie)
}
