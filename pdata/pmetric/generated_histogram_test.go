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
	"testing"

	"github.com/stretchr/testify/assert"

	otlpmetrics "go.opentelemetry.io/collector/pdata/internal/data/protogen/metrics/v1"
)

func TestHistogram_MoveTo(t *testing.T) {
	ms := generateTestHistogram()
	dest := NewHistogram()
	ms.MoveTo(dest)
	assert.Equal(t, NewHistogram(), ms)
	assert.Equal(t, generateTestHistogram(), dest)
}

func TestHistogram_CopyTo(t *testing.T) {
	ms := NewHistogram()
	orig := NewHistogram()
	orig.CopyTo(ms)
	assert.Equal(t, orig, ms)
	orig = generateTestHistogram()
	orig.CopyTo(ms)
	assert.Equal(t, orig, ms)
}

func TestHistogram_AggregationTemporality(t *testing.T) {
	ms := NewHistogram()
	assert.Equal(t, AggregationTemporality(otlpmetrics.AggregationTemporality(0)), ms.AggregationTemporality())
	testValAggregationTemporality := AggregationTemporality(otlpmetrics.AggregationTemporality(1))
	ms.SetAggregationTemporality(testValAggregationTemporality)
	assert.Equal(t, testValAggregationTemporality, ms.AggregationTemporality())
}

func TestHistogram_DataPoints(t *testing.T) {
	ms := NewHistogram()
	assert.Equal(t, NewHistogramDataPointSlice().getOrig(), ms.DataPoints().getOrig())
	fillTestHistogramDataPointSlice(ms.DataPoints())
	assert.Equal(t, generateTestHistogramDataPointSlice().getOrig(), ms.DataPoints().getOrig())
}

func generateTestHistogram() Histogram {
	tv := NewHistogram()
	fillTestHistogram(tv)
	return tv
}

func fillTestHistogram(tv Histogram) {
	tv.getOrig().AggregationTemporality = otlpmetrics.AggregationTemporality(1)
	fillTestHistogramDataPointSlice(newHistogramDataPointSliceFromParent(tv))
}
