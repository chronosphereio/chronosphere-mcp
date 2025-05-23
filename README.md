# MCP Server
MCP server for Chronosphere. Serves tools for fetching logs, metrics, traces, events as well as select entities.

This project is a work in progress. Many features are not yet implemented and features may be added, changed or removed without warning.

## Prototyping/manual testing with libra chat
See [chat/README.md]

## Agent
The agent is currently copied from a template. It doesn't do much yet, but we plan to use it for testing and building some examples of agentic behavior using the server.

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

## Using the Agent
To use the agent, first create a YAML file `agent.yaml`:

```
agent:
  openAIAPIKey: your-api-key
```

Then you can run:

```
make run-agent
```

with the MCP server running at 0.0.0.0:8080.

To augment the inputs that are provided to the model, edit `agent/resources/inputs.txt`; each non-blank line will correspond to a separate question to the agent.
