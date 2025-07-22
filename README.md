# Chronosphere MCP Server
MCP server for Chronosphere. Serves tools for fetching logs, metrics, traces, events as well as select entities.

This project uses [semver](https://semver.org/) for release versions. We have not graduated to 1.0, so breaking
changes may occur for minor version bumps.

## MCP config with popular hosts (claude desktop, cursor)
### Remote Server
The easiest way to use the MCP server is using our remote hosted server:

```json
{
    "mcpServers": {
        "chronosphere": {
            "url": "https://<org name>.chronosphere.io/api/mcp/mcp",
            "headers": {
                "Authorization": "Bearer <chronosphere api token>"
            }
        }
    }
}
```

### Building from source
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
