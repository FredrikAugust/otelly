package main

import (
	"context"
	"io"
	"log/slog"
	"os"

	"github.com/fredrikaugust/otelly/telemetry"
	"github.com/fredrikaugust/otelly/ui"
)

func main() {
	logFile, err := os.OpenFile("debug.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		slog.Error("could not set up logger", "error", err)
		os.Exit(1)
	}
	defer logFile.Close()

	logger := slog.New(
		slog.NewTextHandler(
			io.MultiWriter(
				logFile,
				os.Stdout,
			),
			&slog.HandlerOptions{
				AddSource: false,
				Level:     slog.LevelDebug,
			},
		),
	)
	slog.SetDefault(logger)
	slog.Info("logger initialized")

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		if err = telemetry.Start(ctx); err != nil {
			slog.Error("failed to start receiver", "error", err)
			cancel()
		}
	}()

	if err = ui.Start(ctx); err != nil {
		slog.Error("failed to start ui", "error", err)
	}

	cancel()

	<-ctx.Done()
	slog.Info("application quit successfully")
}
