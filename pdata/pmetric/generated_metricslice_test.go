// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

// Code generated by "pdata/internal/cmd/pdatagen/main.go". DO NOT EDIT.
// To regenerate this file run "make genpdata".

package pmetric

import (
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"

	"go.opentelemetry.io/collector/pdata/internal"
	otlpmetrics "go.opentelemetry.io/collector/pdata/internal/data/protogen/metrics/v1"
)

func TestMetricSlice(t *testing.T) {
	es := NewMetricSlice()
	assert.Equal(t, 0, es.Len())
	state := internal.StateMutable
	es = newMetricSlice(&[]*otlpmetrics.Metric{}, &state)
	assert.Equal(t, 0, es.Len())

	emptyVal := NewMetric()
	testVal := generateTestMetric()
	for i := 0; i < 7; i++ {
		el := es.AppendEmpty()
		assert.Equal(t, emptyVal, es.At(i))
		fillTestMetric(el)
		assert.Equal(t, testVal, es.At(i))
	}
	assert.Equal(t, 7, es.Len())
}

func TestMetricSliceReadOnly(t *testing.T) {
	sharedState := internal.StateReadOnly
	es := newMetricSlice(&[]*otlpmetrics.Metric{}, &sharedState)
	assert.Equal(t, 0, es.Len())
	assert.Panics(t, func() { es.AppendEmpty() })
	assert.Panics(t, func() { es.EnsureCapacity(2) })
	es2 := NewMetricSlice()
	es.CopyTo(es2)
	assert.Panics(t, func() { es2.CopyTo(es) })
	assert.Panics(t, func() { es.MoveAndAppendTo(es2) })
	assert.Panics(t, func() { es2.MoveAndAppendTo(es) })
}

func TestMetricSlice_CopyTo(t *testing.T) {
	dest := NewMetricSlice()
	// Test CopyTo to empty
	NewMetricSlice().CopyTo(dest)
	assert.Equal(t, NewMetricSlice(), dest)

	// Test CopyTo larger slice
	generateTestMetricSlice().CopyTo(dest)
	assert.Equal(t, generateTestMetricSlice(), dest)

	// Test CopyTo same size slice
	generateTestMetricSlice().CopyTo(dest)
	assert.Equal(t, generateTestMetricSlice(), dest)
}

func TestMetricSlice_EnsureCapacity(t *testing.T) {
	es := generateTestMetricSlice()

	// Test ensure smaller capacity.
	const ensureSmallLen = 4
	es.EnsureCapacity(ensureSmallLen)
	assert.Less(t, ensureSmallLen, es.Len())
	assert.Equal(t, es.Len(), cap(*es.orig))
	assert.Equal(t, generateTestMetricSlice(), es)

	// Test ensure larger capacity
	const ensureLargeLen = 9
	es.EnsureCapacity(ensureLargeLen)
	assert.Less(t, generateTestMetricSlice().Len(), ensureLargeLen)
	assert.Equal(t, ensureLargeLen, cap(*es.orig))
	assert.Equal(t, generateTestMetricSlice(), es)
}

func TestMetricSlice_MoveAndAppendTo(t *testing.T) {
	// Test MoveAndAppendTo to empty
	expectedSlice := generateTestMetricSlice()
	dest := NewMetricSlice()
	src := generateTestMetricSlice()
	src.MoveAndAppendTo(dest)
	assert.Equal(t, generateTestMetricSlice(), dest)
	assert.Equal(t, 0, src.Len())
	assert.Equal(t, expectedSlice.Len(), dest.Len())

	// Test MoveAndAppendTo empty slice
	src.MoveAndAppendTo(dest)
	assert.Equal(t, generateTestMetricSlice(), dest)
	assert.Equal(t, 0, src.Len())
	assert.Equal(t, expectedSlice.Len(), dest.Len())

	// Test MoveAndAppendTo not empty slice
	generateTestMetricSlice().MoveAndAppendTo(dest)
	assert.Equal(t, 2*expectedSlice.Len(), dest.Len())
	for i := 0; i < expectedSlice.Len(); i++ {
		assert.Equal(t, expectedSlice.At(i), dest.At(i))
		assert.Equal(t, expectedSlice.At(i), dest.At(i+expectedSlice.Len()))
	}
}

func TestMetricSlice_RemoveIf(t *testing.T) {
	// Test RemoveIf on empty slice
	emptySlice := NewMetricSlice()
	emptySlice.RemoveIf(func(el Metric) bool {
		t.Fail()
		return false
	})

	// Test RemoveIf
	filtered := generateTestMetricSlice()
	pos := 0
	filtered.RemoveIf(func(el Metric) bool {
		pos++
		return pos%3 == 0
	})
	assert.Equal(t, 5, filtered.Len())
}

func TestMetricSlice_Sort(t *testing.T) {
	es := generateTestMetricSlice()
	es.Sort(func(a, b Metric) bool {
		return uintptr(unsafe.Pointer(a.orig)) < uintptr(unsafe.Pointer(b.orig))
	})
	for i := 1; i < es.Len(); i++ {
		assert.True(t, uintptr(unsafe.Pointer(es.At(i-1).orig)) < uintptr(unsafe.Pointer(es.At(i).orig)))
	}
	es.Sort(func(a, b Metric) bool {
		return uintptr(unsafe.Pointer(a.orig)) > uintptr(unsafe.Pointer(b.orig))
	})
	for i := 1; i < es.Len(); i++ {
		assert.True(t, uintptr(unsafe.Pointer(es.At(i-1).orig)) > uintptr(unsafe.Pointer(es.At(i).orig)))
	}
}

func TestMetricSlice_ForEach(t *testing.T) {
	// Test ForEach on empty slice
	emptySlice := NewMetricSlice()
	emptySlice.ForEach(func(el Metric) {
		t.Fail()
	})

	// Test ForEach
	slice := generateTestMetricSlice()
	pos := 0
	slice.ForEach(func(el Metric) {
		pos++
	})
	assert.Equal(t, 7, slice.Len())
}

func generateTestMetricSlice() MetricSlice {
	es := NewMetricSlice()
	fillTestMetricSlice(es)
	return es
}

func fillTestMetricSlice(es MetricSlice) {
	*es.orig = make([]*otlpmetrics.Metric, 7)
	for i := 0; i < 7; i++ {
		(*es.orig)[i] = &otlpmetrics.Metric{}
		fillTestMetric(newMetric((*es.orig)[i], es.state))
	}
}
