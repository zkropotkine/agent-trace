package main

import (
	"context"

	"github.com/zkropotkine/agent-trace/assembler"
	"github.com/zkropotkine/agent-trace/config"
	"github.com/zkropotkine/agent-trace/internal/router"
	"github.com/zkropotkine/agent-trace/pkg/logger"
)

func main() {
	cfg := config.LoadConfig()

	// Create logger and context manager
	baseLogger := logger.NewLogRusLogger(config.Log{
		Level:  cfg.Log.Level,
		Format: cfg.Log.Format,
	})

	ctx := context.Background()

	logger.DefaultLogger = baseLogger
	ctx = logger.WithLogger(ctx, baseLogger)

	// Build app dependencies
	registry := assembler.BuildApp(ctx, cfg)

	// Setup router
	app := router.SetupRouter(ctx, *registry)

	baseLogger.Infof("AgentTrace running on %s", cfg.Port)
	if err := app.Run(cfg.Port); err != nil {
		baseLogger.Fatalf("failed to start server: %v", err)
	}
}
