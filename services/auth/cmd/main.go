package main

import (
	"context"
	"fmt"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/auth"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/my_logger"
	reposhare "github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/repository"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/services/auth/internal/pkg/session_manager/delivery"
	"google.golang.org/grpc"
	"net"
	"os"
)

func main() {
	lis, err := net.Listen("tcp", ":8082")
	if err != nil {
		fmt.Printf("can`t listen port", err)

		return
	}

	server := grpc.NewServer()

	baseCtx := context.Background()

	_, err = my_logger.New([]string{"stdout"}, []string{"stderr"})
	if err != nil {
		fmt.Println(err)

		return
	}

	pool, err := reposhare.NewPgxPool(baseCtx, os.Getenv("URL_DATABASE_AUTH"))
	if err != nil {
		fmt.Println(err)

		return
	}

	sessionManager, err := delivery.NewSessionManager(pool)
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
