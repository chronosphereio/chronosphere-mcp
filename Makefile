tools_bin_path            := $(abspath ./_tools/bin)
server_bin_path           := $(abspath ./bin/chronomcp)
agent_bin_path           := $(abspath ./bin/agent)

CONFIG_FILE ?= config.yaml
AGENT_CONFIG_FILE ?= agent.yaml
ENV_FILE ?= .env
LIBRECHAT_CONFIG ?= librechat.yaml
AGENT_INPUTS_FILE ?= agent/resources/inputs.txt
LDFLAGS ?= -ldflags="-X github.com/chronosphereio/chronosphere-mcp/pkg/version.Version=$(shell git describe --tags --always --dirty) -X github.com/chronosphereio/chronosphere-mcp/pkg/version.BuildDate=$(shell date -u +%Y-%m-%dT%H:%M:%SZ) -X github.com/chronosphereio/chronosphere-mcp/pkg/version.GitCommit=$(shell git rev-parse HEAD)"

.PHONY: install-tools
install-tools: go-version-check
	cd tools && GOBIN=$(tools_bin_path) go install  github.com/go-swagger/go-swagger/cmd/swagger  github.com/golangci/golangci-lint/cmd/golangci-lint

.PHONY: go-version-check
go-version-check:
	# make sure you're running the right version of Go, otherwise builds/codegen/tests
	# may have inconsistent results that are hard to debug.
	go version | grep go1.24 || (echo "Error: you must be running go1.24.x" && exit 1)

.PHONY: swagger-gen
swagger-gen: install-tools
	rm -rf generated/dataunstable/dataunstable
	rm -rf generated/dataunstable/models
	$(tools_bin_path)/swagger generate client -f generated/dataunstable/spec.json -t generated/dataunstable -c dataunstable
	rm -rf generated/configv1/dataunstable
	rm -rf generated/configv1/models
	$(tools_bin_path)/swagger generate client -f generated/configv1/spec.json -t generated/configv1 -c configv1
	rm -rf generated/stateunstable/stateunstable
	rm -rf generated/stateunstable/models
	$(tools_bin_path)/swagger generate client -f generated/stateunstable/spec.json -t generated/stateunstable -c stateunstable

.PHONY: swagger-serve-dataunstable
swagger-serve-dataunstable:
	$(tools_bin_path)/swagger serve mcp-server/pkg/generated/clients/dataunstable.swagger.json

.PHONY: run-chat
run-chat:
	@(cd chat && \
		docker-compose down && \
		if [ ! -f $(ENV_FILE) ]; then \
			echo "Env file $(ENV_FILE) not found"; \
			exit 1; \
		fi && \
		if [ ! -f $(LIBRECHAT_CONFIG) ]; then \
			echo "Librechat config file $(LIBRECHAT_CONFIG) not found"; \
			exit 1; \
		fi && \
		echo  "LibreChat should be on localhost:3080 once container up (check docker-compose ps)" && \
		docker-compose up -d)

.PHONY: run-server build-server chronomcp run-chronomcp
chronomcp:
	go build $(LDFLAGS) -o $(server_bin_path) ./mcp-server

run-chronomcp: chronomcp
	CONFIG_FILE=$(CONFIG_FILE) ./scripts/run-chronomcp.sh $(server_bin_path)

build-server: chronomcp # alias for backwards compatibility

run-server: run-chronomcp # alias for backwards compatibility


.PHONY: run-agent
run-agent: build-agent
	if [ ! -f $(AGENT_CONFIG_FILE) ]; then \
		echo "Config file $(AGENT_CONFIG_FILE) not found"; \
		echo "Make a agent.yml"; \
		exit 1; \
	fi
	@echo "Starting agent..."
	$(agent_bin_path) -f $(AGENT_CONFIG_FILE) -i $(AGENT_INPUTS_FILE)

.PHONY: build-agent
build-agent:
	go build $(LDFLAGS) -o $(agent_bin_path) ./agent

.PHONY: build-mcpgen
build-mcpgen:
	cd tools && go build $(LDFLAGS) -o $(tools_bin_path)/mcpgen ./cmd/mcpgen

.PHONY: tools-mcpgen
tools-gen: build-mcpgen
	$(tools_bin_path)/mcpgen -spec ./generated/configv1/spec.json -pkg configv1 -target ./mcp-server/pkg/generated/tools/configv1 -allowed-entities \
		monitors,dashboards,slos

.PHONY: lint
lint: install-tools
	@echo "--- :golang: linting code"
	GOFLAGS=$(GOFLAGS) $(tools_bin_path)/golangci-lint run

.PHONY: all-gen test-all-gen tidy
tidy:
	go mod tidy

all-gen: swagger-gen tools-gen
	@echo "--- :golang: all codegen complete"

test-all-gen: tidy all-gen
	./scripts/check-branch.sh
