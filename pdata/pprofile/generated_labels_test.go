// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

// Code generated by "pdata/internal/cmd/pdatagen/main.go". DO NOT EDIT.
// To regenerate this file run "make genpdata".

package pprofile

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"go.opentelemetry.io/collector/pdata/internal"
	otlpprofiles "go.opentelemetry.io/collector/pdata/internal/data/protogen/profiles/v1experimental"
)

func TestLabels(t *testing.T) {
	es := NewLabels()
	assert.Equal(t, 0, es.Len())
	state := internal.StateMutable
	es = newLabels(&[]otlpprofiles.Label{}, &state)
	assert.Equal(t, 0, es.Len())

	emptyVal := NewLabel()
	testVal := generateTestLabel()
	for i := 0; i < 7; i++ {
		el := es.AppendEmpty()
		assert.Equal(t, emptyVal, es.At(i))
		fillTestLabel(el)
		assert.Equal(t, testVal, es.At(i))
	}
	assert.Equal(t, 7, es.Len())
}

func TestLabelsReadOnly(t *testing.T) {
	sharedState := internal.StateReadOnly
	es := newLabels(&[]otlpprofiles.Label{}, &sharedState)
	assert.Equal(t, 0, es.Len())
	assert.Panics(t, func() { es.AppendEmpty() })
	assert.Panics(t, func() { es.EnsureCapacity(2) })
	es2 := NewLabels()
	es.CopyTo(es2)
	assert.Panics(t, func() { es2.CopyTo(es) })
	assert.Panics(t, func() { es.MoveAndAppendTo(es2) })
	assert.Panics(t, func() { es2.MoveAndAppendTo(es) })
}

func TestLabels_CopyTo(t *testing.T) {
	dest := NewLabels()
	// Test CopyTo to empty
	NewLabels().CopyTo(dest)
	assert.Equal(t, NewLabels(), dest)

	// Test CopyTo larger slice
	generateTestLabels().CopyTo(dest)
	assert.Equal(t, generateTestLabels(), dest)

	// Test CopyTo same size slice
	generateTestLabels().CopyTo(dest)
	assert.Equal(t, generateTestLabels(), dest)
}

func TestLabels_EnsureCapacity(t *testing.T) {
	es := generateTestLabels()

	// Test ensure smaller capacity.
	const ensureSmallLen = 4
	es.EnsureCapacity(ensureSmallLen)
	assert.Less(t, ensureSmallLen, es.Len())
	assert.Equal(t, es.Len(), cap(*es.orig))
	assert.Equal(t, generateTestLabels(), es)

	// Test ensure larger capacity
	const ensureLargeLen = 9
	es.EnsureCapacity(ensureLargeLen)
	assert.Less(t, generateTestLabels().Len(), ensureLargeLen)
	assert.Equal(t, ensureLargeLen, cap(*es.orig))
	assert.Equal(t, generateTestLabels(), es)
}

func TestLabels_MoveAndAppendTo(t *testing.T) {
	// Test MoveAndAppendTo to empty
	expectedSlice := generateTestLabels()
	dest := NewLabels()
	src := generateTestLabels()
	src.MoveAndAppendTo(dest)
	assert.Equal(t, generateTestLabels(), dest)
	assert.Equal(t, 0, src.Len())
	assert.Equal(t, expectedSlice.Len(), dest.Len())

	// Test MoveAndAppendTo empty slice
	src.MoveAndAppendTo(dest)
	assert.Equal(t, generateTestLabels(), dest)
	assert.Equal(t, 0, src.Len())
	assert.Equal(t, expectedSlice.Len(), dest.Len())

	// Test MoveAndAppendTo not empty slice
	generateTestLabels().MoveAndAppendTo(dest)
	assert.Equal(t, 2*expectedSlice.Len(), dest.Len())
	for i := 0; i < expectedSlice.Len(); i++ {
		assert.Equal(t, expectedSlice.At(i), dest.At(i))
		assert.Equal(t, expectedSlice.At(i), dest.At(i+expectedSlice.Len()))
	}
}

func TestLabels_RemoveIf(t *testing.T) {
	// Test RemoveIf on empty slice
	emptySlice := NewLabels()
	emptySlice.RemoveIf(func(el Label) bool {
		t.Fail()
		return false
	})

	// Test RemoveIf
	filtered := generateTestLabels()
	pos := 0
	filtered.RemoveIf(func(el Label) bool {
		pos++
		return pos%3 == 0
	})
	assert.Equal(t, 5, filtered.Len())
}

func generateTestLabels() Labels {
	es := NewLabels()
	fillTestLabels(es)
	return es
}

func fillTestLabels(es Labels) {
	*es.orig = make([]otlpprofiles.Label, 7)
	for i := 0; i < 7; i++ {
		(*es.orig)[i] = otlpprofiles.Label{}
		fillTestLabel(newLabel(&(*es.orig)[i], es.state))
	}
}
