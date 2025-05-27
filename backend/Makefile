build:
	@go build -o bin/tcc cmd/main.go

test:
	@go test -v ./...

run: build
	@./bin/tcc

db-version:
	@goose --version

db-create:
	@if [ -z "$(word 2, $(MAKECMDGOALS))" ]; then \
		echo "Error: migration name is required. Use 'make db-create name_of_migration'"; \
		exit 1; \
	fi
	@goose create $(word 2, $(MAKECMDGOALS)) sql --dir db/migrations

db-up:
	@goose postgres -dir db/migrations "postgres://$(DATABASE_USER):$(DATABASE_PASSWORD)@$(DATABASE_HOST):$(DATABASE_PORT)/$(DATABASE_NAME)?sslmode=disable" up

db-down:
	@goose postgres -dir db/migrations "postgres://$(DATABASE_USER):$(DATABASE_PASSWORD)@$(DATABASE_HOST):$(DATABASE_PORT)/$(DATABASE_NAME)?sslmode=disable" down

db-redo:
	@goose postgres -dir db/migrations "postgres://$(DATABASE_USER):$(DATABASE_PASSWORD)@$(DATABASE_HOST):$(DATABASE_PORT)/$(DATABASE_NAME)?sslmode=disable" redo

db-status:
	@goose postgres -dir db/migrations "postgres://$(DATABASE_USER):$(DATABASE_PASSWORD)@$(DATABASE_HOST):$(DATABASE_PORT)/$(DATABASE_NAME)?sslmode=disable" status
