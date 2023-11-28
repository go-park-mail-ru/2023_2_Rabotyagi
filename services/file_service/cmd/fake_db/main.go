package main

import (
	"fmt"
	"os"

	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/my_logger"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/services/file_service/internal/fake_db"

	"github.com/brianvoe/gofakeit/v6"
)

func main() {
	const baseCount = 10

	gofakeit.Seed(0)

	if len(os.Args) != 3 {
		fmt.Println(`command should be format ./fake_db postgres://postgres:postgres@localhost:5432/youla?sslmode=disable .
where first arg is url_db and second path to root where find static/img`)

		return
	}

	urlDB := os.Args[1]
	pathRoot := os.Args[2]

	logger, err := my_logger.New([]string{"stdout"}, []string{"stderr"})
	if err != nil {
		fmt.Println(err)

		return
	}
	defer logger.Sync()

	err = fake_db.RunScriptFillDB(urlDB,
		logger, baseCount, pathRoot)
	if err != nil {
		logger.Error(err)
	}
}
