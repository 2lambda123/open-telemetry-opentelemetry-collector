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

	"go.opentelemetry.io/collector/pdata/internal"
	otlpmetrics "go.opentelemetry.io/collector/pdata/internal/data/protogen/metrics/v1"
	"go.opentelemetry.io/collector/pdata/pcommon"
)

func TestExponentialHistogramDataPoint_MoveTo(t *testing.T) {
	ms := generateTestExponentialHistogramDataPoint()
	dest := NewExponentialHistogramDataPoint()
	ms.MoveTo(dest)
	assert.Equal(t, NewExponentialHistogramDataPoint(), ms)
	assert.Equal(t, generateTestExponentialHistogramDataPoint(), dest)
}

func TestExponentialHistogramDataPoint_CopyTo(t *testing.T) {
	ms := NewExponentialHistogramDataPoint()
	orig := NewExponentialHistogramDataPoint()
	orig.CopyTo(ms)
	assert.Equal(t, orig, ms)
	orig = generateTestExponentialHistogramDataPoint()
	orig.CopyTo(ms)
	assert.Equal(t, orig, ms)
}

func TestExponentialHistogramDataPoint_Attributes(t *testing.T) {
	ms := NewExponentialHistogramDataPoint()
	assert.Equal(t, internal.Map(pcommon.NewMap()).GetOrig(), internal.Map(ms.Attributes()).GetOrig())
	internal.FillTestMap(internal.Map(ms.Attributes()))
	assert.Equal(t, internal.GenerateTestMap().GetOrig(), internal.Map(ms.Attributes()).GetOrig())
}

func TestExponentialHistogramDataPoint_StartTimestamp(t *testing.T) {
	ms := NewExponentialHistogramDataPoint()
	assert.Equal(t, pcommon.Timestamp(0), ms.StartTimestamp())
	testValStartTimestamp := pcommon.Timestamp(1234567890)
	ms.SetStartTimestamp(testValStartTimestamp)
	assert.Equal(t, testValStartTimestamp, ms.StartTimestamp())
}

func TestExponentialHistogramDataPoint_Timestamp(t *testing.T) {
	ms := NewExponentialHistogramDataPoint()
	assert.Equal(t, pcommon.Timestamp(0), ms.Timestamp())
	testValTimestamp := pcommon.Timestamp(1234567890)
	ms.SetTimestamp(testValTimestamp)
	assert.Equal(t, testValTimestamp, ms.Timestamp())
}

func TestExponentialHistogramDataPoint_Count(t *testing.T) {
	ms := NewExponentialHistogramDataPoint()
	assert.Equal(t, uint64(0), ms.Count())
	ms.SetCount(uint64(17))
	assert.Equal(t, uint64(17), ms.Count())
}

func TestExponentialHistogramDataPoint_Sum(t *testing.T) {
	ms := NewExponentialHistogramDataPoint()
	assert.Equal(t, float64(0.0), ms.Sum())
	ms.SetSum(float64(17.13))
	assert.True(t, ms.HasSum())
	assert.Equal(t, float64(17.13), ms.Sum())
	ms.RemoveSum()
	assert.False(t, ms.HasSum())
}

func TestExponentialHistogramDataPoint_Scale(t *testing.T) {
	ms := NewExponentialHistogramDataPoint()
	assert.Equal(t, int32(0), ms.Scale())
	ms.SetScale(int32(4))
	assert.Equal(t, int32(4), ms.Scale())
}

func TestExponentialHistogramDataPoint_ZeroCount(t *testing.T) {
	ms := NewExponentialHistogramDataPoint()
	assert.Equal(t, uint64(0), ms.ZeroCount())
	ms.SetZeroCount(uint64(201))
	assert.Equal(t, uint64(201), ms.ZeroCount())
}

