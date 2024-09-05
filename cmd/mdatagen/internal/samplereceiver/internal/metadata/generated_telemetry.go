// Code generated by mdatagen. DO NOT EDIT.

package metadata

import (
	"context"
	"errors"

	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config/configtelemetry"
)

// Deprecated: [v0.108.0] use LeveledMeter instead.
func Meter(settings component.TelemetrySettings) metric.Meter {
	return settings.MeterProvider.Meter("go.opentelemetry.io/collector/internal/receiver/samplereceiver")
}

func LeveledMeter(settings component.TelemetrySettings, level configtelemetry.Level) metric.Meter {
	return settings.LeveledMeterProvider(level).Meter("go.opentelemetry.io/collector/internal/receiver/samplereceiver")
}

func Tracer(settings component.TelemetrySettings) trace.Tracer {
	return settings.TracerProvider.Tracer("go.opentelemetry.io/collector/internal/receiver/samplereceiver")
}

// TelemetryBuilder provides an interface for components to report telemetry
// as defined in metadata and user config.
type TelemetryBuilder struct {
	meter                                metric.Meter
	BatchSizeTriggerSend                 metric.Int64Counter
	DisabledMetric                       metric.Int64Gauge
	ProcessRuntimeTotalAllocBytes        metric.Int64ObservableCounter
	observeProcessRuntimeTotalAllocBytes func(context.Context, metric.Observer) error
	QueueLength                          metric.Int64ObservableGauge
	RequestDuration                      metric.Float64Histogram
	RequestDurationNoAttrs               metric.Float64Histogram
	meters                               map[configtelemetry.Level]metric.Meter
}

// TelemetryBuilderOption applies changes to default builder.
type TelemetryBuilderOption interface {
	apply(*TelemetryBuilder)
}

type telemetryBuilderOptionFunc func(mb *TelemetryBuilder)

func (tbof telemetryBuilderOptionFunc) apply(mb *TelemetryBuilder) {
	tbof(mb)
}

// WithProcessRuntimeTotalAllocBytesCallback sets callback for observable ProcessRuntimeTotalAllocBytes metric.
func WithProcessRuntimeTotalAllocBytesCallback(cb func() int64, opts ...metric.ObserveOption) TelemetryBuilderOption {
	return telemetryBuilderOptionFunc(func(builder *TelemetryBuilder) {
		builder.observeProcessRuntimeTotalAllocBytes = func(_ context.Context, o metric.Observer) error {
			o.ObserveInt64(builder.ProcessRuntimeTotalAllocBytes, cb(), opts...)
			return nil
		}
	})
}

// InitQueueLength configures the QueueLength metric.
func (builder *TelemetryBuilder) InitQueueLength(cb func() int64, opts ...metric.ObserveOption) error {
	var err error
	builder.QueueLength, err = builder.meters[configtelemetry.LevelBasic].Int64ObservableGauge(
		"otelcol_queue_length",
		metric.WithDescription("This metric is optional and therefore not initialized in NewTelemetryBuilder."),
		metric.WithUnit("{items}"),
	)
	if err != nil {
		return err
	}
	_, err = builder.meters[configtelemetry.LevelBasic].RegisterCallback(func(_ context.Context, o metric.Observer) error {
		o.ObserveInt64(builder.QueueLength, cb(), opts...)
		return nil
	}, builder.QueueLength)
	return err
}

// NewTelemetryBuilder provides a struct with methods to update all internal telemetry
// for a component
func NewTelemetryBuilder(settings component.TelemetrySettings, options ...TelemetryBuilderOption) (*TelemetryBuilder, error) {
	builder := TelemetryBuilder{meters: map[configtelemetry.Level]metric.Meter{}}
	for _, op := range options {
		op.apply(&builder)
	}
	builder.meters[configtelemetry.LevelBasic] = LeveledMeter(settings, configtelemetry.LevelBasic)
	var err, errs error
	builder.BatchSizeTriggerSend, err = builder.meters[configtelemetry.LevelBasic].Int64Counter(
		"otelcol_batch_size_trigger_send",
		metric.WithDescription("Number of times the batch was sent due to a size trigger [deprecated since v0.110.0]"),
		metric.WithUnit("{times}"),
	)
	errs = errors.Join(errs, err)
	builder.DisabledMetric, err = builder.meters[configtelemetry.LevelBasic].Int64Gauge(
		"otelcol_disabled_metric",
		metric.WithDescription("This metric is disabled by default"),
		metric.WithUnit("{items}"),
	)
	errs = errors.Join(errs, err)
	builder.ProcessRuntimeTotalAllocBytes, err = builder.meters[configtelemetry.LevelBasic].Int64ObservableCounter(
		"otelcol_process_runtime_total_alloc_bytes",
		metric.WithDescription("Cumulative bytes allocated for heap objects (see 'go doc runtime.MemStats.TotalAlloc')"),
		metric.WithUnit("By"),
	)
	errs = errors.Join(errs, err)
	_, err = builder.meters[configtelemetry.LevelBasic].RegisterCallback(builder.observeProcessRuntimeTotalAllocBytes, builder.ProcessRuntimeTotalAllocBytes)
	errs = errors.Join(errs, err)
	builder.RequestDuration, err = builder.meters[configtelemetry.LevelBasic].Float64Histogram(
		"otelcol_request_duration",
		metric.WithDescription("Duration of request [alpha]"),
		metric.WithUnit("s"),
		metric.WithExplicitBucketBoundaries([]float64{1, 10, 100}...),
	)
	errs = errors.Join(errs, err)
	builder.RequestDurationNoAttrs, err = builder.meters[configtelemetry.LevelBasic].Float64Histogram(
		"otelcol_request_duration_no_attrs",
		metric.WithDescription("Duration of request. Test mdatagen without any attributes"),
		metric.WithUnit("s"), metric.WithExplicitBucketBoundaries([]float64{1, 10, 100}...),
	)
	errs = errors.Join(errs, err)
	return &builder, errs
}
