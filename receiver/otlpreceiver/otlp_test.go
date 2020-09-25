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

package otlpreceiver

import (
	"bytes"
	"compress/gzip"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/gogo/protobuf/jsonpb"
	json "github.com/json-iterator/go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	spb "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/component/componenttest"
	"go.opentelemetry.io/collector/config/configgrpc"
	"go.opentelemetry.io/collector/config/confighttp"
	"go.opentelemetry.io/collector/config/configmodels"
	"go.opentelemetry.io/collector/config/confignet"
	"go.opentelemetry.io/collector/config/configtls"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/consumer/pdata"
	"go.opentelemetry.io/collector/exporter/exportertest"
	collectortrace "go.opentelemetry.io/collector/internal/data/opentelemetry-proto-gen/collector/trace/v1"
	otlpcommon "go.opentelemetry.io/collector/internal/data/opentelemetry-proto-gen/common/v1"
	otlpresource "go.opentelemetry.io/collector/internal/data/opentelemetry-proto-gen/resource/v1"
	otlptrace "go.opentelemetry.io/collector/internal/data/opentelemetry-proto-gen/trace/v1"
	"go.opentelemetry.io/collector/internal/data/testdata"
	"go.opentelemetry.io/collector/obsreport/obsreporttest"
	"go.opentelemetry.io/collector/testutil"
	"go.opentelemetry.io/collector/translator/conventions"
)

const otlpReceiverName = "otlp_receiver_test"

var traceJSON = []byte(`
	{
	  "resource_spans": [
		{
		  "resource": {
			"attributes": [
			  {
				"key": "host.hostname",
				"value": { "stringValue": "testHost" }
			  }
			]
		  },
		  "instrumentation_library_spans": [
			{
			  "spans": [
				{
				  "trace_id": "5B8EFFF798038103D269B633813FC60C",
				  "span_id": "7uGbfsPBsXM=",
				  "name": "testSpan",
				  "start_time_unix_nano": 1544712660000000000,
				  "end_time_unix_nano": 1544712661000000000,
				  "attributes": [
					{
					  "key": "attr1",
					  "value": { "intValue": 55 }
					}
				  ]
				}
			  ]
			}
		  ]
		}
	  ]
	}`)

var resourceSpansOtlp = otlptrace.ResourceSpans{

	Resource: &otlpresource.Resource{
		Attributes: []*otlpcommon.KeyValue{
			{
				Key:   conventions.AttributeHostHostname,
				Value: &otlpcommon.AnyValue{Value: &otlpcommon.AnyValue_StringValue{StringValue: "testHost"}},
			},
		},
	},
	InstrumentationLibrarySpans: []*otlptrace.InstrumentationLibrarySpans{
		{
			Spans: []*otlptrace.Span{
				{
					TraceId:           otlpcommon.NewTraceID([]byte{0x5B, 0x8E, 0xFF, 0xF7, 0x98, 0x3, 0x81, 0x3, 0xD2, 0x69, 0xB6, 0x33, 0x81, 0x3F, 0xC6, 0xC}),
					SpanId:            []byte{0xEE, 0xE1, 0x9B, 0x7E, 0xC3, 0xC1, 0xB1, 0x73},
					Name:              "testSpan",
					StartTimeUnixNano: 1544712660000000000,
					EndTimeUnixNano:   1544712661000000000,
					Attributes: []*otlpcommon.KeyValue{
						{
							Key:   "attr1",
							Value: &otlpcommon.AnyValue{Value: &otlpcommon.AnyValue_IntValue{IntValue: 55}},
						},
					},
				},
			},
		},
	},
}

var traceOtlp = pdata.TracesFromOtlp([]*otlptrace.ResourceSpans{&resourceSpansOtlp})

