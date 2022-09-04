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

package plog

import (
	otlplogs "go.opentelemetry.io/collector/pdata/internal/data/protogen/logs/v1"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"testing"

	"github.com/stretchr/testify/assert"
)

var logsOTLP = func() Logs {
	ld := NewLogs()
	rl := ld.ResourceLogs().AppendEmpty()
	rl.Resource().Attributes().UpsertString("host.name", "testHost")
	rl.SetSchemaUrl("testSchemaURL")
	il := rl.ScopeLogs().AppendEmpty()
	il.Scope().SetName("name")
	il.Scope().SetVersion("version")
	il.Scope().SetDroppedAttributesCount(1)
	lg := il.LogRecords().AppendEmpty()
	lg.SetSeverityNumber(SeverityNumber(otlplogs.SeverityNumber_SEVERITY_NUMBER_ERROR))
	lg.SetSeverityText("Error")
	lg.SetDroppedAttributesCount(1)
	lg.SetFlags(LogRecordFlags(otlplogs.LogRecordFlags_LOG_RECORD_FLAG_UNSPECIFIED))
	traceID := pcommon.TraceID([16]byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10})
	spanID := pcommon.SpanID([8]byte{0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18})
	lg.SetTraceID(traceID)
	lg.SetSpanID(spanID)
	return ld
}()

func TestLogsJSON(t *testing.T) {
	encoder := NewJSONMarshaler()
	jsonBuf, err := encoder.MarshalLogs(logsOTLP)
	assert.NoError(t, err)

	decoder := NewJSONUnmarshaler()
	var got interface{}
	got, err = decoder.UnmarshalLogs(jsonBuf)
	assert.NoError(t, err)

	assert.EqualValues(t, logsOTLP, got)
}

var logsJSON = `{"resourceLogs":[{"resource":{"attributes":[{"key":"host.name","value":{"stringValue":"testHost"}}]},"scopeLogs":[{"scope":{"name":"name","version":"version"},"logRecords":[{"severityText":"Error","body":{},"traceId":"","spanId":""}]}]}]}`

func TestLogsJSON_Marshal(t *testing.T) {
	encoder := NewJSONMarshaler()
	jsonBuf, err := encoder.MarshalLogs(logsOTLP)
	assert.NoError(t, err)
	assert.Equal(t, logsJSON, string(jsonBuf))
}
