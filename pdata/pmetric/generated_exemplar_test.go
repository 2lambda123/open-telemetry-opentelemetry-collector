// Copyright The OpenTelemetry Authors
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

// Code generated by "pdata/internal/cmd/pdatagen/main.go". DO NOT EDIT.
// To regenerate this file run "make genpdata".

package pmetric

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"go.opentelemetry.io/collector/pdata/internal"
	"go.opentelemetry.io/collector/pdata/internal/data"
	"go.opentelemetry.io/collector/pdata/pcommon"
)

func TestExemplar_MoveTo(t *testing.T) {
	ms := generateTestExemplar()
	dest := NewMutableExemplar()
	ms.MoveTo(dest)
	assert.Equal(t, NewMutableExemplar(), ms)
	assert.Equal(t, generateTestExemplar(), dest)
}

func TestExemplar_CopyTo(t *testing.T) {
	ms := NewMutableExemplar()
	orig := NewMutableExemplar()
	orig.CopyTo(ms)
	assert.Equal(t, orig, ms)
	orig = generateTestExemplar()
	orig.CopyTo(ms)
	assert.Equal(t, orig, ms)
}

func TestExemplar_Timestamp(t *testing.T) {
	ms := NewMutableExemplar()
	assert.Equal(t, pcommon.Timestamp(0), ms.Timestamp())
	testValTimestamp := pcommon.Timestamp(1234567890)
	ms.SetTimestamp(testValTimestamp)
	assert.Equal(t, testValTimestamp, ms.Timestamp())
}

func TestExemplar_ValueType(t *testing.T) {
	tv := NewMutableExemplar()
	assert.Equal(t, ExemplarValueTypeEmpty, tv.ValueType())
}

func TestExemplar_DoubleValue(t *testing.T) {
	ms := NewMutableExemplar()
	assert.Equal(t, float64(0.0), ms.DoubleValue())
	ms.SetDoubleValue(float64(17.13))
	assert.Equal(t, float64(17.13), ms.DoubleValue())
	assert.Equal(t, ExemplarValueTypeDouble, ms.ValueType())
}

func TestExemplar_IntValue(t *testing.T) {
	ms := NewMutableExemplar()
	assert.Equal(t, int64(0), ms.IntValue())
	ms.SetIntValue(int64(17))
	assert.Equal(t, int64(17), ms.IntValue())
	assert.Equal(t, ExemplarValueTypeInt, ms.ValueType())
}

func TestExemplar_FilteredAttributes(t *testing.T) {
	ms := NewMutableExemplar()
	assert.Equal(t, pcommon.NewMutableMap(), ms.FilteredAttributes())
	internal.FillTestMap(internal.MutableMap(ms.FilteredAttributes()))
	assert.Equal(t, pcommon.MutableMap(internal.GenerateTestMap()), ms.FilteredAttributes())
}

func TestExemplar_TraceID(t *testing.T) {
	ms := NewMutableExemplar()
	assert.Equal(t, pcommon.TraceID(data.TraceID([16]byte{})), ms.TraceID())
	testValTraceID := pcommon.TraceID(data.TraceID([16]byte{1, 2, 3, 4, 5, 6, 7, 8, 8, 7, 6, 5, 4, 3, 2, 1}))
	ms.SetTraceID(testValTraceID)
	assert.Equal(t, testValTraceID, ms.TraceID())
}

func TestExemplar_SpanID(t *testing.T) {
	ms := NewMutableExemplar()
	assert.Equal(t, pcommon.SpanID(data.SpanID([8]byte{})), ms.SpanID())
	testValSpanID := pcommon.SpanID(data.SpanID([8]byte{8, 7, 6, 5, 4, 3, 2, 1}))
	ms.SetSpanID(testValSpanID)
	assert.Equal(t, testValSpanID, ms.SpanID())
}
