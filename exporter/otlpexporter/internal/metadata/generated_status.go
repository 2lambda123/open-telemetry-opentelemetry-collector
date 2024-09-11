// Code generated by mdatagen. DO NOT EDIT.

package metadata

import (
	"go.opentelemetry.io/collector/component"
)

var (
	Type      = component.MustNewType("otlp")
	ScopeName = "go.opentelemetry.io/collector/exporter/otlpexporter"
)

const (
	LogsStability     = component.StabilityLevelBeta
	ProfilesStability = component.StabilityLevelBeta
	TracesStability   = component.StabilityLevelStable
	MetricsStability  = component.StabilityLevelStable
)
