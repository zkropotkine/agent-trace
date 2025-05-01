package assembler

import (
	"log"

	"github.com/zkropotkine/agent-trace/config"
	"github.com/zkropotkine/agent-trace/internal/db"
	"github.com/zkropotkine/agent-trace/internal/handler"
	"github.com/zkropotkine/agent-trace/internal/repository"
	"github.com/zkropotkine/agent-trace/internal/router"
)

func BuildApp(cfg *config.Config) *router.RouteRegistry {
	client, err := db.NewMongoClient(cfg.MongoURI)
	if err != nil {
		log.Fatalf("failed to connect to MongoDB: %v", err)
	}

	database := client.Database("agentTrace")
	traceRepo := repository.NewMongoTraceRepository(database)
	traceHandler := handler.NewTraceHandler(traceRepo)

	registry := &router.RouteRegistry{
		TraceHandler: traceHandler,
	}

	return registry
}
