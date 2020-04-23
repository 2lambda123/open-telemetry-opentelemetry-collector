// Copyright  OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package internal

import (
	"github.com/open-telemetry/opentelemetry-collector/consumer/pdata"
	"github.com/open-telemetry/opentelemetry-collector/internal/data"
)

// Initializes a metric with a metric slice and returns it.
func InitializeMetricSlice(metricData data.MetricData) pdata.MetricSlice {
	rms := metricData.ResourceMetrics()
	rms.Resize(1)
	rm := rms.At(0)
	ilms := rm.InstrumentationLibraryMetrics()
	ilms.Resize(1)
	ilm := ilms.At(0)
	return ilm.Metrics()
}

// AddNewMetric appends an empty metric to the metric slice, resizing
// the slice by 1, and returns the new metric.
func AddNewMetric(metrics pdata.MetricSlice) pdata.Metric {
	len := metrics.Len()
	metrics.Resize(len + 1)
	return metrics.At(len)
}
