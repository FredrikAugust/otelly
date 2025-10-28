package db

import (
	"context"
	"fmt"

	"go.opentelemetry.io/collector/pdata/pcommon"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
)

func (d *Database) GetResource(id string) (*Resource, error) {
	var res Resource

	err := d.sqlDB.Get(
		&res,
		`
		SELECT
			*
		FROM
			resource
		WHERE id = ?`,
		id,
	)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

// InsertResource returns the resource ID and error
func (d *Database) InsertResource(ctx context.Context, res pcommon.Resource) (string, error) {
	// See Database struct definition for why we do this
	d.resourceLock.Lock()
	defer d.resourceLock.Unlock()

	resName, exists := res.Attributes().Get(string(semconv.ServiceNameKey))
	if !exists {
		resName = pcommon.NewValueStr("unknown")
	}
	resNamespace, exists := res.Attributes().Get(string(semconv.ServiceNamespaceKey))
	if !exists {
		resNamespace = pcommon.NewValueStr("unknown")
	}
	resID := fmt.Sprintf("%s:%s", resName.Str(), resNamespace.Str())

	_, err := d.ExecContext(ctx, `INSERT OR REPLACE INTO resource VALUES (?, ?, ?)`,
		resID,
		resName.Str(),
		resNamespace.Str(),
	)

	return resID, err
}