func TestJsonHttp(t *testing.T) {
	tests := []struct {
		name     string
		encoding string
		err      error
	}{
		{
			name:     "JSONUncompressed",
			encoding: "",
		},
		{
			name:     "JSONGzipCompressed",
			encoding: "gzip",
		},
		{
			name:     "NotGRPCError",
			encoding: "",
			err:      errors.New("my error"),
		},
		{
			name:     "GRPCError",
			encoding: "",
			err:      status.New(codes.Internal, "").Err(),
		},
	}
	addr := testutil.GetAvailableLocalAddress(t)

	// Set the buffer count to 1 to make it flush the test span immediately.
	sink := new(exportertest.SinkTraceExporter)
	ocr := newHTTPReceiver(t, addr, sink, nil)

	require.NoError(t, ocr.Start(context.Background(), componenttest.NewNopHost()), "Failed to start trace receiver")
	defer ocr.Shutdown(context.Background())

	// TODO(nilebox): make starting server deterministic
	// Wait for the servers to start
	<-time.After(10 * time.Millisecond)

	url := fmt.Sprintf("http://%s/v1/trace", addr)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var buf *bytes.Buffer
			var err error
			switch test.encoding {
			case "gzip":
				buf, err = compressGzip(traceJSON)
				require.NoError(t, err, "Error while gzip compressing trace: %v", err)
			default:
				buf = bytes.NewBuffer(traceJSON)
			}
			sink.SetConsumeTraceError(test.err)
			req, err := http.NewRequest("POST", url, buf)
			require.NoError(t, err, "Error creating trace POST request: %v", err)
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Content-Encoding", test.encoding)

			client := &http.Client{}
			resp, err := client.Do(req)
			require.NoError(t, err, "Error posting trace to grpc-gateway server: %v", err)

			respBytes, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("Error reading response from trace grpc-gateway, %v", err)
			}
			respStr := string(respBytes)
			err = resp.Body.Close()
			if err != nil {
				t.Errorf("Error closing response body, %v", err)
			}

			if test.err == nil {
				assert.Equal(t, 200, resp.StatusCode)
				var respJSON map[string]interface{}
				assert.NoError(t, json.Unmarshal([]byte(respStr), &respJSON))
				assert.Len(t, respJSON, 0, "Got unexpected response from trace grpc-gateway")
			} else {
				errStatus := &spb.Status{}
				assert.NoError(t, json.Unmarshal([]byte(respStr), errStatus))
				if s, ok := status.FromError(test.err); ok {
					assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
					assert.True(t, proto.Equal(errStatus, s.Proto()))
				} else {
					assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
					assert.True(t, proto.Equal(errStatus, &spb.Status{Code: int32(codes.Unknown), Message: "my error"}))
				}
			}

			got := sink.AllTraces()[0]
			assert.EqualValues(t, got, traceOtlp)
		})
	}
}

func TestJsonMarshaling(t *testing.T) {
	m := jsonpb.Marshaler{}
	json, err := m.MarshalToString(&resourceSpansOtlp)
	assert.NoError(t, err)

	var resourceSpansOtlp2 otlptrace.ResourceSpans
	err = jsonpb.UnmarshalString(json, &resourceSpansOtlp2)
	assert.NoError(t, err)

	assert.EqualValues(t, resourceSpansOtlp, resourceSpansOtlp2)
}

func TestJsonUnmarshaling(t *testing.T) {
	var resourceSpansOtlp2 otlptrace.ResourceSpans
	err := jsonpb.UnmarshalString(`
		{
		  "instrumentation_library_spans": [
			{
			  "spans": [
				{
				}
			  ]
			}
		  ]
		}`, &resourceSpansOtlp2)
	assert.NoError(t, err)
	assert.EqualValues(t, otlpcommon.TraceID{}, resourceSpansOtlp2.InstrumentationLibrarySpans[0].Spans[0].TraceId)

	tests := []struct {
		name  string
		json  string
		bytes []byte
	}{
		{
			name:  "empty string trace id",
			json:  `""`,
			bytes: nil,
		},
		{
			name:  "zero bytes trace id",
			json:  `"00000000000000000000000000000000"`,
			bytes: []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var resourceSpansOtlp2 otlptrace.ResourceSpans
			jsonStr := fmt.Sprintf(`
			{
			  "instrumentation_library_spans": [
				{
				  "spans": [
					{
					  "trace_id": %v
					}
				  ]
				}
			  ]
			}`, test.json)
			err := jsonpb.UnmarshalString(jsonStr, &resourceSpansOtlp2)
			assert.NoError(t, err)
			assert.EqualValues(t, otlpcommon.NewTraceID(test.bytes), resourceSpansOtlp2.InstrumentationLibrarySpans[0].Spans[0].TraceId)
		})
	}
}

