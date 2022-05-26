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

package pmetric

import (
	"testing"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/collector/pdata/internal"
)

var metricsOTLP = func() Metrics {
	md := NewMetrics()
	rm := md.ResourceMetrics().AppendEmpty()
	rm.Resource().Attributes().UpsertString("host.name", "testHost")
	il := rm.ScopeMetrics().AppendEmpty()
	il.Scope().SetName("name")
	il.Scope().SetVersion("version")
	il.Metrics().AppendEmpty().SetName("testMetric")
	return md
}()

var metricsJSON = `{"resourceMetrics":[{"resource":{"attributes":[{"key":"host.name","value":{"stringValue":"testHost"}}]},"scopeMetrics":[{"scope":{"name":"name","version":"version"},"metrics":[{"name":"testMetric"}]}]}]}`

func TestMetricsJSON(t *testing.T) {
	encoder := NewJSONMarshaler()
	jsonBuf, err := encoder.MarshalMetrics(metricsOTLP)
	assert.NoError(t, err)

	decoder := NewJSONUnmarshaler()
	var got interface{}
	got, err = decoder.UnmarshalMetrics(jsonBuf)
	assert.NoError(t, err)

	assert.EqualValues(t, metricsOTLP, got)
}

func TestMetricsJSON_Marshal(t *testing.T) {
	encoder := NewJSONMarshaler()
	jsonBuf, err := encoder.MarshalMetrics(metricsOTLP)
	assert.NoError(t, err)
	assert.Equal(t, metricsJSON, string(jsonBuf))
}

func TestMetricsNil(t *testing.T) {
	jsonBuf := `{
"resourceMetrics": [
	{
	"resource": {
		"attributes": [
		{
			"key": "service.name",
			"value": {
			"stringValue": "unknown_service:node"
			}
		},
		{
			"key": "telemetry.sdk.language",
			"value": {
			"stringValue": "nodejs"
			}
		},
		{
			"key": "telemetry.sdk.name",
			"value": {
			"stringValue": "opentelemetry"
			}
		},
		{
			"key": "telemetry.sdk.version",
			"value": {
			"stringValue": "0.24.0"
			}
		}
		],
		"droppedAttributesCount": 0
	},
	"instrumentationLibraryMetrics": [
		{
		"metrics": [
			{
			"name": "metric_name",
			"description": "Example of a UpDownCounter",
			"unit": "1",
			"doubleSum": {
				"dataPoints": [
				{
					"labels": [
					{
						"key": "pid",
						"value": "50712"
					}
					],
					"value": 1,
					"startTimeUnixNano": 1631056185376000000,
					"timeUnixNano": 1631056185378763800
				}
				],
				"isMonotonic": false,
				"aggregationTemporality": 2
			}
			},
			{
			"name": "your_metric_name",
			"description": "Example of a sync observer with callback",
			"unit": "1",
			"doubleGauge": {
				"dataPoints": [
				{
					"labels": [
					{
						"key": "label",
						"value": "1"
					}
					],
					"value": 0.07604853280317792,
					"startTimeUnixNano": 1631056185376000000,
					"timeUnixNano": 1631056189394600700
				}
				]
			}
			},
			{
			"name": "your_metric_name",
			"description": "Example of a sync observer with callback",
			"unit": "1",
			"doubleGauge": {
				"dataPoints": [
				{
					"labels": [
					{
						"key": "label",
						"value": "2"
					}
					],
					"value": 0.9332005145656965,
					"startTimeUnixNano": 1631056185376000000,
					"timeUnixNano": 1631056189394630400
				}
				]
			}
			}
		],
		"instrumentationLibrary": {
			"name": "example-meter"
		}
		}
	]
	}
]
}`
	decoder := NewJSONUnmarshaler()
	var got interface{}
	got, err := decoder.UnmarshalMetrics([]byte(jsonBuf))
	assert.Error(t, err)
	assert.EqualValues(t, Metrics{}, got)
}

