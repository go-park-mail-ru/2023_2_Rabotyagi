package delivery

import (
	"context"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/auth"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/my_logger"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/services/auth/internal/pkg/session_manager/usecases"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ IAuthService = (*usecases.AuthService)(nil)

type IAuthService interface {
	AddUser(ctx context.Context, email string, password string) (string, error)
	GetUserRawJWT(ctx context.Context, email string, password string) (string, error)
	Delete(ctx context.Context, rawJwt string) (string, error)
	Check(ctx context.Context, rawJwt string) bool
}

type SessionManager struct {
	auth.UnimplementedSessionMangerServer

	service IAuthService
	pool    *pgxpool.Pool
	logger  *zap.SugaredLogger
}

func NewSessionManager(pool *pgxpool.Pool, authService IAuthService) (*SessionManager, error) {
	logger, err := my_logger.Get()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "logger failed get")
	}

	return &SessionManager{
		service: authService,
		pool:    pool,
		logger:  logger,
	}, nil
}

func (s *SessionManager) Check(ctx context.Context, sessionUser *auth.Session) (*auth.SessionStatus, error) {
	if sessionUser == nil {
		return nil, status.Errorf(codes.InvalidArgument, "sessionUser == nil")
	}

	correct := s.service.Check(ctx, sessionUser.AccessToken)

	return &auth.SessionStatus{Correct: correct}, nil
}

func (s *SessionManager) Create(ctx context.Context, user *auth.User) (*auth.Session, error) {
	if user == nil {
		return nil, status.Errorf(codes.InvalidArgument, "user == nil")
	}

	rawJWT, err := s.service.AddUser(ctx, user.Email, user.Password)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "in Create")
	}

	return &auth.Session{AccessToken: rawJWT}, nil
}

func (s *SessionManager) Login(ctx context.Context, user *auth.User) (*auth.Session, error) {
	if user == nil {
		return nil, status.Errorf(codes.InvalidArgument, "user == nil")
	}

	rawJWT, err := s.service.GetUserRawJWT(ctx, user.Email, user.Password)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "in Login")
	}

	return &auth.Session{AccessToken: rawJWT}, nil
}

func (s *SessionManager) Delete(ctx context.Context, sessionUser *auth.Session) (*auth.Session, error) {
	if sessionUser == nil {
		return nil, status.Errorf(codes.InvalidArgument, "sessionUser == nil")
	}

	rawJwt, err := s.service.Delete(ctx, sessionUser.AccessToken)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "in Delete")
	}

	return &auth.Session{AccessToken: rawJwt}, nil
}
