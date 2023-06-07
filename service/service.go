// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package service // import "go.opentelemetry.io/collector/service"

import (
	"context"
	"fmt"
	"runtime"

	"go.uber.org/multierr"
	"go.uber.org/zap"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config/configtelemetry"
	"go.opentelemetry.io/collector/connector"
	"go.opentelemetry.io/collector/exporter"
	"go.opentelemetry.io/collector/extension"
	"go.opentelemetry.io/collector/internal/obsreportconfig"
	"go.opentelemetry.io/collector/processor"
	"go.opentelemetry.io/collector/receiver"
	"go.opentelemetry.io/collector/service/extensions"
	"go.opentelemetry.io/collector/service/internal/graph"
	"go.opentelemetry.io/collector/service/internal/proctelemetry"
	"go.opentelemetry.io/collector/service/telemetry"
)

// Settings holds configuration for building a new service.
type Settings struct {
	// BuildInfo provides collector start information.
	BuildInfo component.BuildInfo

	// Receivers builder for receivers.
	Receivers *receiver.Builder

	// Processors builder for processors.
	Processors *processor.Builder

	// Exporters builder for exporters.
	Exporters *exporter.Builder

	// Connectors builder for connectors.
	Connectors *connector.Builder

	// Extensions builder for extensions.
	Extensions *extension.Builder

	// AsyncErrorChannel is the channel that is used to report fatal errors.
	AsyncErrorChannel chan error

	// LoggingOptions provides a way to change behavior of zap logging.
	LoggingOptions []zap.Option
}

// Service represents the implementation of a component.Host.
type Service struct {
	buildInfo         component.BuildInfo
	telemetry         *telemetry.Telemetry
	telemetrySettings component.TelemetrySettings
	host              *serviceHost
}

func New(ctx context.Context, set Settings, cfg Config) (*Service, error) {
	srv := &Service{
		buildInfo: set.BuildInfo,
		host: &serviceHost{
			receivers:         set.Receivers,
			processors:        set.Processors,
			exporters:         set.Exporters,
			connectors:        set.Connectors,
			extensions:        set.Extensions,
			buildInfo:         set.BuildInfo,
			asyncErrorChannel: set.AsyncErrorChannel,
		},
	}
	var err error
	srv.telemetry, err = telemetry.New(ctx, telemetry.Settings{BuildInfo: set.BuildInfo, ZapOptions: set.LoggingOptions}, cfg.Telemetry)
	if err != nil {
		return nil, fmt.Errorf("failed to get logger: %w", err)
	}

	srv.telemetrySettings = component.TelemetrySettings{
		Logger:         srv.telemetry.Logger(),
		TracerProvider: srv.telemetry.TracerProvider(),
		MeterProvider:  srv.telemetry.MeterProvider(),
		Resource:       srv.telemetry.Resource(),
	}

	// process the configuration and initialize the pipeline
	if err = srv.initExtensionsAndPipeline(ctx, set, cfg); err != nil {
		// If pipeline initialization fails then shut down the telemetry.
		if shutdownErr := srv.telemetry.Shutdown(ctx); shutdownErr != nil {
			err = multierr.Append(err, fmt.Errorf("failed to shutdown collector telemetry: %w", shutdownErr))
		}

		return nil, err
	}

	if cfg.Telemetry.Metrics.Level != configtelemetry.LevelNone && cfg.Telemetry.Metrics.Address != "" {
		// The process telemetry initialization requires the ballast size, which is available after the extensions are initialized.
		if err = proctelemetry.RegisterProcessMetrics(srv.telemetryInitializer.ocRegistry, srv.telemetryInitializer.mp, obsreportconfig.UseOtelForInternalMetricsfeatureGate.IsEnabled(), getBallastSize(srv.host)); err != nil {
			return nil, fmt.Errorf("failed to register process metrics: %w", err)
		}
	}

	return srv, nil
}

