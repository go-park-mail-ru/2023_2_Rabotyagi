FROM golang:1.21.1-alpine3.18

WORKDIR /var/backend

COPY cmd cmd
COPY db/images_for_fake_db db/images_for_fake_db
COPY internal internal
COPY go.mod .
COPY go.sum .

RUN go mod tidy
RUN go mod download

ENV ALLOW_ORIGIN=localhost:3000
ENV PORT_BACKEND=8080
ENV POSTGRES_DB=youla
ENV POSTGRES_USER=postgres
ENV POSTGRES_PASSWORD=postgres
ENV POSTGRES_PORT=5432
ENV POSTGRES_ADDRESS=localhost
ENV PATH_TO_ROOT=/var/backend

ENV URL_DATA_BASE=postgres://postgres:postgres@localhost/youla?sslmode=disable
ENV SCHEMA=http://

EXPOSE 8080

ENTRYPOINT go run  /var/backend/cmd/app/main.go
