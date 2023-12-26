FROM golang:1.21.1-alpine3.18 as build

WORKDIR /var/backend

RUN apk update
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

COPY cmd cmd
COPY internal internal
COPY pkg pkg
COPY go.mod .
COPY go.sum .

RUN go mod tidy
RUN go build -o main ./cmd/app/main.go

#=========================================================================================
FROM alpine:3.18 as production

WORKDIR /var/backend
COPY --from=build /var/backend/main main
COPY --from=build /go/bin/migrate migrate

RUN mkdir -p /var/log/backend
COPY db/migrations db/migrations

ENV ENVIRONMENT=development
ENV SERVICE_NAME=backend
ENV SCHEMA=http://
ENV ALLOW_ORIGIN=localhost:3000
ENV PORT_BACKEND=8080
ENV ADDRESS_FS_GRPC=backend-fs:8011
ENV ADDRESS_AUTH_GRPC=backend-auth:8012
ENV URL_DATA_BASE=postgres://postgres:postgres@localhost/youla?sslmode=disable
ENV PREMIUM_SHOP_ID=297668
ENV PREMIUM_SHOP_SECRET=test_qlRvNM1Btl6h3upjYaWEJSxfzjqyI6CdsrbcPsFS_3M
ENV PATH_CERT_FILE=/etc/ssl/goods-galaxy.ru.crt
ENV PATH_KEY_FILE=/etc/ssl/goods-galaxy.ru.key
ENV OUTPUT_LOG_PATH=/var/log/backend/logs.json
ENV ERROR_OUTPUT_LOG_PATH=/var/log/backend/err_logs.json

EXPOSE 8080

ENTRYPOINT ./main