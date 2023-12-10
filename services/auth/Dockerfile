FROM golang:1.21.1-alpine3.18 as build

WORKDIR /var/auth

RUN apk update
RUN apk add git

COPY pkg pkg
COPY services/auth/cmd services/auth/cmd
COPY services/auth/internal services/auth/internal
COPY go.mod .
COPY go.sum .

RUN go mod tidy
RUN go mod download
RUN go build -o main ./services/auth/cmd/main.go

#=========================================================================================
FROM alpine:3.18 as production

WORKDIR /var/auth
COPY --from=build /var/auth/main main

RUN mkdir -p /var/log/backend

ENV ENVIRONMENT=development
ENV SERVICE_NAME=backend_auth
ENV PORT_AUTH=8082
ENV ADDRESS_AUTH_GRPC=:8012
ENV URL_DATA_BASE=postgres://postgres:postgres@postgres:5432/youla?sslmode=disable
ENV OUTPUT_LOG_PATH=/var/log/backend/logs_auth.json
ENV ERROR_OUTPUT_LOG_PATH=/var/log/backend/err_logs_auth.json

EXPOSE 8082

ENTRYPOINT ./main