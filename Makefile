build:
	docker-compose build twitter-clone

run:
	docker-compose up twitter-clone

test:
	go test -v ./...

swag:
	swag init -g cmd/main.go