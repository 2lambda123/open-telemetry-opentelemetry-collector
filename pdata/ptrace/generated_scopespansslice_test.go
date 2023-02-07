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

func TestScopeSpansSlice(t *testing.T) {
	es := NewMutableScopeSpansSlice()
	assert.Equal(t, 0, es.Len())
	es = newMutableScopeSpansSliceFromOrig(&[]*otlptrace.ScopeSpans{})
	assert.Equal(t, 0, es.Len())

	emptyVal := NewMutableScopeSpans()
	testVal := generateTestScopeSpans()
	for i := 0; i < 7; i++ {
		el := es.AppendEmpty()
		assert.Equal(t, emptyVal, es.At(i))
		fillTestScopeSpans(el)
		assert.Equal(t, testVal, es.At(i))
	}
	assert.Equal(t, 7, es.Len())
}

func TestScopeSpansSlice_CopyTo(t *testing.T) {
	dest := NewMutableScopeSpansSlice()
	// Test CopyTo to empty
	NewMutableScopeSpansSlice().CopyTo(dest)
	assert.Equal(t, NewMutableScopeSpansSlice(), dest)

	// Test CopyTo larger slice
	generateTestScopeSpansSlice().CopyTo(dest)
	assert.Equal(t, generateTestScopeSpansSlice(), dest)

	// Test CopyTo same size slice
	generateTestScopeSpansSlice().CopyTo(dest)
	assert.Equal(t, generateTestScopeSpansSlice(), dest)
}

func TestScopeSpansSlice_EnsureCapacity(t *testing.T) {
	es := generateTestScopeSpansSlice()

	// Test ensure smaller capacity.
	const ensureSmallLen = 4
	es.EnsureCapacity(ensureSmallLen)
	assert.Less(t, ensureSmallLen, es.Len())
	assert.Equal(t, es.Len(), cap(*es.orig))
	assert.Equal(t, generateTestScopeSpansSlice(), es)

	// Test ensure larger capacity
	const ensureLargeLen = 9
	es.EnsureCapacity(ensureLargeLen)
	assert.Less(t, generateTestScopeSpansSlice().Len(), ensureLargeLen)
	assert.Equal(t, ensureLargeLen, cap(*es.orig))
	assert.Equal(t, generateTestScopeSpansSlice(), es)
}

func TestScopeSpansSlice_MoveAndAppendTo(t *testing.T) {
	// Test MoveAndAppendTo to empty
	expectedSlice := generateTestScopeSpansSlice()
	dest := NewMutableScopeSpansSlice()
	src := generateTestScopeSpansSlice()
	src.MoveAndAppendTo(dest)
	assert.Equal(t, generateTestScopeSpansSlice(), dest)
	assert.Equal(t, 0, src.Len())
	assert.Equal(t, expectedSlice.Len(), dest.Len())

	// Test MoveAndAppendTo empty slice
	src.MoveAndAppendTo(dest)
	assert.Equal(t, generateTestScopeSpansSlice(), dest)
	assert.Equal(t, 0, src.Len())
	assert.Equal(t, expectedSlice.Len(), dest.Len())

	// Test MoveAndAppendTo not empty slice
	generateTestScopeSpansSlice().MoveAndAppendTo(dest)
	assert.Equal(t, 2*expectedSlice.Len(), dest.Len())
	for i := 0; i < expectedSlice.Len(); i++ {
		assert.Equal(t, expectedSlice.At(i), dest.At(i))
		assert.Equal(t, expectedSlice.At(i), dest.At(i+expectedSlice.Len()))
	}
}

func TestScopeSpansSlice_RemoveIf(t *testing.T) {
	// Test RemoveIf on empty slice
	emptySlice := NewMutableScopeSpansSlice()
	emptySlice.RemoveIf(func(el MutableScopeSpans) bool {
		t.Fail()
		return false
	})

	// Test RemoveIf
	filtered := generateTestScopeSpansSlice()
	pos := 0
	filtered.RemoveIf(func(el MutableScopeSpans) bool {
		pos++
		return pos%3 == 0
	})
	assert.Equal(t, 5, filtered.Len())
}

func TestScopeSpansSlice_Sort(t *testing.T) {
	es := generateTestScopeSpansSlice()
	es.Sort(func(a, b MutableScopeSpans) bool {
		return uintptr(unsafe.Pointer(a.orig)) < uintptr(unsafe.Pointer(b.orig))
	})
	for i := 1; i < es.Len(); i++ {
		assert.True(t, uintptr(unsafe.Pointer(es.At(i-1).orig)) < uintptr(unsafe.Pointer(es.At(i).orig)))
	}
	es.Sort(func(a, b MutableScopeSpans) bool {
		return uintptr(unsafe.Pointer(a.orig)) > uintptr(unsafe.Pointer(b.orig))
	})
	for i := 1; i < es.Len(); i++ {
		assert.True(t, uintptr(unsafe.Pointer(es.At(i-1).orig)) > uintptr(unsafe.Pointer(es.At(i).orig)))
	}
}
