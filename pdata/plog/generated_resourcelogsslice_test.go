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

func TestResourceLogsSlice(t *testing.T) {
	es := NewResourceLogsSlice()
	assert.Equal(t, 0, es.Len())

	emptyVal := NewResourceLogs()
	testVal := generateTestResourceLogs()
	for i := 0; i < 7; i++ {
		el := es.AppendEmpty()
		assert.Equal(t, emptyVal.getOrig(), es.At(i).getOrig())
		fillTestResourceLogs(el)
		assert.Equal(t, testVal.getOrig(), es.At(i).getOrig())
	}
	assert.Equal(t, 7, es.Len())
}

func TestResourceLogsSlice_CopyTo(t *testing.T) {
	dest := NewResourceLogsSlice()
	// Test CopyTo to empty
	NewResourceLogsSlice().CopyTo(dest)
	assert.Equal(t, NewResourceLogsSlice(), dest)

	// Test CopyTo larger slice
	generateTestResourceLogsSlice().CopyTo(dest)
	assert.Equal(t, generateTestResourceLogsSlice(), dest)

	// Test CopyTo same size slice
	generateTestResourceLogsSlice().CopyTo(dest)
	assert.Equal(t, generateTestResourceLogsSlice(), dest)
}

func TestResourceLogsSlice_EnsureCapacity(t *testing.T) {
	es := generateTestResourceLogsSlice()

	// Test ensure smaller capacity.
	const ensureSmallLen = 4
	es.EnsureCapacity(ensureSmallLen)
	assert.Less(t, ensureSmallLen, es.Len())
	assert.Equal(t, es.Len(), cap(*es.getOrig()))
	assert.Equal(t, generateTestResourceLogsSlice(), es)

	// Test ensure larger capacity
	const ensureLargeLen = 9
	es.EnsureCapacity(ensureLargeLen)
	assert.Less(t, generateTestResourceLogsSlice().Len(), ensureLargeLen)
	assert.Equal(t, ensureLargeLen, cap(*es.getOrig()))
	assert.Equal(t, generateTestResourceLogsSlice(), es)
}

func TestResourceLogsSlice_MoveAndAppendTo(t *testing.T) {
	// Test MoveAndAppendTo to empty
	expectedSlice := generateTestResourceLogsSlice()
	dest := NewResourceLogsSlice()
	src := generateTestResourceLogsSlice()
	src.MoveAndAppendTo(dest)
	assert.Equal(t, generateTestResourceLogsSlice(), dest)
	assert.Equal(t, 0, src.Len())
	assert.Equal(t, expectedSlice.Len(), dest.Len())

	// Test MoveAndAppendTo empty slice
	src.MoveAndAppendTo(dest)
	assert.Equal(t, generateTestResourceLogsSlice(), dest)
	assert.Equal(t, 0, src.Len())
	assert.Equal(t, expectedSlice.Len(), dest.Len())

	// Test MoveAndAppendTo not empty slice
	generateTestResourceLogsSlice().MoveAndAppendTo(dest)
	assert.Equal(t, 2*expectedSlice.Len(), dest.Len())
	for i := 0; i < expectedSlice.Len(); i++ {
		assert.Equal(t, expectedSlice.At(i).getOrig(), dest.At(i).getOrig())
		assert.Equal(t, expectedSlice.At(i).getOrig(), dest.At(i+expectedSlice.Len()).getOrig())
	}
}

func TestResourceLogsSlice_RemoveIf(t *testing.T) {
	// Test RemoveIf on empty slice
	emptySlice := NewResourceLogsSlice()
	emptySlice.RemoveIf(func(el ResourceLogs) bool {
		t.Fail()
		return false
	})

	// Test RemoveIf
	filtered := generateTestResourceLogsSlice()
	pos := 0
	filtered.RemoveIf(func(el ResourceLogs) bool {
		pos++
		return pos%3 == 0
	})
	assert.Equal(t, 5, filtered.Len())
}

func TestResourceLogsSlice_Sort(t *testing.T) {
	es := generateTestResourceLogsSlice()
	es.Sort(func(a, b ResourceLogs) bool {
		return uintptr(unsafe.Pointer(a.getOrig())) < uintptr(unsafe.Pointer(b.getOrig()))
	})
	for i := 1; i < es.Len(); i++ {
		assert.True(t, uintptr(unsafe.Pointer(es.At(i-1).getOrig())) < uintptr(unsafe.Pointer(es.At(i).getOrig())))
	}
	es.Sort(func(a, b ResourceLogs) bool {
		return uintptr(unsafe.Pointer(a.getOrig())) > uintptr(unsafe.Pointer(b.getOrig()))
	})
	for i := 1; i < es.Len(); i++ {
		assert.True(t, uintptr(unsafe.Pointer(es.At(i-1).getOrig())) > uintptr(unsafe.Pointer(es.At(i).getOrig())))
	}
}

func generateTestResourceLogsSlice() ResourceLogsSlice {
	es := NewResourceLogsSlice()
	fillTestResourceLogsSlice(es)
	return es
}

func fillTestResourceLogsSlice(es ResourceLogsSlice) {
	*es.getOrig() = make([]*otlplogs.ResourceLogs, 7)
	for i := 0; i < 7; i++ {
		(*es.getOrig())[i] = &otlplogs.ResourceLogs{}
		fillTestResourceLogs(newResourceLogs((*es.getOrig())[i], es, i))
	}
}
