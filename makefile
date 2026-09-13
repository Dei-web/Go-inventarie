.PHONY: up down build restart logs

up:
	docker compose up -d

down:
	docker compose down

build:
	docker compose build

restart:
	docker compose down
	docker compose up -d

logs:
	docker compose logs -f
