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
	"go.opentelemetry.io/collector/pdata/internal"
	otlplogs "go.opentelemetry.io/collector/pdata/internal/data/protogen/logs/v1"
	"go.opentelemetry.io/collector/pdata/pcommon"
)

// ResourceLogs is a collection of logs from a Resource.
//
// This is a reference type, if passed by value and callee modifies it the
// caller will see the modification.
//
// Must use NewMutableResourceLogs function to create new instances.
// Important: zero-initialized instance is not valid for use.
type ResourceLogs struct {
	commonResourceLogs
}

type MutableResourceLogs struct {
	commonResourceLogs
	preventConversion struct{} // nolint:unused
}

type commonResourceLogs struct {
	orig *otlplogs.ResourceLogs
}

func newResourceLogsFromOrig(orig *otlplogs.ResourceLogs) ResourceLogs {
	return ResourceLogs{commonResourceLogs{orig}}
}

func newMutableResourceLogsFromOrig(orig *otlplogs.ResourceLogs) MutableResourceLogs {
	return MutableResourceLogs{commonResourceLogs: commonResourceLogs{orig}}
}

// NewMutableResourceLogs creates a new empty ResourceLogs.
//
// This must be used only in testing code. Users should use "AppendEmpty" when part of a Slice,
// OR directly access the member if this is embedded in another struct.
func NewMutableResourceLogs() MutableResourceLogs {
	return newMutableResourceLogsFromOrig(&otlplogs.ResourceLogs{})
}

// nolint:unused
func (ms ResourceLogs) asMutable() MutableResourceLogs {
	return MutableResourceLogs{commonResourceLogs: commonResourceLogs{orig: ms.orig}}
}

func (ms MutableResourceLogs) AsImmutable() ResourceLogs {
	return ResourceLogs{commonResourceLogs{orig: ms.orig}}
}

// MoveTo moves all properties from the current struct overriding the destination and
// resetting the current instance to its zero value
func (ms MutableResourceLogs) MoveTo(dest MutableResourceLogs) {
	*dest.orig = *ms.orig
	*ms.orig = otlplogs.ResourceLogs{}
}

// Resource returns the resource associated with this ResourceLogs.
func (ms ResourceLogs) Resource() pcommon.Resource {
	return internal.NewResourceFromOrig(&ms.orig.Resource)
}

// Resource returns the resource associated with this ResourceLogs.
func (ms MutableResourceLogs) Resource() pcommon.MutableResource {
	return internal.NewMutableResourceFromOrig(&ms.orig.Resource)
}

// SchemaUrl returns the schemaurl associated with this ResourceLogs.
func (ms commonResourceLogs) SchemaUrl() string {
	return ms.orig.SchemaUrl
}

// SetSchemaUrl replaces the schemaurl associated with this ResourceLogs.
func (ms MutableResourceLogs) SetSchemaUrl(v string) {
	ms.orig.SchemaUrl = v
}

// ScopeLogs returns the ScopeLogs associated with this ResourceLogs.
func (ms ResourceLogs) ScopeLogs() ScopeLogsSlice {
	return newScopeLogsSliceFromOrig(&ms.orig.ScopeLogs)
}

func (ms MutableResourceLogs) ScopeLogs() MutableScopeLogsSlice {
	return newMutableScopeLogsSliceFromOrig(&ms.orig.ScopeLogs)
}

// CopyTo copies all properties from the current struct overriding the destination.
func (ms commonResourceLogs) CopyTo(dest MutableResourceLogs) {
	ResourceLogs{ms}.Resource().CopyTo(dest.Resource())
	dest.SetSchemaUrl(ms.SchemaUrl())
	ResourceLogs{ms}.ScopeLogs().CopyTo(dest.ScopeLogs())
}

func generateTestResourceLogs() MutableResourceLogs {
	tv := NewMutableResourceLogs()
	fillTestResourceLogs(tv)
	return tv
}

func fillTestResourceLogs(tv MutableResourceLogs) {
	internal.FillTestResource(internal.NewMutableResourceFromOrig(&tv.orig.Resource))
	tv.orig.SchemaUrl = "https://opentelemetry.io/schemas/1.5.0"
	fillTestScopeLogsSlice(newMutableScopeLogsSliceFromOrig(&tv.orig.ScopeLogs))
}
