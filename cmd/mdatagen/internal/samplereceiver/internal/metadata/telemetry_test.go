// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package metadata

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/component/componenttest"
)

// TODO: The goal of this functionality is to have it be generated by mdatagen. This is meant to be an example
// of what the output of mdatagen should be.
// This method will be expanded when generated by mdatagen to run tests on each telemetry metric attribute for its given type.
func TestTelemetryMetrics(t *testing.T) {
	// Current status: Shows how to properly test gauge, sum, histogram, optional, and async metrics.
	// Also tests with and without attributes on a telemetry metric.
	// This is not a comprehensive test of the sample receiver's telemetry metrics, but should cover each of its cases.

	// Component telemetry setup should be uniform across components
	tt, err := componenttest.SetupTelemetry(component.MustNewID(Type.String()))
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, tt.Shutdown(context.Background())) })

	var tb *TelemetryBuilder
	tb, err = NewTelemetryBuilder(tt.TelemetrySettings(),
		// Need default callback function for each async metric
		WithProcessRuntimeTotalAllocBytesCallback(func() int64 {
			return 1
		}))
	require.NoError(t, err)

	// Initialize default values for all attributes used by telemetry metrics. Must be used or we'll get lint errors.
	stringAttr := attribute.String("string_attr", "value")
	stateAttr := attribute.Int64("state", 1)
	// Choose first value in array for enums
	enumAttr := attribute.String("enum_attr", "red")
	booleanAttr := attribute.String("boolean_attr", "false")

	batchSizeTriggerSendAttrs := attribute.NewSet(stringAttr, stateAttr, enumAttr, booleanAttr)

	// Ensure async metric is set upon setup.
	ts, err := tt.GetIntSumDataPoint("otelcol_process_runtime_total_alloc_bytes", attribute.NewSet())
	require.NoError(t, err)
	require.Equal(t, int64(1), ts.Value)

	tb.BatchSizeTriggerSend.Add(context.Background(), 1, metric.WithAttributeSet(batchSizeTriggerSendAttrs))
	ts, err = tt.GetIntSumDataPoint("otelcol_batch_size_trigger_send", batchSizeTriggerSendAttrs)
	require.NoError(t, err)
	require.Equal(t, int64(1), ts.Value)

	// Ensure different attributes means resetting expected value here.
	tb.BatchSizeTriggerSend.Add(context.Background(), 1, metric.WithAttributeSet(attribute.NewSet()))
	ts, err = tt.GetIntSumDataPoint("otelcol_batch_size_trigger_send", batchSizeTriggerSendAttrs)
	require.NoError(t, err)
	require.Equal(t, int64(1), ts.Value)

	// For histograms, use metadata buckets to set information here. Want to include values in different buckets,
	// and then check to make sure they're in the right buckets.
	tb.RequestDuration.Record(context.Background(), 1, metric.WithAttributeSet(attribute.NewSet(stateAttr)))
	tb.RequestDuration.Record(context.Background(), 15, metric.WithAttributeSet(attribute.NewSet(stateAttr)))
	th, err := tt.GetFloatHistogramDataPoint("otelcol_request_duration", attribute.NewSet(stateAttr))
	require.NoError(t, err)
	require.Equal(t, uint64(1), th.BucketCounts[0])
	require.Equal(t, uint64(0), th.BucketCounts[1])
	require.Equal(t, uint64(1), th.BucketCounts[2])
	require.Equal(t, uint64(0), th.BucketCounts[3])

	// For optional metrics, we should check to make sure they are not recorded before being initialized, even if
	// another metric has been recorded, which would trigger the callback.
	ts, err = tt.GetIntGaugeDataPoint("otelcol_queue_length", attribute.NewSet())
	require.Error(t, err)

	// Init optional metric and trigger record from another metric.
	// We can check the other metric to ensure it's correctly updated as well.
	tb.InitQueueLength(func() int64 { return 1 })
	tb.BatchSizeTriggerSend.Add(context.Background(), 1, metric.WithAttributeSet(batchSizeTriggerSendAttrs))
	ts, err = tt.GetIntGaugeDataPoint("otelcol_queue_length", attribute.NewSet())
	require.NoError(t, err)
	require.Equal(t, int64(1), ts.Value)
	ts, err = tt.GetIntSumDataPoint("otelcol_batch_size_trigger_send", batchSizeTriggerSendAttrs)
	require.NoError(t, err)
	require.Equal(t, int64(2), ts.Value)
}
