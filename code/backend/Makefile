.PHONY: build run test lint migrate tidy clean

APP_NAME := xyfamily
BUILD_DIR := bin

build:
	mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(APP_NAME) ./cmd/xyfamily/

run:
	go run ./cmd/xyfamily/

test:
	go test -v -race ./...

lint:
	golangci-lint run ./...

tidy:
	go mod tidy

migrate-up:
	psql -h localhost -U xyfamily -d xyfamily -f migrations/V001__init.sql

migrate:
	psql -h localhost -U xyfamily -d xyfamily -f migrations/V001__init.sql

clean:
	rm -rf $(BUILD_DIR)

docker-up:
	docker compose -f docker-compose.dev.yml up -d

docker-down:
	docker compose -f docker-compose.dev.yml down
