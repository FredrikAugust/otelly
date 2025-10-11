package telemetry

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/confmap"
	"go.opentelemetry.io/collector/confmap/provider/fileprovider"
	"go.opentelemetry.io/collector/otelcol"
)

func Start(ctx context.Context) error {
	col, err := otelcol.NewCollector(otelcol.CollectorSettings{
		Factories: createCollectorFactories,
		BuildInfo: component.NewDefaultBuildInfo(),
		ConfigProviderSettings: otelcol.ConfigProviderSettings{
			ResolverSettings: confmap.ResolverSettings{
				URIs: []string{"./telemetry/config.yml"},
				ProviderFactories: []confmap.ProviderFactory{
					fileprovider.NewFactory(),
				},
			},
		},
	})
	if err != nil {
		return err
	}

	return col.Run(ctx)
}
