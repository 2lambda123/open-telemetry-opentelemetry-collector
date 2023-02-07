// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package consumer

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"go.opentelemetry.io/collector/pdata/ptrace"
)

func TestNilFuncTraces(t *testing.T) {
	_, err := NewTraces(nil)
	assert.Equal(t, errNilFunc, err)
}

func TestConsumeTraces(t *testing.T) {
	consumeCalled := false
	cp, err := NewTraces(func(context.Context, ptrace.Traces) error { consumeCalled = true; return nil })
	assert.NoError(t, err)
	assert.NoError(t, cp.ConsumeTraces(context.Background(), ptrace.NewTraces()))
	assert.True(t, consumeCalled)
}

func TestConsumeTraces_ReturnError(t *testing.T) {
	want := errors.New("my_error")
	cp, err := NewTraces(func(context.Context, ptrace.Traces) error { return want })
	assert.NoError(t, err)
	assert.Equal(t, want, cp.ConsumeTraces(context.Background(), ptrace.NewTraces()))
}
