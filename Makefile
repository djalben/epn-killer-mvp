.PHONY: build run test clean

build:
	go build -o bin/server ./cmd

run: build
	./bin/server

test:
	go test ./internal/...

clean:
	rm -rf bin/
