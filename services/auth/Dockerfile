FROM golang:1.21.1-alpine3.18 as build

WORKDIR /var/auth

COPY cmd cmd
COPY internal internal
COPY go.mod .
COPY go.sum .

RUN go mod tidy
RUN go mod download
RUN go build -o main ./cmd/app/main.go

#=========================================================================================
FROM alpine:3.18 as production

WORKDIR /var/auth
COPY --from=build /var/auth/main main

RUN mkdir -p /var/log/backend

ENV ENVIRONMENT=development
ENV ADDRESS_AUTH_GRPC=127.0.0.1:8082
ENV SCHEMA=http://
ENV ALLOW_ORIGIN=localhost:3000
ENV PORT_BACKEND=8080
ENV PATH_TO_ROOT=/var/auth
ENV OUTPUT_LOG_PATH=stdout /var/log/backend/logs_auth.json
ENV ERROR_OUTPUT_LOG_PATH=stderr /var/log/backend/err_logs_auth.json
EXPOSE 8082
EXPOSE 8019

ENTRYPOINT ./main