package tui

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDateInput(t *testing.T) {
	node := newDateNode()
	node.Init()
	assert.Empty(t, node.input.Value())

	node.input.SetValue("tomorrow")
	if !assert.NoError(t, node.input.Err) {
		return
	}
	date := node.Request().(time.Time)
	assert.False(t, date.IsZero())
	assert.Equal(t, "tomorrow", node.input.Value())
	_, view := node.View()
	assert.Contains(t, view, " // "+date.Format("January 2, 2006 15:04"))

	node.input.SetValue("not a date")
	assert.Error(t, node.input.Err)
	assert.Equal(t, 2, node.Height())

	node.input.SetValue("tomorrow")
	assert.NoError(t, node.input.Err)
	assert.Equal(t, 1, node.Height())
}
