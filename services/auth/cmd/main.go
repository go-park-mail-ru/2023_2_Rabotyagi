package main

import (
	"fmt"

	authconfig "github.com/go-park-mail-ru/2023_2_Rabotyagi/services/auth/internal/config"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/services/auth/internal/server"
)

func main() {
	configServer := authconfig.New()

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
