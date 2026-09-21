OAPI_CODEGEN_VERSION := v2.8.0
OAPI_CODEGEN_BIN := bin/oapi-codegen
OAPI_CODEGEN := ./$(OAPI_CODEGEN_BIN)
GOOSE_VERSION := v3.26.0
GOOSE_BIN := bin/goose
COMPOSE := docker compose --env-file .env -f deploy/docker-compose.yml

export GOBIN := $(CURDIR)/bin

.PHONY: tools generate generate-api run up migrate-up migrate-status

tools: $(OAPI_CODEGEN_BIN)

$(OAPI_CODEGEN_BIN):
	@echo "START: installing oapi-codegen $(OAPI_CODEGEN_VERSION)"
	@go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@$(OAPI_CODEGEN_VERSION)
	@echo "DONE: oapi-codegen installed at $(OAPI_CODEGEN_BIN)"

generate: generate-api

generate-api: $(OAPI_CODEGEN_BIN)
	@echo "START: creating generated-code directory"
	@mkdir -p gen
	@echo "DONE: generated-code directory is ready"
	@echo "START: generating Go API from api/contract.yaml"
	@$(OAPI_CODEGEN) --config api/openapi.codegen.yaml api/contract.yaml
	@echo "DONE: Go API generated at gen/openapi.gen.go"

run: generate
	@echo "START: running Contract Service"
	@set -a; if [ -f .env ]; then . ./.env || exit $$?; fi; set +a; \
	trap 'status=$$?; trap - EXIT INT TERM; echo "DONE: Contract Service stopped (exit code $$status)"; exit $$status' EXIT INT TERM; go run ./cmd

$(GOOSE_BIN):
	@go install github.com/pressly/goose/v3/cmd/goose@$(GOOSE_VERSION)

up:
	@$(COMPOSE) up -d --wait --wait-timeout 60 postgres

migrate-up: $(GOOSE_BIN)
	@GOOSE_DRIVER=postgres GOOSE_DBSTRING='connect_timeout=5' ./$(GOOSE_BIN) -env .env -dir migrations up

migrate-status: $(GOOSE_BIN)
	@GOOSE_DRIVER=postgres GOOSE_DBSTRING='connect_timeout=5' ./$(GOOSE_BIN) -env .env -dir migrations status
