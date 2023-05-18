// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

// Code generated by "pdata/internal/cmd/pdatagen/main.go". DO NOT EDIT.
// To regenerate this file run "make genpdata".

package plog

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"go.opentelemetry.io/collector/pdata/internal"
	"go.opentelemetry.io/collector/pdata/internal/data"
	otlplogs "go.opentelemetry.io/collector/pdata/internal/data/protogen/logs/v1"
	"go.opentelemetry.io/collector/pdata/pcommon"
)

func TestLogRecord_MoveTo(t *testing.T) {
	ms := generateTestLogRecord()
	dest := NewLogRecord()
	ms.MoveTo(dest)
	assert.Equal(t, NewLogRecord(), ms)
	assert.Equal(t, generateTestLogRecord(), dest)
}

func TestLogRecord_CopyTo(t *testing.T) {
	ms := NewLogRecord()
	orig := NewLogRecord()
	orig.CopyTo(ms)
	assert.Equal(t, orig, ms)
	orig = generateTestLogRecord()
	orig.CopyTo(ms)
	assert.Equal(t, orig, ms)
}

func TestLogRecord_ObservedTimestamp(t *testing.T) {
	ms := NewLogRecord()
	assert.Equal(t, pcommon.Timestamp(0), ms.ObservedTimestamp())
	testValObservedTimestamp := pcommon.Timestamp(1234567890)
	ms.SetObservedTimestamp(testValObservedTimestamp)
	assert.Equal(t, testValObservedTimestamp, ms.ObservedTimestamp())
}

func TestLogRecord_Timestamp(t *testing.T) {
	ms := NewLogRecord()
	assert.Equal(t, pcommon.Timestamp(0), ms.Timestamp())
	testValTimestamp := pcommon.Timestamp(1234567890)
	ms.SetTimestamp(testValTimestamp)
	assert.Equal(t, testValTimestamp, ms.Timestamp())
}

func TestLogRecord_TraceID(t *testing.T) {
	ms := NewLogRecord()
	assert.Equal(t, pcommon.TraceID(data.TraceID([16]byte{})), ms.TraceID())
	testValTraceID := pcommon.TraceID(data.TraceID([16]byte{1, 2, 3, 4, 5, 6, 7, 8, 8, 7, 6, 5, 4, 3, 2, 1}))
	ms.SetTraceID(testValTraceID)
	assert.Equal(t, testValTraceID, ms.TraceID())
}

func TestLogRecord_SpanID(t *testing.T) {
	ms := NewLogRecord()
	assert.Equal(t, pcommon.SpanID(data.SpanID([8]byte{})), ms.SpanID())
	testValSpanID := pcommon.SpanID(data.SpanID([8]byte{8, 7, 6, 5, 4, 3, 2, 1}))
	ms.SetSpanID(testValSpanID)
	assert.Equal(t, testValSpanID, ms.SpanID())
}

func TestLogRecord_Flags(t *testing.T) {
	ms := NewLogRecord()
	assert.Equal(t, LogRecordFlags(0), ms.Flags())
	testValFlags := LogRecordFlags(1)
	ms.SetFlags(testValFlags)
	assert.Equal(t, testValFlags, ms.Flags())
}

func TestLogRecord_SeverityText(t *testing.T) {
	ms := NewLogRecord()
	assert.Equal(t, "", ms.SeverityText())
	ms.SetSeverityText("INFO")
	assert.Equal(t, "INFO", ms.SeverityText())
}

func TestLogRecord_SeverityNumber(t *testing.T) {
	ms := NewLogRecord()
	assert.Equal(t, SeverityNumber(otlplogs.SeverityNumber(0)), ms.SeverityNumber())
	testValSeverityNumber := SeverityNumber(otlplogs.SeverityNumber(5))
	ms.SetSeverityNumber(testValSeverityNumber)
	assert.Equal(t, testValSeverityNumber, ms.SeverityNumber())
}

func TestLogRecord_Body(t *testing.T) {
	ms := NewLogRecord()
	internal.FillTestValue(internal.Value(ms.Body()))
	assert.Equal(t, pcommon.Value(internal.GenerateTestValue()), ms.Body())
}

func TestLogRecord_Attributes(t *testing.T) {
	ms := NewLogRecord()
	assert.Equal(t, pcommon.NewMap(), ms.Attributes())
	internal.FillTestMap(internal.Map(ms.Attributes()))
	assert.Equal(t, pcommon.Map(internal.GenerateTestMap()), ms.Attributes())
}

func TestLogRecord_DroppedAttributesCount(t *testing.T) {
	ms := NewLogRecord()
	assert.Equal(t, uint32(0), ms.DroppedAttributesCount())
	ms.SetDroppedAttributesCount(uint32(17))
	assert.Equal(t, uint32(17), ms.DroppedAttributesCount())
}

func generateTestLogRecord() LogRecord {
	tv := NewLogRecord()
	fillTestLogRecord(tv)
	return tv
}

func fillTestLogRecord(tv LogRecord) {
	tv.orig.ObservedTimeUnixNano = 1234567890
	tv.orig.TimeUnixNano = 1234567890
	tv.orig.TraceId = data.TraceID([16]byte{1, 2, 3, 4, 5, 6, 7, 8, 8, 7, 6, 5, 4, 3, 2, 1})
	tv.orig.SpanId = data.SpanID([8]byte{8, 7, 6, 5, 4, 3, 2, 1})
	tv.orig.Flags = 1
	tv.orig.SeverityText = "INFO"
	tv.orig.SeverityNumber = otlplogs.SeverityNumber(5)
	internal.FillTestValue(internal.NewValue(&tv.orig.Body))
	internal.FillTestMap(internal.NewMap(&tv.orig.Attributes))
	tv.orig.DroppedAttributesCount = uint32(17)
}
