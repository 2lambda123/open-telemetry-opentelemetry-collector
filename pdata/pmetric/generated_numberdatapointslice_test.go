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

func TestNumberDataPointSlice(t *testing.T) {
	es := NewNumberDataPointSlice()
	assert.Equal(t, 0, es.Len())
	state := internal.StateMutable
	es = newNumberDataPointSlice(&[]*otlpmetrics.NumberDataPoint{}, &state)
	assert.Equal(t, 0, es.Len())

	emptyVal := NewNumberDataPoint()
	testVal := generateTestNumberDataPoint()
	for i := 0; i < 7; i++ {
		el := es.AppendEmpty()
		assert.Equal(t, emptyVal, es.At(i))
		fillTestNumberDataPoint(el)
		assert.Equal(t, testVal, es.At(i))
	}
	assert.Equal(t, 7, es.Len())
}

func TestNumberDataPointSliceReadOnly(t *testing.T) {
	sharedState := internal.StateReadOnly
	es := newNumberDataPointSlice(&[]*otlpmetrics.NumberDataPoint{}, &sharedState)
	assert.Equal(t, 0, es.Len())
	assert.Panics(t, func() { es.AppendEmpty() })
	assert.Panics(t, func() { es.EnsureCapacity(2) })
	es2 := NewNumberDataPointSlice()
	es.CopyTo(es2)
	assert.Panics(t, func() { es2.CopyTo(es) })
	assert.Panics(t, func() { es.MoveAndAppendTo(es2) })
	assert.Panics(t, func() { es2.MoveAndAppendTo(es) })
}

func TestNumberDataPointSlice_CopyTo(t *testing.T) {
	dest := NewNumberDataPointSlice()
	// Test CopyTo to empty
	NewNumberDataPointSlice().CopyTo(dest)
	assert.Equal(t, NewNumberDataPointSlice(), dest)

	// Test CopyTo larger slice
	generateTestNumberDataPointSlice().CopyTo(dest)
	assert.Equal(t, generateTestNumberDataPointSlice(), dest)

	// Test CopyTo same size slice
	generateTestNumberDataPointSlice().CopyTo(dest)
	assert.Equal(t, generateTestNumberDataPointSlice(), dest)
}

func TestNumberDataPointSlice_EnsureCapacity(t *testing.T) {
	es := generateTestNumberDataPointSlice()

	// Test ensure smaller capacity.
	const ensureSmallLen = 4
	es.EnsureCapacity(ensureSmallLen)
	assert.Less(t, ensureSmallLen, es.Len())
	assert.Equal(t, es.Len(), cap(*es.orig))
	assert.Equal(t, generateTestNumberDataPointSlice(), es)

	// Test ensure larger capacity
	const ensureLargeLen = 9
	es.EnsureCapacity(ensureLargeLen)
	assert.Less(t, generateTestNumberDataPointSlice().Len(), ensureLargeLen)
	assert.Equal(t, ensureLargeLen, cap(*es.orig))
	assert.Equal(t, generateTestNumberDataPointSlice(), es)
}

func TestNumberDataPointSlice_MoveAndAppendTo(t *testing.T) {
	// Test MoveAndAppendTo to empty
	expectedSlice := generateTestNumberDataPointSlice()
	dest := NewNumberDataPointSlice()
	src := generateTestNumberDataPointSlice()
	src.MoveAndAppendTo(dest)
	assert.Equal(t, generateTestNumberDataPointSlice(), dest)
	assert.Equal(t, 0, src.Len())
	assert.Equal(t, expectedSlice.Len(), dest.Len())

	// Test MoveAndAppendTo empty slice
	src.MoveAndAppendTo(dest)
	assert.Equal(t, generateTestNumberDataPointSlice(), dest)
	assert.Equal(t, 0, src.Len())
	assert.Equal(t, expectedSlice.Len(), dest.Len())

	// Test MoveAndAppendTo not empty slice
	generateTestNumberDataPointSlice().MoveAndAppendTo(dest)
	assert.Equal(t, 2*expectedSlice.Len(), dest.Len())
	for i := 0; i < expectedSlice.Len(); i++ {
		assert.Equal(t, expectedSlice.At(i), dest.At(i))
		assert.Equal(t, expectedSlice.At(i), dest.At(i+expectedSlice.Len()))
	}
}

func TestNumberDataPointSlice_RemoveIf(t *testing.T) {
	// Test RemoveIf on empty slice
	emptySlice := NewNumberDataPointSlice()
	emptySlice.RemoveIf(func(el NumberDataPoint) bool {
		t.Fail()
		return false
	})

	// Test RemoveIf
	filtered := generateTestNumberDataPointSlice()
	pos := 0
	filtered.RemoveIf(func(el NumberDataPoint) bool {
		pos++
		return pos%3 == 0
	})
	assert.Equal(t, 5, filtered.Len())
}

func TestNumberDataPointSlice_Sort(t *testing.T) {
	es := generateTestNumberDataPointSlice()
	es.Sort(func(a, b NumberDataPoint) bool {
		return uintptr(unsafe.Pointer(a.orig)) < uintptr(unsafe.Pointer(b.orig))
	})
	for i := 1; i < es.Len(); i++ {
		assert.True(t, uintptr(unsafe.Pointer(es.At(i-1).orig)) < uintptr(unsafe.Pointer(es.At(i).orig)))
	}
	es.Sort(func(a, b NumberDataPoint) bool {
		return uintptr(unsafe.Pointer(a.orig)) > uintptr(unsafe.Pointer(b.orig))
	})
	for i := 1; i < es.Len(); i++ {
		assert.True(t, uintptr(unsafe.Pointer(es.At(i-1).orig)) > uintptr(unsafe.Pointer(es.At(i).orig)))
	}
}

func TestNumberDataPointSlice_ForEach(t *testing.T) {
	// Test ForEach on empty slice
	emptySlice := NewNumberDataPointSlice()
	emptySlice.ForEach(func(el NumberDataPoint) {
		t.Fail()
	})

	// Test ForEach
	slice := generateTestNumberDataPointSlice()
	pos := 0
	slice.ForEach(func(el NumberDataPoint) {
		pos++
	})
	assert.Equal(t, 7, slice.Len())
}

func generateTestNumberDataPointSlice() NumberDataPointSlice {
	es := NewNumberDataPointSlice()
	fillTestNumberDataPointSlice(es)
	return es
}

func fillTestNumberDataPointSlice(es NumberDataPointSlice) {
	*es.orig = make([]*otlpmetrics.NumberDataPoint, 7)
	for i := 0; i < 7; i++ {
		(*es.orig)[i] = &otlpmetrics.NumberDataPoint{}
		fillTestNumberDataPoint(newNumberDataPoint((*es.orig)[i], es.state))
	}
}
