// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package envprovider // import "go.opentelemetry.io/collector/confmap/provider/envprovider"

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"os"
	"strings"

	"go.opentelemetry.io/collector/confmap"
	"go.opentelemetry.io/collector/confmap/provider/internal"
)

const schemeName = "env"

type provider struct {
	logger *zap.Logger
}

// NewWithSettings returns a new confmap.Provider that reads the configuration from the given environment variable.
//
// This Provider supports "env" scheme, and can be called with a selector:
// `env:NAME_OF_ENVIRONMENT_VARIABLE`
func NewWithSettings(a confmap.ProviderSettings) confmap.Provider {
	return &provider{
		logger: a.Logger,
	}
}

func (emp *provider) Retrieve(_ context.Context, uri string, _ confmap.WatcherFunc) (*confmap.Retrieved, error) {
	if !strings.HasPrefix(uri, schemeName+":") {
		return nil, fmt.Errorf("%q uri is not supported by %q provider", uri, schemeName)
	}
	envVarName := uri[len(schemeName)+1:]
	val, exists := os.LookupEnv(envVarName)
	if !exists {
		emp.logger.Warn(fmt.Sprintf("Environment variable %s is used in configuration but is unset", envVarName))
	} else if len(val) == 0 {
		emp.logger.Warn(fmt.Sprintf("Environment variable %s is used in configuration but is empty", envVarName))
	}

	return internal.NewRetrievedFromYAML([]byte(val))
}

func (*provider) Scheme() string {
	return schemeName
}

func (*provider) Shutdown(context.Context) error {
	return nil
}
