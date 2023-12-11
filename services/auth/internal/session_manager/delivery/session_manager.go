package delivery

import (
	"context"
	"fmt"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/auth"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/my_logger"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/myerrors"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/services/auth/internal/session_manager/usecases"
	"github.com/jackc/pgx/v5/pgxpool"
)

var _ IAuthService = (*usecases.AuthService)(nil)

type IAuthService interface {
	AddUser(ctx context.Context, email string, password string) (string, error)
	LoginUser(ctx context.Context, email string, password string) (string, error)
	Delete(ctx context.Context, rawJwt string) (string, error)
	Check(ctx context.Context, rawJwt string) (uint64, error)
}

type SessionManager struct {
	auth.UnimplementedSessionMangerServer

	service IAuthService
	pool    *pgxpool.Pool
	logger  *my_logger.MyLogger
}

func NewSessionManager(pool *pgxpool.Pool, authService IAuthService) (*SessionManager, error) {
	logger, err := my_logger.Get()
	if err != nil {
		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return &SessionManager{ //nolint:exhaustruct
		service: authService,
		pool:    pool,
		logger:  logger,
	}, nil
}

func (s *SessionManager) Check(ctx context.Context, sessionUser *auth.Session) (*auth.UserID, error) {
	if sessionUser == nil {
		return nil, myerrors.NewErrorInternal("sessionUser == nil")
	}

	userID, err := s.service.Check(ctx, sessionUser.GetAccessToken())
	if err != nil {
		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return &auth.UserID{UserId: userID}, nil
}

func (s *SessionManager) Create(ctx context.Context, user *auth.User) (*auth.Session, error) {
	if user == nil {
		return nil, myerrors.NewErrorInternal("user == nil")
	}

	rawJWT, err := s.service.AddUser(ctx, user.GetEmail(), user.GetPassword())
	if err != nil {
		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return &auth.Session{AccessToken: rawJWT}, nil
}

func (s *SessionManager) Login(ctx context.Context, user *auth.User) (*auth.Session, error) {
	if user == nil {
		return nil, myerrors.NewErrorInternal("user == nil")
	}

	rawJWT, err := s.service.LoginUser(ctx, user.GetEmail(), user.GetPassword())
	if err != nil {
		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return &auth.Session{AccessToken: rawJWT}, nil
}

func (s *SessionManager) Delete(ctx context.Context, sessionUser *auth.Session) (*auth.Session, error) {
	if sessionUser == nil {
		return nil, myerrors.NewErrorInternal("sessionUser == nil")
	}

	rawJwt, err := s.service.Delete(ctx, sessionUser.GetAccessToken())
	if err != nil {
		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return &auth.Session{AccessToken: rawJwt}, nil
}
