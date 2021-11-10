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

/*
Package client contains generic representations of clients connecting to different receivers. Components,
such as processors or exporters, can make use of this information to make decisions related to
grouping of batches, tenancy, load balancing, tagging, among others.

The structs defined here are typically used within the context that is propagated down the pipeline,
with the values being produced by authenticators or receivers, and consumed by processors and exporters.

Producers

Authenticators are responsible for obtaining a client.Info from the current context and enhancing the client.Info
with an implementation of client.AuthData. The attribute names should be documented with their return types and
considered part of the public API for the authenticator.

Receivers are responsible for obtaining a client.Info from the current context and enhancing the client.Info
with the net.Addr from the peer. For HTTP requests, this is typically the IP address of the client. This new context
should then be propagated down the pipeline when calling the next consumer. Typically, however, receivers would delegate
this processing to helpers such as the confighttp or configgrpc packages: both contain interceptors that will enhance
the context with the client.Info, such that no actions are needed by receivers that are built using confighttp.HTTPServerSettings
or configgrpc.GRPCServerSettings.

Consumers

Processors and exporters have access to the client.Info via client.FromContext. Among other usages, this data can be used to:

- annotate data points with authentication data (username, tenant, ...)

- route data points based on authentication data

- rate limit client calls based on IP addresses

Neither processors nor exporters are expected to change data from the client.Info contained in the context, although
there might be reasons to do so. In that case, this should be clearly documented as part of the component's README file.

Processors and exporters relying on the existence of data from the client.Info, especially client.AuthData, should
clearly document this as part of the component's README file. The expected pattern for consuming data is to allow
users to specify the attribute name to use in the component. The expected data type should also be communicated to
users, who should then compare this with the authenticators that are part of the pipeline.

For example, assuming that the OIDC authenticator pushes a subject string attribute and that we have a hypothetical
authprinter processor that prints the username to the console, this is how an OpenTelemetry Collector configuration
would look like:

  extensions:
    oidc:
      issuer_url: http://localhost:8080/auth/realms/opentelemetry
      audience: collector

  receivers:
    otlp:
      protocols:
        grpc:
          auth:
            authenticator: oidc

  processors:
    authprinter:
      attribute: subject

  exporters:
    logging:

  service:
    extensions: [oidc]
    pipelines:
      traces:
        receivers: [otlp]
        processors: [authprinter]
        exporters: [logging]
*/
package client // import "go.opentelemetry.io/collector/client"

import (
	"context"
	"net"
)

type ctxKey struct{}

// Info contains data related to the clients connecting to receivers.
type Info struct {
	// Addr for the client connecting to this collector. Available in a best-effort basis, and generally reliable
	// for receivers making use of confighttp.ToServer and configgrpc.ToServerOption.
	Addr net.Addr

	// Auth information from the incoming request as provided by configauth.ServerAuthenticator implementations
	// tied to the receiver for this connection.
	Auth AuthData
}

// AuthData represents the authentication data as seen by authenticators tied to the receivers.
type AuthData interface {
	// GetAttribute returns the value for the given attribute. Authenticator implementations might define different
	// data types for different attributes. While "string" is used most of the time, a key named "membership" might
	// return a list of strings.
	GetAttribute(string) interface{}

	// GetAttributes returns the names of all attributes in this authentication data.
	GetAttributeNames() []string
}

// NewContext takes an existing context and derives a new context with the client.Info value stored on it.
func NewContext(ctx context.Context, c *Info) context.Context {
	return context.WithValue(ctx, ctxKey{}, c)
}

// FromContext takes a context and returns a ClientInfo from it. When a ClientInfo isn't present, a new
// empty one is returned.
func FromContext(ctx context.Context) *Info {
	c, ok := ctx.Value(ctxKey{}).(*Info)
	if !ok {
		c = &Info{}
	}
	return c
}
