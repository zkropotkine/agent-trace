package main

import (
	"log"

	"github.com/zkropotkine/agent-trace/assembler"
	"github.com/zkropotkine/agent-trace/config"
	"github.com/zkropotkine/agent-trace/internal/router"
)

func main() {
	cfg := config.LoadConfig()

	// Build app dependencies
	registry := assembler.BuildApp(cfg)

	// Setup router
	app := router.SetupRouter(*registry)

	log.Printf("AgentTrace running on %s", cfg.Port)
	if err := app.Run(cfg.Port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
