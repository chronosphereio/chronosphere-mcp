The agent is currently copied from a template. It doesn't do much yet, but we plan to use it for testing and building some examples of agentic behavior using the server.

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
