package db

import (
	"context"
	"database/sql"
	"encoding/json"

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
				return err
			}
		}
	}

	return tx.Commit()
}

func (d *Database) ClearLogs() error {
	_, err := d.sqlDB.Exec(`TRUNCATE TABLE log`)
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
