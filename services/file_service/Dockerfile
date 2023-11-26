FROM golang:1.21.1-alpine3.18 as build

WORKDIR /var/file_service

COPY cmd cmd
COPY internal internal
COPY go.mod .
COPY go.sum .

RUN go mod tidy
RUN go mod download
RUN go build -o main ./cmd/app/main.go

#=========================================================================================
FROM alpine:3.18 as production

WORKDIR /var/file_service
COPY --from=build /var/file_service/main main

RUN mkdir -p /var/log/backend
RUN mkdir -p static/img
COPY services/file_service/static/images_for_fake_db static/images_for_fake_db

ENV ENVIRONMENT=development
ENV ADDRESS_FS_GRPC=127.0.0.1:8081
ENV SCHEMA=http://
ENV ALLOW_ORIGIN=localhost:3000
ENV PORT_BACKEND=8080
ENV PATH_TO_ROOT=/var/file_service
ENV FILE_SERVICE_DIR=/var/file_service/static/img
ENV OUTPUT_LOG_PATH=stdout /var/log/backend/logs_fs.json
ENV ERROR_OUTPUT_LOG_PATH=stderr /var/log/backend/err_logs_fs.json

EXPOSE 8081
EXPOSE 8018

ENTRYPOINT ./main