func TestExponentialHistogramDataPoint_Positive(t *testing.T) {
	ms := NewExponentialHistogramDataPoint()
	fillTestExponentialHistogramDataPointBuckets(ms.Positive())
	assert.Equal(t, generateTestExponentialHistogramDataPointBuckets().getOrig(), ms.Positive().getOrig())
}

func TestExponentialHistogramDataPoint_Negative(t *testing.T) {
	ms := NewExponentialHistogramDataPoint()
	fillTestExponentialHistogramDataPointBuckets(ms.Negative())
	assert.Equal(t, generateTestExponentialHistogramDataPointBuckets().getOrig(), ms.Negative().getOrig())
}

func TestExponentialHistogramDataPoint_Exemplars(t *testing.T) {
	ms := NewExponentialHistogramDataPoint()
	assert.Equal(t, NewExemplarSlice().getOrig(), ms.Exemplars().getOrig())
	fillTestExemplarSlice(ms.Exemplars())
	assert.Equal(t, generateTestExemplarSlice().getOrig(), ms.Exemplars().getOrig())
}

func TestExponentialHistogramDataPoint_Flags(t *testing.T) {
	ms := NewExponentialHistogramDataPoint()
	assert.Equal(t, DataPointFlags(0), ms.Flags())
	testValFlags := DataPointFlags(1)
	ms.SetFlags(testValFlags)
	assert.Equal(t, testValFlags, ms.Flags())
}

func TestExponentialHistogramDataPoint_Min(t *testing.T) {
	ms := NewExponentialHistogramDataPoint()
	assert.Equal(t, float64(0.0), ms.Min())
	ms.SetMin(float64(9.23))
	assert.True(t, ms.HasMin())
	assert.Equal(t, float64(9.23), ms.Min())
	ms.RemoveMin()
	assert.False(t, ms.HasMin())
}

func TestExponentialHistogramDataPoint_Max(t *testing.T) {
	ms := NewExponentialHistogramDataPoint()
	assert.Equal(t, float64(0.0), ms.Max())
	ms.SetMax(float64(182.55))
	assert.True(t, ms.HasMax())
	assert.Equal(t, float64(182.55), ms.Max())
	ms.RemoveMax()
	assert.False(t, ms.HasMax())
}

func generateTestExponentialHistogramDataPoint() ExponentialHistogramDataPoint {
	tv := NewExponentialHistogramDataPoint()
	fillTestExponentialHistogramDataPoint(tv)
	return tv
}

func fillTestExponentialHistogramDataPoint(tv ExponentialHistogramDataPoint) {
	internal.FillTestMap(internal.NewMapFromParent(wrappedExponentialHistogramDataPointAttributes{ExponentialHistogramDataPoint: tv}))
	tv.getOrig().StartTimeUnixNano = 1234567890
	tv.getOrig().TimeUnixNano = 1234567890
	tv.getOrig().Count = uint64(17)
	tv.getOrig().Sum_ = &otlpmetrics.ExponentialHistogramDataPoint_Sum{Sum: float64(17.13)}
	tv.getOrig().Scale = int32(4)
	tv.getOrig().ZeroCount = uint64(201)
	fillTestExponentialHistogramDataPointBuckets(newExponentialHistogramDataPointBucketsFromParent(wrappedExponentialHistogramDataPointPositive{ExponentialHistogramDataPoint: tv}))
	fillTestExponentialHistogramDataPointBuckets(newExponentialHistogramDataPointBucketsFromParent(wrappedExponentialHistogramDataPointNegative{ExponentialHistogramDataPoint: tv}))
	fillTestExemplarSlice(newExemplarSliceFromParent(wrappedExponentialHistogramDataPointExemplars{ExponentialHistogramDataPoint: tv}))
	tv.getOrig().Flags = 1
	tv.getOrig().Min_ = &otlpmetrics.ExponentialHistogramDataPoint_Min{Min: float64(9.23)}
	tv.getOrig().Max_ = &otlpmetrics.ExponentialHistogramDataPoint_Max{Max: float64(182.55)}
}
