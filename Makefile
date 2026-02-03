.PHONY: build run test clean

build:
	go build -o bin/server ./cmd/main.go

run:
	go run ./cmd/main.go

run-race:
	go run --race ./cmd/main.go

lint:
	golangci-lint run ./... -v

generate:
	go generate ./...

format:
	gofmt -s -w .
	goimports -w -d .

test:
	go test ./internal/...

clean:
	rm -rf bin/