func TestProtoHttp(t *testing.T) {
	tests := []struct {
		name     string
		encoding string
		err      error
	}{
		{
			name:     "ProtoUncompressed",
			encoding: "",
		},
		{
			name:     "ProtoGzipCompressed",
			encoding: "gzip",
		},
		{
			name:     "NotGRPCError",
			encoding: "",
			err:      errors.New("my error"),
		},
		{
			name:     "GRPCError",
			encoding: "",
			err:      status.New(codes.Internal, "").Err(),
		},
	}
	addr := testutil.GetAvailableLocalAddress(t)

	// Set the buffer count to 1 to make it flush the test span immediately.
	tSink := new(exportertest.SinkTraceExporter)
	mSink := new(exportertest.SinkMetricsExporter)
	ocr := newHTTPReceiver(t, addr, tSink, mSink)

	require.NoError(t, ocr.Start(context.Background(), componenttest.NewNopHost()), "Failed to start trace receiver")
	defer ocr.Shutdown(context.Background())

	// TODO(nilebox): make starting server deterministic
	// Wait for the servers to start
	<-time.After(10 * time.Millisecond)

	url := fmt.Sprintf("http://%s/v1/trace", addr)

	wantOtlp := pdata.TracesToOtlp(testdata.GenerateTraceDataOneSpan())
	traceProto := collectortrace.ExportTraceServiceRequest{
		ResourceSpans: wantOtlp,
	}
	traceBytes, err := traceProto.Marshal()
	if err != nil {
		t.Errorf("Error marshaling protobuf: %v", err)
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var buf *bytes.Buffer
			var err error
			switch test.encoding {
			case "gzip":
				buf, err = compressGzip(traceBytes)
				require.NoError(t, err, "Error while gzip compressing trace: %v", err)
			default:
				buf = bytes.NewBuffer(traceBytes)
			}
			tSink.SetConsumeTraceError(test.err)
			req, err := http.NewRequest("POST", url, buf)
			require.NoError(t, err, "Error creating trace POST request: %v", err)
			req.Header.Set("Content-Type", "application/x-protobuf")
			req.Header.Set("Content-Encoding", test.encoding)

			client := &http.Client{}
			resp, err := client.Do(req)
			require.NoError(t, err, "Error posting trace to grpc-gateway server: %v", err)

			respBytes, err := ioutil.ReadAll(resp.Body)
			require.NoError(t, err, "Error reading response from trace grpc-gateway")
			require.NoError(t, resp.Body.Close(), "Error closing response body")

			if test.err == nil {
				require.Equal(t, 200, resp.StatusCode, "Unexpected return status")
				tmp := &collectortrace.ExportTraceServiceResponse{}
				err = tmp.Unmarshal(respBytes)
				require.NoError(t, err, "Unable to unmarshal response to ExportTraceServiceResponse proto")
			} else {
				errStatus := &spb.Status{}
				assert.NoError(t, proto.Unmarshal(respBytes, errStatus))
				if s, ok := status.FromError(test.err); ok {
					assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
					assert.True(t, proto.Equal(errStatus, s.Proto()))
				} else {
					assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
					assert.True(t, proto.Equal(errStatus, &spb.Status{Code: int32(codes.Unknown), Message: "my error"}))
				}
			}

			require.Equal(t, "application/x-protobuf", resp.Header.Get("Content-Type"), "Unexpected response Content-Type")

			gotOtlp := pdata.TracesToOtlp(tSink.AllTraces()[0])

			if len(gotOtlp) != len(wantOtlp) {
				t.Fatalf("len(traces):\nGot: %d\nWant: %d\n", len(gotOtlp), len(wantOtlp))
			}

			got := gotOtlp[0]
			want := wantOtlp[0]

			if !assert.EqualValues(t, got, want) {
				t.Errorf("Sending trace proto over http failed\nGot:\n%v\nWant:\n%v\n",
					got.String(),
					want.String())
			}
		})
	}
}

