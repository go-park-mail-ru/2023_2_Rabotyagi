# 2023_2_Rabotyagi
Backend репозиторий команды Работяги

### Наши контакты:

Владислав Ильинский: https://github.com/Vilinvil и тг https://t.me/Vilin0

Никита Демирев: 'https://github.com/NickDemiman' и тг https://t.me/NikDemiman

Алексей Красноперов: 'https://github.com/SanExpett' и тг https://t.me/SanExpet

Таня Емельянова: 'https://github.com/TanyaEmka' и тг https://t.me/jupi_abri

### Репа фронт
https://github.com/frontend-park-mail-ru/2023_2_Rabotyagi/tree/minimal-front

### Фигма
https://www.figma.com/file/YLSZ9uY9gVn6bMDJchrEzD?node-id=23:2127&mode=design#567544444

### Приложение
http://84.23.53.28/

### Запуск локально

`go run cmd/app/main.go`

### Тестирование 

`mkdir -p bin && go test -coverprofile=bin/cover.out ./internal/... && go tool cover -html=bin/cover.out -o=bin/cover.html && go tool cover --func bin/cover.out`

## Документация
 Ссылка https://app.swaggerhub.com/apis/IVN15072002/yula-project_api/1.0
 Также посмотреть информацию по ручками api можно в docs/swagger.yaml 

### Сгенерировать swagger документацию

`swag init -g cmd/app/main.go`


## Docker image build

### Local

Из корня проекта прописываем
```shell
docker build -t rabotyagi/backend .
```

Далее, чтобы убедиться что image забилдился, прописываем:
```shell
docker images
```

Должны увидеть следующее:
```shell
REPOSITORY          TAG       IMAGE ID       CREATED          SIZE
rabotyagi/backend   latest    25dbaeeef1af   50 seconds ago   307MB
```

### Запуск контейнера 

`docker run -p 8080:8080 rabotyagi/backend`


### Локальная установка тула для миграций
```shell
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

### Пример команды, чтобы накатить миграцию
```shell
 migrate -database postgres://postgres:password@localhost:5432/youla?sslmode=disable -path db/migrations up
```
### Пример команды, чтобы отменить миграцию
```shell
 migrate -database postgres://postgres:password@localhost:5432/youla?sslmode=disable -path db/migrations up
```
