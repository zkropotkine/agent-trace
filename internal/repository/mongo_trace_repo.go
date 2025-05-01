package repository

import (
	"context"

	"github.com/zkropotkine/agent-trace/internal/model"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoTraceRepository struct {
	collection *mongo.Collection
}

func NewMongoTraceRepository(db *mongo.Database) TraceRepository {
	return &mongoTraceRepository{
		collection: db.Collection("traces"),
	}
}

func (r *mongoTraceRepository) InsertTrace(ctx context.Context, trace model.Trace) error {
	_, err := r.collection.InsertOne(ctx, trace)
	return err
}
