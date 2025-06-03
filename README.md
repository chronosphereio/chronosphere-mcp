# MCP Server
MCP server for Chronosphere. Serves tools for fetching logs, metrics, traces, events as well as select entities.

This project is a work in progress. Many features are not yet implemented and features may be added, changed or removed without warning.

## MCP config with popular hosts (claude desktop, cursor)

First build the binary
```sh
make chronomcp
```

```json
{
  "mcpServers": {
    "chronosphere-mcp": {
      "command": "<PATH/TO/REPO>/bin/chronomcp",
      "args": [
        "-c",
        "<PATH/TO/REPO>/config.yaml"
      ],
      "env": {
        "CHRONOSPHERE_ORG_NAME": "<your org here>",
        "CHRONOSPHERE_API_TOKEN": "<your api token here>"
      }
    }
  }
}
```


## Developing
### Running the server
#### Authentication to Chronosphere

This MCP server uses the same authentication methods as chronoctl. By default, the Makefile expects the API token to be stored in `.chronosphere_api_token`.

#### Run the mcp server
```sh
make run-chronomcp CHRONOSPHERE_ORG_NAME=<your org here> CHRONOSPHERE_API_TOKEN=<your api token here>
```

### Debugging MCP Tools

The MCP project provides an inspector useful for directly calling tools APIs. To use:

1. Start the MCP server with streamable http transport `make run-chronomcp CONFIG_FILE=./config.http.yaml CHRONOSPHERE_ORG_NAME=<your org here>`
1. Run `npx @modelcontextprotocol/inspector node build/index.js`.
1. Open http://localhost:6274/#resources , fill in `http://0.0.0.0:8081/mcp` in the URL, with transport type Streamable HTTP.

## Librechat agent (experimental)
See [chat/README.md](chat/README.md)

## Agent (experimental)
See [agent/README.md](chat/README.md)


### Releases
We use [goreleaser](https://goreleaser.com/) to manage releases.

You'll need a [github token](https://github.com/settings/personal-access-tokens/) and put it in a .github_release_token file.
The token needs at least the following [permissions](https://goreleaser.com/ci/actions/#token-permissions)
- `content: write`
- `issues: write`

To create a new release, first create a tag:
```sh
git tag vX.Y.Z
git push origin vX.Y.Z
```

Then run the following command to perform a dry run of the release:

```
```sh
make release-dry-run
# verify the release looks good, then run:
make release
```
