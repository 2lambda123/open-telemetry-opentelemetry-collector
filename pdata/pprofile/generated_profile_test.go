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
	"go.opentelemetry.io/collector/pdata/pcommon"
)

func TestProfile_MoveTo(t *testing.T) {
	ms := generateTestProfile()
	dest := NewProfile()
	ms.MoveTo(dest)
	assert.Equal(t, NewProfile(), ms)
	assert.Equal(t, generateTestProfile(), dest)
	sharedState := internal.StateReadOnly
	assert.Panics(t, func() { ms.MoveTo(newProfile(&otlpprofiles.Profile{}, &sharedState)) })
	assert.Panics(t, func() { newProfile(&otlpprofiles.Profile{}, &sharedState).MoveTo(dest) })
}

func TestProfile_CopyTo(t *testing.T) {
	ms := NewProfile()
	orig := NewProfile()
	orig.CopyTo(ms)
	assert.Equal(t, orig, ms)
	orig = generateTestProfile()
	orig.CopyTo(ms)
	assert.Equal(t, orig, ms)
	sharedState := internal.StateReadOnly
	assert.Panics(t, func() { ms.CopyTo(newProfile(&otlpprofiles.Profile{}, &sharedState)) })
}

func TestProfile_SampleType(t *testing.T) {
	ms := NewProfile()
	assert.Equal(t, NewValueTypes(), ms.SampleType())
	fillTestValueTypes(ms.SampleType())
	assert.Equal(t, generateTestValueTypes(), ms.SampleType())
}

func TestProfile_Sample(t *testing.T) {
	ms := NewProfile()
	assert.Equal(t, NewSamples(), ms.Sample())
	fillTestSamples(ms.Sample())
	assert.Equal(t, generateTestSamples(), ms.Sample())
}

func TestProfile_StartTime(t *testing.T) {
	ms := NewProfile()
	assert.Equal(t, pcommon.Timestamp(0), ms.StartTime())
	testValStartTime := pcommon.Timestamp(1234567890)
	ms.SetStartTime(testValStartTime)
	assert.Equal(t, testValStartTime, ms.StartTime())
}

func generateTestProfile() Profile {
	tv := NewProfile()
	fillTestProfile(tv)
	return tv
}

func fillTestProfile(tv Profile) {
	fillTestValueTypes(newValueTypes(&tv.orig.SampleType, tv.state))
	fillTestSamples(newSamples(&tv.orig.Sample, tv.state))
	tv.orig.TimeNanos = 1234567890
}
