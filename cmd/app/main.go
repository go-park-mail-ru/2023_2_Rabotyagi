package main

import (
	"log"

	rabotyagi "github.com/go-park-mail-ru/2023_2_Rabotyagi"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/transport"
)

func main() {
	handler := new(handler.AuthHandler)

	srv := new(rabotyagi.Server)
	if err := srv.Run("8080", handler.InitRoutes()); err != nil {
		log.Fatalf("AAAAAAAAAAAAAA %s", err.Error())
	}
}