var metricsSumOTLPFull = func() Metrics {
	metric := NewMetrics()
	rs := metric.ResourceMetrics().AppendEmpty()
	rs.SetSchemaUrl("schemaURL")
	// Add resource.
	rs.Resource().Attributes().UpsertString("host.name", "testHost")
	rs.Resource().Attributes().UpsertString("service.name", "testService")
	rs.Resource().SetDroppedAttributesCount(1)
	// Add InstrumentationLibraryMetrics.
	m := rs.ScopeMetrics().AppendEmpty()
	m.Scope().SetName("instrumentation name")
	m.Scope().SetVersion("instrumentation version")
	m.SetSchemaUrl("schemaURL")
	// Add Metric
	sumData := m.Metrics().AppendEmpty()
	sumData.SetName("test sum")
	sumData.SetDescription("test sum")
	sumData.SetDataType(internal.MetricDataTypeSum)
	sumData.SetUnit("unit")
	sumData.Sum().SetAggregationTemporality(internal.MetricAggregationTemporalityCumulative)
	sumData.Sum().SetIsMonotonic(true)
	datapoint := sumData.Sum().DataPoints().AppendEmpty()
	datapoint.SetStartTimestamp(internal.NewTimestampFromTime(time.Now()))
	datapoint.SetFlags(internal.MetricDataPointFlagsNone)
	datapoint.SetIntVal(100)
	datapoint.Attributes().UpsertString("string", "value")
	datapoint.Attributes().UpsertBool("bool", true)
	datapoint.Attributes().UpsertInt("int", 1)
	datapoint.Attributes().UpsertDouble("double", 1.1)
	datapoint.Attributes().UpsertMBytes("bytes", []byte("foo"))
	exemplar := datapoint.Exemplars().AppendEmpty()
	exemplar.SetDoubleVal(99.3)
	exemplar.SetTimestamp(internal.NewTimestampFromTime(time.Now()))
	traceID := internal.NewTraceID([16]byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10})
	spanID := internal.NewSpanID([8]byte{0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18})
	exemplar.SetSpanID(spanID)
	exemplar.SetTraceID(traceID)
	exemplar.FilteredAttributes().UpsertString("service.name", "testService")
	datapoint.SetTimestamp(internal.NewTimestampFromTime(time.Now()))
	return metric
}

var metricsGaugeOTLPFull = func() Metrics {
	metric := NewMetrics()
	rs := metric.ResourceMetrics().AppendEmpty()
	rs.SetSchemaUrl("schemaURL")
	// Add resource.
	rs.Resource().Attributes().UpsertString("host.name", "testHost")
	rs.Resource().Attributes().UpsertString("service.name", "testService")
	rs.Resource().SetDroppedAttributesCount(1)
	// Add InstrumentationLibraryMetrics.
	m := rs.ScopeMetrics().AppendEmpty()
	m.Scope().SetName("instrumentation name")
	m.Scope().SetVersion("instrumentation version")
	m.SetSchemaUrl("schemaURL")
	// Add Metric
	gaugeData := m.Metrics().AppendEmpty()
	gaugeData.SetName("test gauge")
	gaugeData.SetDescription("test gauge")
	gaugeData.SetDataType(internal.MetricDataTypeGauge)
	gaugeData.SetUnit("unit")
	datapoint := gaugeData.Gauge().DataPoints().AppendEmpty()
	datapoint.SetStartTimestamp(internal.NewTimestampFromTime(time.Now()))
	datapoint.SetFlags(internal.MetricDataPointFlagsNone)
	datapoint.SetDoubleVal(10.2)
	datapoint.Attributes().UpsertString("string", "value")
	datapoint.Attributes().UpsertBool("bool", true)
	datapoint.Attributes().UpsertInt("int", 1)
	datapoint.Attributes().UpsertDouble("double", 1.1)
	datapoint.Attributes().UpsertMBytes("bytes", []byte("foo"))
	exemplar := datapoint.Exemplars().AppendEmpty()
	exemplar.SetDoubleVal(99.3)
	exemplar.SetTimestamp(internal.NewTimestampFromTime(time.Now()))
	traceID := internal.NewTraceID([16]byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10})
	spanID := internal.NewSpanID([8]byte{0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18})
	exemplar.SetSpanID(spanID)
	exemplar.SetTraceID(traceID)
	exemplar.FilteredAttributes().UpsertString("service.name", "testService")
	datapoint.SetTimestamp(internal.NewTimestampFromTime(time.Now()))
	return metric
}

