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
	"go.opentelemetry.io/collector/pdata/pcommon"
)

func TestScopeLogs_MoveTo(t *testing.T) {
	ms := generateTestScopeLogs()
	dest := NewMutableScopeLogs()
	ms.MoveTo(dest)
	assert.Equal(t, NewMutableScopeLogs(), ms)
	assert.Equal(t, generateTestScopeLogs(), dest)
}

func TestScopeLogs_CopyTo(t *testing.T) {
	ms := NewMutableScopeLogs()
	orig := NewMutableScopeLogs()
	orig.CopyTo(ms)
	assert.Equal(t, orig, ms)
	orig = generateTestScopeLogs()
	orig.CopyTo(ms)
	assert.Equal(t, orig, ms)
}

func TestScopeLogs_Scope(t *testing.T) {
	ms := NewMutableScopeLogs()
	internal.FillTestInstrumentationScope(internal.MutableInstrumentationScope(ms.Scope()))
	assert.Equal(t, pcommon.MutableInstrumentationScope(internal.GenerateTestInstrumentationScope()), ms.Scope())
}

func TestScopeLogs_SchemaUrl(t *testing.T) {
	ms := NewMutableScopeLogs()
	assert.Equal(t, "", ms.SchemaUrl())
	ms.SetSchemaUrl("https://opentelemetry.io/schemas/1.5.0")
	assert.Equal(t, "https://opentelemetry.io/schemas/1.5.0", ms.SchemaUrl())
}

func TestScopeLogs_LogRecords(t *testing.T) {
	ms := NewMutableScopeLogs()
	assert.Equal(t, NewMutableLogRecordSlice(), ms.LogRecords())
	fillTestLogRecordSlice(ms.LogRecords())
	assert.Equal(t, generateTestLogRecordSlice(), ms.LogRecords())
}
