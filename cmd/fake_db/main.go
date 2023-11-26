package main

import (
	"fmt"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/my_logger"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/config"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/fake_db"
)

func main() {
	const baseCount = 10

	gofakeit.Seed(0)

	configServer := config.New()

	logger, err := my_logger.New([]string{"stdout"}, []string{"stderr"})
	if err != nil {
		fmt.Println(err)

		return
	}
	defer logger.Sync()

	err = fake_db.RunScriptFillDB(configServer.URLDataBase,
		logger, baseCount, configServer.PathToRoot)
	if err != nil {
		logger.Error(err)
	}
}
