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

package plog

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"go.opentelemetry.io/collector/pdata/internal"
)

func TestResourceLogs_MoveTo(t *testing.T) {
	ms := generateTestResourceLogs()
	dest := NewResourceLogs()
	ms.MoveTo(dest)
	assert.Equal(t, NewResourceLogs(), ms)
	assert.Equal(t, generateTestResourceLogs(), dest)
}

func TestResourceLogs_CopyTo(t *testing.T) {
	ms := NewResourceLogs()
	orig := NewResourceLogs()
	orig.CopyTo(ms)
	assert.Equal(t, orig, ms)
	orig = generateTestResourceLogs()
	orig.CopyTo(ms)
	assert.Equal(t, orig, ms)
}

func TestResourceLogs_Resource(t *testing.T) {
	ms := NewResourceLogs()
	internal.FillTestResource(internal.Resource(ms.Resource()))
	assert.Equal(t, internal.GenerateTestResource().GetOrig(), internal.Resource(ms.Resource()).GetOrig())
}

func TestResourceLogs_SchemaUrl(t *testing.T) {
	ms := NewResourceLogs()
	assert.Equal(t, "", ms.SchemaUrl())
	ms.SetSchemaUrl("https://opentelemetry.io/schemas/1.5.0")
	assert.Equal(t, "https://opentelemetry.io/schemas/1.5.0", ms.SchemaUrl())
}

func TestResourceLogs_ScopeLogs(t *testing.T) {
	ms := NewResourceLogs()
	assert.Equal(t, NewScopeLogsSlice().getOrig(), ms.ScopeLogs().getOrig())
	fillTestScopeLogsSlice(ms.ScopeLogs())
	assert.Equal(t, generateTestScopeLogsSlice().getOrig(), ms.ScopeLogs().getOrig())
}

func generateTestResourceLogs() ResourceLogs {
	tv := NewResourceLogs()
	fillTestResourceLogs(tv)
	return tv
}

func fillTestResourceLogs(tv ResourceLogs) {
	internal.FillTestResource(internal.NewResource(&tv.getOrig().Resource,
		wrappedResourceLogsResource{ResourceLogs: tv}))
	tv.getOrig().SchemaUrl = "https://opentelemetry.io/schemas/1.5.0"
	fillTestScopeLogsSlice(newScopeLogsSliceFromParent(tv))
}
