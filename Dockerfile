FROM golang:1.21.1-alpine3.18

WORKDIR /var/backend

COPY cmd cmd
COPY internal internal
COPY go.mod .
COPY go.sum .

RUN go mod tidy
RUN go mod download

EXPOSE 8080

ENTRYPOINT go run  /var/backend/cmd/app/main.go
