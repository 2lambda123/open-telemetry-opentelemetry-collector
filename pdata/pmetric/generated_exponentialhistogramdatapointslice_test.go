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
	"unsafe"

	"github.com/stretchr/testify/assert"

	otlpmetrics "go.opentelemetry.io/collector/pdata/internal/data/protogen/metrics/v1"
)

func TestExponentialHistogramDataPointSlice(t *testing.T) {
	es := NewExponentialHistogramDataPointSlice()
	assert.Equal(t, 0, es.Len())

	emptyVal := NewExponentialHistogramDataPoint()
	testVal := generateTestExponentialHistogramDataPoint()
	for i := 0; i < 7; i++ {
		el := es.AppendEmpty()
		assert.Equal(t, emptyVal.getOrig(), es.At(i).getOrig())
		fillTestExponentialHistogramDataPoint(el)
		assert.Equal(t, testVal.getOrig(), es.At(i).getOrig())
	}
	assert.Equal(t, 7, es.Len())
}

func TestExponentialHistogramDataPointSlice_CopyTo(t *testing.T) {
	dest := NewExponentialHistogramDataPointSlice()
	// Test CopyTo to empty
	NewExponentialHistogramDataPointSlice().CopyTo(dest)
	assert.Equal(t, NewExponentialHistogramDataPointSlice(), dest)

	// Test CopyTo larger slice
	generateTestExponentialHistogramDataPointSlice().CopyTo(dest)
	assert.Equal(t, generateTestExponentialHistogramDataPointSlice(), dest)

	// Test CopyTo same size slice
	generateTestExponentialHistogramDataPointSlice().CopyTo(dest)
	assert.Equal(t, generateTestExponentialHistogramDataPointSlice(), dest)
}

func TestExponentialHistogramDataPointSlice_EnsureCapacity(t *testing.T) {
	es := generateTestExponentialHistogramDataPointSlice()

	// Test ensure smaller capacity.
	const ensureSmallLen = 4
	es.EnsureCapacity(ensureSmallLen)
	assert.Less(t, ensureSmallLen, es.Len())
	assert.Equal(t, es.Len(), cap(*es.getOrig()))
	assert.Equal(t, generateTestExponentialHistogramDataPointSlice(), es)

	// Test ensure larger capacity
	const ensureLargeLen = 9
	es.EnsureCapacity(ensureLargeLen)
	assert.Less(t, generateTestExponentialHistogramDataPointSlice().Len(), ensureLargeLen)
	assert.Equal(t, ensureLargeLen, cap(*es.getOrig()))
	assert.Equal(t, generateTestExponentialHistogramDataPointSlice(), es)
}

func TestExponentialHistogramDataPointSlice_MoveAndAppendTo(t *testing.T) {
	// Test MoveAndAppendTo to empty
	expectedSlice := generateTestExponentialHistogramDataPointSlice()
	dest := NewExponentialHistogramDataPointSlice()
	src := generateTestExponentialHistogramDataPointSlice()
	src.MoveAndAppendTo(dest)
	assert.Equal(t, generateTestExponentialHistogramDataPointSlice(), dest)
	assert.Equal(t, 0, src.Len())
	assert.Equal(t, expectedSlice.Len(), dest.Len())

	// Test MoveAndAppendTo empty slice
	src.MoveAndAppendTo(dest)
	assert.Equal(t, generateTestExponentialHistogramDataPointSlice(), dest)
	assert.Equal(t, 0, src.Len())
	assert.Equal(t, expectedSlice.Len(), dest.Len())

	// Test MoveAndAppendTo not empty slice
	generateTestExponentialHistogramDataPointSlice().MoveAndAppendTo(dest)
	assert.Equal(t, 2*expectedSlice.Len(), dest.Len())
	for i := 0; i < expectedSlice.Len(); i++ {
		assert.Equal(t, expectedSlice.At(i).getOrig(), dest.At(i).getOrig())
		assert.Equal(t, expectedSlice.At(i).getOrig(), dest.At(i+expectedSlice.Len()).getOrig())
	}
}

func TestExponentialHistogramDataPointSlice_RemoveIf(t *testing.T) {
	// Test RemoveIf on empty slice
	emptySlice := NewExponentialHistogramDataPointSlice()
	emptySlice.RemoveIf(func(el ExponentialHistogramDataPoint) bool {
		t.Fail()
		return false
	})

	// Test RemoveIf
	filtered := generateTestExponentialHistogramDataPointSlice()
	pos := 0
	filtered.RemoveIf(func(el ExponentialHistogramDataPoint) bool {
		pos++
		return pos%3 == 0
	})
	assert.Equal(t, 5, filtered.Len())
}

func TestExponentialHistogramDataPointSlice_Sort(t *testing.T) {
	es := generateTestExponentialHistogramDataPointSlice()
	es.Sort(func(a, b ExponentialHistogramDataPoint) bool {
		return uintptr(unsafe.Pointer(a.getOrig())) < uintptr(unsafe.Pointer(b.getOrig()))
	})
	for i := 1; i < es.Len(); i++ {
		assert.True(t, uintptr(unsafe.Pointer(es.At(i-1).getOrig())) < uintptr(unsafe.Pointer(es.At(i).getOrig())))
	}
	es.Sort(func(a, b ExponentialHistogramDataPoint) bool {
		return uintptr(unsafe.Pointer(a.getOrig())) > uintptr(unsafe.Pointer(b.getOrig()))
	})
	for i := 1; i < es.Len(); i++ {
		assert.True(t, uintptr(unsafe.Pointer(es.At(i-1).getOrig())) > uintptr(unsafe.Pointer(es.At(i).getOrig())))
	}
}

func generateTestExponentialHistogramDataPointSlice() ExponentialHistogramDataPointSlice {
	es := NewExponentialHistogramDataPointSlice()
	fillTestExponentialHistogramDataPointSlice(es)
	return es
}

func fillTestExponentialHistogramDataPointSlice(es ExponentialHistogramDataPointSlice) {
	*es.getOrig() = make([]*otlpmetrics.ExponentialHistogramDataPoint, 7)
	for i := 0; i < 7; i++ {
		(*es.getOrig())[i] = &otlpmetrics.ExponentialHistogramDataPoint{}
		fillTestExponentialHistogramDataPoint(newExponentialHistogramDataPointFromOrig((*es.getOrig())[i]))
	}
}
