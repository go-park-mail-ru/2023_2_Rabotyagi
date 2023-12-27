package main

import (
	"fmt"
	"strings"

	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/mylogger"
	authconfig "github.com/go-park-mail-ru/2023_2_Rabotyagi/services/auth/internal/config"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/services/auth/internal/server"
)

func main() {
	configLogger, err := authconfig.NewConfigLogger()
	if err != nil {
		fmt.Printf("ошибка создания конфига логгера %+v", err)

		return
	}

	logger, err := mylogger.New(strings.Split(configLogger.OutputLogPath, " "),
		strings.Split(configLogger.ErrorOutputLogPath, " "))
	if err != nil {
		fmt.Printf("ошибка создания логгера %+v", err)

		return
	}

	configServer, err := authconfig.New()
	if err != nil {
		fmt.Printf("ошибка создания конфига %+v", err)

		return
	}

	chErrHTTP := make(chan error)
	go func() {
		err := <-chErrHTTP
		if err != nil {
			fmt.Printf("Error in http server: %s", err)

			return
		}
	}()

	srv := &server.Server{}
	if err := srv.RunFull(configServer, logger, chErrHTTP); err != nil {
		fmt.Printf("Error in grpc server: %s", err.Error())
	}
}
