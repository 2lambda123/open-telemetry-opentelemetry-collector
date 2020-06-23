// Copyright The OpenTelemetry Authors
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

package metrics

import (
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/collector/internal/data"
	"go.opentelemetry.io/collector/internal/goldendataset"
	"go.opentelemetry.io/collector/testbed/correctness"
	"go.opentelemetry.io/collector/testbed/testbed"
)

func TestMetricsGoldenData(t *testing.T) {
	tests, err := correctness.LoadPictOutputPipelineDefs("../testdata/generated_pict_pairs_metrics_pipeline.txt")
	require.NoError(t, err)
	for _, test := range tests {
		test.TestName = fmt.Sprintf("%s-%s", test.Receiver, test.Exporter)
		test.DataSender = correctness.ConstructMetricsSender(t, test.Receiver)
		test.DataReceiver = correctness.ConstructReceiver(t, test.Exporter)
		t.Run(test.TestName, func(t *testing.T) {
			testWithMetricsGoldenDataset(t, test.DataSender, test.DataReceiver)
		})
	}
}

func testWithMetricsGoldenDataset(t *testing.T, sender testbed.DataSender, receiver testbed.DataReceiver) {
	mds := getTestMetrics(t)
	accumulator := newDiffAccumulator()
	h := NewTestHarness(
		t,
		NewMetricSupplier(mds),
		NewMetricIndex(mds),
		NewDualSender(sender),
		accumulator,
	)
	tc := NewCorrectnessTestCase(t, sender, receiver, h)

	tc.StartTestbedReceiver()
	tc.StartCollector()
	tc.StartTestbedSender()

	tc.SendFirstMetric()
	tc.WaitForAllMetrics()

	tc.StopTestbedReceiver()
	tc.StopCollector()

	// if accumulator.foundDiffs {
	// 	t.Fail()
	// }
}

type diffAccumulator struct {
	foundDiffs bool
}

func newDiffAccumulator() *diffAccumulator {
	return &diffAccumulator{}
}

func (d *diffAccumulator) Accept(metricName string, diffs []*MetricDiff) {
	if len(diffs) > 0 {
		d.foundDiffs = true
	}
	log.Printf("Found diffs for [%v]\n%v", metricName, diffs)
}

func getTestMetrics(t *testing.T) []data.MetricData {
	const file = "../../../internal/goldendataset/testdata/generated_pict_pairs_metrics.txt"
	mds, err := goldendataset.GenerateMetricDatas(file)
	require.NoError(t, err)
	return mds
}

func _TestProm(t *testing.T) {
	testbedSender := testbed.NewPrometheusDataSender(testbed.DefaultHost, testbed.GetAvailablePort(t))
	// testbedSender := testbed.NewOTLPMetricDataSender(testbed.DefaultHost, testbed.GetAvailablePort(t))
	testbedReceiver := testbed.NewPrometheusDataReceiver(testbed.GetAvailablePort(t))
	// testbedReceiver := testbed.NewOTLPDataReceiver(testbed.GetAvailablePort(t))

	testData := getTestData(t)
	testHarness := NewTestHarness(
		t,
		NewMetricSupplier(testData),
		NewMetricIndex(testData),
		NewDualSender(testbedSender),
		newDiffAccumulator(),
	)

	tc := NewCorrectnessTestCase(t, testbedSender, testbedReceiver, testHarness)
	tc.StartCollector()
	tc.StartTestbedSender()
	tc.StartTestbedReceiver()
	tc.SendFirstMetric()
	tc.WaitForAllMetrics()
}

func getTestData(t *testing.T) []data.MetricData {
	const file = "../../../internal/goldendataset/testdata/generated_pict_pairs_metrics.txt"
	mds, err := goldendataset.GenerateMetricDatas(file)
	require.NoError(t, err)
	return mds
}
