package db_test

import (
	"database/sql"
	"testing"

	"github.com/fredrikaugust/otelly/db"
	"github.com/stretchr/testify/assert"
)

func TestFilterRootSpans(t *testing.T) {
	t.Run("filters root spans", func(t *testing.T) {
		spans := []db.Span{
			{ParentSpanID: sql.NullString{Valid: false}, ID: "test"},
			{ParentSpanID: sql.NullString{Valid: true, String: "test"}, ID: "test1"},
			{ParentSpanID: sql.NullString{Valid: true, String: "test1"}, ID: "test2"},
			{ParentSpanID: sql.NullString{Valid: true, String: "test2"}, ID: "test3"},
		}

		rootSpans := db.FilterRootSpans(spans)

		assert.Equal(t, "test", rootSpans[0].ID)
		assert.Len(t, rootSpans, 1)
	})
}
