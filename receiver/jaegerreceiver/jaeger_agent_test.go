// Copyright 2019, OpenTelemetry Authors
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

package jaegerreceiver

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"contrib.go.opencensus.io/exporter/jaeger"
	commonpb "github.com/census-instrumentation/opencensus-proto/gen-go/agent/common/v1"
	tracepb "github.com/census-instrumentation/opencensus-proto/gen-go/trace/v1"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"go.opencensus.io/trace"
	"go.uber.org/zap"

	"github.com/open-telemetry/opentelemetry-collector/consumer/consumerdata"
	"github.com/open-telemetry/opentelemetry-collector/exporter/exportertest"
	"github.com/open-telemetry/opentelemetry-collector/internal"
	"github.com/open-telemetry/opentelemetry-collector/internal/testutils"
	"github.com/open-telemetry/opentelemetry-collector/receiver/receivertest"
)

func TestJaegerAgentUDP_ThriftCompact_6831(t *testing.T) {
	port := 6831
	addrForClient := fmt.Sprintf(":%d", port)
	testJaegerAgent(t, addrForClient, &Configuration{
		AgentCompactThriftPort: port,
	})
}

func TestJaegerAgentUDP_ThriftBinary_6832(t *testing.T) {
	t.Skipf("Unfortunately due to Jaeger internal versioning, OpenCensus-Go's Thrift seems to conflict with ours")

	port := 6832
	addrForClient := fmt.Sprintf(":%d", port)
	testJaegerAgent(t, addrForClient, &Configuration{
		AgentBinaryThriftPort: port,
	})
}

func TestJaegerHTTP(t *testing.T) {
	port := int(testutils.GetAvailablePort(t))
	config := &Configuration{
		AgentHTTPPort: port,
	}
	jr, err := New(context.Background(), config, nil, zap.NewNop())
	if err != nil {
		t.Fatalf("Failed to create new Jaeger Receiver: %v", err)
	}
	defer jr.StopTraceReception()

	mh := receivertest.NewMockHost()
	if err := jr.StartTraceReception(mh); err != nil {
		t.Fatalf("StartTraceReception failed: %v", err)
	}

	// allow http server to start
	<-time.After(100 * time.Millisecond)

	// this functionality is just stubbed out at the moment.  just confirm they 200.
	testURL := fmt.Sprintf("http://localhost:%d/sampling?service=test", port)
	resp, err := http.Get(testURL)
	assert.NoError(t, err, "should not have failed to make request")
	if resp != nil {
		assert.Equal(t, 200, resp.StatusCode, "should have returned 200")
	}

	testURL = fmt.Sprintf("http://localhost:%d/sampling?service=test", port)
	resp, err = http.Get(testURL)
	assert.NoError(t, err, "should not have failed to make request")
	if resp != nil {
		assert.Equal(t, 200, resp.StatusCode, "should have returned 200")
	}
}

