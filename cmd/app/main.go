package main

import (
	"fmt"

	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/config"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/server"
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
	if err := srv.Run(configServer); err != nil {
		fmt.Printf("Error in server: %s", err.Error())
	}
}
