%:
	@:

.PHONY: server client

ifneq (,$(wildcard ./.env))
    include .env
    export
endif


client:
	@go run ./cmd/client.go

worker:
	@go run ./cmd/worker.go

submit:
	@go run ./cmd/submitter.go 

server:
	@go run ./cmd/server.go 

migrationDir = ./server/migrations/

dev:
	air

generate:
	@echo "Generating sqlc repositoy..."
	@sqlc generate

migrate:
	@-mkdir ${migrationDir}
	$(eval args=$(filter-out $@,$(MAKECMDGOALS)))
	@goose -dir=${migrationDir} create ${args}

up:
	@goose -dir=${migrationDir} sqlite render_box.db up

down:
	@goose -dir=${migrationDir} sqlite render_box.db down
