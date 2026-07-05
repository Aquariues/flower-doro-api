.PHONY: dev prod build test tidy migrate-up docker-build

dev:
	docker compose -f docker-compose.dev.yml up --build

prod:
	docker compose -f docker-compose.prod.yml up --build -d

build:
	go build ./cmd/server

test:
	go test ./...

tidy:
	go mod tidy

docker-build:
	docker build -t flower-doro-api:local .