var metricsHistogramOTLPFull = func() Metrics {
	metric := NewMetrics()
	rs := metric.ResourceMetrics().AppendEmpty()
	rs.SetSchemaUrl("schemaURL")
	// Add resource.
	rs.Resource().Attributes().UpsertString("host.name", "testHost")
	rs.Resource().Attributes().UpsertString("service.name", "testService")
	rs.Resource().SetDroppedAttributesCount(1)
	// Add InstrumentationLibraryMetrics.
	m := rs.ScopeMetrics().AppendEmpty()
	m.Scope().SetName("instrumentation name")
	m.Scope().SetVersion("instrumentation version")
	m.SetSchemaUrl("schemaURL")
	// Add Metric
	histogramData := m.Metrics().AppendEmpty()
	histogramData.SetName("test Histogram")
	histogramData.SetDescription("test Histogram")
	histogramData.SetDataType(internal.MetricDataTypeHistogram)
	histogramData.SetUnit("unit")
	histogramData.Histogram().SetAggregationTemporality(MetricAggregationTemporalityCumulative)
	datapoint := histogramData.Histogram().DataPoints().AppendEmpty()
	datapoint.SetStartTimestamp(internal.NewTimestampFromTime(time.Now()))
	datapoint.SetFlags(internal.MetricDataPointFlagsNone)
	datapoint.Attributes().UpsertString("string", "value")
	datapoint.Attributes().UpsertBool("bool", true)
	datapoint.Attributes().UpsertInt("int", 1)
	datapoint.Attributes().UpsertDouble("double", 1.1)
	datapoint.Attributes().UpsertMBytes("bytes", []byte("foo"))
	datapoint.SetCount(4)
	datapoint.SetSum(345)
	datapoint.SetMBucketCounts([]uint64{1, 1, 2})
	datapoint.SetMExplicitBounds([]float64{10, 100})
	exemplar := datapoint.Exemplars().AppendEmpty()
	exemplar.SetDoubleVal(99.3)
	exemplar.SetTimestamp(internal.NewTimestampFromTime(time.Now()))
	datapoint.SetMin(float64(time.Now().Unix()))
	traceID := internal.NewTraceID([16]byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10})
	spanID := internal.NewSpanID([8]byte{0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18})
	exemplar.SetSpanID(spanID)
	exemplar.SetTraceID(traceID)
	exemplar.FilteredAttributes().UpsertString("service.name", "testService")
	datapoint.SetMax(float64(time.Now().Unix()))
	datapoint.SetTimestamp(internal.NewTimestampFromTime(time.Now()))
	return metric
}

var metricsExponentialHistogramOTLPFull = func() Metrics {
	metric := NewMetrics()
	rs := metric.ResourceMetrics().AppendEmpty()
	rs.SetSchemaUrl("schemaURL")
	// Add resource.
	rs.Resource().Attributes().UpsertString("host.name", "testHost")
	rs.Resource().Attributes().UpsertString("service.name", "testService")
	rs.Resource().SetDroppedAttributesCount(1)
	// Add InstrumentationLibraryMetrics.
	m := rs.ScopeMetrics().AppendEmpty()
	m.Scope().SetName("instrumentation name")
	m.Scope().SetVersion("instrumentation version")
	m.SetSchemaUrl("schemaURL")
	// Add Metric
	histogramData := m.Metrics().AppendEmpty()
	histogramData.SetName("test ExponentialHistogram")
	histogramData.SetDescription("test ExponentialHistogram")
	histogramData.SetDataType(internal.MetricDataTypeExponentialHistogram)
	histogramData.SetUnit("unit")
	histogramData.ExponentialHistogram().SetAggregationTemporality(MetricAggregationTemporalityCumulative)
	datapoint := histogramData.ExponentialHistogram().DataPoints().AppendEmpty()
	datapoint.SetStartTimestamp(internal.NewTimestampFromTime(time.Now()))
	datapoint.SetFlags(internal.MetricDataPointFlagsNone)
	datapoint.Attributes().UpsertString("string", "value")
	datapoint.Attributes().UpsertBool("bool", true)
	datapoint.Attributes().UpsertInt("int", 1)
	datapoint.Attributes().UpsertDouble("double", 1.1)
	datapoint.Attributes().UpsertMBytes("bytes", []byte("foo"))
	datapoint.SetCount(4)
	datapoint.SetSum(345)
	datapoint.Positive().SetMBucketCounts([]uint64{1, 1, 2})
	datapoint.Positive().SetOffset(2)
	exemplar := datapoint.Exemplars().AppendEmpty()
	exemplar.SetDoubleVal(99.3)
	exemplar.SetTimestamp(internal.NewTimestampFromTime(time.Now()))
	datapoint.SetMin(float64(time.Now().Unix()))
	traceID := internal.NewTraceID([16]byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10})
	spanID := internal.NewSpanID([8]byte{0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18})
	exemplar.SetSpanID(spanID)
	exemplar.SetTraceID(traceID)
	exemplar.FilteredAttributes().UpsertString("service.name", "testService")
	datapoint.SetMax(float64(time.Now().Unix()))
	datapoint.SetTimestamp(internal.NewTimestampFromTime(time.Now()))
	return metric
}

