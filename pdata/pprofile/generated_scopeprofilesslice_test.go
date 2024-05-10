// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

// Code generated by "pdata/internal/cmd/pdatagen/main.go". DO NOT EDIT.
// To regenerate this file run "make genpdata".

package pprofile

import (
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"

	otlpprofiles "go.opentelemetry.io/collector/pdata/internal/data/protogen/profiles/v1experimental"
)

func TestScopeProfilesSlice(t *testing.T) {
	es := NewScopeProfilesSlice()
	assert.Equal(t, 0, es.Len())
	es = newScopeProfilesSlice(&[]*otlpprofiles.ScopeProfiles{})
	assert.Equal(t, 0, es.Len())

	emptyVal := NewScopeProfiles()
	testVal := generateTestScopeProfiles()
	for i := 0; i < 7; i++ {
		el := es.AppendEmpty()
		assert.Equal(t, emptyVal, es.At(i))
		fillTestScopeProfiles(el)
		assert.Equal(t, testVal, es.At(i))
	}
	assert.Equal(t, 7, es.Len())
}

func TestScopeProfilesSlice_CopyTo(t *testing.T) {
	dest := NewScopeProfilesSlice()
	// Test CopyTo to empty
	NewScopeProfilesSlice().CopyTo(dest)
	assert.Equal(t, NewScopeProfilesSlice(), dest)

	// Test CopyTo larger slice
	generateTestScopeProfilesSlice().CopyTo(dest)
	assert.Equal(t, generateTestScopeProfilesSlice(), dest)

	// Test CopyTo same size slice
	generateTestScopeProfilesSlice().CopyTo(dest)
	assert.Equal(t, generateTestScopeProfilesSlice(), dest)
}

func TestScopeProfilesSlice_EnsureCapacity(t *testing.T) {
	es := generateTestScopeProfilesSlice()

	// Test ensure smaller capacity.
	const ensureSmallLen = 4
	es.EnsureCapacity(ensureSmallLen)
	assert.Less(t, ensureSmallLen, es.Len())
	assert.Equal(t, es.Len(), cap(*es.orig))
	assert.Equal(t, generateTestScopeProfilesSlice(), es)

	// Test ensure larger capacity
	const ensureLargeLen = 9
	es.EnsureCapacity(ensureLargeLen)
	assert.Less(t, generateTestScopeProfilesSlice().Len(), ensureLargeLen)
	assert.Equal(t, ensureLargeLen, cap(*es.orig))
	assert.Equal(t, generateTestScopeProfilesSlice(), es)
}

func TestScopeProfilesSlice_MoveAndAppendTo(t *testing.T) {
	// Test MoveAndAppendTo to empty
	expectedSlice := generateTestScopeProfilesSlice()
	dest := NewScopeProfilesSlice()
	src := generateTestScopeProfilesSlice()
	src.MoveAndAppendTo(dest)
	assert.Equal(t, generateTestScopeProfilesSlice(), dest)
	assert.Equal(t, 0, src.Len())
	assert.Equal(t, expectedSlice.Len(), dest.Len())

	// Test MoveAndAppendTo empty slice
	src.MoveAndAppendTo(dest)
	assert.Equal(t, generateTestScopeProfilesSlice(), dest)
	assert.Equal(t, 0, src.Len())
	assert.Equal(t, expectedSlice.Len(), dest.Len())

	// Test MoveAndAppendTo not empty slice
	generateTestScopeProfilesSlice().MoveAndAppendTo(dest)
	assert.Equal(t, 2*expectedSlice.Len(), dest.Len())
	for i := 0; i < expectedSlice.Len(); i++ {
		assert.Equal(t, expectedSlice.At(i), dest.At(i))
		assert.Equal(t, expectedSlice.At(i), dest.At(i+expectedSlice.Len()))
	}
}

func TestScopeProfilesSlice_RemoveIf(t *testing.T) {
	// Test RemoveIf on empty slice
	emptySlice := NewScopeProfilesSlice()
	emptySlice.RemoveIf(func(el ScopeProfiles) bool {
		t.Fail()
		return false
	})

	// Test RemoveIf
	filtered := generateTestScopeProfilesSlice()
	pos := 0
	filtered.RemoveIf(func(el ScopeProfiles) bool {
		pos++
		return pos%3 == 0
	})
	assert.Equal(t, 5, filtered.Len())
}

func TestScopeProfilesSlice_Sort(t *testing.T) {
	es := generateTestScopeProfilesSlice()
	es.Sort(func(a, b ScopeProfiles) bool {
		return uintptr(unsafe.Pointer(a.orig)) < uintptr(unsafe.Pointer(b.orig))
	})
	for i := 1; i < es.Len(); i++ {
		assert.True(t, uintptr(unsafe.Pointer(es.At(i-1).orig)) < uintptr(unsafe.Pointer(es.At(i).orig)))
	}
	es.Sort(func(a, b ScopeProfiles) bool {
		return uintptr(unsafe.Pointer(a.orig)) > uintptr(unsafe.Pointer(b.orig))
	})
	for i := 1; i < es.Len(); i++ {
		assert.True(t, uintptr(unsafe.Pointer(es.At(i-1).orig)) > uintptr(unsafe.Pointer(es.At(i).orig)))
	}
}

func generateTestScopeProfilesSlice() ScopeProfilesSlice {
	es := NewScopeProfilesSlice()
	fillTestScopeProfilesSlice(es)
	return es
}

func fillTestScopeProfilesSlice(es ScopeProfilesSlice) {
	*es.orig = make([]*otlpprofiles.ScopeProfiles, 7)
	for i := 0; i < 7; i++ {
		(*es.orig)[i] = &otlpprofiles.ScopeProfiles{}
		fillTestScopeProfiles(newScopeProfiles((*es.orig)[i]))
	}
}