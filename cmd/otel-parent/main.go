package main

import (
	"context"
	"log"
	"os"
	"os/exec"

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
	ctx, span := otel.Tracer(instrumentationName).Start(context.Background(), "parent")
	defer span.End()

	log.Printf("parent: SpanID: %s", trace.SpanContextFromContext(ctx).SpanID())

	cmd := execCommandWithSpanContext(ctx, "otel-child")

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		panic(err)
	}
}

func execCommandWithSpanContext(ctx context.Context, name string, args ...string) *exec.Cmd {
	cmd := exec.CommandContext(ctx, "otel-child")
	cmd.Env = os.Environ()

	envCarrier := propagation.MapCarrier{}
	otel.GetTextMapPropagator().Inject(ctx, envCarrier)
	for _, key := range envCarrier.Keys() {
		cmd.Env = append(cmd.Env, key+"="+envCarrier.Get(key))
	}
	return cmd
}
