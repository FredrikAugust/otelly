package telemetry

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/confmap"
	"go.opentelemetry.io/collector/confmap/provider/fileprovider"
	"go.opentelemetry.io/collector/otelcol"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Start(ctx context.Context) error {
	col, err := otelcol.NewCollector(otelcol.CollectorSettings{
		Factories: createCollectorFactories,
		BuildInfo: component.NewDefaultBuildInfo(),
		LoggingOptions: []zap.Option{
			zap.WrapCore(func(c zapcore.Core) zapcore.Core {
				return zap.L().Core()
			}),
		},
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
