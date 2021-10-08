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

package featuregate

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegistry(t *testing.T) {
	r := registry{gates: map[string]Gate{}}

	gate := Gate{
		ID:          "foo",
		Description: "Test Gate",
		Enabled:     true,
	}

	assert.Empty(t, r.list())
	assert.False(t, r.isEnabled(gate.ID))

	assert.NoError(t, r.add(gate))
	assert.Len(t, r.list(), 1)
	assert.True(t, r.isEnabled(gate.ID))

	r.apply(map[string]bool{gate.ID: false})
	assert.False(t, r.isEnabled(gate.ID))

	assert.Error(t, r.add(gate))
}

func TestGlobalRegistry(t *testing.T) {
	gate := Gate{
		ID:          "feature_gate_test.foo",
		Description: "Test Gate",
		Enabled:     true,
	}

	assert.NotContains(t, List(), gate)
	assert.False(t, IsEnabled(gate.ID))

	assert.NotPanics(t, func() { Register(gate) })
	assert.Contains(t, List(), gate)
	assert.True(t, IsEnabled(gate.ID))

	Apply(map[string]bool{gate.ID: false})
	assert.False(t, IsEnabled(gate.ID))

	assert.Panics(t, func() { Register(gate) })
	reg.Lock()
	delete(reg.gates, gate.ID)
	reg.Unlock()
}