func TestOTLPReceiverInvalidContentEncoding(t *testing.T) {
	tests := []struct {
		name        string
		content     string
		encoding    string
		reqBodyFunc func() (*bytes.Buffer, error)
		resBodyFunc func() ([]byte, error)
		status      int
	}{
		{
			name:     "JsonGzipUncompressed",
			content:  "application/json",
			encoding: "gzip",
			reqBodyFunc: func() (*bytes.Buffer, error) {
				return bytes.NewBuffer([]byte(`{"key": "value"}`)), nil
			},
			resBodyFunc: func() ([]byte, error) {
				return json.Marshal(status.New(codes.InvalidArgument, "gzip: invalid header").Proto())
			},
			status: 400,
		},
		{
			name:     "ProtoGzipUncompressed",
			content:  "application/x-protobuf",
			encoding: "gzip",
			reqBodyFunc: func() (*bytes.Buffer, error) {
				return bytes.NewBuffer([]byte(`{"key": "value"}`)), nil
			},
			resBodyFunc: func() ([]byte, error) {
				return proto.Marshal(status.New(codes.InvalidArgument, "gzip: invalid header").Proto())
			},
			status: 400,
		},
	}
	addr := testutil.GetAvailableLocalAddress(t)

	// Set the buffer count to 1 to make it flush the test span immediately.
	tSink := new(exportertest.SinkTraceExporter)
	mSink := new(exportertest.SinkMetricsExporter)
	ocr := newHTTPReceiver(t, addr, tSink, mSink)

	require.NoError(t, ocr.Start(context.Background(), componenttest.NewNopHost()), "Failed to start trace receiver")
	defer ocr.Shutdown(context.Background())

	url := fmt.Sprintf("http://%s/v1/trace", addr)

	// Wait for the servers to start
	<-time.After(10 * time.Millisecond)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			body, err := test.reqBodyFunc()
			require.NoError(t, err, "Error creating request body: %v", err)

			req, err := http.NewRequest("POST", url, body)
			require.NoError(t, err, "Error creating trace POST request: %v", err)
			req.Header.Set("Content-Type", test.content)
			req.Header.Set("Content-Encoding", test.encoding)

			client := &http.Client{}
			resp, err := client.Do(req)
			require.NoError(t, err, "Error posting trace to grpc-gateway server: %v", err)

			respBytes, err := ioutil.ReadAll(resp.Body)
			require.NoError(t, err, "Error reading response from trace grpc-gateway")
			exRespBytes, err := test.resBodyFunc()
			require.NoError(t, err, "Error creating expecting response body")
			require.NoError(t, resp.Body.Close(), "Error closing response body")

			require.Equal(t, test.status, resp.StatusCode, "Unexpected return status")
			require.Equal(t, test.content, resp.Header.Get("Content-Type"), "Unexpected response Content-Type")
			require.Equal(t, exRespBytes, respBytes, "Unexpected response content")
		})
	}
}

func TestGRPCNewPortAlreadyUsed(t *testing.T) {
	addr := testutil.GetAvailableLocalAddress(t)
	ln, err := net.Listen("tcp", addr)
	require.NoError(t, err, "failed to listen on %q: %v", addr, err)
	defer ln.Close()

	r := newGRPCReceiver(t, otlpReceiverName, addr, new(exportertest.SinkTraceExporter), new(exportertest.SinkMetricsExporter))
	require.NotNil(t, r)

	require.Error(t, r.Start(context.Background(), componenttest.NewNopHost()))
}

func TestHTTPNewPortAlreadyUsed(t *testing.T) {
	addr := testutil.GetAvailableLocalAddress(t)
	ln, err := net.Listen("tcp", addr)
	require.NoError(t, err, "failed to listen on %q: %v", addr, err)
	defer ln.Close()

	r := newHTTPReceiver(t, addr, new(exportertest.SinkTraceExporter), new(exportertest.SinkMetricsExporter))
	require.NotNil(t, r)

	require.Error(t, r.Start(context.Background(), componenttest.NewNopHost()))
}

