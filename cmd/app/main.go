package main

import (
	"log"

	handler "github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/transport/handlers"
	rabotyagi "github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/transport/server"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/config"
)

//	@title      YULA project API
//	@version    1.0
//	@description  This is a server of YULA server.
//
// @Schemes http
// @host	84.23.53.28:8080
// @BasePath  /api/v1
func main() {
	configServer := config.New()

	srv := new(rabotyagi.Server)
	if err := srv.Run(configServer, handler.NewMux(configServer.AllowOrigin)); err != nil {
		log.Fatalf("AAAAAAAAAAAAAA %s", err.Error())
	}
}
