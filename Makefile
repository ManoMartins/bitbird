include .env
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

test:
	go clean -testcache
	go test ./... -v

watch:
	reflex -r '\.go$$' -s -- sh -c 'make test'

migration_create:
	migrate create -ext sql -dir internal/database/migrations -seq $(name)

migration_up:
	migrate -path=internal/database/migrations -database "postgresql://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=disable" -verbose up

migration_down:
	migrate -path=internal/database/migrations -database "postgresql://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=disable" -verbose down
