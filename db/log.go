package db

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"go.opentelemetry.io/collector/pdata/plog"
	"go.uber.org/zap"
)

func (d *Database) InsertResourceLogs(ctx context.Context, logs plog.ResourceLogs) error {
	resID, err := d.InsertResource(ctx, logs.Resource())
	if err != nil {
		return err
	}

	tx, err := d.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, scopeLogs := range logs.ScopeLogs().All() {
		for _, logRecord := range scopeLogs.LogRecords().All() {
			attrs, err := json.Marshal(logRecord.Attributes().AsRaw())
			if err != nil {
				zap.L().Warn("could not serialize log attributes to JSON", zap.Error(err))
				attrs = []byte("{}")
			}

			zap.L().Debug("inserting new log record", zap.String("body", logRecord.Body().Str()))

			_, err = tx.ExecContext(
				ctx,
				`INSERT INTO log VALUES (?, ?, ?, ?, ?, ?, ?)`,
				sql.NullString{String: logRecord.SpanID().String(), Valid: !logRecord.SpanID().IsEmpty()},
				logRecord.Body().Str(),
				logRecord.Timestamp().AsTime(),
				logRecord.SeverityNumber(),
				logRecord.SeverityText(),
				resID,
				attrs,
			)
			if err != nil {
				zap.L().Warn("failed to create log", zap.String("body", logRecord.Body().Str()))
			}
		}
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction log: %w", err)
	}

	return nil
}

func (d *Database) ClearLogs(ctx context.Context) error {
	_, err := d.ExecContext(ctx, `TRUNCATE TABLE log`)
	if err != nil {
		return err
	}

	return nil
}

func (d *Database) GetLogs(ctx context.Context) ([]Log, error) {
	logs := make([]Log, 0)
	err := d.sqlDB.SelectContext(
		ctx,
		&logs,
		`
		SELECT
			*
		FROM
			log
		ORDER BY
			timestamp DESC`,
	)
	if err != nil {
		return logs, err
	}

	return logs, nil
}
