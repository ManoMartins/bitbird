.PHONY: up down stop

COMPOSE_FILE=build/compose.yaml

dev: up
	go run cmd/main.go

up:
	docker-compose -f $(COMPOSE_FILE) up -d

down:
	docker-compose -f $(COMPOSE_FILE) down

stop:
	docker-compose -f $(COMPOSE_FILE) stop