FROM golang:1.21.1-alpine3.18 as build

WORKDIR /var/file_service

RUN apk update
RUN apk add git

COPY pkg pkg
COPY services/file_service/cmd services/file_service/cmd
COPY services/file_service/internal services/file_service/internal
COPY go.mod .
COPY go.sum .

RUN go mod tidy
RUN go mod download
RUN go build -o main ./services/file_service/cmd/app/main.go
RUN go build -o fake_db ./services/file_service/cmd/fake_db/main.go

#=========================================================================================
FROM alpine:3.18 as production

WORKDIR /var/file_service
COPY --from=build /var/file_service/main main
COPY --from=build /var/file_service/fake_db fake_db

RUN mkdir -p /var/log/backend
RUN mkdir -p static/img

ENV ENVIRONMENT=development
ENV ADDRESS_FS_GRPC=127.0.0.1:8081
ENV SCHEMA=http://
ENV ALLOW_ORIGIN=localhost:3000
ENV PORT_BACKEND=8080
ENV PATH_TO_ROOT=/var/file_service
ENV FILE_SERVICE_DIR=/var/file_service/static/img
ENV OUTPUT_LOG_PATH=/var/log/backend/logs_fs.json
ENV ERROR_OUTPUT_LOG_PATH=/var/log/backend/err_logs_fs.json

EXPOSE 8081
EXPOSE 8018

ENTRYPOINT ./main