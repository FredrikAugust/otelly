// Package telemetry deals with receiving OTEL traces, doing any
// transformations and storing them in a local store.
package telemetry

import (
	"github.com/fredrikaugust/otelly/bus"
	"github.com/fredrikaugust/otelly/db"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/otlpjsonfilereceiver"
	"go.opentelemetry.io/collector/otelcol"
	"go.opentelemetry.io/collector/receiver/otlpreceiver"
)

// Builds factories which amount to ghe components needed
// to set up a collector.
func createCollectorFactories(bus *bus.TransportBus, db *db.Database) (otelcol.Factories, error) {
	var err error

	factories := otelcol.Factories{}

	factories.Receivers, err = otelcol.MakeFactoryMap(
		otlpreceiver.NewFactory(),
		otlpjsonfilereceiver.NewFactory(),
	)
	if err != nil {
		return otelcol.Factories{}, err
	}

	factories.Exporters, err = otelcol.MakeFactoryMap(
		createOtellyExporter(bus, db),
	)
	if err != nil {
		return otelcol.Factories{}, err
	}

	return factories, nil
}