func TestGRPCStartWithoutConsumers(t *testing.T) {
	addr := testutil.GetAvailableLocalAddress(t)
	r := newGRPCReceiver(t, otlpReceiverName, addr, nil, nil)
	require.NotNil(t, r)
	require.Error(t, r.Start(context.Background(), componenttest.NewNopHost()))
}

func TestHTTPStartWithoutConsumers(t *testing.T) {
	addr := testutil.GetAvailableLocalAddress(t)
	r := newHTTPReceiver(t, addr, nil, nil)
	require.NotNil(t, r)
	require.Error(t, r.Start(context.Background(), componenttest.NewNopHost()))
}

// TestOTLPReceiverTrace_HandleNextConsumerResponse checks if the trace receiver
// is returning the proper response (return and metrics) when the next consumer
// in the pipeline reports error. The test changes the responses returned by the
// next trace consumer, checks if data was passed down the pipeline and if
// proper metrics were recorded. It also uses all endpoints supported by the
// trace receiver.
func TestOTLPReceiverTrace_HandleNextConsumerResponse(t *testing.T) {
	type ingestionStateTest struct {
		okToIngest   bool
		expectedCode codes.Code
	}
	tests := []struct {
		name                         string
		expectedReceivedBatches      int
		expectedIngestionBlockedRPCs int
		ingestionStates              []ingestionStateTest
	}{
		{
			name:                         "IngestTest",
			expectedReceivedBatches:      2,
			expectedIngestionBlockedRPCs: 1,
			ingestionStates: []ingestionStateTest{
				{
					okToIngest:   true,
					expectedCode: codes.OK,
				},
				{
					okToIngest:   false,
					expectedCode: codes.Unknown,
				},
				{
					okToIngest:   true,
					expectedCode: codes.OK,
				},
			},
		},
	}

	addr := testutil.GetAvailableLocalAddress(t)
	req := &collectortrace.ExportTraceServiceRequest{
		ResourceSpans: []*otlptrace.ResourceSpans{
			{
				InstrumentationLibrarySpans: []*otlptrace.InstrumentationLibrarySpans{
					{
						Spans: []*otlptrace.Span{
							{
								TraceId: otlpcommon.NewTraceID(
									[]byte{
										0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08,
										0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10,
									},
								),
							},
						},
					},
				},
			},
		},
	}

	exportBidiFn := func(
		t *testing.T,
		cc *grpc.ClientConn,
		msg *collectortrace.ExportTraceServiceRequest) error {

		acc := collectortrace.NewTraceServiceClient(cc)
		_, err := acc.Export(context.Background(), req)

		return err
	}

	exporters := []struct {
		receiverTag string
		exportFn    func(
			t *testing.T,
			cc *grpc.ClientConn,
			msg *collectortrace.ExportTraceServiceRequest) error
	}{
		{
			receiverTag: "otlp_trace",
			exportFn:    exportBidiFn,
		},
	}
	for _, exporter := range exporters {
		for _, tt := range tests {
			t.Run(tt.name+"/"+exporter.receiverTag, func(t *testing.T) {
				doneFn, err := obsreporttest.SetupRecordedMetricsTest()
				require.NoError(t, err)
				defer doneFn()

				sink := new(exportertest.SinkTraceExporter)

				ocr := newGRPCReceiver(t, exporter.receiverTag, addr, sink, nil)
				require.NotNil(t, ocr)
				require.NoError(t, ocr.Start(context.Background(), componenttest.NewNopHost()))
				defer ocr.Shutdown(context.Background())

				cc, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithBlock())
				require.NoError(t, err)
				defer cc.Close()

				for _, ingestionState := range tt.ingestionStates {
					if ingestionState.okToIngest {
						sink.SetConsumeTraceError(nil)
					} else {
						sink.SetConsumeTraceError(fmt.Errorf("%q: consumer error", tt.name))
					}

					err = exporter.exportFn(t, cc, req)

					status, ok := status.FromError(err)
					require.True(t, ok)
					assert.Equal(t, ingestionState.expectedCode, status.Code())
				}

				require.Equal(t, tt.expectedReceivedBatches, len(sink.AllTraces()))

				obsreporttest.CheckReceiverTracesViews(t, exporter.receiverTag, "grpc", int64(tt.expectedReceivedBatches), int64(tt.expectedIngestionBlockedRPCs))
			})
		}
	}
}

