tools_bin_path            := $(abspath ./_tools/bin)
server_bin_path           := $(abspath ./bin/server)
agent_bin_path           := $(abspath ./bin/agent)

CONFIG_FILE ?= config.yaml
AGENT_CONFIG_FILE ?= agent.yaml
ENV_FILE ?= .env
LIBRECHAT_CONFIG ?= librechat.yaml
AGENT_INPUTS_FILE ?= agent/resources/inputs.txt

.PHONY: install-tools
install-tools: go-version-check
	cd tools && GOBIN=$(tools_bin_path) go install \
		github.com/go-swagger/go-swagger/cmd/swagger

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

.PHONY: swagger-serve-dataunstable
swagger-serve-dataunstable:
	$(tools_bin_path)/swagger serve mcp-server/pkg/generated/clients/dataunstable.swagger.json

.PHONY: run-client
run-client:
	docker-compose down
	if [ ! -f $(ENV_FILE) ]; then \
		echo "Env file $(ENV_FILE) not found"; \
		exit 1; \
	fi
	if [ ! -f $(LIBRECHAT_CONFIG) ]; then \
		echo "Librechat config file $(LIBRECHAT_CONFIG) not found"; \
		exit 1; \
	fi
	docker-compose up -d


.PHONY: run-server
run-server: build-server
	if [ ! -f $(CONFIG_FILE) ]; then \
		echo "Config file $(CONFIG_FILE) not found"; \
		exit 1; \
	fi
	@echo "Starting MCP server..."
	@echo  "LibreChat should be on localhost:3080 once container up (check docker-compose ps)"
	$(server_bin_path) -c $(CONFIG_FILE) --org-name meta --api-token-filename .chronosphere_api_token --verbose

build-server:
	go build -o $(server_bin_path) ./mcp-server

.PHONY: run-agent
run-agent: build-agent
	if [ ! -f $(AGENT_CONFIG_FILE) ]; then \
		echo "Config file $(AGENT_CONFIG_FILE) not found"; \
		echo "Make a agent.yml"; \
		exit 1; \
	fi
	@echo "Starting agent..."
	$(agent_bin_path) -f $(AGENT_CONFIG_FILE) -i $(AGENT_INPUTS_FILE)

build-agent:
	go build -o $(agent_bin_path) ./agent

build-mcpgen:
	cd tools && go build -o $(tools_bin_path)/mcpgen ./cmd/mcpgen

tools-gen: build-mcpgen
	$(tools_bin_path)/mcpgen -spec ./generated/configv1/spec.json -pkg configv1 -target ./mcp-server/pkg/generated/tools/configv1 -allowed-entities \
		monitors,dashboards,slos

