// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package envprovider

import (
	"context"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/collector/confmap"
	"go.opentelemetry.io/collector/confmap/confmaptest"
)

const envSchemePrefix = schemeName + ":"

const validYAML = `
processors:
  batch:
exporters:
  otlp:
    endpoint: "localhost:4317"
`

func TestValidateProviderScheme(t *testing.T) {
	assert.NoError(t, confmaptest.ValidateProviderScheme(NewWithSettings(confmap.NewProviderSettingsNoopLogger())))
}

func TestEmptyName(t *testing.T) {
	env := NewWithSettings(confmap.NewProviderSettingsNoopLogger())
	_, err := env.Retrieve(context.Background(), "", nil)
	require.Error(t, err)
	assert.NoError(t, env.Shutdown(context.Background()))
}

func TestUnsupportedScheme(t *testing.T) {
	env := NewWithSettings(confmap.NewProviderSettingsNoopLogger())
	_, err := env.Retrieve(context.Background(), "https://", nil)
	assert.Error(t, err)
	assert.NoError(t, env.Shutdown(context.Background()))
}

func TestInvalidYAML(t *testing.T) {
	const envName = "invalid-yaml"
	t.Setenv(envName, "[invalid,")
	env := NewWithSettings(confmap.NewProviderSettingsNoopLogger())
	_, err := env.Retrieve(context.Background(), envSchemePrefix+envName, nil)
	assert.Error(t, err)
	assert.NoError(t, env.Shutdown(context.Background()))
}

func TestEnv(t *testing.T) {
	const envName = "default-config"
	t.Setenv(envName, validYAML)

	env := NewWithSettings(confmap.NewProviderSettingsNoopLogger())
	ret, err := env.Retrieve(context.Background(), envSchemePrefix+envName, nil)
	require.NoError(t, err)
	retMap, err := ret.AsConf()
	assert.NoError(t, err)
	expectedMap := confmap.NewFromStringMap(map[string]any{
		"processors::batch":         nil,
		"exporters::otlp::endpoint": "localhost:4317",
	})
	assert.Equal(t, expectedMap.ToStringMap(), retMap.ToStringMap())

	assert.NoError(t, env.Shutdown(context.Background()))
}

func TestEnvWithLogger(t *testing.T) {
	const envName = "default-config"
	t.Setenv(envName, validYAML)
	core, ol := observer.New(zap.WarnLevel)
	logger := zap.New(core)

	env := NewWithSettings(confmap.ProviderSettings{Logger: logger})
	ret, err := env.Retrieve(context.Background(), envSchemePrefix+envName, nil)
	require.NoError(t, err)
	retMap, err := ret.AsConf()
	assert.NoError(t, err)
	expectedMap := confmap.NewFromStringMap(map[string]any{
		"processors::batch":         nil,
		"exporters::otlp::endpoint": "localhost:4317",
	})
	assert.Equal(t, expectedMap.ToStringMap(), retMap.ToStringMap())

	assert.NoError(t, env.Shutdown(context.Background()))
	assert.Equal(t, 0, ol.Len())
}

func TestUnsetEnvWithLoggerWarn(t *testing.T) {
	const envName = "default-config"
	core, ol := observer.New(zap.WarnLevel)
	logger := zap.New(core)

	env := NewWithSettings(confmap.ProviderSettings{Logger: logger})
	ret, err := env.Retrieve(context.Background(), envSchemePrefix+envName, nil)
	require.NoError(t, err)
	retMap, err := ret.AsConf()
	assert.NoError(t, err)
	expectedMap := confmap.NewFromStringMap(map[string]any{})
	assert.Equal(t, expectedMap.ToStringMap(), retMap.ToStringMap())

	assert.NoError(t, env.Shutdown(context.Background()))

	assert.Equal(t, 1, ol.Len())
	logLine := ol.All()[0]
	assert.Equal(t, "Configuration references unset environment variable", logLine.Message)
	assert.Equal(t, zap.WarnLevel, logLine.Level)
	assert.Equal(t, envName, logLine.Context[0].String)
}

func TestEmptyEnvWithLoggerWarn(t *testing.T) {
	const envName = "default-config"
	t.Setenv(envName, "")

	core, ol := observer.New(zap.InfoLevel)
	logger := zap.New(core)

	env := NewWithSettings(confmap.ProviderSettings{Logger: logger})
	ret, err := env.Retrieve(context.Background(), envSchemePrefix+envName, nil)
	require.NoError(t, err)
	retMap, err := ret.AsConf()
	assert.NoError(t, err)
	expectedMap := confmap.NewFromStringMap(map[string]any{})
	assert.Equal(t, expectedMap.ToStringMap(), retMap.ToStringMap())

	assert.NoError(t, env.Shutdown(context.Background()))

	assert.Equal(t, 1, ol.Len())
	logLine := ol.All()[0]
	assert.Equal(t, "Configuration references empty environment variable", logLine.Message)
	assert.Equal(t, zap.InfoLevel, logLine.Level)
	assert.Equal(t, envName, logLine.Context[0].String)
}
