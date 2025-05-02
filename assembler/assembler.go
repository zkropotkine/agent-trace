package assembler

import (
	"context"
	"log"

	"github.com/zkropotkine/agent-trace/config"
	"github.com/zkropotkine/agent-trace/internal/db"
	"github.com/zkropotkine/agent-trace/internal/handler"
	"github.com/zkropotkine/agent-trace/internal/repository"
	"github.com/zkropotkine/agent-trace/internal/router"
)

func BuildApp(ctx context.Context, cfg *config.Config) *router.RouteRegistry {
	mongo := cfg.Mongo

	client, err := db.NewMongoClient(mongo.URI)
	if err != nil {
		log.Fatalf("failed to connect to MongoDB: %v", err)
	}

	dbClient := client.Database(mongo.DB)
	collection := dbClient.Collection(mongo.Collection)

	traceRepo := repository.NewMongoTraceRepository(collection)
	traceHandler := handler.NewTraceHandler(traceRepo)

	registry := &router.RouteRegistry{
		TraceHandler: traceHandler,
	}

	return registry
}
