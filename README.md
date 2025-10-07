# Chronosphere MCP Server
MCP server for [Chronosphere](https://chronosphere.io). Serves tools for fetching logs, metrics, traces, events as well as select entities.

This project uses [semver](https://semver.org/) for release versions. We have not graduated to 1.0, so breaking
changes may occur for minor version bumps.

## MCP config with popular hosts (claude desktop, cursor)
### Remote Server
The easiest way to use the MCP server is using our remote hosted server:

### Auth
You can use either a [Chronosphere API token](https://docs.chronosphere.io/tooling/api-info#create-an-api-token) or OAuth
with the Chronosphere MCP server. To use MCP with OAuth, the MCP client must support OAuth. 

OAuth support is new and have not been tested with all clients. If it does not work for you, please raise the issue to
Chronosphere support in slack with the following information:
1. MCP client you're using (e.g. VS code, codex, cursor, etc)
2. What steps you took to attempt authentication
3. What error you're seeing.

#### Cursor/VSCode
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

This configuration should work for Cursor and VSCode. Leave out the `headers` section to use OAuth instead of a Chronosphere API token.

More details for VSCode [here](https://code.visualstudio.com/docs/copilot/customization/mcp-servers) and Cursor [here](https://cursor.com/docs/context/mcp)

#### Claude code
Adding chronosphere MCP server to claude code
```shell
claude mcp add -t http -H "Authorization: ${CHRONOSPHERE_API_TOKEN}" chronosphere "https://${CHRONOSPHERE_ORG_NAME}.chronosphere.io/api/mcp/mcp"
```

You can leave out the Authorization header if you are using OAuth. Once you're in claude type `/mcp` and select the server to login to trigger the OAuth flow.

More details [here](https://docs.claude.com/en/docs/claude-code/mcp)

#### Codex CLI
```
experimental_use_rmcp_client = true
[mcp_servers.chronosphere]
url = "https://<org_name>.chronosphere.io/api/mcp/mcp"
bearer_token = "<chronosphere api token>"
```

For oauth login, you must enable `experimental_use_rmcp_client = true` and then run `codex mcp login chronosphere`

More details [here](https://github.com/openai/codex/blob/main/docs/config.md#mcp_servers)

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

## Available Tools

| Group | Tool Name | Description |
|-------|-----------|-------------|
| configapi | get_classic_dashboard | Get classic-dashboards resource |
| configapi | get_dashboard | Get dashboards resource |
| configapi | get_drop_rule | Get drop-rules resource |
| configapi | get_mapping_rule | Get mapping-rules resource |
| configapi | get_monitor | Get monitors resource |
| configapi | get_notification_policy | Get notification-policies resource |
| configapi | get_recording_rule | Get recording-rules resource |
| configapi | get_rollup_rule | Get rollup-rules resource |
| configapi | get_slo | Get slos resource |
| configapi | list_classic_dashboards | List classic-dashboards resources |
| configapi | list_dashboards | List dashboards resources |
| configapi | list_drop_rules | List drop-rules resources |
| configapi | list_mapping_rules | List mapping-rules resources |
| configapi | list_monitors | List monitors resources |
| configapi | list_notification_policies | List notification-policies resources |
| configapi | list_recording_rules | List recording-rules resources |
| configapi | list_rollup_rules | List rollup-rules resources |
| configapi | list_slos | List slos resources |
| events | get_events_metadata | List properties you can query on events |
| events | list_events | List events from a given query |
| events | list_events_label_values | List values for a given label name |
| logs | get_log | Get a full log message by its ID. The ID is the unique identifier for the log. |
| logs | get_log_histogram | Get histogram of logs from a given query |
| logs | list_log_field_names | List field names of logs |
| logs | list_log_field_values | List field values of logs |
| logs | query_logs_range | Execute a range query for logs. This endpoint returns logs as either timeSeries or gridData. It may return a large amount of data, so be careful putting the result of this direction into context. U... |
| logscale | query_logscale | Query LogScale repository with a given query string. LogScale uses its own query language for searching and analyzing logs. |
| metrics | list_prometheus_label_names | Returns the list of label names (keys) available on metrics that match the given selectors. Use this tool when you need to discover what labels are available on specific metrics or services. Exampl... |
| metrics | list_prometheus_label_values | Returns the list of values for a specific label name, optionally filtered by selectors. Use this tool when you know the label name and want to discover what values it has across your metrics. Commo... |
| metrics | list_prometheus_series | Returns the complete time series (full label sets with all key-value pairs) that match the given selectors. Each result shows the exact combination of labels for an active time series. Use this too... |
| metrics | list_prometheus_series_metadata |  |
| metrics | query_prometheus_instant | Evaluates a Prometheus instant query at a single point in time |
| metrics | query_prometheus_range | Executes a Prometheus PromQL query over a specified time range and returns time series data points as JSON. Supports standard PromQL syntax plus Chronosphere custom functions: - cardinality_estimat... |
| metrics | render_prometheus_range_query | Evaluates a Prometheus expression query over a range of time and renders it as a PNG image. |
| monitors | list_monitor_statuses | Lists the current status of monitors in Chronosphere. Returns monitor statuses with alert states and optional signal and series details. |
| traces | list_traces | List traces from a given query |

*Note: To regenerate this table after tool updates, run: `make tools-gen && go run scripts/generate-tools-table.go`*

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
