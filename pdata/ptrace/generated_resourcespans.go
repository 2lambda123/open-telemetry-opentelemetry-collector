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
	"go.opentelemetry.io/collector/pdata/internal"
	otlpresource "go.opentelemetry.io/collector/pdata/internal/data/protogen/resource/v1"
	otlptrace "go.opentelemetry.io/collector/pdata/internal/data/protogen/trace/v1"
	"go.opentelemetry.io/collector/pdata/pcommon"
)

// ResourceSpans is a collection of spans from a Resource.
//
// This is a reference type, if passed by value and callee modifies it the
// caller will see the modification.
//
// Must use NewResourceSpans function to create new instances.
// Important: zero-initialized instance is not valid for use.
type ResourceSpans struct {
	*pResourceSpans
}

type pResourceSpans struct {
	orig   *otlptrace.ResourceSpans
	state  *internal.State
	parent ResourceSpansSlice
	idx    int
}

func (ms ResourceSpans) getOrig() *otlptrace.ResourceSpans {
	if *ms.state == internal.StateDirty {
		ms.orig, ms.state = ms.parent.refreshElementOrigState(ms.idx)
	}
	return ms.orig
}

func (ms ResourceSpans) ensureMutability() {
	if *ms.state == internal.StateShared {
		ms.parent.ensureMutability()
	}
}

func (ms ResourceSpans) getState() *internal.State {
	return ms.state
}

type wrappedResourceSpansResource struct {
	ResourceSpans
}

func (es wrappedResourceSpansResource) RefreshOrigState() (*otlpresource.Resource, *internal.State) {
	return &es.getOrig().Resource, es.getState()
}

func (es wrappedResourceSpansResource) EnsureMutability() {
	es.ensureMutability()
}

func (es wrappedResourceSpansResource) GetState() *internal.State {
	return es.getState()
}

func (ms ResourceSpans) getScopeSpansOrig() *[]*otlptrace.ScopeSpans {
	return &ms.getOrig().ScopeSpans
}

func newResourceSpansFromScopeSpansOrig(childOrig *[]*otlptrace.ScopeSpans) ResourceSpans {
	state := internal.StateExclusive
	return ResourceSpans{&pResourceSpans{
		state: &state,
		orig: &otlptrace.ResourceSpans{
			ScopeSpans: *childOrig,
		},
	}}
}

func newResourceSpans(orig *otlptrace.ResourceSpans, parent ResourceSpansSlice, idx int) ResourceSpans {
	return ResourceSpans{&pResourceSpans{
		orig:   orig,
		state:  parent.getState(),
		parent: parent,
		idx:    idx,
	}}
}

// NewResourceSpans creates a new empty ResourceSpans.
//
// This must be used only in testing code. Users should use "AppendEmpty" when part of a Slice,
// OR directly access the member if this is embedded in another struct.
func NewResourceSpans() ResourceSpans {
	state := internal.StateExclusive
	return ResourceSpans{&pResourceSpans{orig: &otlptrace.ResourceSpans{}, state: &state}}
}

// MoveTo moves all properties from the current struct overriding the destination and
// resetting the current instance to its zero value
func (ms ResourceSpans) MoveTo(dest ResourceSpans) {
	ms.ensureMutability()
	dest.ensureMutability()
	*dest.getOrig() = *ms.getOrig()
	*ms.getOrig() = otlptrace.ResourceSpans{}
}

// Resource returns the resource associated with this ResourceSpans.
func (ms ResourceSpans) Resource() pcommon.Resource {
	return pcommon.Resource(internal.NewResource(&ms.getOrig().Resource, wrappedResourceSpansResource{ResourceSpans: ms}))
}

// SchemaUrl returns the schemaurl associated with this ResourceSpans.
func (ms ResourceSpans) SchemaUrl() string {
	return ms.getOrig().SchemaUrl
}

// SetSchemaUrl replaces the schemaurl associated with this ResourceSpans.
func (ms ResourceSpans) SetSchemaUrl(v string) {
	ms.ensureMutability()
	ms.getOrig().SchemaUrl = v
}

// ScopeSpans returns the <no value> associated with this ResourceSpans.
func (ms ResourceSpans) ScopeSpans() ScopeSpansSlice {
	return newScopeSpansSliceFromParent(ms)
}

// CopyTo copies all properties from the current struct overriding the destination.
func (ms ResourceSpans) CopyTo(dest ResourceSpans) {
	dest.ensureMutability()
	ms.Resource().CopyTo(dest.Resource())
	dest.SetSchemaUrl(ms.SchemaUrl())
	ms.ScopeSpans().CopyTo(dest.ScopeSpans())
}
