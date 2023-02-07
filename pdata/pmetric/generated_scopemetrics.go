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
	"go.opentelemetry.io/collector/pdata/pcommon"
)

// ScopeMetrics is a collection of metrics from a LibraryInstrumentation.
//
// This is a reference type, if passed by value and callee modifies it the
// caller will see the modification.
//
// Must use NewMutableScopeMetrics function to create new instances.
// Important: zero-initialized instance is not valid for use.
type ScopeMetrics struct {
	commonScopeMetrics
}

type MutableScopeMetrics struct {
	commonScopeMetrics
	preventConversion struct{} // nolint:unused
}

type commonScopeMetrics struct {
	orig *otlpmetrics.ScopeMetrics
}

func newScopeMetricsFromOrig(orig *otlpmetrics.ScopeMetrics) ScopeMetrics {
	return ScopeMetrics{commonScopeMetrics{orig}}
}

func newMutableScopeMetricsFromOrig(orig *otlpmetrics.ScopeMetrics) MutableScopeMetrics {
	return MutableScopeMetrics{commonScopeMetrics: commonScopeMetrics{orig}}
}

// NewMutableScopeMetrics creates a new empty ScopeMetrics.
//
// This must be used only in testing code. Users should use "AppendEmpty" when part of a Slice,
// OR directly access the member if this is embedded in another struct.
func NewMutableScopeMetrics() MutableScopeMetrics {
	return newMutableScopeMetricsFromOrig(&otlpmetrics.ScopeMetrics{})
}

// nolint:unused
func (ms ScopeMetrics) asMutable() MutableScopeMetrics {
	return MutableScopeMetrics{commonScopeMetrics: commonScopeMetrics{orig: ms.orig}}
}

func (ms MutableScopeMetrics) AsImmutable() ScopeMetrics {
	return ScopeMetrics{commonScopeMetrics{orig: ms.orig}}
}

// MoveTo moves all properties from the current struct overriding the destination and
// resetting the current instance to its zero value
func (ms MutableScopeMetrics) MoveTo(dest MutableScopeMetrics) {
	*dest.orig = *ms.orig
	*ms.orig = otlpmetrics.ScopeMetrics{}
}

// Scope returns the scope associated with this ScopeMetrics.
func (ms ScopeMetrics) Scope() pcommon.InstrumentationScope {
	return internal.NewInstrumentationScopeFromOrig(&ms.orig.Scope)
}

// Scope returns the scope associated with this ScopeMetrics.
func (ms MutableScopeMetrics) Scope() pcommon.MutableInstrumentationScope {
	return internal.NewMutableInstrumentationScopeFromOrig(&ms.orig.Scope)
}

// SchemaUrl returns the schemaurl associated with this ScopeMetrics.
func (ms commonScopeMetrics) SchemaUrl() string {
	return ms.orig.SchemaUrl
}

// SetSchemaUrl replaces the schemaurl associated with this ScopeMetrics.
func (ms MutableScopeMetrics) SetSchemaUrl(v string) {
	ms.orig.SchemaUrl = v
}

// Metrics returns the Metrics associated with this ScopeMetrics.
func (ms ScopeMetrics) Metrics() MetricSlice {
	return newMetricSliceFromOrig(&ms.orig.Metrics)
}

func (ms MutableScopeMetrics) Metrics() MutableMetricSlice {
	return newMutableMetricSliceFromOrig(&ms.orig.Metrics)
}

// CopyTo copies all properties from the current struct overriding the destination.
func (ms commonScopeMetrics) CopyTo(dest MutableScopeMetrics) {
	ScopeMetrics{ms}.Scope().CopyTo(dest.Scope())
	dest.SetSchemaUrl(ms.SchemaUrl())
	ScopeMetrics{ms}.Metrics().CopyTo(dest.Metrics())
}

func generateTestScopeMetrics() MutableScopeMetrics {
	tv := NewMutableScopeMetrics()
	fillTestScopeMetrics(tv)
	return tv
}

func fillTestScopeMetrics(tv MutableScopeMetrics) {
	internal.FillTestInstrumentationScope(internal.NewMutableInstrumentationScopeFromOrig(&tv.orig.Scope))
	tv.orig.SchemaUrl = "https://opentelemetry.io/schemas/1.5.0"
	fillTestMetricSlice(newMutableMetricSliceFromOrig(&tv.orig.Metrics))
}
