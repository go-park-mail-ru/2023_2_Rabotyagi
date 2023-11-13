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

### Документация
[Посмотреть здесь](docs/swagger.yaml)

### Сгенерировать swagger документацию

```shell
swag init -ot yaml --parseDependency --parseInternal -g cmd/app/main.go
```

## Локальное поднятие бека, бд, pgadmin вместе
1. Запускаем  все
```shell
docker compose -f deployments/docker-compose.yml up --build -d 
```
2. Далее ждем пока поднимется бек. Команда ниже должна дать вывод как ниже 
```shell
docker compose -f  deployments/docker-compose.yml logs backend
```
Вот такой вывод примерно
```
deployments-backend-1  | {"level":"info","ts":1699520968.4875963,"caller":"server/server.go:55","msg":"Start server:8080"}
```
3. Далее накатываем миграции
```shell
docker exec -it deployments-backend-1 ./migrate -database postgres://postgres:postgres@postgres:5432/youla?sslmode=disable -path db/migrations up
```
4. Далее заполняем бд данными.
```shell
docker exec -it deployments-backend-1 ./fake_db
```
Если произошли какие-то проблемы во время заполнения бд. То откатываем миграции и накатываем еще раз(шаг 3 только в конце up заменяем на down, потом опять вызов с up в конце)

Если все окей, то увидите что-то такое
```
{"level":"info","ts":1699521811.2572942,"caller":"repository/fake_storage.go:305","msg":"end filling favourites\n"}
```
Это все бек + бд + pgadmin запущены
## Запуск локально из терминала / ide

1. Поднимаем бд
```shell
docker compose -f deployments/docker-compose.yml up -d backend
```
2. Прописываем env или соответствующая настройка в ide
```shell
export URL_DATA_BASE=postgres://postgres:postgres@localhost:5432/youla?sslmode=disable
```
3. Запускаем бек
```shell 
go run cmd/app/main.go
```
4. Накатить миграции
```shell
migrate -database postgres://postgres:postgres@localhost:5432/youla?sslmode=disable -path db/migrations up
```
5. Заполнить бд
```shell
sudo go run cmd/fake_db/main.go
```

### Локальная установка тула для миграций
```shell
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

### Тестирование 

```shell
mkdir -p bin && go test -coverprofile=bin/cover.out ./internal/... && go tool cover -html=bin/cover.out -o=bin/cover.html && go tool cover --func bin/cover.out
```
