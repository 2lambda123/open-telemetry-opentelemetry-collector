// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

// Code generated by "pdata/internal/cmd/pdatagen/main.go". DO NOT EDIT.
// To regenerate this file run "make genpdata".

package pprofile

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"go.opentelemetry.io/collector/pdata/internal"
	"go.opentelemetry.io/collector/pdata/pcommon"
)

func TestScopeProfiles_MoveTo(t *testing.T) {
	ms := generateTestScopeProfiles()
	dest := NewScopeProfiles()
	ms.MoveTo(dest)
	assert.Equal(t, NewScopeProfiles(), ms)
	assert.Equal(t, generateTestScopeProfiles(), dest)
}

func TestScopeProfiles_CopyTo(t *testing.T) {
	ms := NewScopeProfiles()
	orig := NewScopeProfiles()
	orig.CopyTo(ms)
	assert.Equal(t, orig, ms)
	orig = generateTestScopeProfiles()
	orig.CopyTo(ms)
	assert.Equal(t, orig, ms)
}

func TestScopeProfiles_Scope(t *testing.T) {
	ms := NewScopeProfiles()
	internal.FillTestInstrumentationScope(internal.InstrumentationScope(ms.Scope()))
	assert.Equal(t, pcommon.InstrumentationScope(internal.GenerateTestInstrumentationScope()), ms.Scope())
}

func TestScopeProfiles_SchemaUrl(t *testing.T) {
	ms := NewScopeProfiles()
	assert.Equal(t, "", ms.SchemaUrl())
	ms.SetSchemaUrl("https://opentelemetry.io/schemas/1.5.0")
	assert.Equal(t, "https://opentelemetry.io/schemas/1.5.0", ms.SchemaUrl())
}

func TestScopeProfiles_Profiles(t *testing.T) {
	ms := NewScopeProfiles()
	assert.Equal(t, NewProfileSlice(), ms.Profiles())
	fillTestProfileSlice(ms.Profiles())
	assert.Equal(t, generateTestProfileSlice(), ms.Profiles())
}

func generateTestScopeProfiles() ScopeProfiles {
	tv := NewScopeProfiles()
	fillTestScopeProfiles(tv)
	return tv
}

func fillTestScopeProfiles(tv ScopeProfiles) {
	internal.FillTestInstrumentationScope(internal.NewInstrumentationScope(&tv.orig.Scope))
	tv.orig.SchemaUrl = "https://opentelemetry.io/schemas/1.5.0"
	fillTestProfileSlice(newProfileSlice(&tv.orig.Profiles))
}