.PHONY: all all-prod

all: update-env compose-db-up compose-frontend-up go-mod-tidy test swag run
all-down: compose-db-down compose-frontend-down
all-without-front: compose-db-up go-mod-tidy test swag run

compose-frontend-up:
	docker compose -f docker-compose.yml up -d frontend

compose-frontend-down:
	docker compose -f docker-compose.yml down frontend

# prod
compose-full-up: update-env
	docker compose -f docker-compose.yml up backend postgres pgadmin --build -d

compose-full-down:
	docker compose -f docker-compose.yml down backend postgres pgadmin

compose-logs:
	docker compose -f docker-compose.yml logs

migrate-docker-up:
	docker exec -it 2023_2_rabotyagi-backend-1 ./migrate -database postgres://postgres:postgres@postgres:5432/youla?sslmode=disable -path db/migrations up

migrate-docker-down:
	docker exec -it 2023_2_rabotyagi-backend-1  ./migrate -database postgres://postgres:postgres@postgres:5432/youla?sslmode=disable -path db/migrations down

fill-db-docker: migrate-docker-up
	docker exec -it 2023_2_rabotyagi-backend-1  ./fake_db

refill-db-docker: migrate-docker-down fill-db-docker

# dev
compose-db-up:
	docker compose -f docker-compose.yml up -d postgres

compose-db-down:
	docker compose -f docker-compose.yml down postgres

swag:
	swag init -ot yaml --parseDependency --parseInternal -g cmd/app/main.go

go-mod-tidy:
	go mod tidy

mkdir-bin:
	mkdir -p bin

test: mkdir-bin
	 go test -coverprofile=bin/cover.out ./internal/... && go tool cover -html=bin/cover.out -o=bin/cover.html && go tool cover --func bin/cover.out

build: mkdir-bin
	go build -o bin/main cmd/app/main.go

run: build
	sudo ./bin/main

migrate-up:
	migrate -database postgres://postgres:postgres@localhost:5432/youla?sslmode=disable -path db/migrations up

migrate-down:
	migrate -database postgres://postgres:postgres@localhost:5432/youla?sslmode=disable -path db/migrations down

build-fill-db: mkdir-bin
	go build -o bin/fill-db cmd/fake_db/main.go

fill-db: build-fill-db migrate-up
	sudo ./bin/fill-db

refill-db: migrate-down migrate-up fill-db

logs:
	cat /var/log/backend/logs.json | jq

compose-pull:
	docker compose down
	docker compose pull

update-env:
	mkdir -p ".env"
	cp .env.example/.env.backend .env/.env.backend
	cp .env.example/.env.pgadmin .env/.env.pgadmin
	cp .env.example/.env.postgres .env/.env.postgres