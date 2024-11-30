package prom

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

func RegisterGlobal(
	ctx context.Context,
	serviceName,
	serviceVersion,
	namespace string,
) (shutdown func(context.Context) error, err error) {
	r, err := resource.Merge(resource.Default(),
		resource.NewWithAttributes(semconv.SchemaURL,
			semconv.ServiceName(serviceName),
			semconv.ServiceVersion(serviceVersion),
		))
	if err != nil {
		return nil, err
	}

	exporter, err := prometheus.New(
		prometheus.WithNamespace(namespace),
	)
	if err != nil {
		return nil, err
	}

	provider := metric.NewMeterProvider(
		metric.WithReader(exporter),
		metric.WithResource(r),
	)

	otel.SetMeterProvider(provider)

	return func(ctx context.Context) error {
		return provider.Shutdown(ctx)
	}, nil
}
