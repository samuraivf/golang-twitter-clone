build:
	docker-compose build todo-app

run:
	docker-compose up todo-app

test:
	go test -v ./...

swag:
	swag init -g cmd/main.go