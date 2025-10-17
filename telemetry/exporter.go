package telemetry

import (
	"context"

	"github.com/fredrikaugust/otelly/bus"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/exporter"
	"go.opentelemetry.io/collector/exporter/exporterhelper"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/ptrace"
)

const ExporterName = "otelly"

type traceConfig struct {
	bus *bus.TransportBus
}

func createOtellyExporter(bus *bus.TransportBus) exporter.Factory {
	return exporter.NewFactory(
		component.MustNewType(ExporterName),
		func() component.Config {
			return &traceConfig{
				bus: bus,
			}
		},
		exporter.WithTraces(createTraces, component.StabilityLevelDevelopment),
		exporter.WithLogs(createLogs, component.StabilityLevelDevelopment),
	)
}

func createTraces(ctx context.Context, set exporter.Settings, cfg component.Config) (exporter.Traces, error) {
	return exporterhelper.NewTraces(
		ctx,
		set,
		cfg,
		func(ctx context.Context, td ptrace.Traces) error {
			bus := cfg.(*traceConfig).bus

			return traceReceiver(ctx, td, bus)
		},
	)
}

func createLogs(ctx context.Context, set exporter.Settings, cfg component.Config) (exporter.Logs, error) {
	return exporterhelper.NewLogs(
		ctx,
		set,
		cfg,
		func(ctx context.Context, ld plog.Logs) error {
			bus := cfg.(*traceConfig).bus

			return logReceiver(ctx, ld, bus)
		},
	)
}

func traceReceiver(_ context.Context, td ptrace.Traces, bus *bus.TransportBus) error {
	for _, resourceSpans := range td.ResourceSpans().All() {
		bus.TraceBus <- resourceSpans
	}

	return nil
}

func logReceiver(_ context.Context, ld plog.Logs, bus *bus.TransportBus) error {
	for _, resourceLogs := range ld.ResourceLogs().All() {
		bus.LogBus <- resourceLogs
	}

	return nil
}
