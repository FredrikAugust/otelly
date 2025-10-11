// Package telemetry deals with receiving OTEL traces, doing any
// transformations and storing them in a local store.
package telemetry

import (
	"context"

	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/otlpjsonfilereceiver"
	"go.opentelemetry.io/collector/otelcol"
	"go.opentelemetry.io/collector/receiver/otlpreceiver"
)

type Store struct{}

func (s *Store) Receive(ctx context.Context) error {
	return nil
}

// Builds factories which amount to ghe components needed
// to set up a collector.
func createCollectorFactories() (otelcol.Factories, error) {
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
		createOtellyExporter(),
	)
	if err != nil {
		return otelcol.Factories{}, err
	}

	return factories, nil
}
