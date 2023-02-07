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

package ptrace

import (
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"

	otlptrace "go.opentelemetry.io/collector/pdata/internal/data/protogen/trace/v1"
)

func TestResourceSpansSlice(t *testing.T) {
	es := NewMutableResourceSpansSlice()
	assert.Equal(t, 0, es.Len())
	es = newMutableResourceSpansSliceFromOrig(&[]*otlptrace.ResourceSpans{})
	assert.Equal(t, 0, es.Len())

	emptyVal := NewMutableResourceSpans()
	testVal := generateTestResourceSpans()
	for i := 0; i < 7; i++ {
		el := es.AppendEmpty()
		assert.Equal(t, emptyVal, es.At(i))
		fillTestResourceSpans(el)
		assert.Equal(t, testVal, es.At(i))
	}
	assert.Equal(t, 7, es.Len())
}

func TestResourceSpansSlice_CopyTo(t *testing.T) {
	dest := NewMutableResourceSpansSlice()
	// Test CopyTo to empty
	NewMutableResourceSpansSlice().CopyTo(dest)
	assert.Equal(t, NewMutableResourceSpansSlice(), dest)

	// Test CopyTo larger slice
	generateTestResourceSpansSlice().CopyTo(dest)
	assert.Equal(t, generateTestResourceSpansSlice(), dest)

	// Test CopyTo same size slice
	generateTestResourceSpansSlice().CopyTo(dest)
	assert.Equal(t, generateTestResourceSpansSlice(), dest)
}

func TestResourceSpansSlice_EnsureCapacity(t *testing.T) {
	es := generateTestResourceSpansSlice()

	// Test ensure smaller capacity.
	const ensureSmallLen = 4
	es.EnsureCapacity(ensureSmallLen)
	assert.Less(t, ensureSmallLen, es.Len())
	assert.Equal(t, es.Len(), cap(*es.orig))
	assert.Equal(t, generateTestResourceSpansSlice(), es)

	// Test ensure larger capacity
	const ensureLargeLen = 9
	es.EnsureCapacity(ensureLargeLen)
	assert.Less(t, generateTestResourceSpansSlice().Len(), ensureLargeLen)
	assert.Equal(t, ensureLargeLen, cap(*es.orig))
	assert.Equal(t, generateTestResourceSpansSlice(), es)
}

func TestResourceSpansSlice_MoveAndAppendTo(t *testing.T) {
	// Test MoveAndAppendTo to empty
	expectedSlice := generateTestResourceSpansSlice()
	dest := NewMutableResourceSpansSlice()
	src := generateTestResourceSpansSlice()
	src.MoveAndAppendTo(dest)
	assert.Equal(t, generateTestResourceSpansSlice(), dest)
	assert.Equal(t, 0, src.Len())
	assert.Equal(t, expectedSlice.Len(), dest.Len())

	// Test MoveAndAppendTo empty slice
	src.MoveAndAppendTo(dest)
	assert.Equal(t, generateTestResourceSpansSlice(), dest)
	assert.Equal(t, 0, src.Len())
	assert.Equal(t, expectedSlice.Len(), dest.Len())

	// Test MoveAndAppendTo not empty slice
	generateTestResourceSpansSlice().MoveAndAppendTo(dest)
	assert.Equal(t, 2*expectedSlice.Len(), dest.Len())
	for i := 0; i < expectedSlice.Len(); i++ {
		assert.Equal(t, expectedSlice.At(i), dest.At(i))
		assert.Equal(t, expectedSlice.At(i), dest.At(i+expectedSlice.Len()))
	}
}

func TestResourceSpansSlice_RemoveIf(t *testing.T) {
	// Test RemoveIf on empty slice
	emptySlice := NewMutableResourceSpansSlice()
	emptySlice.RemoveIf(func(el MutableResourceSpans) bool {
		t.Fail()
		return false
	})

	// Test RemoveIf
	filtered := generateTestResourceSpansSlice()
	pos := 0
	filtered.RemoveIf(func(el MutableResourceSpans) bool {
		pos++
		return pos%3 == 0
	})
	assert.Equal(t, 5, filtered.Len())
}

func TestResourceSpansSlice_Sort(t *testing.T) {
	es := generateTestResourceSpansSlice()
	es.Sort(func(a, b MutableResourceSpans) bool {
		return uintptr(unsafe.Pointer(a.orig)) < uintptr(unsafe.Pointer(b.orig))
	})
	for i := 1; i < es.Len(); i++ {
		assert.True(t, uintptr(unsafe.Pointer(es.At(i-1).orig)) < uintptr(unsafe.Pointer(es.At(i).orig)))
	}
	es.Sort(func(a, b MutableResourceSpans) bool {
		return uintptr(unsafe.Pointer(a.orig)) > uintptr(unsafe.Pointer(b.orig))
	})
	for i := 1; i < es.Len(); i++ {
		assert.True(t, uintptr(unsafe.Pointer(es.At(i-1).orig)) > uintptr(unsafe.Pointer(es.At(i).orig)))
	}
}
