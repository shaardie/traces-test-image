package main

import (
	"context"
	"log"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
)

const name = "traces-test-image"

func init() {
	ctx := context.TODO()

	// Create otlp exporter
	exporter, err := otlptracegrpc.New(ctx)
	if err != nil {
		log.Fatalf("Failed to create exporter, %v", err)
	}

	// Create Resource
	rs, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceVersionKey.String(name),
			semconv.ServiceVersionKey.String("v0.1.0"),
		),
	)
	if err != nil {
		log.Fatalf("Failed to create resource, %v", err)
	}

	// Create tracer
	tracer := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(rs),
	)

	otel.SetTracerProvider(tracer)
}

func loop(ctx context.Context) {
	_, span := otel.Tracer(name).Start(ctx, "Poll")
	defer span.End()
	log.Printf("Loop, traceID=%v", span.SpanContext().TraceID())
}

func main() {
	for {
		loop(context.Background())
		time.Sleep(10 * time.Second)
	}

}
