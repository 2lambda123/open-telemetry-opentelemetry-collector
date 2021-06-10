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

package internal_test

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/gogo/protobuf/proto"
	"github.com/golang/snappy"
	"github.com/prometheus/prometheus/pkg/value"
	"github.com/prometheus/prometheus/prompb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/exporter/prometheusremotewriteexporter"
	"go.opentelemetry.io/collector/processor/batchprocessor"
	"go.opentelemetry.io/collector/receiver/prometheusreceiver"
	"go.opentelemetry.io/collector/service"
	"go.opentelemetry.io/collector/service/parserprovider"
)

// Test that staleness markers are emitted for timeseries that intermittently disappear.
// This test runs the entire collector and end-to-end scrapes then checks with the
// Prometheus remotewrite exporter that staleness markers are emitted per timeseries.
// See https://github.com/open-telemetry/opentelemetry-collector/issues/3413
func TestStalenessMarkersEndToEnd(t *testing.T) {
	if testing.Short() {
		t.Skip("This test can take a long time")
	}

	ctx, cancel := context.WithCancel(context.Background())

	// 1. Setup the server that sends series that intermittently appear and disappear.
	var n uint64
	scrapeServer := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Increment the scrape count atomically per scrape.
		i := atomic.AddUint64(&n, 1)

		select {
		case <-ctx.Done():
			return
		default:
		}

		// Alternate metrics per scrape so that every one of
		// them will be reported as stale.
		if i%2 == 0 {
			fmt.Fprintf(rw, `
# HELP jvm_memory_bytes_used Used bytes of a given JVM memory area.
# TYPE jvm_memory_bytes_used gauge
jvm_memory_bytes_used{area="heap"} %.1f`, float64(i))
		} else {
			fmt.Fprintf(rw, `
# HELP jvm_memory_pool_bytes_used Used bytes of a given JVM memory pool.
# TYPE jvm_memory_pool_bytes_used gauge
jvm_memory_pool_bytes_used{pool="CodeHeap 'non-nmethods'"} %.1f`, float64(i))
		}
	}))
	defer scrapeServer.Close()

	serverURL, err := url.Parse(scrapeServer.URL)
	require.Nil(t, err)

	// 2. Set up the Prometheus RemoteWrite endpoint.
	prweUploads := make(chan *prompb.WriteRequest)
	prweServer := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Snappy decode the uploads.
		payload, rerr := ioutil.ReadAll(req.Body)
		if err != nil {
			panic(rerr)
		}
		recv := make([]byte, len(payload))
		decoded, derr := snappy.Decode(recv, payload)
		if err != nil {
			panic(derr)
		}

		writeReq := new(prompb.WriteRequest)
		if uerr := proto.Unmarshal(decoded, writeReq); uerr != nil {
			panic(uerr)
		}

		select {
		case <-ctx.Done():
			return
		case prweUploads <- writeReq:
		}
	}))
	defer prweServer.Close()

	// 3. Set the OpenTelemetry Prometheus receiver.
	config := fmt.Sprintf(`
receivers:
  prometheus:
    config:
      scrape_configs:
        - job_name: 'test'
          scrape_interval: 5ms
          static_configs:
            - targets: [%q]

processors:
  batch:
exporters:
  prometheusremotewrite:
    endpoint: %q
    insecure: true

service:
  pipelines:
    metrics:
      receivers: [prometheus]
      processors: [batch]
      exporters: [prometheusremotewrite]`, serverURL.Host, prweServer.URL)

	// 4. Run the OpenTelemetry Collector.
	receivers, err := component.MakeReceiverFactoryMap(prometheusreceiver.NewFactory())
	require.Nil(t, err)
	exporters, err := component.MakeExporterFactoryMap(prometheusremotewriteexporter.NewFactory())
	require.Nil(t, err)
	processors, err := component.MakeProcessorFactoryMap(batchprocessor.NewFactory())
	require.Nil(t, err)

	factories := component.Factories{
		Receivers:  receivers,
		Exporters:  exporters,
		Processors: processors,
	}

	appSettings := service.CollectorSettings{
		Factories:      factories,
		ParserProvider: parserprovider.NewInMemory(strings.NewReader(config)),
		BuildInfo: component.BuildInfo{
			Command:     "otelcol",
			Description: "OpenTelemetry Collector",
			Version:     "tests",
		},
		LoggingOptions: []zap.Option{
			// Turn off the verbose logging from the collector.
			zap.WrapCore(func(zapcore.Core) zapcore.Core {
				return zapcore.NewNopCore()
			}),
		},
	}
	app, err := service.New(appSettings)
	require.Nil(t, err)
	go func() {
		if err := app.Run(); err != nil {
			t.Error(err)
		}
	}()
	defer app.Shutdown()

	// 5. Let's wait on at least 8 fetches.
	var wReqL []*prompb.WriteRequest
	for i := 0; i < 10; i++ {
		wReqL = append(wReqL, <-prweUploads)
	}
	defer cancel()

	// Assert that we encounter the stale markers aka special NaNs for every time series.
	staleMarkerCount := 0
	totalSamples := 0
	for i, wReq := range wReqL {
		t.Run(fmt.Sprintf("WriteRequest#%d", i), func(t *testing.T) {
			require.True(t, len(wReq.Timeseries) > 0, "Expecting at least 1 timeSeries")
			for j, ts := range wReq.Timeseries {
				t.Run(fmt.Sprintf("TimeSeries#%d", j), func(t *testing.T) {
					assert.True(t, len(ts.Samples) > 0, "Expected at least 1 Sample")
					for _, sample := range ts.Samples {
						totalSamples++
						if value.IsStaleNaN(sample.Value) {
							staleMarkerCount++
						}
					}
				})
			}
		})
	}

	require.True(t, totalSamples > 0, "Expected at least 1 sample")
	// Expect at least 40% chance of stale markers being emitted.
	chance := float64(staleMarkerCount) / float64(totalSamples)
	require.True(t, chance > 0.4, fmt.Sprintf("Expected at least one stale marker: %.3f", chance))
}
