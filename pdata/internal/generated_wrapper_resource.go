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

package internal

import (
	otlpcommon "go.opentelemetry.io/collector/pdata/internal/data/protogen/common/v1"
	otlpresource "go.opentelemetry.io/collector/pdata/internal/data/protogen/resource/v1"
)

type Resource struct {
	*pResource
}

type pResource struct {
	orig   *otlpresource.Resource
	state  *State
	parent Parent[*otlpresource.Resource]
}

func (ms Resource) GetOrig() *otlpresource.Resource {
	if *ms.state == StateDirty {
		ms.orig, ms.state = ms.parent.RefreshOrigState()
	}
	return ms.orig
}

func (ms Resource) EnsureMutability() {
	if *ms.state == StateShared {
		ms.parent.EnsureMutability()
	}
}

func (ms Resource) GetState() *State {
	return ms.state
}

func NewResource(orig *otlpresource.Resource, parent Parent[*otlpresource.Resource]) Resource {
	if parent == nil {
		state := StateExclusive
		return Resource{&pResource{orig: orig, state: &state}}
	}
	return Resource{&pResource{orig: orig, state: parent.GetState(), parent: parent}}
}

type WrappedResourceAttributes struct {
	Resource
}

func (es WrappedResourceAttributes) GetChildOrig() *[]otlpcommon.KeyValue {
	return &es.GetOrig().Attributes
}

func GenerateTestResource() Resource {
	orig := otlpresource.Resource{}
	tv := NewResource(&orig, nil)
	FillTestResource(tv)
	return tv
}

func FillTestResource(tv Resource) {
	FillTestMap(NewMapFromParent(WrappedResourceAttributes{Resource: tv}))
	tv.GetOrig().DroppedAttributesCount = uint32(17)
}
