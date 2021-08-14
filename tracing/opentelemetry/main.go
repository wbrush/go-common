package opentelemetry

/*
Example of http wrapping at:
    https://github.com/open-telemetry/opentelemetry-go-contrib/blob/master/instrumentation/net/http
*/

import (
    "context"
    "github.com/sirupsen/logrus"

    texporter "github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace"
    "go.opentelemetry.io/otel/api/global"
    "go.opentelemetry.io/otel/api/trace"
    sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

var otTracer *Tracer

type Tracer struct {
    ProjectId   string
    SyncAlways  bool
    TracerCtx   context.Context
}

type Span struct {
    TraceName string
    SpanName  string
    Span      trace.Span
    SpanCtx   context.Context
}

func NewTracer(projectId string, syncAlwaysFlag bool) (t *Tracer) {
    if otTracer != nil {
        return otTracer
    }

    // Create exporter.
    exporter, err := texporter.NewExporter(texporter.WithProjectID(projectId))
    if err != nil {
        logrus.Errorf("texporter.NewExporter: %v", err)
        return nil
    }

    // Create trace provider with the exporter.
    //
    // By default it uses AlwaysSample() which samples all traces.
    // In a production environment or high QPS setup please use
    // ProbabilitySampler set at the desired probability.
    var tp *sdktrace.Provider
    probability := 1.0
    if !syncAlwaysFlag {
        probability = 0.5
    }
    config := sdktrace.Config{DefaultSampler:sdktrace.ProbabilitySampler(probability)}
    tp, err = sdktrace.NewProvider(sdktrace.WithSyncer(exporter),sdktrace.WithConfig(config))
    if err != nil {
        logrus.Errorf("Error creating probability sampler tracer: %s",err.Error())
        return nil
    }
    global.SetTraceProvider(tp)

    otTracer = &Tracer{
        ProjectId:  projectId,
        SyncAlways: syncAlwaysFlag,
        TracerCtx:  context.Background(),
    }

    logrus.Infof("Tracing enabled with %f probability sampler", probability)
    return otTracer
}

func GetTracer() (tracer *Tracer) {
    return otTracer
}

func (tracer *Tracer) GetTracerCtx() (context.Context) {
    return otTracer.TracerCtx
}

/*
   // Create custom span.
   tracer := global.TraceProvider().Tracer("example.com/trace")
   err = func(ctx context.Context) error {
       ctx, span := tracer.Start(ctx, "foo")
       defer span.End()

       // Do some work.

       return nil
   }(ctx)
*/

func (tracer *Tracer) StartSpan(ctx context.Context, traceName, spanName string) (sp *Span) {

    otTracer := global.TraceProvider().Tracer(traceName)
    ctx, span := otTracer.Start(ctx, spanName)

    sp = &Span{
        TraceName:traceName,
        SpanName: spanName,
        Span:    span,
        SpanCtx: ctx,
    }

    return
}

func (span *Span) GetSpanCtx() (context.Context){
    return span.SpanCtx
}

func (span *Span) EndSpan() {
    (*span).Span.End()
}
