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

	otlpmetrics "go.opentelemetry.io/collector/pdata/internal/data/protogen/metrics/v1"
)

func TestScopeMetricsSlice(t *testing.T) {
	es := NewScopeMetricsSlice()
	assert.Equal(t, 0, es.Len())
	es = newScopeMetricsSlice(&[]*otlpmetrics.ScopeMetrics{})
	assert.Equal(t, 0, es.Len())

	emptyVal := NewScopeMetrics()
	testVal := generateTestScopeMetrics()
	for i := 0; i < 7; i++ {
		el := es.AppendEmpty()
		assert.Equal(t, emptyVal, es.At(i))
		fillTestScopeMetrics(el)
		assert.Equal(t, testVal, es.At(i))
	}
	assert.Equal(t, 7, es.Len())
}

func TestScopeMetricsSlice_CopyTo(t *testing.T) {
	dest := NewScopeMetricsSlice()
	// Test CopyTo to empty
	NewScopeMetricsSlice().CopyTo(dest)
	assert.Equal(t, NewScopeMetricsSlice(), dest)

	// Test CopyTo larger slice
	generateTestScopeMetricsSlice().CopyTo(dest)
	assert.Equal(t, generateTestScopeMetricsSlice(), dest)

	// Test CopyTo same size slice
	generateTestScopeMetricsSlice().CopyTo(dest)
	assert.Equal(t, generateTestScopeMetricsSlice(), dest)
}

func TestScopeMetricsSlice_EnsureCapacity(t *testing.T) {
	es := generateTestScopeMetricsSlice()

	// Test ensure smaller capacity.
	const ensureSmallLen = 4
	es.EnsureCapacity(ensureSmallLen)
	assert.Less(t, ensureSmallLen, es.Len())
	assert.Equal(t, es.Len(), cap(*es.orig))
	assert.Equal(t, generateTestScopeMetricsSlice(), es)

	// Test ensure larger capacity
	const ensureLargeLen = 9
	es.EnsureCapacity(ensureLargeLen)
	assert.Less(t, generateTestScopeMetricsSlice().Len(), ensureLargeLen)
	assert.Equal(t, ensureLargeLen, cap(*es.orig))
	assert.Equal(t, generateTestScopeMetricsSlice(), es)
}

func TestScopeMetricsSlice_MoveAndAppendTo(t *testing.T) {
	// Test MoveAndAppendTo to empty
	expectedSlice := generateTestScopeMetricsSlice()
	dest := NewScopeMetricsSlice()
	src := generateTestScopeMetricsSlice()
	src.MoveAndAppendTo(dest)
	assert.Equal(t, generateTestScopeMetricsSlice(), dest)
	assert.Equal(t, 0, src.Len())
	assert.Equal(t, expectedSlice.Len(), dest.Len())

	// Test MoveAndAppendTo empty slice
	src.MoveAndAppendTo(dest)
	assert.Equal(t, generateTestScopeMetricsSlice(), dest)
	assert.Equal(t, 0, src.Len())
	assert.Equal(t, expectedSlice.Len(), dest.Len())

	// Test MoveAndAppendTo not empty slice
	generateTestScopeMetricsSlice().MoveAndAppendTo(dest)
	assert.Equal(t, 2*expectedSlice.Len(), dest.Len())
	for i := 0; i < expectedSlice.Len(); i++ {
		assert.Equal(t, expectedSlice.At(i), dest.At(i))
		assert.Equal(t, expectedSlice.At(i), dest.At(i+expectedSlice.Len()))
	}
}

func TestScopeMetricsSlice_RemoveIf(t *testing.T) {
	// Test RemoveIf on empty slice
	emptySlice := NewScopeMetricsSlice()
	emptySlice.RemoveIf(func(el ScopeMetrics) bool {
		t.Fail()
		return false
	})

	// Test RemoveIf
	filtered := generateTestScopeMetricsSlice()
	pos := 0
	filtered.RemoveIf(func(el ScopeMetrics) bool {
		pos++
		return pos%3 == 0
	})
	assert.Equal(t, 5, filtered.Len())
}

func generateTestScopeMetricsSlice() ScopeMetricsSlice {
	es := NewScopeMetricsSlice()
	fillTestScopeMetricsSlice(es)
	return es
}

func fillTestScopeMetricsSlice(es ScopeMetricsSlice) {
	*es.orig = make([]*otlpmetrics.ScopeMetrics, 7)
	for i := 0; i < 7; i++ {
		(*es.orig)[i] = &otlpmetrics.ScopeMetrics{}
		fillTestScopeMetrics(newScopeMetrics((*es.orig)[i]))
	}
}
