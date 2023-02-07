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

	"go.opentelemetry.io/collector/pdata/pmetric"
)

func TestNilFuncMetrics(t *testing.T) {
	_, err := NewMetrics(nil)
	assert.Equal(t, errNilFunc, err)
}

func TestConsumeMetrics(t *testing.T) {
	consumeCalled := false
	cp, err := NewMetrics(func(context.Context, pmetric.Metrics) error { consumeCalled = true; return nil })
	assert.NoError(t, err)
	assert.NoError(t, cp.ConsumeMetrics(context.Background(), pmetric.NewMetrics()))
	assert.True(t, consumeCalled)
}

func TestConsumeMetrics_ReturnError(t *testing.T) {
	want := errors.New("my_error")
	cp, err := NewMetrics(func(context.Context, pmetric.Metrics) error { return want })
	assert.NoError(t, err)
	assert.Equal(t, want, cp.ConsumeMetrics(context.Background(), pmetric.NewMetrics()))
}
