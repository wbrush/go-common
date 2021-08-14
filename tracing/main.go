package tracing

import (
	"context"
	"github.com/wbrush/go-common/tracing/opentelemetry"
)

type (
	TracerBase interface {
		NewTracer(projectId string, syncAlwaysFlag bool) (t opentelemetry.Tracer)
		GetTracer() (tracer *opentelemetry.Tracer)
	}
	Tracer interface {
		GetTracerCtx() context.Context
		StartSpan(traceName, spanName string) (sp *opentelemetry.Span)
	}
	Spanner interface {
		GetSpanCtx() context.Context
		EndSpan()
	}
)
