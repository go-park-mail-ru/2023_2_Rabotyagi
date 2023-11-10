FROM golang:1.21.1-alpine3.18 as build

WORKDIR /var/backend

RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
COPY cmd cmd
COPY internal internal
COPY go.mod .
COPY go.sum .

RUN go mod tidy
RUN go mod download
RUN go build -o main ./cmd/app/main.go
RUN go build -o fake_db ./cmd/fake_db/main.go

#=========================================================================================
FROM alpine:3.18 as production

WORKDIR /var/backend
COPY --from=build /var/backend/main main
COPY --from=build /var/backend/fake_db fake_db
COPY --from=build /go/bin/migrate migrate

RUN mkdir -p static/img
COPY static/images_for_fake_db static/images_for_fake_db
COPY db/migrations db/migrations

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

ENTRYPOINT ./main