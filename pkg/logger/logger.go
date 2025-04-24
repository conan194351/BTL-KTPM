package logger

import (
	"context"
)

var (
	LOG_METADATA = "log_metadata"
)

type Logger interface {
	Info(msg string, ctx context.Context)
	Debug(msg string, ctx context.Context)
	Warn(msg string, ctx context.Context)
	Fatal(err error, msg string, ctx context.Context)
	Panic(err error, msg string, ctx context.Context)
	Error(err error, msg string, ctx context.Context)
	CtxWithMetadata(ctx context.Context, value interface{}) context.Context
}
