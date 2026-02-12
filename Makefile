include .env
export

run:
	go run cmd/main.go

build:
	docker compose build

up:
	docker compose up -d

stop:
	docker compose down

rebuild:
	make stop
	make build
	make up

restart:
	make stop
	make build
	make up

logs:
	docker compose logs -f app

ps:
	docker compose ps

init:
	cp .env.example .env