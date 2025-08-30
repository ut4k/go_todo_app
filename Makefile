.PHONY: help build build-local up down logs ps test
.DEFAULT_GOAL := help

DOCKER_TAG := latest
build:
	docker build -t ut4k/gotodo:${DOCKER_TAG} \
		-- target deploy ./

build-local:
	docker compose build --no-cache

up:
	docker compose up -docker

down:
	docker compose down

logs:
	docker compose logs -f

ps:
	docker compose ps

test:
	go test -race -shuffle=on ./...

