include .env
export

run:
	go run cmd/main.go

build:
	docker compose build

up:
	docker compose up -d

restart: stop build up

logs:
	docker compose logs -f app

ps:
	docker compose ps

stop:
	docker compose down

init:
	cp .env.example .env