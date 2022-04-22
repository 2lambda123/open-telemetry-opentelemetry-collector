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

package service // import "go.opentelemetry.io/collector/service"

import (
	"context"
	"fmt"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config"
	"go.opentelemetry.io/collector/internal/configunmarshaler"
	"go.opentelemetry.io/collector/service/internal/mapresolver"
)

// ConfigProvider provides the service configuration.
//
// The typical usage is the following:
//
//		cfgProvider.Get(...)
//		cfgProvider.Watch() // wait for an event.
//		cfgProvider.Get(...)
//		cfgProvider.Watch() // wait for an event.
//		// repeat Get/Watch cycle until it is time to shut down the Collector process.
//		cfgProvider.Shutdown()
type ConfigProvider interface {
	// Get returns the service configuration, or error otherwise.
	//
	// Should never be called concurrently with itself, Watch or Shutdown.
	Get(ctx context.Context, factories component.Factories) (*config.Config, error)

	// Watch blocks until any configuration change was detected or an unrecoverable error
	// happened during monitoring the configuration changes.
	//
	// Error is nil if the configuration is changed and needs to be re-fetched. Any non-nil
	// error indicates that there was a problem with watching the config changes.
	//
	// Should never be called concurrently with itself or Get.
	Watch() <-chan error

	// Shutdown signals that the provider is no longer in use and the that should close
	// and release any resources that it may have created.
	//
	// This function must terminate the Watch channel.
	//
	// Should never be called concurrently with itself or Get.
	Shutdown(ctx context.Context) error
}

type configProvider struct {
	mapResolver       *mapresolver.MapResolver
	configUnmarshaler configunmarshaler.ConfigUnmarshaler
}

// ConfigProviderSettings are the settings to configure the behavior of the ConfigProvider.
type ConfigProviderSettings struct {
	// Locations from where the config.Map is retrieved, and merged in the given order.
	// It is required to have at least one location.
	Locations []string

	// MapProviders is a map of pairs <scheme, config.MapProvider>.
	// It is required to have at least one config.MapProvider.
	MapProviders map[string]config.MapProvider

	// MapConverters is a slice of config.MapConverterFunc.
	MapConverters []config.MapConverterFunc

	// Deprecated: [v0.50.0] because providing custom ConfigUnmarshaler is not necessary since users can wrap/implement
	// ConfigProvider if needed to change the resulted config. This functionality will be kept for at least 2 minor versions,
	// and if nobody express a need for it will be removed.
	Unmarshaler configunmarshaler.ConfigUnmarshaler
}

func newDefaultConfigProviderSettings(locations []string) ConfigProviderSettings {
	set := mapresolver.NewDefaultSettings(locations)
	return ConfigProviderSettings{
		Locations:     set.Locations,
		MapProviders:  set.MapProviders,
		MapConverters: set.MapConverters,
		Unmarshaler:   configunmarshaler.NewDefault(),
	}
}

// NewConfigProvider returns a new ConfigProvider that provides the service configuration:
// * Initially it resolves the "configuration map":
//	 * Retrieve the config.Map by merging all retrieved maps from the given `locations` in order.
// 	 * Then applies all the config.MapConverterFunc in the given order.
// * Then unmarshalls the config.Map into the service Config.
func NewConfigProvider(set ConfigProviderSettings) (ConfigProvider, error) {
	mr, err := mapresolver.NewMapResolver(&mapresolver.Settings{
		Locations:     set.Locations,
		MapProviders:  set.MapProviders,
		MapConverters: set.MapConverters,
	})
	if err != nil {
		return nil, err
	}

	unmarshaler := set.Unmarshaler
	if unmarshaler == nil {
		unmarshaler = configunmarshaler.NewDefault()
	}

	return &configProvider{
		mapResolver:       mr,
		configUnmarshaler: unmarshaler,
	}, nil
}

func (cm *configProvider) Get(ctx context.Context, factories component.Factories) (*config.Config, error) {
	retMap, err := cm.mapResolver.Resolve(ctx)
	if err != nil {
		return nil, fmt.Errorf("cannot resolve the configuration: %w", err)
	}

	var cfg *config.Config
	if cfg, err = cm.configUnmarshaler.Unmarshal(retMap, factories); err != nil {
		return nil, fmt.Errorf("cannot unmarshal the configuration: %w", err)
	}

	if err = cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return cfg, nil
}

func (cm *configProvider) Watch() <-chan error {
	return cm.mapResolver.Watch()
}

func (cm *configProvider) Shutdown(ctx context.Context) error {
	return cm.mapResolver.Shutdown(ctx)
}
