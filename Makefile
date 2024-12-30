include .env

COMPOSE_FILE=build/compose.yaml

setup:
	@echo "Installing CLI dependencies..."
	go install github.com/cespare/reflex@latest
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	@echo "All CLI tools installed!"

dev: service_up service_wait_database migration_up
	go run cmd/main.go

dev_ui:
	cd ui && npm run dev

service_up:
	docker compose -f $(COMPOSE_FILE) up -d

service_down:
	docker compose -f $(COMPOSE_FILE) down

service_stop:
	docker compose -f $(COMPOSE_FILE) stop

service_wait_database:
	go run scripts/wait_for_database.go

test: service_up service_wait_database
	go run cmd/main.go & \
	go test ./... -v

watch:
	reflex -r '\.go$$' -s -- sh -c 'go clean -testcache && go test ./... -v'

migration_create:
	migrate create -ext sql -dir internal/database/migrations -seq $(name)

migration_up:
	migrate -path=internal/database/migrations -database "postgresql://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=disable" -verbose up

migration_down:
	migrate -path=internal/database/migrations -database "postgresql://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=disable" -verbose down
