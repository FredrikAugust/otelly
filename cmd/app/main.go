package main

import (
	"context"
	"log/slog"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/fredrikaugust/otelly/bus"
	"github.com/fredrikaugust/otelly/db"
	"github.com/fredrikaugust/otelly/telemetry"
	"github.com/fredrikaugust/otelly/ui"
	slogzap "github.com/samber/slog-zap/v2"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	cleanup := configureLogging()
	defer cleanup()

	ctx, cancel := context.WithCancel(context.Background())

	bus := bus.NewTransportBus()

	db, err := configureDB(ctx)
	if err != nil {
		slog.Error("couldn't configure DB", "error", err)
		return
	}
	defer db.Close()

	go func() {
		if err := telemetry.Start(ctx, bus); err != nil {
			slog.Error("failed to start receiver", "error", err)
			cancel()
		}
	}()

	p := tea.NewProgram(ui.NewEntryModel(), tea.WithAltScreen(), tea.WithContext(ctx))
	if _, err := p.Run(); err != nil {
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

	logCfg := zap.NewDevelopmentConfig()
	logCfg.EncoderConfig.EncodeTime = func(t time.Time, pae zapcore.PrimitiveArrayEncoder) {
		pae.AppendString(t.Format("15:04:05.000"))
	}
	logCfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	logCfg.OutputPaths = []string{
		"./debug.log",
	}
	zapLogger, _ := logCfg.Build()
	zap.ReplaceGlobals(zapLogger)

	slogLogger := slog.New(slogzap.Option{Level: slog.LevelDebug, Logger: zapLogger}.NewZapHandler())
	slog.SetDefault(slogLogger)

	slog.Info("logger initialized")

	return logFile.Close
}

func configureDB(ctx context.Context) (*db.Database, error) {
	db, err := db.NewDB()
	if err != nil {
		return nil, err
	}
	err = db.Migrate(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}
