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

package pmetric

import (
	"go.opentelemetry.io/collector/pdata/internal"
	otlpmetrics "go.opentelemetry.io/collector/pdata/internal/data/protogen/metrics/v1"
)

// Summary represents the type of a metric that is calculated by aggregating as a Summary of all reported double measurements over a time interval.
//
// This is a reference type, if passed by value and callee modifies it the
// caller will see the modification.
//
// Must use NewSummary function to create new instances.
// Important: zero-initialized instance is not valid for use.
type Summary struct {
	*pSummary
}

type pSummary struct {
	orig   *otlpmetrics.Summary
	state  *internal.State
	parent Metric
}

func (ms Summary) getOrig() *otlpmetrics.Summary {
	if *ms.state == internal.StateDirty {
		ms.orig, ms.state = ms.parent.refreshSummaryOrigState()
	}
	return ms.orig
}

func (ms Summary) ensureMutability() {
	if *ms.state == internal.StateShared {
		ms.parent.ensureMutability()
	}
}

func (ms Summary) getState() *internal.State {
	return ms.state
}

func (ms Summary) getDataPointsOrig() *[]*otlpmetrics.SummaryDataPoint {
	return &ms.getOrig().DataPoints
}

func newSummaryFromDataPointsOrig(childOrig *[]*otlpmetrics.SummaryDataPoint) Summary {
	state := internal.StateExclusive
	return Summary{&pSummary{
		state: &state,
		orig: &otlpmetrics.Summary{
			DataPoints: *childOrig,
		},
	}}
}

func newSummary(orig *otlpmetrics.Summary, parent Metric) Summary {
	return Summary{&pSummary{
		orig:   orig,
		state:  parent.getState(),
		parent: parent,
	}}
}

// NewSummary creates a new empty Summary.
//
// This must be used only in testing code. Users should use "AppendEmpty" when part of a Slice,
// OR directly access the member if this is embedded in another struct.
func NewSummary() Summary {
	state := internal.StateExclusive
	return Summary{&pSummary{orig: &otlpmetrics.Summary{}, state: &state}}
}

// MoveTo moves all properties from the current struct overriding the destination and
// resetting the current instance to its zero value
func (ms Summary) MoveTo(dest Summary) {
	ms.ensureMutability()
	dest.ensureMutability()
	*dest.getOrig() = *ms.getOrig()
	*ms.getOrig() = otlpmetrics.Summary{}
}

// DataPoints returns the <no value> associated with this Summary.
func (ms Summary) DataPoints() SummaryDataPointSlice {
	return newSummaryDataPointSliceFromParent(ms)
}

// CopyTo copies all properties from the current struct overriding the destination.
func (ms Summary) CopyTo(dest Summary) {
	dest.ensureMutability()
	ms.DataPoints().CopyTo(dest.DataPoints())
}
