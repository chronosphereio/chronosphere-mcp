# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

If you have things you want to save for future reference, write a file MEMORY.md in the root of the repository.
Feel free to update MEMORY.md with any useful information you find while working on this repository.

## Development Commands

### Environment Setup
- Go version check: `make go-version-check` (requires Go 1.25.1)
- Install tools: `make install-tools`

### Building
- Build server: `make chronomcp` - run this after making changes to the server code to make sure it compiles.
- Generate Swagger clients: `make swagger-gen`
- Generate MCP tools: `make tools-gen`

### Generated Files
- Do not manually edit generated artifacts in `generated/*` or `mcp-server/pkg/generated/*`.
- Do not manually edit or delete Swagger specs (example `generated/statev1/spec.json`).
- If generated outputs need to change, regenerate with `make swagger-gen`, `make tools-gen`, or `make all-gen`.

### Testing
- When writing tests, use the github.com/stretchr/testify package instead of the standard library's asserts.
- Run tests using `go test ./path/to/package` to run specific tests for a package.
- Use table tests for better organization and readability of tests.

### Running
- Run MCP server: `make run-chronomcp CHRONOSPHERE_ORG_NAME=<your_org_name> CHRONOSPHERE_API_TOKEN=<your_api_token>` (uses `config.http.yaml` by default; override with `CONFIG_FILE=...`)

### Linting
- Lint code: `make lint` - always run this after changing Go code and fix lint warnings and errors that are flagged.

### Commit Message Conventions (GoReleaser)
- Use Conventional Commit-style subjects so release notes are grouped correctly.
- Preferred format: `<type>(<optional-scope>)<optional-!>: <summary>`
- Use these commit types for changelog groups:
  - `feat:` -> 🚀 Features
  - `fix:` -> 🐛 Bug Fixes
  - `docs:` -> 📖 Documentation
  - `chore:` -> 🏠 Housekeeping
- Breaking changes: add `!` after type/scope (example: `feat(api)!: remove deprecated endpoint`).
- Commits starting with `test:`, `scripts:`, and `ci:` are excluded from GoReleaser changelog groups.
- Avoid merge-style subjects (`Merge pull request`, `Merge branch`) for authored commits; those are excluded too.

### Ports and Transports
- `config.http.yaml` (default): Streamable HTTP on `0.0.0.0:8081`
- `config.sse.yaml`: SSE on `0.0.0.0:8080` - NOTE: sse transport is deprecated in favor of streamable http.
- `config.yaml`: stdio transport

## Project Architecture

This repository contains a Model Context Protocol (MCP) server for Chronosphere, an observability platform. The server provides tools for fetching logs, metrics, traces, events, and select entities.

### Key Components

1. **MCP Server**: 
   - Main server component that serves the MCP API
   - Located in `mcp-server/` directory
   - Uses Streamable HTTP transport
   - Includes authentication via Chronosphere API tokens

2. **Tools System**:
   - Located in `mcp-server/pkg/tools/`
   - MCPTools interface allows grouping tools together
   - Each tool has metadata and a handler function
   - Tools can be disabled via configuration

3. **Resources**:
   - Resources for logs, metrics, etc.
   - Located in `mcp-server/pkg/resources/`

### Key Files

- **config.yaml**: Main configuration for the MCP server
- **Makefile**: Contains all build/run commands

## Authentication

`CHRONOSPHERE_ORG_NAME` is required when running via `make run-chronomcp`.
`CHRONOSPHERE_API_TOKEN` is optional if using bearer forwarding/OAuth-compatible hosts.

## Debugging

For debugging MCP Tools:
1. Start server with transport config:
   - Streamable HTTP: `make run-chronomcp CONFIG_FILE=config.http.yaml CHRONOSPHERE_ORG_NAME=<org> CHRONOSPHERE_API_TOKEN=<token>`
2. Run `npx @modelcontextprotocol/inspector`
3. Open `http://localhost:6274` and connect using:
   - Streamable HTTP URL: `http://0.0.0.0:8081/mcp`
