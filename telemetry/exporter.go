package telemetry

import (
	"context"
	"errors"

	"github.com/fredrikaugust/otelly/bus"
	"github.com/fredrikaugust/otelly/db"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/exporter"
	"go.opentelemetry.io/collector/exporter/exporterhelper"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/ptrace"
)

const ExporterName = "otelly"

type traceConfig struct {
	bus *bus.TransportBus
	db  *db.Database
}

func createOtellyExporter(bus *bus.TransportBus, db *db.Database) exporter.Factory {
	return exporter.NewFactory(
		component.MustNewType(ExporterName),
		func() component.Config {
			return &traceConfig{
				bus: bus,
				db:  db,
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
			db := cfg.(*traceConfig).db

			return traceReceiver(ctx, td, bus, db)
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
			db := cfg.(*traceConfig).db

			return logReceiver(ctx, ld, bus, db)
		},
	)
}

func traceReceiver(ctx context.Context, td ptrace.Traces, bus *bus.TransportBus, db *db.Database) error {
	var err error

	for _, resourceSpans := range td.ResourceSpans().All() {
		bus.TraceBus <- resourceSpans
		e := db.InsertResourceSpans(ctx, resourceSpans)
		if e != nil {
			err = errors.Join(err, e)
		}
	}

	return err
}

func logReceiver(ctx context.Context, ld plog.Logs, bus *bus.TransportBus, db *db.Database) error {
	var err error

	for _, resourceLogs := range ld.ResourceLogs().All() {
		bus.LogBus <- resourceLogs
		e := db.InsertResourceLogs(ctx, resourceLogs)
		if e != nil {
			err = errors.Join(err, e)
		}
	}

	return err
}
