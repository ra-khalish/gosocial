include .envrc
MIGRATIONS_PATH = ./cmd/migrate/migrations

.PHONY: migrate-create
migration:
	@migrate create -seq -ext sql -dir $(MIGRATIONS_PATH) $(filter-out $@,$(MAKECMDGOALS))

.PHONY: migrate-up
migrate-up:
	@migrate -verbose -source file://cmd/migrate/migrations -database postgres://admin:adminpassword@localhost:5432/social?sslmode=disable up

.PHONY: migrate-down
migrate-down:
	@migrate -verbose -source file://cmd/migrate/migrations -database postgres://admin:adminpassword@localhost:5432/social?sslmode=disable up

.PHONY: seed
seed:
	@go run ./cmd/migrate/seed/main.go