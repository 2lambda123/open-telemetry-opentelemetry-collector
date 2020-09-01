// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package internal

import (
	"time"

	"go.opentelemetry.io/collector/consumer/pdata"
	"go.opentelemetry.io/collector/internal/data"
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

func TimeToUnixNano(t time.Time) pdata.TimestampUnixNano {
	return pdata.TimestampUnixNano(uint64(t.UnixNano()))
}
