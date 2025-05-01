package repository

import (
	"context"

	"github.com/zkropotkine/agent-trace/internal/model"
)

type TraceRepository interface {
	InsertTrace(ctx context.Context, trace model.Trace) error
}
