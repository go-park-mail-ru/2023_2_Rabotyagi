package main

import (
	"fmt"

	"github.com/go-park-mail-ru/2023_2_Rabotyagi/services/file_service/internal/config"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/services/file_service/internal/server"
)

func main() {
	configServer := config.New()

	chErrHTTP := make(chan error)
	go func() {
		err := <-chErrHTTP
		if err != nil {
			fmt.Printf("Error in http server: %s", err)

			return
		}
	}()

	srv := &server.Server{}
	if err := srv.RunFull(configServer, chErrHTTP); err != nil {
		fmt.Printf("Error in grpc server: %s", err.Error())
	}
}
