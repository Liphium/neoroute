package tui

import (
	"fmt"
	"time"

	"github.com/Unbreathable/nowandlater"
)

var _ SchemaNode = &DateNode{}

// DateNode edits a date value using natural language dates (e.g. "tomorrow", "next week", ...)
type DateNode struct {
	ValueNode[time.Time]
}

func newDateNode() *DateNode {
	parser := nowandlater.Parser{}

	return &DateNode{
		prefix: "date(\"",
		suffix: "\")",
		value:  time.Now(),
		comment: func(t time.Time) string {
			return t.Format("January 2, 2006 15:04")
		},
		convert: func(date string) (time.Time, error) {

			t, err := parser.Parse(date)
			if err != nil {
				return time.Time{}, fmt.Errorf("Could not parse date: %w", err)
			}

			// Format the time to a more readable form
			return t, nil
		},
	}
}

func (f *DateNode) Init() {
	f.ValueNode.Init()
	// Display the path, not the byte array's formatted value.
	f.input.SetValue("")
}
