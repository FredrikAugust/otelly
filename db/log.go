package db

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/plog"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
	"go.uber.org/zap"
)

func (d *Database) InsertResourceLogs(logs plog.ResourceLogs) error {
	resName, exists := logs.Resource().Attributes().Get(string(semconv.ServiceNameKey))
	if !exists {
		resName = pcommon.NewValueStr("unknown")
	}
	resNamespace, exists := logs.Resource().Attributes().Get(string(semconv.ServiceNamespaceKey))
	if !exists {
		resNamespace = pcommon.NewValueStr("unknown")
	}
	resID := fmt.Sprintf("%s:%s", resName.Str(), resNamespace.Str())

	tx, err := d.BeginTx(context.Background())
	if err != nil {
		return err
	}

	_, err = tx.Exec(`INSERT OR IGNORE INTO resource VALUES ($1, $2, $3)`,
		resID,
		resName.Str(),
		resNamespace.Str(),
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	for _, scopeLogs := range logs.ScopeLogs().All() {
		for _, logRecord := range scopeLogs.LogRecords().All() {
			attrs, err := json.Marshal(logRecord.Attributes().AsRaw())
			if err != nil {
				attrs = []byte("{}")
			}

			zap.L().Debug("inserting new log record", zap.String("body", logRecord.Body().Str()))

			_, err = tx.Exec(
				`
				INSERT INTO log VALUES (
					$1,
					$2,
					$3,
					$4,
					$5,
					$6,
					$7
				)`,
				sql.NullString{String: logRecord.SpanID().String(), Valid: !logRecord.SpanID().IsEmpty()},
				logRecord.Body().Str(),
				logRecord.Timestamp().AsTime(),
				logRecord.SeverityNumber(),
				logRecord.SeverityText(),
				resID,
				attrs,
			)
			if err != nil {
				tx.Rollback()
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

func (d *Database) GetLogs() ([]Log, error) {
	logs := make([]Log, 0)
	err := d.sqlDB.Select(
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
