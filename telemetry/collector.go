package telemetry

import (
	"context"
	"log/slog"

	"github.com/fredrikaugust/otelly/bus"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/confmap"
	"go.opentelemetry.io/collector/confmap/provider/fileprovider"
	"go.opentelemetry.io/collector/otelcol"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Start(ctx context.Context, bus *bus.TransportBus) error {
	slog.Info("starting collector")
	col, err := otelcol.NewCollector(otelcol.CollectorSettings{
		Factories: func() (otelcol.Factories, error) {
			return createCollectorFactories(bus)
		},
		BuildInfo: component.NewDefaultBuildInfo(),
		LoggingOptions: []zap.Option{
			zap.WrapCore(func(c zapcore.Core) zapcore.Core {
				return zap.NewNop().Core()
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
