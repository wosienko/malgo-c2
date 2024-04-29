package observability

import (
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
)

const TracingServiceName = "c2-gateway"

func ConfigureTraceProvider() *tracesdk.TracerProvider {
	jaegerEndpoint := os.Getenv("JAEGER_ENDPOINT")
	if jaegerEndpoint == "" {
		panic("JAEGER_ENDPOINT env variable is required")
	}

	exp, err := jaeger.New(
		jaeger.WithCollectorEndpoint(
			jaeger.WithEndpoint(jaegerEndpoint),
		),
	)
	if err != nil {
		panic(err)
	}

	tp := tracesdk.NewTracerProvider(
		// WARNING: `tracesdk.WithSyncer` should be not used in production,
		// for prod you should use `tracesdk.WithBatcher`
		tracesdk.WithSyncer(exp), // TODO: add automatic change to `tracesdk.WithBatcher` in the future
		//tracesdk.WithBatcher(exp),
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(TracingServiceName),
		)),
	)

	otel.SetTracerProvider(tp)

	// don't forget about that! lack of that line will cause that trace will not be propagated via messages
	otel.SetTextMapPropagator(propagation.TraceContext{})

	return tp
}
