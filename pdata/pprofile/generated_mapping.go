// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

// Code generated by "pdata/internal/cmd/pdatagen/main.go". DO NOT EDIT.
// To regenerate this file run "make genpdata".

package pprofile

import (
	"go.opentelemetry.io/collector/pdata/internal"
	otlpprofiles "go.opentelemetry.io/collector/pdata/internal/data/protogen/profiles/v1experimental"
	"go.opentelemetry.io/collector/pdata/pcommon"
)

// Sample represents each record value encountered within a profiled program.
//
// This is a reference type, if passed by value and callee modifies it the
// caller will see the modification.
//
// Must use NewMapping function to create new instances.
// Important: zero-initialized instance is not valid for use.
type Mapping struct {
	orig  *otlpprofiles.Mapping
	state *internal.State
}

func newMapping(orig *otlpprofiles.Mapping, state *internal.State) Mapping {
	return Mapping{orig: orig, state: state}
}

// NewMapping creates a new empty Mapping.
//
// This must be used only in testing code. Users should use "AppendEmpty" when part of a Slice,
// OR directly access the member if this is embedded in another struct.
func NewMapping() Mapping {
	state := internal.StateMutable
	return newMapping(&otlpprofiles.Mapping{}, &state)
}

// MoveTo moves all properties from the current struct overriding the destination and
// resetting the current instance to its zero value
func (ms Mapping) MoveTo(dest Mapping) {
	ms.state.AssertMutable()
	dest.state.AssertMutable()
	*dest.orig = *ms.orig
	*ms.orig = otlpprofiles.Mapping{}
}

// ID returns the id associated with this Mapping.
func (ms Mapping) ID() uint64 {
	return ms.orig.Id
}

// SetID replaces the id associated with this Mapping.
func (ms Mapping) SetID(v uint64) {
	ms.state.AssertMutable()
	ms.orig.Id = v
}

// MemoryStart returns the memorystart associated with this Mapping.
func (ms Mapping) MemoryStart() uint64 {
	return ms.orig.MemoryStart
}

// SetMemoryStart replaces the memorystart associated with this Mapping.
func (ms Mapping) SetMemoryStart(v uint64) {
	ms.state.AssertMutable()
	ms.orig.MemoryStart = v
}

// MemoryLimit returns the memorylimit associated with this Mapping.
func (ms Mapping) MemoryLimit() uint64 {
	return ms.orig.MemoryLimit
}

// SetMemoryLimit replaces the memorylimit associated with this Mapping.
func (ms Mapping) SetMemoryLimit(v uint64) {
	ms.state.AssertMutable()
	ms.orig.MemoryLimit = v
}

// FileOffset returns the fileoffset associated with this Mapping.
func (ms Mapping) FileOffset() uint64 {
	return ms.orig.FileOffset
}

// SetFileOffset replaces the fileoffset associated with this Mapping.
func (ms Mapping) SetFileOffset(v uint64) {
	ms.state.AssertMutable()
	ms.orig.FileOffset = v
}

// Filename returns the filename associated with this Mapping.
func (ms Mapping) Filename() int64 {
	return ms.orig.Filename
}

// SetFilename replaces the filename associated with this Mapping.
func (ms Mapping) SetFilename(v int64) {
	ms.state.AssertMutable()
	ms.orig.Filename = v
}

// BuildID returns the buildid associated with this Mapping.
func (ms Mapping) BuildID() int64 {
	return ms.orig.BuildId
}

// SetBuildID replaces the buildid associated with this Mapping.
func (ms Mapping) SetBuildID(v int64) {
	ms.state.AssertMutable()
	ms.orig.BuildId = v
}

// BuildIDKind returns the buildidkind associated with this Mapping.
func (ms Mapping) BuildIDKind() otlpprofiles.BuildIdKind {
	return ms.orig.BuildIdKind
}

// SetBuildIDKind replaces the buildidkind associated with this Mapping.
func (ms Mapping) SetBuildIDKind(v otlpprofiles.BuildIdKind) {
	ms.state.AssertMutable()
	ms.orig.BuildIdKind = v
}

// Attributes returns the Attributes associated with this Mapping.
func (ms Mapping) Attributes() pcommon.UInt64Slice {
	return pcommon.UInt64Slice(internal.NewUInt64Slice(&ms.orig.Attributes, ms.state))
}

// HasFunctions returns the hasfunctions associated with this Mapping.
func (ms Mapping) HasFunctions() bool {
	return ms.orig.HasFunctions
}

// SetHasFunctions replaces the hasfunctions associated with this Mapping.
func (ms Mapping) SetHasFunctions(v bool) {
	ms.state.AssertMutable()
	ms.orig.HasFunctions = v
}

// HasFilenames returns the hasfilenames associated with this Mapping.
func (ms Mapping) HasFilenames() bool {
	return ms.orig.HasFilenames
}

// SetHasFilenames replaces the hasfilenames associated with this Mapping.
func (ms Mapping) SetHasFilenames(v bool) {
	ms.state.AssertMutable()
	ms.orig.HasFilenames = v
}

// HasLineNumbers returns the haslinenumbers associated with this Mapping.
func (ms Mapping) HasLineNumbers() bool {
	return ms.orig.HasLineNumbers
}

// SetHasLineNumbers replaces the haslinenumbers associated with this Mapping.
func (ms Mapping) SetHasLineNumbers(v bool) {
	ms.state.AssertMutable()
	ms.orig.HasLineNumbers = v
}

// HasInlineFrames returns the hasinlineframes associated with this Mapping.
func (ms Mapping) HasInlineFrames() bool {
	return ms.orig.HasInlineFrames
}

// SetHasInlineFrames replaces the hasinlineframes associated with this Mapping.
func (ms Mapping) SetHasInlineFrames(v bool) {
	ms.state.AssertMutable()
	ms.orig.HasInlineFrames = v
}

// CopyTo copies all properties from the current struct overriding the destination.
func (ms Mapping) CopyTo(dest Mapping) {
	dest.state.AssertMutable()
	dest.SetID(ms.ID())
	dest.SetMemoryStart(ms.MemoryStart())
	dest.SetMemoryLimit(ms.MemoryLimit())
	dest.SetFileOffset(ms.FileOffset())
	dest.SetFilename(ms.Filename())
	dest.SetBuildID(ms.BuildID())
	dest.SetBuildIDKind(ms.BuildIDKind())
	ms.Attributes().CopyTo(dest.Attributes())
	dest.SetHasFunctions(ms.HasFunctions())
	dest.SetHasFilenames(ms.HasFilenames())
	dest.SetHasLineNumbers(ms.HasLineNumbers())
	dest.SetHasInlineFrames(ms.HasInlineFrames())
}
