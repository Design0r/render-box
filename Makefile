%:
	@:

ifneq (,$(wildcard ./.env))
    include .env
    export
endif

client:
	@go run ./client.go 

server:
	@go run ./server.go 

migrationDir = ./migrations/

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
	GOOSE_DRIVER=sqlite GOOSE_DBSTRING=render_box.db goose -dir=${migrationDir} up

down:
	@GOOSE_DRIVER=sqlite GOOSE_DBSTRING=render_box.db goose -dir=${migrationDir} down
