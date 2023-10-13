FROM golang:1.21.1-alpine3.18

WORKDIR /var/backend

COPY cmd cmd
COPY internal internal
COPY pkg pkg
COPY go.mod .
COPY go.sum .

RUN go mod tidy
RUN go mod download

EXPOSE 8080
WORKDIR /var/backend/cmd/app

ENTRYPOINT go run main.go
