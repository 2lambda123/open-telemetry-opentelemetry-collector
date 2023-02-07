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
	"go.opentelemetry.io/collector/pdata/pcommon"
)

func TestResourceMetrics_MoveTo(t *testing.T) {
	ms := generateTestResourceMetrics()
	dest := NewMutableResourceMetrics()
	ms.MoveTo(dest)
	assert.Equal(t, NewMutableResourceMetrics(), ms)
	assert.Equal(t, generateTestResourceMetrics(), dest)
}

func TestResourceMetrics_CopyTo(t *testing.T) {
	ms := NewMutableResourceMetrics()
	orig := NewMutableResourceMetrics()
	orig.CopyTo(ms)
	assert.Equal(t, orig, ms)
	orig = generateTestResourceMetrics()
	orig.CopyTo(ms)
	assert.Equal(t, orig, ms)
}

func TestResourceMetrics_Resource(t *testing.T) {
	ms := NewMutableResourceMetrics()
	internal.FillTestResource(internal.MutableResource(ms.Resource()))
	assert.Equal(t, pcommon.MutableResource(internal.GenerateTestResource()), ms.Resource())
}

func TestResourceMetrics_SchemaUrl(t *testing.T) {
	ms := NewMutableResourceMetrics()
	assert.Equal(t, "", ms.SchemaUrl())
	ms.SetSchemaUrl("https://opentelemetry.io/schemas/1.5.0")
	assert.Equal(t, "https://opentelemetry.io/schemas/1.5.0", ms.SchemaUrl())
}

func TestResourceMetrics_ScopeMetrics(t *testing.T) {
	ms := NewMutableResourceMetrics()
	assert.Equal(t, NewMutableScopeMetricsSlice(), ms.ScopeMetrics())
	fillTestScopeMetricsSlice(ms.ScopeMetrics())
	assert.Equal(t, generateTestScopeMetricsSlice(), ms.ScopeMetrics())
}
