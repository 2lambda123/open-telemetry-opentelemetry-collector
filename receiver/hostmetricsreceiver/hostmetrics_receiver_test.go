// Copyright 2020, OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package hostmetricsreceiver

import (
	"context"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/open-telemetry/opentelemetry-collector/component/componenttest"
	"github.com/open-telemetry/opentelemetry-collector/consumer/pdata"
	"github.com/open-telemetry/opentelemetry-collector/exporter/exportertest"
	"github.com/open-telemetry/opentelemetry-collector/receiver/hostmetricsreceiver/internal"
	"github.com/open-telemetry/opentelemetry-collector/receiver/hostmetricsreceiver/internal/scraper/cpuscraper"
)

func TestGatherMetrics_EndToEnd(t *testing.T) {
	sink := &exportertest.SinkMetricsExporter{}

	config := &Config{
		Scrapers: map[string]internal.Config{
			cpuscraper.TypeStr: &cpuscraper.Config{
				ConfigSettings: internal.ConfigSettings{CollectionIntervalValue: 100 * time.Millisecond},
				ReportPerCPU:   true,
			},
		},
	}

	factories := map[string]internal.Factory{
		cpuscraper.TypeStr: &cpuscraper.Factory{},
	}

	receiver, err := NewHostMetricsReceiver(context.Background(), zap.NewNop(), config, factories, sink)

	if runtime.GOOS != "windows" {
		require.Error(t, err, "Expected error when creating a host metrics receiver with cpuscraper collector on a non-windows environment")
		return
	}

	require.NoError(t, err, "Failed to create metrics receiver: %v", err)

	err = receiver.Start(context.Background(), componenttest.NewNopHost())
	require.NoError(t, err, "Failed to start metrics receiver: %v", err)
	defer func() { assert.NoError(t, receiver.Shutdown(context.Background())) }()

	require.Eventually(t, func() bool {
		got := sink.AllMetrics()
		if len(got) == 0 {
			return false
		}

		assertMetricData(t, got)
		return true
	}, time.Second, 10*time.Millisecond, "No metrics were collected")
}

func assertMetricData(t *testing.T, got []pdata.Metrics) {
	metrics := internal.AssertSingleMetricDataAndGetMetricsSlice(t, got)

	// expect 2 metrics
	assert.Equal(t, 2, metrics.Len())

	// for cpu seconds metric, expect 5 timeseries with appropriate labels
	hostCPUTimeMetric := metrics.At(0)
	internal.AssertDescriptorEqual(t, cpuscraper.MetricCPUSecondsDescriptor, hostCPUTimeMetric.MetricDescriptor())
	assert.Equal(t, 4*runtime.NumCPU(), hostCPUTimeMetric.Int64DataPoints().Len())
	internal.AssertInt64MetricLabelExists(t, hostCPUTimeMetric, 0, cpuscraper.CPULabel)
	internal.AssertInt64MetricLabelHasValue(t, hostCPUTimeMetric, 0, cpuscraper.StateLabel, cpuscraper.UserStateLabelValue)
	internal.AssertInt64MetricLabelHasValue(t, hostCPUTimeMetric, 1, cpuscraper.StateLabel, cpuscraper.SystemStateLabelValue)
	internal.AssertInt64MetricLabelHasValue(t, hostCPUTimeMetric, 2, cpuscraper.StateLabel, cpuscraper.IdleStateLabelValue)
	internal.AssertInt64MetricLabelHasValue(t, hostCPUTimeMetric, 3, cpuscraper.StateLabel, cpuscraper.InterruptStateLabelValue)

	// for cpu utilization metric, expect #cores timeseries each with a value < 100
	hostCPUUtilizationMetric := metrics.At(1)
	internal.AssertDescriptorEqual(t, cpuscraper.MetricCPUUtilizationDescriptor, hostCPUUtilizationMetric.MetricDescriptor())
	ddp := hostCPUUtilizationMetric.DoubleDataPoints()
	assert.Equal(t, runtime.NumCPU(), ddp.Len())
	for i := 0; i < ddp.Len(); i++ {
		assert.LessOrEqual(t, ddp.At(i).Value(), float64(100))
	}
}
