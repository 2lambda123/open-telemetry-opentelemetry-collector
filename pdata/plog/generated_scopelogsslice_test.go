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

package plog

import (
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"

	otlplogs "go.opentelemetry.io/collector/pdata/internal/data/protogen/logs/v1"
)

func TestScopeLogsSlice(t *testing.T) {
	es := NewMutableScopeLogsSlice()
	assert.Equal(t, 0, es.Len())
	es = newMutableScopeLogsSliceFromOrig(&[]*otlplogs.ScopeLogs{})
	assert.Equal(t, 0, es.Len())

	emptyVal := NewMutableScopeLogs()
	testVal := generateTestScopeLogs()
	for i := 0; i < 7; i++ {
		el := es.AppendEmpty()
		assert.Equal(t, emptyVal, es.At(i))
		fillTestScopeLogs(el)
		assert.Equal(t, testVal, es.At(i))
	}
	assert.Equal(t, 7, es.Len())
}

func TestScopeLogsSlice_CopyTo(t *testing.T) {
	dest := NewMutableScopeLogsSlice()
	// Test CopyTo to empty
	NewMutableScopeLogsSlice().CopyTo(dest)
	assert.Equal(t, NewMutableScopeLogsSlice(), dest)

	// Test CopyTo larger slice
	generateTestScopeLogsSlice().CopyTo(dest)
	assert.Equal(t, generateTestScopeLogsSlice(), dest)

	// Test CopyTo same size slice
	generateTestScopeLogsSlice().CopyTo(dest)
	assert.Equal(t, generateTestScopeLogsSlice(), dest)
}

func TestScopeLogsSlice_EnsureCapacity(t *testing.T) {
	es := generateTestScopeLogsSlice()

	// Test ensure smaller capacity.
	const ensureSmallLen = 4
	es.EnsureCapacity(ensureSmallLen)
	assert.Less(t, ensureSmallLen, es.Len())
	assert.Equal(t, es.Len(), cap(*es.orig))
	assert.Equal(t, generateTestScopeLogsSlice(), es)

	// Test ensure larger capacity
	const ensureLargeLen = 9
	es.EnsureCapacity(ensureLargeLen)
	assert.Less(t, generateTestScopeLogsSlice().Len(), ensureLargeLen)
	assert.Equal(t, ensureLargeLen, cap(*es.orig))
	assert.Equal(t, generateTestScopeLogsSlice(), es)
}

func TestScopeLogsSlice_MoveAndAppendTo(t *testing.T) {
	// Test MoveAndAppendTo to empty
	expectedSlice := generateTestScopeLogsSlice()
	dest := NewMutableScopeLogsSlice()
	src := generateTestScopeLogsSlice()
	src.MoveAndAppendTo(dest)
	assert.Equal(t, generateTestScopeLogsSlice(), dest)
	assert.Equal(t, 0, src.Len())
	assert.Equal(t, expectedSlice.Len(), dest.Len())

	// Test MoveAndAppendTo empty slice
	src.MoveAndAppendTo(dest)
	assert.Equal(t, generateTestScopeLogsSlice(), dest)
	assert.Equal(t, 0, src.Len())
	assert.Equal(t, expectedSlice.Len(), dest.Len())

	// Test MoveAndAppendTo not empty slice
	generateTestScopeLogsSlice().MoveAndAppendTo(dest)
	assert.Equal(t, 2*expectedSlice.Len(), dest.Len())
	for i := 0; i < expectedSlice.Len(); i++ {
		assert.Equal(t, expectedSlice.At(i), dest.At(i))
		assert.Equal(t, expectedSlice.At(i), dest.At(i+expectedSlice.Len()))
	}
}

func TestScopeLogsSlice_RemoveIf(t *testing.T) {
	// Test RemoveIf on empty slice
	emptySlice := NewMutableScopeLogsSlice()
	emptySlice.RemoveIf(func(el MutableScopeLogs) bool {
		t.Fail()
		return false
	})

	// Test RemoveIf
	filtered := generateTestScopeLogsSlice()
	pos := 0
	filtered.RemoveIf(func(el MutableScopeLogs) bool {
		pos++
		return pos%3 == 0
	})
	assert.Equal(t, 5, filtered.Len())
}

func TestScopeLogsSlice_Sort(t *testing.T) {
	es := generateTestScopeLogsSlice()
	es.Sort(func(a, b MutableScopeLogs) bool {
		return uintptr(unsafe.Pointer(a.orig)) < uintptr(unsafe.Pointer(b.orig))
	})
	for i := 1; i < es.Len(); i++ {
		assert.True(t, uintptr(unsafe.Pointer(es.At(i-1).orig)) < uintptr(unsafe.Pointer(es.At(i).orig)))
	}
	es.Sort(func(a, b MutableScopeLogs) bool {
		return uintptr(unsafe.Pointer(a.orig)) > uintptr(unsafe.Pointer(b.orig))
	})
	for i := 1; i < es.Len(); i++ {
		assert.True(t, uintptr(unsafe.Pointer(es.At(i-1).orig)) > uintptr(unsafe.Pointer(es.At(i).orig)))
	}
}
