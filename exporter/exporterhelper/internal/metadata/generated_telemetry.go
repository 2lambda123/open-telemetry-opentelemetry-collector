// Code generated by mdatagen. DO NOT EDIT.

package metadata

import (
	"errors"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

func Meter(settings component.TelemetrySettings) metric.Meter {
	return settings.MeterProvider.Meter("go.opentelemetry.io/collector/exporter/exporterhelper")
}

func Tracer(settings component.TelemetrySettings) trace.Tracer {
	return settings.TracerProvider.Tracer("go.opentelemetry.io/collector/exporter/exporterhelper")
}

// TelemetryBuilder provides an interface for components to report telemetry
// as defined in metadata and user config.
type TelemetryBuilder struct {
	ExporterEnqueueFailedLogRecords   metric.Int64Counter
	ExporterEnqueueFailedMetricPoints metric.Int64Counter
	ExporterEnqueueFailedSpans        metric.Int64Counter
	ExporterSendFailedLogRecords      metric.Int64Counter
	ExporterSendFailedMetricPoints    metric.Int64Counter
	ExporterSendFailedSpans           metric.Int64Counter
	ExporterSentLogRecords            metric.Int64Counter
	ExporterSentMetricPoints          metric.Int64Counter
	ExporterSentSpans                 metric.Int64Counter
}

// telemetryBuilderOption applies changes to default builder.
type telemetryBuilderOption func(*TelemetryBuilder)

// NewTelemetryBuilder provides a struct with methods to update all internal telemetry
// for a component
func NewTelemetryBuilder(settings component.TelemetrySettings, options ...telemetryBuilderOption) (*TelemetryBuilder, error) {
	builder := TelemetryBuilder{}
	var err, errs error
	meter := Meter(settings)
	builder.ExporterEnqueueFailedLogRecords, err = meter.Int64Counter(
		"exporter_enqueue_failed_log_records",
		metric.WithDescription("Number of log records failed to be added to the sending queue."),
		metric.WithUnit("1"),
	)
	errs = errors.Join(errs, err)
	builder.ExporterEnqueueFailedMetricPoints, err = meter.Int64Counter(
		"exporter_enqueue_failed_metric_points",
		metric.WithDescription("Number of metric points failed to be added to the sending queue."),
		metric.WithUnit("1"),
	)
	errs = errors.Join(errs, err)
	builder.ExporterEnqueueFailedSpans, err = meter.Int64Counter(
		"exporter_enqueue_failed_spans",
		metric.WithDescription("Number of spans failed to be added to the sending queue."),
		metric.WithUnit("1"),
	)
	errs = errors.Join(errs, err)
	builder.ExporterSendFailedLogRecords, err = meter.Int64Counter(
		"exporter_send_failed_log_records",
		metric.WithDescription("Number of log records in failed attempts to send to destination."),
		metric.WithUnit("1"),
	)
	errs = errors.Join(errs, err)
	builder.ExporterSendFailedMetricPoints, err = meter.Int64Counter(
		"exporter_send_failed_metric_points",
		metric.WithDescription("Number of metric points in failed attempts to send to destination."),
		metric.WithUnit("1"),
	)
	errs = errors.Join(errs, err)
	builder.ExporterSendFailedSpans, err = meter.Int64Counter(
		"exporter_send_failed_spans",
		metric.WithDescription("Number of spans in failed attempts to send to destination."),
		metric.WithUnit("1"),
	)
	errs = errors.Join(errs, err)
	builder.ExporterSentLogRecords, err = meter.Int64Counter(
		"exporter_sent_log_records",
		metric.WithDescription("Number of log record successfully sent to destination."),
		metric.WithUnit("1"),
	)
	errs = errors.Join(errs, err)
	builder.ExporterSentMetricPoints, err = meter.Int64Counter(
		"exporter_sent_metric_points",
		metric.WithDescription("Number of metric points successfully sent to destination."),
		metric.WithUnit("1"),
	)
	errs = errors.Join(errs, err)
	builder.ExporterSentSpans, err = meter.Int64Counter(
		"exporter_sent_spans",
		metric.WithDescription("Number of spans successfully sent to destination."),
		metric.WithUnit("1"),
	)
	errs = errors.Join(errs, err)
	return &builder, errs
}
