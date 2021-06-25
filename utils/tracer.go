package utils

import (
	"context"
	"runtime"

	"go.opentelemetry.io/otel/trace"

	"go.opentelemetry.io/otel"
)

var tracer = otel.GetTracerProvider().Tracer("")

func StartSpan(ctx context.Context) (context.Context, trace.Span) {
	pc, _, _, _ := runtime.Caller(1)
	details := runtime.FuncForPC(pc)
	return tracer.Start(ctx, details.Name())
}
