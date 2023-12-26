# 2023_2_Rabotyagi
Backend репозиторий команды Работяги

### Наши контакты:

Владислав Ильинский: https://github.com/Vilinvil и тг https://t.me/Vilin0

Никита Демирев: 'https://github.com/NickDemiman' и тг https://t.me/NikDemiman

Алексей Красноперов: 'https://github.com/SanExpett' и тг https://t.me/SanExpet

Таня Емельянова: 'https://github.com/TanyaEmka' и тг https://t.me/jupi_abri

### Репа фронт
https://github.com/frontend-park-mail-ru/2023_2_Rabotyagi

### Фигма
https://www.figma.com/file/YLSZ9uY9gVn6bMDJchrEzD?node-id=23:2127&mode=design#567544444

### Приложение
https://goods-galaxy.ru/

### Метрики
https://goods-galaxy.ru/grafana/

Логин: ```guest```
Пароль: ```guest```

### Приложение dev стенд
https://dev.goods-galaxy.ru/

### Документация
[Посмотреть здесь](docs/swagger.yaml)

## Команды для разработки
### Локальное поднятие бека, бд, pgadmin вместе
1. Запускаем  все
```shell
make compose-full-up
```

### Если впервые запускаем бек
1. Далее ждем пока поднимется бек. Команда ниже должна дать вывод как ниже 
```shell
make compose-logs
```
Вот такой вывод примерно
```
deployments-backend-1  | {"level":"info","ts":1699520968.4875963,"caller":"server/server.go:55","msg":"Start server:8080"}
```
2. Далее накатываем миграции и заполняем бд
```shell
make fill-db-docker
```

Если все окей, то увидите что-то такое в конце
```
{"level":"info","ts":1699521811.2572942,"caller":"repository/fake_storage.go:305","msg":"end filling favourites\n"}
```
Если произошли какие-то проблемы во время заполнения бд. То это перезапишет данные в бд
```shell
make refill-db-docker
```

Это все бек + бд + pgadmin запущены
## Запуск локально из терминала / ide

1. Запускаем сразу все
```shell
make all
```
### Все без фронта 
```shell
make all-without-front
```

### Если нужно накатить миграции
```shell
make migrate-docker-up backend=local-deploy-backend-1
```
### Если нужно откатить миграции
```shell
make migrate-docker-down backend=local-deploy-backend-1
```
### Если нужно перезаполнить бд
```shell
make refill-db-docker backend=local-deploy-backend-1 backend-fs=local-deploy-backend-fs-1
```

### Локальная установка тула для миграций
```shell
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

### Тестирование
```shell
make test
```

### Сгенерировать swagger документацию
```shell
make swag
```

### Сгенерировать easyjson файл
```shell
easyjson <file.go>
```

### Сгенерировать mock файл
```shell
mockgen --source=<filename> --destination=<filename> --package=mocks
```