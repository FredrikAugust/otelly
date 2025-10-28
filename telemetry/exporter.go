package telemetry

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/fredrikaugust/otelly/bus"
	"github.com/fredrikaugust/otelly/db"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/exporter"
	"go.opentelemetry.io/collector/exporter/exporterhelper"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/ptrace"
	"go.uber.org/zap"
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
	var wg sync.WaitGroup

	zap.L().Debug("received spans", zap.Int("spanCount", td.SpanCount()))

	for _, resourceSpans := range td.ResourceSpans().All() {
		wg.Go(func() {
			err := db.InsertResourceSpans(ctx, resourceSpans)
			if err != nil {
				zap.L().Warn("could not insert resource spans", zap.Error(err))
			}
		})
	}

	wg.Wait()

	// TODO: only send new ones
	spans, err := db.GetSpans(ctx)
	if err != nil {
		return err
	}

	select {
	case bus.SpanBus <- spans:
		return nil
	case <-time.After(1 * time.Second):
		return errors.New("trace receiver timed out after 1 second")
	}
}

func logReceiver(ctx context.Context, ld plog.Logs, bus *bus.TransportBus, db *db.Database) error {
	var wg sync.WaitGroup

	zap.L().Debug("received logs", zap.Int("logRecordCount", ld.LogRecordCount()))

	for _, resourceLogs := range ld.ResourceLogs().All() {
		resourceLogs := resourceLogs // capture loop variable
		wg.Go(func() {
			err := db.InsertResourceLogs(ctx, resourceLogs)
			if err != nil {
				zap.L().Warn("could not insert resource logs", zap.Error(err))
			}
		})
	}

	wg.Wait()

	logs, err := db.GetLogs(ctx)
	if err != nil {
		return err
	}

	select {
	case bus.LogBus <- logs:
		return nil
	case <-time.After(1 * time.Second):
		return errors.New("log receiver timed out after 1 second")
	}
}
