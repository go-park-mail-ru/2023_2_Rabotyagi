package main

import (
	"log"

	handler "github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/transport/handlers"
	rabotyagi "github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/transport/server"
)

//  @title      YULA project API
//  @version    1.0
//  @description  This is a server of YULA server.

// @host    127.0.0.1:8080
// @BasePath  /api/v1
func main() {
	handler := new(handler.Handler)

	srv := new(rabotyagi.Server)
	if err := srv.Run("8080", handler.InitRoutes()); err != nil {
		log.Fatalf("AAAAAAAAAAAAAA %s", err.Error())
	}
}
