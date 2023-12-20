.PHONY: all
all: update-env go-mod-tidy test swag compose-full-up

.PHONY: all-without-front
all-without-front: update-env compose-frontend-up


# for frontend
.PHONY: compose-frontend-up
compose-frontend-up: update-env
	docker compose -f local-deploy/docker-compose-frontend.yml up --build -d

.PHONY: compose-frontend-down
compose-frontend-down:
	docker compose -f local-deploy/docker-compose-frontend.yml down

backend=local-deploy-backend-1
.PHONY: migrate-docker-up
migrate-docker-up:
	docker exec -it ${backend} ./migrate -database postgres://postgres:postgres@postgres:5432/youla?sslmode=disable -path db/migrations up

.PHONY: migrate-docker-down
migrate-docker-down:
	docker exec -it ${backend}  ./migrate -database postgres://postgres:postgres@postgres:5432/youla?sslmode=disable -path db/migrations down

backend-fs=local-deploy-backend-fs-1
.PHONY: fill-db-docker
fill-db-docker: migrate-docker-up
	docker exec -it ${backend-fs}  ./fake_db postgres://postgres:postgres@postgres:5432/youla?sslmode=disable .

.PHONY: refill-db-docker
refill-db-docker: migrate-docker-down fill-db-docker

# dev
.PHONY: compose-up
compose-up:
	docker compose -f docker-compose.yml up -d postgres frontend nginx grafana node-exporter

.PHONY: compose-down
compose-down:
	docker compose -f docker-compose.yml down postgres frontend nginx grafana node-exporter

.PHONY: compose-full-up
compose-full-up: update-env
	docker compose -f docker-compose.yml up --build -d

.PHONY: compose-full-down
compose-full-down:
	docker compose -f docker-compose.yml down


.PHONY: compose-pull
compose-pull:
	docker compose down
	docker compose pull

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
	 go test --race -coverpkg=./... -coverprofile=bin/cover.out ./... \
 	 && cat bin/cover.out | grep -v "mocks" | grep -v "easyjson" | grep -v ".pb" > bin/pure_cover.out \
  	 && go tool cover -html=bin/pure_cover.out -o=bin/cover.html \
  	 && go tool cover --func bin/pure_cover.out

.PHONY: test-actions
test-actions:
	./scripts/test-actions.sh

.PHONY: lint
lint:
	golangci-lint run --timeout=3m ./...

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


.PHONY: update-env
update-env:
	mkdir -p ".env"
	mkdir -p "services/file_service/.env"
	mkdir -p "services/auth/.env"
	cp -r .env.example/. .env
	cp services/file_service/.env.example/.env services/file_service/.env/.env
	cp services/auth/.env.example/.env services/auth/.env/.env

.PHONY: generate
generate:
	protoc --go_out=. --go-grpc_out=. --go-grpc_opt=paths=source_relative --go_opt=paths=source_relative ./pkg/file_service/*.proto