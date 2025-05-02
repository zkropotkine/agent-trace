package repository

import (
	"context"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"

	"github.com/stretchr/testify/assert"
	"github.com/zkropotkine/agent-trace/internal/model"
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

func TestMongoTraceRepository_GetByID(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("success", func(mt *mtest.T) {
		r := NewMongoTraceRepository(mt.Coll)

		id := "64b0c2f4e13c0000aa000000"
		objID, _ := primitive.ObjectIDFromHex(id)

		mt.AddMockResponses(mtest.CreateCursorResponse(1, "agentTrace.traces", mtest.FirstBatch, bson.D{
			{"_id", objID},
			{"traceId", "1"},
			{"agentName", "B"},
			{"timestamp", time.Now().UTC()},
		}))

		res, err := r.GetByID(context.Background(), id)
		assert.NoError(t, err)
		assert.Equal(t, "1", res.TraceID)
	})

	mt.Run("invalid id", func(mt *mtest.T) {
		r := NewMongoTraceRepository(mt.Coll)

		_, err := r.GetByID(context.Background(), "invalid")
		assert.Error(t, err)
		assert.EqualError(t, err, "the provided hex string is not a valid ObjectID")
	})

	mt.Run("not found", func(mt *mtest.T) {
		r := NewMongoTraceRepository(mt.Coll)

		id := "64b0c2f4e13c0000aa001234"
		mt.AddMockResponses(mtest.CreateCursorResponse(0, "agentTrace.traces", mtest.FirstBatch))

		_, err := r.GetByID(context.Background(), id)
		assert.Error(t, err)
	})
}

func TestMongoTraceRepository_GetTraces(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("success with filter", func(mt *mtest.T) {
		r := NewMongoTraceRepository(mt.Coll)

		ns := mt.Coll.Database().Name() + "." + mt.Coll.Name()

		mt.AddMockResponses(
			mtest.CreateCursorResponse(
				1,
				ns,
				mtest.FirstBatch,
				bson.D{
					{"traceId", "1"},
					{"agentName", "A"},
					{"timestamp", time.Now().UTC()},
				},
			),
			mtest.CreateCursorResponse(
				0,
				ns,
				mtest.NextBatch,
			),
		)

		res, err := r.GetTraces(context.Background(), TraceFilter{AgentName: "A", Limit: 10, Offset: 0})
		assert.NoError(t, err)
		assert.Len(t, res, 1)
		assert.Equal(t, "1", res[0].TraceID)
	})
}
