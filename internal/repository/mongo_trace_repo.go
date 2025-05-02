package repository

import (
	"context"

	"github.com/zkropotkine/agent-trace/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoTraceRepository struct {
	collection *mongo.Collection
}

func NewMongoTraceRepository(collection *mongo.Collection) TraceRepository {
	return &mongoTraceRepository{
		collection: collection,
	}
}

func (r *mongoTraceRepository) InsertTrace(ctx context.Context, trace model.Trace) error {
	_, err := r.collection.InsertOne(ctx, trace)
	return err
}

func (r *mongoTraceRepository) GetTraces(ctx context.Context, filter TraceFilter) ([]model.Trace, error) {
	mongoFilter := bson.M{}
	if filter.AgentName != "" {
		mongoFilter["agent_name"] = filter.AgentName
	}
	if filter.From != nil || filter.To != nil {
		timeRange := bson.M{}
		if filter.From != nil {
			timeRange["$gte"] = *filter.From
		}
		if filter.To != nil {
			timeRange["$lte"] = *filter.To
		}
		mongoFilter["timestamp"] = timeRange
	}

	opts := options.Find().SetLimit(filter.Limit).SetSkip(filter.Offset).SetSort(bson.D{{Key: "timestamp", Value: -1}})
	cursor, err := r.collection.Find(ctx, mongoFilter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []model.Trace
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}

func (r *mongoTraceRepository) GetByID(ctx context.Context, id string) (*model.Trace, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var trace model.Trace
	err = r.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&trace)
	if err != nil {
		return nil, err
	}

	return &trace, nil
}
