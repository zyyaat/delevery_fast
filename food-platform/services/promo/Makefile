.PHONY: build run test lint migrate migrate-down create-migration docker-build clean

build:
	go build -o bin/server ./cmd/server

run:
	go run ./cmd/server

test:
	go test -v -race -cover ./...

lint:
	golangci-lint run

migrate:
	migrate -path migrations -database "$(DATABASE_URL)" up

migrate-down:
	migrate -path migrations -database "$(DATABASE_URL)" down 1

create-migration:
	@migrate create -ext sql -dir migrations -seq $(NAME)

docker-build:
	docker build -t food-platform/promo:latest -f deployments/Dockerfile .

clean:
	rm -rf bin/
