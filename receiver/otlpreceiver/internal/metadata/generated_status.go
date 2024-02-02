// Code generated by mdatagen. DO NOT EDIT.

package metadata

import (
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"

	"go.opentelemetry.io/collector/component"
)

var (
	Type    = component.MustNewType("otlp")
	nameSep = "/"
)

const (
	LogsStability    = component.StabilityLevelBeta
	TracesStability  = component.StabilityLevelStable
	MetricsStability = component.StabilityLevelStable
)

func Meter(settings component.TelemetrySettings) metric.Meter {
	return settings.MeterProvider.Meter("otelcol/otlpreceiver")
}

func Tracer(settings component.TelemetrySettings) trace.Tracer {
	return settings.TracerProvider.Tracer("otelcol/otlpreceiver")
}

func CustomMetricName(name string) string {
	return "receiver" + nameSep + "otlp" + nameSep + name
}
