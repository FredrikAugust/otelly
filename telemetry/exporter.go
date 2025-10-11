package telemetry

import (
	"context"
	"log/slog"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/exporter"
	"go.opentelemetry.io/collector/exporter/exporterhelper"
	"go.opentelemetry.io/collector/pdata/ptrace"
)

const ExporterName = "otelly"

func createOtellyExporter() exporter.Factory {
	return exporter.NewFactory(
		component.MustNewType(ExporterName),
		func() component.Config {
			return nil
		},
		exporter.WithTraces(createTraces, component.StabilityLevelDevelopment),
	)
}

func createTraces(ctx context.Context, set exporter.Settings, cfg component.Config) (exporter.Traces, error) {
	return exporterhelper.NewTraces(
		ctx,
		set,
		cfg,
		func(ctx context.Context, td ptrace.Traces) error {
			slog.Debug("new trace", "spanCount", td.SpanCount())
			return nil
		},
	)
}