func TestGRPCInvalidTLSCredentials(t *testing.T) {
	cfg := &Config{
		ReceiverSettings: configmodels.ReceiverSettings{
			NameVal: "IncorrectTLS",
		},
		Protocols: Protocols{
			GRPC: &configgrpc.GRPCServerSettings{
				NetAddr: confignet.NetAddr{
					Endpoint:  "localhost:50000",
					Transport: "tcp",
				},
				TLSSetting: &configtls.TLSServerSetting{
					TLSSetting: configtls.TLSSetting{
						CertFile: "willfail",
					},
				},
			},
		},
	}

	// TLS is resolved during Creation of the receiver for GRPC.
	_, err := createReceiver(cfg)
	assert.EqualError(t, err,
		`failed to load TLS config: for auth via TLS, either both certificate and key must be supplied, or neither`)
}

func TestHTTPInvalidTLSCredentials(t *testing.T) {
	cfg := &Config{
		ReceiverSettings: configmodels.ReceiverSettings{
			NameVal: "IncorrectTLS",
		},
		Protocols: Protocols{
			HTTP: &confighttp.HTTPServerSettings{
				Endpoint: "localhost:50000",
				TLSSetting: &configtls.TLSServerSetting{
					TLSSetting: configtls.TLSSetting{
						CertFile: "willfail",
					},
				},
			},
		},
	}

	// TLS is resolved during Start for HTTP.
	r := newReceiver(t, NewFactory(), cfg, new(exportertest.SinkTraceExporter), new(exportertest.SinkMetricsExporter))
	assert.EqualError(t, r.Start(context.Background(), componenttest.NewNopHost()),
		`failed to load TLS config: for auth via TLS, either both certificate and key must be supplied, or neither`)
}

func newGRPCReceiver(t *testing.T, name string, endpoint string, tc consumer.TraceConsumer, mc consumer.MetricsConsumer) *otlpReceiver {
	factory := NewFactory()
	cfg := factory.CreateDefaultConfig().(*Config)
	cfg.SetName(name)
	cfg.GRPC.NetAddr.Endpoint = endpoint
	cfg.HTTP = nil
	return newReceiver(t, factory, cfg, tc, mc)
}

func newHTTPReceiver(t *testing.T, endpoint string, tc consumer.TraceConsumer, mc consumer.MetricsConsumer) *otlpReceiver {
	factory := NewFactory()
	cfg := factory.CreateDefaultConfig().(*Config)
	cfg.SetName(otlpReceiverName)
	cfg.HTTP.Endpoint = endpoint
	cfg.GRPC = nil
	return newReceiver(t, factory, cfg, tc, mc)
}

func newReceiver(t *testing.T, factory component.ReceiverFactory, cfg *Config, tc consumer.TraceConsumer, mc consumer.MetricsConsumer) *otlpReceiver {
	r, err := createReceiver(cfg)
	require.NoError(t, err)
	if tc != nil {
		params := component.ReceiverCreateParams{}
		_, err := factory.CreateTraceReceiver(context.Background(), params, cfg, tc)
		require.NoError(t, err)
	}
	if mc != nil {
		params := component.ReceiverCreateParams{}
		_, err := factory.CreateMetricsReceiver(context.Background(), params, cfg, mc)
		require.NoError(t, err)
	}
	return r
}

func compressGzip(body []byte) (*bytes.Buffer, error) {
	var buf bytes.Buffer

	gw := gzip.NewWriter(&buf)
	defer gw.Close()

	_, err := gw.Write(body)
	if err != nil {
		return nil, err
	}

	return &buf, nil
}
