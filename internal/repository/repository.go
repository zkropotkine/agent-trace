package repository

import (
	"context"
	"time"

	"github.com/zkropotkine/agent-trace/internal/model"
)

// TraceFilter defines filtering and pagination options for querying traces.
type TraceFilter struct {
	AgentName string
	From      *time.Time
	To        *time.Time
	Limit     int64
	Offset    int64
}

type TraceRepository interface {
	InsertTrace(ctx context.Context, trace model.Trace) error
	GetTraces(ctx context.Context, filter TraceFilter) ([]model.Trace, error)
	GetByID(ctx context.Context, id string) (*model.Trace, error)
}
