package main

import (
	"context"
	"fmt"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/auth"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/config"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/my_logger"
	reposhare "github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/repository"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/services/auth/internal/pkg/session_manager/delivery"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/services/auth/internal/pkg/session_manager/repository"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/services/auth/internal/pkg/session_manager/usecases"
	"google.golang.org/grpc"
	"net"
	"os"
)

var (
	standardOutputLogPath = "stdout"
	standardErrLogPath    = "stderr"
)

func main() {
	lis, err := net.Listen("tcp", ":"+config.StandardAddressAuthGrpc)
	if err != nil {
		fmt.Printf("can`t listen port %s", err)

		return
	}

	server := grpc.NewServer()

	baseCtx := context.Background()

	_, err = my_logger.New([]string{standardOutputLogPath}, []string{standardErrLogPath})
	if err != nil {
		fmt.Println(err)

		return
	}

	pool, err := reposhare.NewPgxPool(baseCtx, os.Getenv("URL_DATABASE_AUTH"))
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

	fmt.Println("starting server at :8082")

	if err := server.Serve(lis); err != nil {
		fmt.Println(err)

		return
	}
}
