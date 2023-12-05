.PHONY: all
all: update-env compose-db-up compose-frontend-up go-mod-tidy test swag run

.PHONY: all-down
all-down: compose-db-down compose-frontend-down

.PHONY: all-without-front
all-without-front: compose-db-up go-mod-tidy test swag run

.PHONY: compose-frontend-up
compose-frontend-up:
	docker compose -f docker-compose.yml up -d frontend

.PHONY: compose-frontend-down
compose-frontend-down:
	docker compose -f docker-compose.yml down frontend

# for frontend
.PHONY: compose-full-up
compose-full-up: update-env
	docker compose -f docker-compose.yml up backend postgres pgadmin --build -d

.PHONY: compose-full-down
compose-full-down:
	docker compose -f docker-compose.yml down backend postgres pgadmin

.PHONY: compose-logs
compose-logs:
	docker compose -f docker-compose.yml logs

.PHONY: migrate-docker-up
migrate-docker-up:
	docker exec -it 2023_2_rabotyagi-backend-1 ./migrate -database postgres://postgres:postgres@postgres:5432/youla?sslmode=disable -path db/migrations up

.PHONY: migrate-docker-down
migrate-docker-down:
	docker exec -it 2023_2_rabotyagi-backend-1  ./migrate -database postgres://postgres:postgres@postgres:5432/youla?sslmode=disable -path db/migrations down

.PHONY: fill-db-docker
fill-db-docker: migrate-docker-up
	docker exec -it 2023_2_rabotyagi-backend-fs-1  ./fake_db postgres://postgres:postgres@postgres:5432/youla?sslmode=disable .

.PHONY: refill-db-docker
refill-db-docker: migrate-docker-down fill-db-docker

# dev
.PHONY: compose-db-up
compose-db-up:
	docker compose -f docker-compose.yml up -d postgres

.PHONY: compose-db-down
compose-db-down:
	docker compose -f docker-compose.yml down postgres

.PHONY: swag
swag:
	swag init -ot yaml --parseDependency --parseInternal -g cmd/app/main.go

.PHONY: go-mod-tidy
go-mod-tidy:
	go mod tidy

.PHONY: mkdir-bin
mkdir-bin:
	mkdir -p bin

.PHONY: test
test: mkdir-bin
	 go test -coverpkg=./... -coverprofile=bin/cover.out ./... \
 	 && cat bin/cover.out | grep -v "mocks" | grep -v ".pb" > bin/pure_cover.out \
  	 && go tool cover -html=bin/pure_cover.out -o=bin/cover.html \
  	 && go tool cover --func bin/pure_cover.out

.PHONY: build
build: mkdir-bin
	go build -o bin/main cmd/app/main.go

.PHONY: run
run: build
	sudo ./bin/main

.PHONY: create-migration
create-migration:
	migrate create -ext sql -dir ./db/migrations $(name)

.PHONY: migrate-up
migrate-up:
	migrate -database postgres://postgres:postgres@localhost:5432/youla?sslmode=disable -path db/migrations up

.PHONY: migrate-down
migrate-down:
	migrate -database postgres://postgres:postgres@localhost:5432/youla?sslmode=disable -path db/migrations down

.PHONY: build-fill-db
build-fill-db: mkdir-bin
	go build -o bin/fill-db services/file_service/cmd/fake_db/main.go

.PHONY: fill-db
fill-db: build-fill-db migrate-up
	sudo ./bin/fill-db postgres://postgres:postgres@localhost:5432/youla?sslmode=disable ./services/file_service

.PHONY: refill-db
refill-db: migrate-down migrate-up fill-db

.PHONY: logs
logs:
	cat /var/log/backend/logs.json | jq

.PHONY: compose-pull
compose-pull:
	docker compose down
	docker compose pull

.PHONY: update-env
update-env:
	mkdir -p ".env"
	mkdir -p "services/file_service/.env"
	mkdir -p "services/auth/.env"
	cp .env.example/.env.backend .env/.env.backend
	cp .env.example/.env.pgadmin .env/.env.pgadmin
	cp .env.example/.env.postgres .env/.env.postgres
	cp services/file_service/.env.example/.env services/file_service/.env/.env
	cp services/auth/.env.example/.env services/auth/.env/.env

.PHONY: generate
generate:
	protoc --go_out=. --go-grpc_out=. --go-grpc_opt=paths=source_relative --go_opt=paths=source_relative ./pkg/file_service/*.proto