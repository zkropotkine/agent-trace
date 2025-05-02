package repository

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/zkropotkine/agent-trace/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

func TestMongoTraceRepository_InsertTrace(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("successful insert", func(mt *mtest.T) {
		insertResponse := bson.D{
			{Key: "ok", Value: 1},
			{Key: "n", Value: 1},
			{Key: "acknowledged", Value: true},
			{Key: "insertedId", Value: primitive.NewObjectID()},
		}
		mt.AddMockResponses(mtest.CreateCursorResponse(1, "foo.bar", mtest.FirstBatch, insertResponse))

		trace := model.Trace{
			TraceID:     "test-trace-123",
			SessionID:   "test-session-123",
			AgentName:   "TestAgent",
			Status:      "success",
			InputPrompt: "Test prompt",
			Output:      "Test output",
			LatencyMS:   100,
			TokenUsage:  model.TokenUsage{Input: 10, Output: 20, Total: 30},
			CreatedAt:   time.Now(),
			SubSteps:    []model.SubStep{},
		}

		repo := NewMongoTraceRepository(mt.Coll)

		err := repo.InsertTrace(context.Background(), trace)

		assert.NoError(t, err, "Expected no error, but got: %v", err)
	})

	mt.Run("insert error", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateCommandErrorResponse(mtest.CommandError{
			Code:    11000, // Duplicate key error
			Message: "duplicate key error",
			Name:    "WriteError",
		}))

		// Sample trace to insert
		trace := model.Trace{
			TraceID:   "error-trace",
			SessionID: "error-session",
			AgentName: "ErrorAgent",
		}

		repo := NewMongoTraceRepository(mt.Coll)

		err := repo.InsertTrace(context.Background(), trace)

		assert.Error(t, err, "Expected an error but got none")
		assert.Contains(t, err.Error(), "duplicate key error")
	})
}
