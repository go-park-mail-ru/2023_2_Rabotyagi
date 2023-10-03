# 2023_2_Rabotyagi
Backend репозиторий команды Работяги

### Наши контакты:

Владислав Ильинский: https://github.com/Vilinvil и тг https://t.me/Vilin0

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

### Push image to remote