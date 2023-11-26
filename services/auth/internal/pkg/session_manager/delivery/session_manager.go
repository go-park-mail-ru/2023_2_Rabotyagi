package delivery

import (
	"context"
	"encoding/hex"
	"fmt"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/server/repository"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/utils"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/auth"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/my_logger"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/myerrors"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/services/auth/internal/pkg/jwt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

type SessionManager struct {
	auth.UnimplementedSessionMangerServer

	pool   *pgxpool.Pool
	logger *zap.SugaredLogger
}

func NewSessionManager(pool *pgxpool.Pool) (*SessionManager, error) {
	logger, err := my_logger.Get()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "logger failed get")
	}

	return &SessionManager{
		pool:   pool,
		logger: logger,
	}, nil
}

func (s *SessionManager) Check(ctx context.Context, sessionUser *auth.Session) (*auth.SessionStatus, error) {
	if sessionUser == nil {
		return nil, status.Errorf(codes.InvalidArgument, "sessionUser == nil")
	}

	_, err := jwt.NewUserJwtPayload(sessionUser.AccessToken, jwt.GetSecret())
	if err != nil {
		return &auth.SessionStatus{Correct: false}, nil
	}

	return &auth.SessionStatus{Correct: true}, nil
}

func (s *SessionManager) Create(ctx context.Context, user *auth.User) (*auth.Session, error) {
	if user == nil {
		return nil, status.Errorf(codes.InvalidArgument, "user == nil")
	}

	jwtPayload := jwt.UserJwtPayload{}

	err := pgx.BeginFunc(ctx, s.pool, func(tx pgx.Tx) error {
		SQLSelectUser := `SELECT email FROM public."user" WHERE email=$1`
		rows := tx.QueryRow(ctx, SQLSelectUser, user.Email)

		if err := rows.Scan(&jwtPayload.Email); err == nil {
			s.logger.Errorf("пользователь с email %s already exist", user.Email)

			return status.Errorf(codes.AlreadyExists, "user with same email already exist")
		}

		SQLInsertUser := `INSERT INTO public."user"(email, password) VALUES ($1, $2)`
		commandTag, err := tx.Exec(ctx, SQLInsertUser, user.Email, user.Password)
		if err != nil {
			s.logger.Errorln(err)

			return status.Errorf(codes.Internal, "can`t insert user")
		}

		if commandTag.RowsAffected() == 0 {
			s.logger.Errorf("Не получилось вставить юзера rowsAffected=0 %+v", user)

			return status.Errorf(codes.Internal, "can`t insert user")
		}

		userID, err := repository.GetLastValSeq(ctx, tx, s.logger, pgx.Identifier{"public", "user"})
		if err != nil {
			s.logger.Errorf("Не получилось взять userID %+v", user)

			return status.Errorf(codes.Internal, "can`t insert user")
		}

		jwtPayload.UserID = userID
		jwtPayload.Email = user.Email
		jwtPayload.Expire = time.Now().Add(jwt.TimeTokenLife).Unix()

		return nil
	})

	if err != nil {
		return nil, err
	}

	rawJwt, err := jwt.GenerateJwtToken(&jwtPayload, jwt.GetSecret())
	if err != nil {
		s.logger.Errorln(err)

		return nil, status.Errorf(codes.Internal, "can`t insert user")
	}

	return &auth.Session{AccessToken: rawJwt}, nil
}

func (s *SessionManager) Login(ctx context.Context, user *auth.User) (*auth.Session, error) {
	if user == nil {
		return nil, status.Errorf(codes.InvalidArgument, "user == nil")
	}

	jwtPayload := jwt.UserJwtPayload{}

	err := pgx.BeginFunc(ctx, s.pool, func(tx pgx.Tx) error {
		var userID uint64

		SQLSelectUser := `SELECT id, email, password FROM public."user" WHERE email=$1`
		rows := tx.QueryRow(ctx, SQLSelectUser, userID, user.Email, user.Password)

		if err := rows.Scan(&jwtPayload.Email); err != nil {
			s.logger.Errorf("пользователь с email %s not exist", user.Email)

			return status.Errorf(codes.NotFound, "user with same email already exist")
		}

		hashPass, err := hex.DecodeString(user.Password)
		if err != nil {
			return fmt.Errorf(myerrors.ErrTemplate, err)
		}

		if !utils.ComparePassAndHash(hashPass, user.Password) {
			return fmt.Errorf(myerrors.ErrTemplate, err)
		}

		jwtPayload.UserID = userID
		jwtPayload.Email = user.Email
		jwtPayload.Expire = time.Now().Add(jwt.TimeTokenLife).Unix()

		return nil
	})

	if err != nil {
		return nil, err
	}

	rawJwt, err := jwt.GenerateJwtToken(&jwtPayload, jwt.GetSecret())
	if err != nil {
		s.logger.Errorln(err)

		return nil, status.Errorf(codes.Internal, "can`t get user")
	}

	return &auth.Session{AccessToken: rawJwt}, nil
}

func (s *SessionManager) Delete(ctx context.Context, sessionUser *auth.Session) (*auth.Session, error) {
	if sessionUser == nil {
		return nil, status.Errorf(codes.InvalidArgument, "sessionUser == nil")
	}

	jwtPayload, err := jwt.NewUserJwtPayload(sessionUser.AccessToken, jwt.GetSecret())
	if err != nil {
		s.logger.Errorln(err)

		return nil, status.Errorf(codes.InvalidArgument, "get incorrect access_token in sessionUser")
	}

	jwtPayload.Expire = time.Now().Unix()

	rawJwt, err := jwt.GenerateJwtToken(jwtPayload, jwt.GetSecret())
	if err != nil {
		s.logger.Errorln(err)

		return nil, status.Errorf(codes.Internal, "can`t delete user")
	}

	return &auth.Session{AccessToken: rawJwt}, nil
}
