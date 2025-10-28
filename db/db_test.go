package db_test

import (
	"testing"

	"github.com/fredrikaugust/otelly/db"
	"github.com/stretchr/testify/assert"
)

func TestDB(t *testing.T) {
	t.Run("creates in-mem db and migrates", func(t *testing.T) {
		db, err := getDB(t)
		if err != nil {
			// This will fail, we just do it like this for the formatting
			assert.Nil(t, err)
		}
		defer db.Close()
	})
}

func getDB(t *testing.T) (*db.Database, error) {
	t.Helper()

	db, err := db.NewDB(":memory:")
	if err != nil {
		return nil, err
	}
	err = db.Migrate(t.Context())
	if err != nil {
		return nil, err
	}
	return db, nil
}
