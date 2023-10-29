package main

import (
	"context"
	"log"
	"os"

	"github.com/yamamoto-febc/otel-env-context/otelsetup"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

const instrumentationName = "github.com/yamamoto-febc/otel-env-context"

func main() {
	shutdown, err := otelsetup.Init(context.Background(), "otel-env-context", "0.0.1")
	if err != nil {
		panic(err)
	}
	defer shutdown(context.Background())

	doSomething()
}

func doSomething() {
	envCarrier := propagation.MapCarrier{
		"traceparent": os.Getenv("traceparent"),
		"tracestate":  os.Getenv("tracestate"),
	}
	parentCtx := otel.GetTextMapPropagator().Extract(context.Background(), envCarrier)
	ctx, span := otel.Tracer(instrumentationName).Start(parentCtx, "child")
	defer span.End()

	log.Printf("child: SpanID: %s", trace.SpanContextFromContext(ctx).SpanID())
}
