package model

import tea "charm.land/bubbletea/v2"

// Allows you to return multiple messages with help of a batch msg
func Multiple(msgs ...tea.Msg) tea.Msg {
	batch := make(tea.BatchMsg, len(msgs))
	for i, msg := range msgs {
		batch[i] = func() tea.Msg { return msg }
	}
	return batch
}
