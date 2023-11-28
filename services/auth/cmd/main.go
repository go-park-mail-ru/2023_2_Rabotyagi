package main

import (
	"context"
	"fmt"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/auth"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/my_logger"
	reposhare "github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/repository"
	authconfig "github.com/go-park-mail-ru/2023_2_Rabotyagi/services/auth/internal/config"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/services/auth/internal/jwt"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/services/auth/internal/session_manager/delivery"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/services/auth/internal/session_manager/repository"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/services/auth/internal/session_manager/usecases"
	"google.golang.org/grpc"
	"net"
	"strings"
)

func main() {
	configServer := authconfig.New()

	logger, err := my_logger.New(strings.Split(configServer.OutputLogPath, " "),
		strings.Split(configServer.ErrorOutputLogPath, " "))
	if err != nil {
		fmt.Printf("can`t create logger %s", err)

		return
	}

	lis, err := net.Listen("tcp", configServer.AddressAuthServiceGrpc)
	if err != nil {
		logger.Errorf("can`t listen port %s", err)

		return
	}

	server := grpc.NewServer()

	baseCtx := context.Background()

	pool, err := reposhare.NewPgxPool(baseCtx, configServer.URLDataBase)
	if err != nil {
		logger.Errorln(err)

		return
	}

	storage, err := repository.NewAuthStorage(pool)
	if err != nil {
		logger.Errorln(err)

		return
	}

	service, err := usecases.NewAuthService(storage)
	if err != nil {
		logger.Errorln(err)

		return
	}

	sessionManager, err := delivery.NewSessionManager(pool, service)
	if err != nil {
		logger.Errorln(err)

		return
	}

	auth.RegisterSessionMangerServer(server, sessionManager)

	logger.Infof("starting server at: %s", configServer.AddressAuthServiceGrpc)

	chCloseRefreshing := make(chan struct{})

	// don`t want use chCloseRefreshing secret now
	jwt.StartRefreshingSecret(jwt.TimeTokenLife, chCloseRefreshing)

	if err := server.Serve(lis); err != nil {
		logger.Errorln(err)

		return
	}
}
