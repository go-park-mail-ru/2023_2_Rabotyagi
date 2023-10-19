package main

import (
	"log"

	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/config"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/server"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/server/delivery/mux"
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

	srv := new(server.Server)
	if err := srv.Run(configServer, mux.NewMux(configServer.AllowOrigin)); err != nil {
		log.Fatalf("Error in server: %s", err.Error())
	}
}
