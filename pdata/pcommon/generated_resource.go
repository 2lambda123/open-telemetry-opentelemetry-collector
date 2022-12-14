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

// Code generated by "model/internal/cmd/pdatagen/main.go". DO NOT EDIT.
// To regenerate this file run "go run model/internal/cmd/pdatagen/main.go".

package pcommon

import (
	"go.opentelemetry.io/collector/pdata/internal"
	otlpresource "go.opentelemetry.io/collector/pdata/internal/data/protogen/resource/v1"
)

// Resource is a message representing the resource information.
//
// This is a reference type, if passed by value and callee modifies it the
// caller will see the modification.
//
// Must use NewResource function to create new instances.
// Important: zero-initialized instance is not valid for use.

type Resource internal.Resource

type MutableResource internal.MutableResource

func newResource(orig *otlpresource.Resource) Resource {
	return Resource(internal.NewResource(orig))
}

func newMutableResource(orig *otlpresource.Resource) MutableResource {
	return MutableResource(internal.NewResource(orig))
}

func (ms Resource) getOrig() *otlpresource.Resource {
	return internal.GetOrigResource(internal.Resource(ms))
}

func (ms MutableResource) getOrig() *otlpresource.Resource {
	return internal.GetMutableOrigResource(internal.MutableResource(ms))
}

// NewResource creates a new empty Resource.
//
// This must be used only in testing code. Users should use "AppendEmpty" when part of a Slice,
// OR directly access the member if this is embedded in another struct.
func NewResource() MutableResource {
	return newMutableResource(&otlpresource.Resource{})
}

// MoveTo moves all properties from the current struct overriding the destination and
// resetting the current instance to its zero value
func (ms MutableResource) MoveTo(dest MutableResource) {
	*dest.getOrig() = *ms.getOrig()
	*ms.getOrig() = otlpresource.Resource{}
}

// Attributes returns the Attributes associated with this Resource.
func (ms Resource) Attributes() Map {
	return Map(internal.NewMap(&ms.getOrig().Attributes))
}

// Attributes returns the Attributes associated with this Resource.
func (ms MutableResource) Attributes() MutableMap {
	return MutableMap(internal.NewMutableMap(&ms.getOrig().Attributes))
}

// DroppedAttributesCount returns the droppedattributescount associated with this Resource.
func (ms Resource) DroppedAttributesCount() uint32 {
	return ms.getOrig().DroppedAttributesCount
}

// MutableDroppedAttributesCount returns the droppedattributescount associated with this Resource.
func (ms MutableResource) DroppedAttributesCount() uint32 {
	return ms.getOrig().DroppedAttributesCount
}

// SetDroppedAttributesCount replaces the droppedattributescount associated with this Resource.
func (ms MutableResource) SetDroppedAttributesCount(v uint32) {
	ms.getOrig().DroppedAttributesCount = v
}

// CopyTo copies all properties from the current struct overriding the destination.
func (ms Resource) CopyTo(dest MutableResource) {
	ms.Attributes().CopyTo(dest.Attributes())
	dest.SetDroppedAttributesCount(ms.DroppedAttributesCount())
}

// CopyTo copies all properties from the current struct overriding the destination.
func (ms MutableResource) CopyTo(dest MutableResource) {
	newResource(ms.getOrig()).CopyTo(dest)
}
