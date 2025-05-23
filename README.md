# MCP Server
MCP server for Chronosphere. Serves tools for fetching logs, metrics, traces, events as well as select entities.

This project is a work in progress. Many features are not yet implemented and features may be added, changed or removed without warning.

## Developing
### Running the server
#### Authentication to Chronosphere

This MCP server uses the same authentication methods as chronoctl. By default, the Makefile expects the API token to be stored in `.chronosphere_api_token`.

#### Run the mcp server
```sh
make run-server ORG_NAME=<your org here>
```

### Debugging MCP Tools

The MCP project provides an inspector useful for directly calling tools APIs. To use:

1. Change baseURL in config.yaml to be 0.0.0.0/sse instead of docker
1. Run `npx @modelcontextprotocol/inspector node build/index.js`.
1. Open http://127.0.0.1:6274/#resources , fill in `http://0.0.0.0:8080/sse` in the URL, with transport type SSE.

## Librechat agent (experimental)
See [chat/README.md]

## Agent (experimental)
See [agent/README.md]
