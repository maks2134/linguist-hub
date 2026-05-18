.PHONY: run build test clean swagger

run:
	go run ./cmd/api

build:
	go build -o bin/api ./cmd/api

test:
	go test -v ./...

clean:
	rm -rf bin/
	rm -rf docs/

deps:
	go mod download
	go mod tidy

swagger:
	swag init -g cmd/api/docs.go -o docs --parseDependency --parseInternal

dev:
	go run ./cmd/api

gen: swagger
