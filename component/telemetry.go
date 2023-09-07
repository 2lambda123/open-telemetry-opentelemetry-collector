// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package component // import "go.opentelemetry.io/collector/component"

import (
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"

	"go.opentelemetry.io/collector/config/configtelemetry"
	"go.opentelemetry.io/collector/pdata/pcommon"
)

type TelemetrySettings struct {
	// Logger that the factory can use during creation and can pass to the created
	// component to be used later as well.
	Logger *zap.Logger

	// GetSampledLoggerFunction is a function to get the sampled Logger
	GetSampledLoggerFunction func() *zap.Logger

	// TracerProvider that the factory can pass to other instrumented third-party libraries.
	TracerProvider trace.TracerProvider

	// MeterProvider that the factory can pass to other instrumented third-party libraries.
	MeterProvider metric.MeterProvider

	// MetricsLevel controls the level of detail for metrics emitted by the collector.
	// Experimental: *NOTE* this field is experimental and may be changed or removed.
	MetricsLevel configtelemetry.Level

	// Resource contains the resource attributes for the collector's telemetry.
	Resource pcommon.Resource
}

// SampledLogger is built from the logger. It is passed to the created component.
// It will be used to avoid flooding the logs with messages that are repeated frequently.
// It will be built the first time used.
func (t *TelemetrySettings) SampledLogger() *zap.Logger {
	return t.GetSampledLoggerFunction()
}
