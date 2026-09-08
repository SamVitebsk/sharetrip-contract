OAPI_CODEGEN_VERSION := v2.8.0
OAPI_CODEGEN_BIN := bin/oapi-codegen
OAPI_CODEGEN := ./$(OAPI_CODEGEN_BIN)

export GOBIN := $(CURDIR)/bin

.PHONY: tools generate generate-api run

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
	@trap 'status=$$?; trap - EXIT INT TERM; echo "DONE: Contract Service stopped (exit code $$status)"; exit $$status' EXIT INT TERM; go run ./cmd
