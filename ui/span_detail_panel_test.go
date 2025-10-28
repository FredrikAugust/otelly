package ui_test

import (
	"testing"

	"github.com/fredrikaugust/otelly/db"
	"github.com/fredrikaugust/otelly/ui"
	"github.com/stretchr/testify/assert"
)

func TestUpdateSpan(t *testing.T) {
	t.Run("set nil span when empty", func(t *testing.T) {
		m := ui.NewSpanDetailPanelModel(nil)

		_, cmd := m.UpdateSpan(nil)
		assert.Nil(t, cmd)
	})

	t.Run("set span", func(t *testing.T) {
		m := ui.NewSpanDetailPanelModel(nil)

		_, cmd := m.UpdateSpan(&db.Span{ID: "test-id"})

		assert.NotNil(t, cmd)
		assert.IsType(t, ui.MsgLoadTrace{}, cmd())
	})
}
