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

type UInt64Slice struct {
	parent Parent[*[]uint64]
}

type stubUInt64SliceParent struct {
	orig *[]uint64
}

func (vp stubUInt64SliceParent) EnsureMutability() {}

func (vp stubUInt64SliceParent) GetChildOrig() *[]uint64 {
	return vp.orig
}

var _ Parent[*[]uint64] = (*stubUInt64SliceParent)(nil)

func (ms UInt64Slice) GetOrig() *[]uint64 {
	return ms.parent.GetChildOrig()
}

func (ms UInt64Slice) EnsureMutability() {
	ms.parent.EnsureMutability()
}

func NewUInt64SliceFromOrig(orig *[]uint64) UInt64Slice {
	return UInt64Slice{parent: &stubUInt64SliceParent{orig: orig}}
}

func NewUInt64SliceFromParent(parent Parent[*[]uint64]) UInt64Slice {
	return UInt64Slice{parent: parent}
}
