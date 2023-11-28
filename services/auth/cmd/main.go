package main

import (
	"context"
	"fmt"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/auth"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/config"
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

	lis, err := net.Listen("tcp", config.StandardAddressAuthGrpc)
	if err != nil {
		fmt.Printf("can`t listen port %s", err)

		return
	}

	server := grpc.NewServer()

	baseCtx := context.Background()

	_, err = my_logger.New([]string{configServer.OutputLogPath}, []string{configServer.ErrorOutputLogPath})
	if err != nil {
		fmt.Println(err)

		return
	}

	pool, err := reposhare.NewPgxPool(baseCtx, config.StandardURLDataBase)
	if err != nil {
		fmt.Println(err)

		return
	}

	storage, err := repository.NewAuthStorage(pool)
	if err != nil {
		fmt.Println(err)

		return
	}

	service, err := usecases.NewAuthService(storage)
	if err != nil {
		fmt.Println(err)

		return
	}

	sessionManager, err := delivery.NewSessionManager(pool, service)
	if err != nil {
		fmt.Println(err)

		return
	}

	auth.RegisterSessionMangerServer(server, sessionManager)

	logger.Infof("starting server at :8082")

	chCloseRefreshing := make(chan struct{})

	// don`t want use chCloseRefreshing secret now
	jwt.StartRefreshingSecret(jwt.TimeTokenLife, chCloseRefreshing)

	if err := server.Serve(lis); err != nil {
		fmt.Println(err)

		return
	}
}
