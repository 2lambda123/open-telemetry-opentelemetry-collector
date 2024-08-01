// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

// Package configtelemetry defines various telemetry level for configuration.
// It enables every component to have access to telemetry level
// to enable metrics only when necessary.
//
// This document provides guidance on which telemetry level to adopt for Collector metrics.
// When adopting a telemetry level, the developer is expected to rely on this guidance to
// justify their choice of telemetry level.
//
// 1. configtelemetry.None
//
// No telemetry data should be collected.
//
// 2. configtelemetry.Basic
//
// Signals associated with this level cover the essential coverage of the component telemetry.
//
// This is the default level recommended when running the Collector.
//
// Signals using this telemetry level can use this guidance:
// * The signals associated with this level must show low cardinality, the number of combinations of dimension values.
// * The signals associated with this level must represent a small data volume. Examples follow:
//   - A max cardinality (total possible combinations of dimension values) of 50.
//   - At most a span actively recording simultaneously, covering the critical path.
//
// * Not all signals defined in the component telemetry are active.
//
// 3. configtelemetry.Normal
//
// Signals associated with this level cover the complete coverage of the component telemetry.
//
// Signals using this telemetry level can use this guidance:
//   - The signals associated with this level must control cardinality.
//     It is acceptable at this level for cardinality to scale linearly with the monitored resources.
//   - The signals associated with this level must represent a controlled data volume. Examples follow:
//   - A max cardinality (total possible combinations of dimension values) of 500.
//   - At most 5 spans actively recording simultaneously.
//
// * All signals defined in the component telemetry are active.
//
// 4. configtelemetry.Detailed
//
// Signals associated with this level cover the complete coverage of the component telemetry.
//
// The signals associated with this level may exhibit high cardinality.
//
// There is no limit on data volume.
//
// All signals defined in the component telemetry are active.
package configtelemetry // import "go.opentelemetry.io/collector/config/configtelemetry"
