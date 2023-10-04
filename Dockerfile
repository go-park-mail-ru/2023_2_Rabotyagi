FROM golang:1.21.1-alpine3.18

WORKDIR /var/backend

COPY . .

RUN go mod tidy
RUN go mod download

WORKDIR /var/backend/cmd/app
ENTRYPOINT go run main.go

EXPOSE 8080