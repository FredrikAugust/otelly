package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/fredrikaugust/otelly/telemetry"
	"github.com/fredrikaugust/otelly/ui"
	slogzap "github.com/samber/slog-zap/v2"
	"go.uber.org/zap"
)

func main() {
	configureLogging()

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		if err := telemetry.Start(ctx); err != nil {
			slog.Error("failed to start receiver", "error", err)
			cancel()
		}
	}()

	if err := ui.Start(ctx); err != nil {
		slog.Error("failed to start ui", "error", err)
	}

	cancel()

	<-ctx.Done()
	slog.Info("application quit successfully")
}

func configureLogging() func() error {
	logFile, err := os.OpenFile("debug.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		slog.Error("could not set up logger", "error", err)
		os.Exit(1)
	}
	defer logFile.Close()

	logCfg := zap.NewDevelopmentConfig()
	logCfg.OutputPaths = []string{
		"./debug.log",
	}
	zapLogger, _ := logCfg.Build()

	restoreGlobalZapLogger := zap.ReplaceGlobals(zapLogger)
	defer restoreGlobalZapLogger()

	slogLogger := slog.New(slogzap.Option{Level: slog.LevelDebug, Logger: zapLogger}.NewZapHandler())
	slog.SetDefault(slogLogger)

	slog.Info("logger initialized")

	return logFile.Close
}
