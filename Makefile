include .env

build:
	go build -o bin/server ./cmd/main.go

run:
	go run ./cmd/main.go

run-race:
	go run --race ./cmd/main.go

clean:
	rm -rf bin/

lint:
	golangci-lint run ./... -v

generate:
	go generate ./...

format:
	gofmt -s -w .
	goimports -w -d .

test:
	go test ./internal/...

bin-deps:
	go install github.com/pressly/goose/v3/cmd/goose@latest

docker-up:
	docker-compose -f ./deployment/docker-compose.yaml --env-file ./deployment/.env up --build

docker-down:
	docker-compose -f ./deployment/docker-compose.yaml down

migrate-create:
	goose postgres  "$(POSTGRES_DSN)" create "$(filter-out $@, $(MAKECMDGOALS))" sql -dir ./migrations/postgres 
 
migrate-force:
	goose postgres "$(POSTGRES_DSN)" -dir ./migrations/postgres  down-to "$(filter-out $@, $(MAKECMDGOALS))"
 
migrate-up:
	goose postgres "$(POSTGRES_DSN)" -dir ./migrations/postgres  up
 
migrate-down:
	goose postgres "$(POSTGRES_DSN)" -dir ./migrations/postgres  down