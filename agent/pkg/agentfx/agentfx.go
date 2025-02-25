package agentfx

import (
	"context"
	"fmt"
	"log"
	"os"
	"path"

	"github.com/mark3labs/mcp-go/client"
	"github.com/tmc/langchaingo/agents"
	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/llms/openai"
	"go.uber.org/config"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"gopkg.in/validator.v2"

	pkgagent "github.com/chronosphereio/mcp-server/agent/pkg/agent"
	"github.com/chronosphereio/mcp-server/agent/pkg/configfx"
	"github.com/chronosphereio/mcp-server/agent/pkg/mcp"
)

// Module registers the server.
var Module = fx.Invoke(invoke)

const configKey = "agent"

type params struct {
	fx.In
	LifeCycle fx.Lifecycle

	ConfigProvider config.Provider
	Inputs         configfx.Inputs
}

type Config struct {
	OpenAIAPIKey string `yaml:"openAIAPIKey"`
}

func invoke(p params) error {
	var cfg Config
	if err := p.ConfigProvider.Get(configKey).Populate(&cfg); err != nil {
		return err
	}
	if err := validator.Validate(cfg); err != nil {
		return err
	}

	tmpDir := os.TempDir()
	loggerCfg := zap.NewProductionConfig()
	loggerCfg.OutputPaths = []string{"stdout", path.Join(tmpDir, "agent.log")}
	loggerCfg.ErrorOutputPaths = []string{"stderr", path.Join(tmpDir, "agent_error.log")}
	logger, err := loggerCfg.Build()
	if err != nil {
		return fmt.Errorf("failed to create logger: %v", err)
	}

	logger.Info("starting agent with logger",
		zap.Strings("log", loggerCfg.OutputPaths),
		zap.Strings("errorLog", loggerCfg.ErrorOutputPaths))

	p.LifeCycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return runAgents(ctx, p.Inputs, cfg)
		},
	})
	return nil
}

func runAgents(ctx context.Context, inputs []string, cfg Config) error {
	for _, input := range inputs {
		if err := runAgent(ctx, input, cfg); err != nil {
			return err
		}
	}
	return nil
}

func runAgent(ctx context.Context, input string, cfg Config) error {
	llm, err := openai.New(openai.WithToken(cfg.OpenAIAPIKey))
	if err != nil {
		return err
	}

	mcpClient, err := client.NewSSEMCPClient(
		"http://0.0.0.0:8080/sse",
	)
	if err != nil {
		log.Fatalf("Failed to create MCP client: %v", err)
	}
	defer mcpClient.Close()
	if err := mcpClient.Start(ctx); err != nil {
		return err
	}

	// Create the adapter
	adapter, err := mcp.NewAdapter(mcpClient)
	if err != nil {
		log.Fatalf("Failed to create adapter: %v", err)
	}

	// Get all tools from MCP server
	mcpTools, err := adapter.Tools()
	if err != nil {
		log.Fatalf("Failed to get tools: %v", err)
	}

	// Create a agent with the tools
	agent := pkgagent.NewMetricsAgent(
		llm,
		mcpTools,
	)
	executor := agents.NewExecutor(agent, agents.WithMaxIterations(20))

	fmt.Printf("=== Executing question ===\n")
	fmt.Println(input)
	fmt.Printf("\n=== Thoughts & responses ===\n")
	answer, err := chains.Run(context.Background(), executor, input)
	if err != nil {
		log.Fatalf("failed to call executor: %v", err)
	}
	fmt.Println(answer)
	return nil
}