// Start starts the extensions and pipelines. If Start fails Shutdown should be called to ensure a clean state.
func (srv *Service) Start(ctx context.Context) error {
	srv.telemetrySettings.Logger.Info("Starting "+srv.buildInfo.Command+"...",
		zap.String("Version", srv.buildInfo.Version),
		zap.Int("NumCPU", runtime.NumCPU()),
	)

	if err := srv.host.serviceExtensions.Start(ctx, srv.host); err != nil {
		return fmt.Errorf("failed to start extensions: %w", err)
	}

	if err := srv.host.pipelines.StartAll(ctx, srv.host); err != nil {
		return fmt.Errorf("cannot start pipelines: %w", err)
	}

	if err := srv.host.serviceExtensions.NotifyPipelineReady(); err != nil {
		return err
	}

	srv.telemetrySettings.Logger.Info("Everything is ready. Begin running and processing data.")
	return nil
}

func (srv *Service) Shutdown(ctx context.Context) error {
	// Accumulate errors and proceed with shutting down remaining components.
	var errs error

	// Begin shutdown sequence.
	srv.telemetrySettings.Logger.Info("Starting shutdown...")

	if err := srv.host.serviceExtensions.NotifyPipelineNotReady(); err != nil {
		errs = multierr.Append(errs, fmt.Errorf("failed to notify that pipeline is not ready: %w", err))
	}

	if err := srv.host.pipelines.ShutdownAll(ctx); err != nil {
		errs = multierr.Append(errs, fmt.Errorf("failed to shutdown pipelines: %w", err))
	}

	if err := srv.host.serviceExtensions.Shutdown(ctx); err != nil {
		errs = multierr.Append(errs, fmt.Errorf("failed to shutdown extensions: %w", err))
	}

	srv.telemetrySettings.Logger.Info("Shutdown complete.")

	if err := srv.telemetry.Shutdown(ctx); err != nil {
		errs = multierr.Append(errs, fmt.Errorf("failed to shutdown service telemetry: %w", err))
	}

	return errs
}

func (srv *Service) initExtensionsAndPipeline(ctx context.Context, set Settings, cfg Config) error {
	var err error
	srv.host.serviceExtensions, err = extensions.New(
		ctx,
		extensions.Settings{Telemetry: srv.telemetrySettings, BuildInfo: srv.buildInfo, Extensions: srv.host.extensions},
		cfg.Extensions)
	if err != nil {
		return fmt.Errorf("failed to build extensions: %w", err)
	}

	graphPipelinesConfigs := make(map[component.ID]*graph.PipelineConfig, len(cfg.Pipelines))
	for id, cfg := range cfg.Pipelines {
		graphPipelinesConfigs[id] = &graph.PipelineConfig{
			Receivers:  cfg.Receivers,
			Processors: cfg.Processors,
			Exporters:  cfg.Exporters,
		}
	}

	pSet := graph.Settings{
		Telemetry:        srv.telemetrySettings,
		BuildInfo:        srv.buildInfo,
		ReceiverBuilder:  set.Receivers,
		ProcessorBuilder: set.Processors,
		ExporterBuilder:  set.Exporters,
		ConnectorBuilder: set.Connectors,
		PipelineConfigs:  graphPipelinesConfigs,
	}

	if srv.host.pipelines, err = graph.Build(ctx, pSet); err != nil {
		return fmt.Errorf("failed to build pipelines: %w", err)
	}

	return nil
}

// Logger returns the logger created for this service.
// This is a temporary API that may be removed soon after investigating how the collector should record different events.
func (srv *Service) Logger() *zap.Logger {
	return srv.telemetrySettings.Logger
}

func getBallastSize(host component.Host) uint64 {
	for _, ext := range host.GetExtensions() {
		if bExt, ok := ext.(interface{ GetBallastSize() uint64 }); ok {
			return bExt.GetBallastSize()
		}
	}
	return 0
}