func testJaegerAgent(t *testing.T, agentEndpoint string, receiverConfig *Configuration) {
	// 1. Create the Jaeger receiver aka "server"
	sink := new(exportertest.SinkTraceExporter)
	jr, err := New(context.Background(), receiverConfig, sink, zap.NewNop())
	if err != nil {
		t.Fatalf("Failed to create new Jaeger Receiver: %v", err)
	}
	defer jr.StopTraceReception()

	mh := receivertest.NewMockHost()
	if err := jr.StartTraceReception(mh); err != nil {
		t.Fatalf("StartTraceReception failed: %v", err)
	}

	now := time.Unix(1542158650, 536343000).UTC()
	nowPlus10min := now.Add(10 * time.Minute)
	nowPlus10min2sec := now.Add(10 * time.Minute).Add(2 * time.Second)

	// 2. Then with a "live application", send spans to the Jaeger exporter.
	jexp, err := jaeger.NewExporter(jaeger.Options{
		AgentEndpoint: agentEndpoint,
		ServiceName:   "TestingAgentUDP",
		Process: jaeger.Process{
			ServiceName: "issaTest",
			Tags: []jaeger.Tag{
				jaeger.BoolTag("bool", true),
				jaeger.StringTag("string", "yes"),
				jaeger.Int64Tag("int64", 1e7),
			},
		},
	})
	if err != nil {
		t.Fatalf("Failed to create the Jaeger OpenCensus exporter for the live application: %v", err)
	}

	// 3. Now finally send some spans
	spandata := []*trace.SpanData{
		{
			SpanContext: trace.SpanContext{
				TraceID: trace.TraceID{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x80},
				SpanID:  trace.SpanID{0xAF, 0xAE, 0xAD, 0xAC, 0xAB, 0xAA, 0xA9, 0xA8},
			},
			ParentSpanID: trace.SpanID{0x1F, 0x1E, 0x1D, 0x1C, 0x1B, 0x1A, 0x19, 0x18},
			Name:         "DBSearch",
			StartTime:    now,
			EndTime:      nowPlus10min,
			Status: trace.Status{
				Code:    trace.StatusCodeNotFound,
				Message: "Stale indices",
			},
			Links: []trace.Link{
				{
					TraceID: trace.TraceID{0xF1, 0xF2, 0xF3, 0xF4, 0xF5, 0xF6, 0xF7, 0xF8, 0xF9, 0xFA, 0xFB, 0xFC, 0xFD, 0xFE, 0xFF, 0x80},
					SpanID:  trace.SpanID{0xCF, 0xCE, 0xCD, 0xCC, 0xCB, 0xCA, 0xC9, 0xC8},
					Type:    trace.LinkTypeParent,
				},
			},
		},
		{
			SpanContext: trace.SpanContext{
				TraceID: trace.TraceID{0xF1, 0xF2, 0xF3, 0xF4, 0xF5, 0xF6, 0xF7, 0xF8, 0xF9, 0xFA, 0xFB, 0xFC, 0xFD, 0xFE, 0xFF, 0x80},
				SpanID:  trace.SpanID{0xCF, 0xCE, 0xCD, 0xCC, 0xCB, 0xCA, 0xC9, 0xC8},
			},
			Name:      "ProxyFetch",
			StartTime: nowPlus10min,
			EndTime:   nowPlus10min2sec,
			Status: trace.Status{
				Code:    trace.StatusCodeInternal,
				Message: "Frontend crash",
			},
			Links: []trace.Link{
				{
					TraceID: trace.TraceID{0xF1, 0xF2, 0xF3, 0xF4, 0xF5, 0xF6, 0xF7, 0xF8, 0xF9, 0xFA, 0xFB, 0xFC, 0xFD, 0xFE, 0xFF, 0x80},
					SpanID:  trace.SpanID{0xAF, 0xAE, 0xAD, 0xAC, 0xAB, 0xAA, 0xA9, 0xA8},
					Type:    trace.LinkTypeChild,
				},
			},
		},
	}

	for _, sd := range spandata {
		jexp.ExportSpan(sd)
	}
	jexp.Flush()

	// Simulate and account for network latency but also the reception process on the server.
	<-time.After(500 * time.Millisecond)

	for i := 0; i < 10; i++ {
		jexp.Flush()
		<-time.After(60 * time.Millisecond)
	}

	got := sink.AllTraces()

	want := []consumerdata.TraceData{
		{
			Node: &commonpb.Node{
				ServiceInfo: &commonpb.ServiceInfo{Name: "issaTest"},
				LibraryInfo: &commonpb.LibraryInfo{},
				Identifier:  &commonpb.ProcessIdentifier{},
				Attributes: map[string]string{
					"bool":   "true",
					"string": "yes",
					"int64":  "10000000",
				},
			},

			Spans: []*tracepb.Span{
				{
					TraceId:      []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x80},
					SpanId:       []byte{0xAF, 0xAE, 0xAD, 0xAC, 0xAB, 0xAA, 0xA9, 0xA8},
					ParentSpanId: []byte{0x1F, 0x1E, 0x1D, 0x1C, 0x1B, 0x1A, 0x19, 0x18},
					Name:         &tracepb.TruncatableString{Value: "DBSearch"},
					StartTime:    internal.TimeToTimestamp(now),
					EndTime:      internal.TimeToTimestamp(nowPlus10min),
					Status: &tracepb.Status{
						Code:    trace.StatusCodeNotFound,
						Message: "Stale indices",
					},
					Attributes: &tracepb.Span_Attributes{
						AttributeMap: map[string]*tracepb.AttributeValue{
							"error": {
								Value: &tracepb.AttributeValue_BoolValue{BoolValue: true},
							},
						},
					},
					Links: &tracepb.Span_Links{
						Link: []*tracepb.Span_Link{
							{
								TraceId: []byte{0xF1, 0xF2, 0xF3, 0xF4, 0xF5, 0xF6, 0xF7, 0xF8, 0xF9, 0xFA, 0xFB, 0xFC, 0xFD, 0xFE, 0xFF, 0x80},
								SpanId:  []byte{0xCF, 0xCE, 0xCD, 0xCC, 0xCB, 0xCA, 0xC9, 0xC8},
								Type:    tracepb.Span_Link_PARENT_LINKED_SPAN,
							},
						},
					},
				},
				{
					TraceId:   []byte{0xF1, 0xF2, 0xF3, 0xF4, 0xF5, 0xF6, 0xF7, 0xF8, 0xF9, 0xFA, 0xFB, 0xFC, 0xFD, 0xFE, 0xFF, 0x80},
					SpanId:    []byte{0xCF, 0xCE, 0xCD, 0xCC, 0xCB, 0xCA, 0xC9, 0xC8},
					Name:      &tracepb.TruncatableString{Value: "ProxyFetch"},
					StartTime: internal.TimeToTimestamp(nowPlus10min),
					EndTime:   internal.TimeToTimestamp(nowPlus10min2sec),
					Status: &tracepb.Status{
						Code:    trace.StatusCodeInternal,
						Message: "Frontend crash",
					},
					Attributes: &tracepb.Span_Attributes{
						AttributeMap: map[string]*tracepb.AttributeValue{
							"error": {
								Value: &tracepb.AttributeValue_BoolValue{BoolValue: true},
							},
						},
					},
					Links: &tracepb.Span_Links{
						Link: []*tracepb.Span_Link{
							{
								TraceId: []byte{0xF1, 0xF2, 0xF3, 0xF4, 0xF5, 0xF6, 0xF7, 0xF8, 0xF9, 0xFA, 0xFB, 0xFC, 0xFD, 0xFE, 0xFF, 0x80},
								SpanId:  []byte{0xAF, 0xAE, 0xAD, 0xAC, 0xAB, 0xAA, 0xA9, 0xA8},
								// TODO: (@pjanotti, @odeke-em) contact the Jaeger maintains to inquire about
								// Parent_Linked_Spans as currently they've only got:
								// * Child_of
								// * Follows_from
								// yet OpenCensus has Parent too but Jaeger uses a zero-value for LinkCHILD.
								Type: tracepb.Span_Link_PARENT_LINKED_SPAN,
							},
						},
					},
				},
			},
			SourceFormat: "jaeger",
		},
	}

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Mismatched responses\n-Got +Want:\n\t%s", diff)
	}
}
