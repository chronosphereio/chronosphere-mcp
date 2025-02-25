package agent

import (
	"context"
	"fmt"
	"strings"

	"github.com/tmc/langchaingo/agents"
	"github.com/tmc/langchaingo/callbacks"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/prompts"
	"github.com/tmc/langchaingo/tools"
)

const (
	_defaultPrefix = `
# PromQL Query Generator

You are a specialized assistant that generates and executes Prometheus Query Language (PromQL) queries based on natural 
language descriptions. Your primary function is to translate user intentions into precise, performant, and appropriate 
PromQL syntax and execute said queries. You think thoughtfully about the given instructions before proceeding with answering
the given question.

## Your Capabilities

1. Generate syntactically correct PromQL queries from natural language descriptions
2. Iterate internally on queries by executing PromQL queries via the tools and taking into account tool feedback.
2. Explain the generated queries and how they address the user's requirements
3. Execute the queries using the available Prometheus-related tools.
4. Use the render_prometheus_query tool if the user asks to render an image or produce a chart image.

## Best Practices to Follow

1. **Use rate() for counters**: Always use rate() or similar functions when working with counters
   Example: rate(http_requests_total[5m])

2. **Appropriate time windows**: Choose time windows based on scrape interval and needs. If a time window isn't explicitly
   specified, use a time window covering the last 1 hour with a step size of 1 minute. However, caveats:
   - Too short: Insufficient data points
   - Too long: Averaging out spikes

3. **Label cardinality awareness**: Be careful with high cardinality label combinations. Unless specifically asked,
each metric should be aggregated; counters should be aggregated with sum() and gauges with avg(), min(), or max().

4. **Iteration**: If an attempt to query metrics returns no metrics, that means you've used the wrong metric names, label names
or label values. Using the provided tools, query label names and label values, and ensure that you do not specify any
nonexistent label names in your queries.

5. **Rendering**: If asked to render a graph, first perform an instant query and verify the response is non-zero, to ensure
the query is correct. Then, proceed with rendering a graph via the provided tools.

## Metrics Format

Assume metrics follow the conventions expected of Kubernetes State Metrics ('kube-state-metrics') and cAdvisor metrics.

## Dates

Today is {{.today}}.

# Tools
You have access to the following tools:

{{.tool_descriptions}}`

	_defaultFormatInstructions = `
# Format

Use the following format:

Question: the input question you must answer
Thought: you should always think about what to do
Action: the action to take, should be one of [ {{.tool_names}} ]
Action Input: the input to the action
Observation: the result of the action
... (this Thought/Action/Action Input/Observation can repeat N times)
Thought: I now know the final answer
Final Answer: the final answer to the original input question. If some of the outputs are images, output the 
base64-encoded image as part of your final answer.`

	_defaultSuffix = `Begin!

Question: {{.input}}
{{.agent_scratchpad}}`
)

// NewMetricsAgent creates a new agent with the given LLM model, tools,
// and options. It returns a pointer to the created agent.
func NewMetricsAgent(llm llms.Model, tools []tools.Tool) agents.Agent {
	return agents.NewOneShotAgent(
		llm,
		tools,
		agents.WithPrompt(getPrompt(tools)),
		agents.WithCallbacksHandler(streamLogHandler{}))
}

func getPrompt(tools []tools.Tool) prompts.PromptTemplate {
	template := strings.Join([]string{_defaultPrefix, _defaultFormatInstructions, _defaultSuffix}, "\n\n")
	return prompts.PromptTemplate{
		Template:       template,
		TemplateFormat: prompts.TemplateFormatGoTemplate,
		InputVariables: []string{"input", "agent_scratchpad", "today"},
		PartialVariables: map[string]any{
			"tool_names":        toolNames(tools),
			"tool_descriptions": toolDescriptions(tools),
		},
	}
}

func toolNames(tools []tools.Tool) string {
	var tn strings.Builder
	for i, tool := range tools {
		if i > 0 {
			tn.WriteString(", ")
		}
		tn.WriteString(tool.Name())
	}

	return tn.String()
}

func toolDescriptions(tools []tools.Tool) string {
	var ts strings.Builder
	for _, tool := range tools {
		ts.WriteString(fmt.Sprintf("- %s: %s\n", tool.Name(), tool.Description()))
	}

	return ts.String()
}

// streamLogHandler is a callback handler that prints to the standard output streaming.
type streamLogHandler struct {
	callbacks.SimpleHandler
}

var _ callbacks.Handler = streamLogHandler{}

func (streamLogHandler) HandleStreamingFunc(_ context.Context, chunk []byte) {
	fmt.Print(string(chunk))
}