var metricsSummaryOTLPFull = func() Metrics {
	metric := NewMetrics()
	rs := metric.ResourceMetrics().AppendEmpty()
	rs.SetSchemaUrl("schemaURL")
	// Add resource.
	rs.Resource().Attributes().UpsertString("host.name", "testHost")
	rs.Resource().Attributes().UpsertString("service.name", "testService")
	rs.Resource().SetDroppedAttributesCount(1)
	// Add InstrumentationLibraryMetrics.
	m := rs.ScopeMetrics().AppendEmpty()
	m.Scope().SetName("instrumentation name")
	m.Scope().SetVersion("instrumentation version")
	m.SetSchemaUrl("schemaURL")
	// Add Metric
	sumData := m.Metrics().AppendEmpty()
	sumData.SetName("test summary")
	sumData.SetDescription("test summary")
	sumData.SetDataType(internal.MetricDataTypeSummary)
	sumData.SetUnit("unit")
	datapoint := sumData.Summary().DataPoints().AppendEmpty()
	datapoint.SetStartTimestamp(internal.NewTimestampFromTime(time.Now()))
	datapoint.SetFlags(internal.MetricDataPointFlagsNone)
	datapoint.SetCount(100)
	datapoint.SetSum(100)
	quantile := datapoint.QuantileValues().AppendEmpty()
	quantile.SetQuantile(0.5)
	quantile.SetValue(1.2)
	datapoint.Attributes().UpsertString("string", "value")
	datapoint.Attributes().UpsertBool("bool", true)
	datapoint.Attributes().UpsertInt("int", 1)
	datapoint.Attributes().UpsertDouble("double", 1.1)
	datapoint.Attributes().UpsertMBytes("bytes", []byte("foo"))
	datapoint.SetTimestamp(internal.NewTimestampFromTime(time.Now()))
	return metric
}

func TestMetricsSumDataJSONFull(t *testing.T) {
	m := metricsSumOTLPFull()
	encoder := NewJSONMarshaler()
	jsonBuf, err := encoder.MarshalMetrics(m)
	assert.NoError(t, err)
	decoder := NewJSONUnmarshaler()
	got, err := decoder.UnmarshalMetrics(jsonBuf)
	assert.NoError(t, err)
	assert.EqualValues(t, m, got)
}

func TestMetricsGaugeDataJSONFull(t *testing.T) {
	m := metricsGaugeOTLPFull()
	encoder := NewJSONMarshaler()
	jsonBuf, err := encoder.MarshalMetrics(m)
	assert.NoError(t, err)
	decoder := NewJSONUnmarshaler()
	got, err := decoder.UnmarshalMetrics(jsonBuf)
	assert.NoError(t, err)
	assert.EqualValues(t, m, got)
}

func TestMetricsHistogramDataJSONFull(t *testing.T) {
	m := metricsHistogramOTLPFull()
	encoder := NewJSONMarshaler()
	jsonBuf, err := encoder.MarshalMetrics(m)
	assert.NoError(t, err)
	decoder := NewJSONUnmarshaler()
	got, err := decoder.UnmarshalMetrics(jsonBuf)
	assert.NoError(t, err)
	assert.EqualValues(t, m, got)
}

func TestMetricsExponentialHistogramDataJSONFull(t *testing.T) {
	m := metricsExponentialHistogramOTLPFull()
	encoder := NewJSONMarshaler()
	jsonBuf, err := encoder.MarshalMetrics(m)
	assert.NoError(t, err)
	decoder := NewJSONUnmarshaler()
	got, err := decoder.UnmarshalMetrics(jsonBuf)
	assert.NoError(t, err)
	assert.EqualValues(t, m, got)
}

func TestMetricsSummaryDataJSONFull(t *testing.T) {
	m := metricsSummaryOTLPFull()
	encoder := NewJSONMarshaler()
	jsonBuf, err := encoder.MarshalMetrics(m)
	assert.NoError(t, err)
	decoder := NewJSONUnmarshaler()
	got, err := decoder.UnmarshalMetrics(jsonBuf)
	assert.NoError(t, err)
	assert.EqualValues(t, m, got)
}

func TestReadKvlistValueUnknownField(t *testing.T) {
	jsonStr := `{"extra":""}`
	iter := jsoniter.ConfigFastest.BorrowIterator([]byte(jsonStr))
	defer jsoniter.ConfigFastest.ReturnIterator(iter)
	readKvlistValue(iter)
	if assert.Error(t, iter.Error) {
		assert.Contains(t, iter.Error.Error(), "unknown field")
	}
}

func TestReadArrayUnknownField(t *testing.T) {
	jsonStr := `{"extra":""}`
	iter := jsoniter.ConfigFastest.BorrowIterator([]byte(jsonStr))
	defer jsoniter.ConfigFastest.ReturnIterator(iter)
	readArray(iter)
	if assert.Error(t, iter.Error) {
		assert.Contains(t, iter.Error.Error(), "unknown field")
	}
}
