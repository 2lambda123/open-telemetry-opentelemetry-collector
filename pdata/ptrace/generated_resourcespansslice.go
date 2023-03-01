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
	"sort"

	"go.opentelemetry.io/collector/pdata/internal"
	otlptrace "go.opentelemetry.io/collector/pdata/internal/data/protogen/trace/v1"
)

// ResourceSpansSlice logically represents a slice of ResourceSpans.
//
// This is a reference type. If passed by value and callee modifies it, the
// caller will see the modification.
//
// Must use NewResourceSpansSlice function to create new instances.
// Important: zero-initialized instance is not valid for use.
type ResourceSpansSlice struct {
	parent Traces
}

func newResourceSpansSliceFromOrig(orig *[]*otlptrace.ResourceSpans) ResourceSpansSlice {
	return ResourceSpansSlice{parent: newTracesFromResourceSpansOrig(orig)}
}

func newResourceSpansSliceFromParent(parent Traces) ResourceSpansSlice {
	return ResourceSpansSlice{parent: parent}
}

func (es ResourceSpansSlice) getOrig() *[]*otlptrace.ResourceSpans {
	return es.parent.getResourceSpansOrig()
}

func (es ResourceSpansSlice) ensureMutability() {
	es.parent.ensureMutability()
}

func (es ResourceSpansSlice) getState() *internal.State {
	return es.parent.getState()
}

func (es ResourceSpansSlice) refreshElementOrigState(i int) (*otlptrace.ResourceSpans, *internal.State) {
	return (*es.getOrig())[i], es.getState()
}

// NewResourceSpansSlice creates a ResourceSpansSlice with 0 elements.
// Can use "EnsureCapacity" to initialize with a given capacity.
func NewResourceSpansSlice() ResourceSpansSlice {
	orig := []*otlptrace.ResourceSpans(nil)
	return newResourceSpansSliceFromOrig(&orig)
}

// Len returns the number of elements in the slice.
//
// Returns "0" for a newly instance created with "NewResourceSpansSlice()".
func (es ResourceSpansSlice) Len() int {
	return len(*es.getOrig())
}

// At returns the element at the given index.
//
// This function is used mostly for iterating over all the values in the slice:
//
//	for i := 0; i < es.Len(); i++ {
//	    e := es.At(i)
//	    ... // Do something with the element
//	}
func (es ResourceSpansSlice) At(i int) ResourceSpans {
	return newResourceSpans((*es.getOrig())[i], es, i)
}

// EnsureCapacity is an operation that ensures the slice has at least the specified capacity.
// 1. If the newCap <= cap then no change in capacity.
// 2. If the newCap > cap then the slice capacity will be expanded to equal newCap.
//
// Here is how a new ResourceSpansSlice can be initialized:
//
//	es := NewResourceSpansSlice()
//	es.EnsureCapacity(4)
//	for i := 0; i < 4; i++ {
//	    e := es.AppendEmpty()
//	    // Here should set all the values for e.
//	}
func (es ResourceSpansSlice) EnsureCapacity(newCap int) {
	es.ensureMutability()
	oldCap := cap(*es.getOrig())
	if newCap <= oldCap {
		return
	}

	newOrig := make([]*otlptrace.ResourceSpans, len(*es.getOrig()), newCap)
	copy(newOrig, *es.getOrig())
	*es.getOrig() = newOrig
}

// AppendEmpty will append to the end of the slice an empty ResourceSpans.
// It returns the newly added ResourceSpans.
func (es ResourceSpansSlice) AppendEmpty() ResourceSpans {
	es.ensureMutability()
	*es.getOrig() = append(*es.getOrig(), &otlptrace.ResourceSpans{})
	return es.At(es.Len() - 1)
}

// MoveAndAppendTo moves all elements from the current slice and appends them to the dest.
// The current slice will be cleared.
func (es ResourceSpansSlice) MoveAndAppendTo(dest ResourceSpansSlice) {
	es.ensureMutability()
	dest.ensureMutability()
	if *dest.getOrig() == nil {
		// We can simply move the entire vector and avoid any allocations.
		*dest.getOrig() = *es.getOrig()
	} else {
		*dest.getOrig() = append(*dest.getOrig(), *es.getOrig()...)
	}
	*es.getOrig() = nil
}

// RemoveIf calls f sequentially for each element present in the slice.
// If f returns true, the element is removed from the slice.
func (es ResourceSpansSlice) RemoveIf(f func(ResourceSpans) bool) {
	es.ensureMutability()
	newLen := 0
	for i := 0; i < len(*es.getOrig()); i++ {
		if f(es.At(i)) {
			continue
		}
		if newLen == i {
			// Nothing to move, element is at the right place.
			newLen++
			continue
		}
		(*es.getOrig())[newLen] = (*es.getOrig())[i]
		newLen++
	}
	// TODO: Prevent memory leak by erasing truncated values.
	*es.getOrig() = (*es.getOrig())[:newLen]
}

// CopyTo copies all elements from the current slice overriding the destination.
func (es ResourceSpansSlice) CopyTo(dest ResourceSpansSlice) {
	dest.ensureMutability()
	srcLen := es.Len()
	destCap := cap(*dest.getOrig())
	exclState := internal.StateExclusive
	if srcLen <= destCap {
		(*dest.getOrig()) = (*dest.getOrig())[:srcLen:destCap]
		for i := range *es.getOrig() {
			srcResourceSpans := ResourceSpans{&pResourceSpans{orig: (*es.getOrig())[i], state: &exclState}}
			destResourceSpans := ResourceSpans{&pResourceSpans{orig: (*dest.getOrig())[i],
				state: &exclState}}
			srcResourceSpans.CopyTo(destResourceSpans)
		}
		return
	}
	origs := make([]otlptrace.ResourceSpans, srcLen)
	wrappers := make([]*otlptrace.ResourceSpans, srcLen)
	for i := range *es.getOrig() {
		wrappers[i] = &origs[i]
		srcResourceSpans := ResourceSpans{&pResourceSpans{orig: (*es.getOrig())[i], state: &exclState}}
		destResourceSpans := ResourceSpans{&pResourceSpans{orig: wrappers[i], state: &exclState}}
		srcResourceSpans.CopyTo(destResourceSpans)
	}
	*dest.getOrig() = wrappers
}

// Sort sorts the ResourceSpans elements within ResourceSpansSlice given the
// provided less function so that two instances of ResourceSpansSlice
// can be compared.
func (es ResourceSpansSlice) Sort(less func(a, b ResourceSpans) bool) {
	es.ensureMutability()
	sort.SliceStable(*es.getOrig(), func(i, j int) bool { return less(es.At(i), es.At(j)) })
}
