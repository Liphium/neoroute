package tui

import (
	tea "charm.land/bubbletea/v2"
)

type SchemaNode interface {
	keyProvider // All of them have their own keys though

	// This should be called to initialize all of the selection functions.
	Init()

	// Should start the selection at the first possible position.
	SelectFromTop()

	// Should start the selection at the last possible position.
	SelectFromBottom()

	// Should take a callback that is called when the upper boundary of the selection is reached and the selection should go up.
	//
	// This should also unselect this node at the same time.
	OnUp(func())

	// Should take a callback that is called when the lower boundary of the selection is reached and the seelction should go down.
	//
	// This should also unselect this node at the same time.
	OnDown(func())

	// Should set a suffix for the entire thing
	SetSuffix(suffix string)

	// Should return the height of the node together with its children
	Height() int

	// Should return the selected position, within the height it returned (0 = not selected).
	Selected() int

	// Should update the node based on the message and return a cursor in case there is a text field (or sth else) in this one.
	Update(msg tea.Msg) tea.Cmd

	// Should render the entire view.
	View() (*tea.Cursor, string)
}

// A struct that can be embedded to implement basic functions for the schema interface.
type basicNode struct {
	Suffix string
	Down   func()
	Up     func()
}

func (b basicNode) GoUp() {
	if b.Up == nil {
		panic("no up function set") // This is important because it should ALWAYS BE SET
	}
	b.Up()
}

func (b basicNode) GoDown() {
	if b.Down == nil {
		panic("no up function set") // This is important because it should ALWAYS BE SET
	}
	b.Down()
}

func (b *basicNode) OnUp(up func()) {
	b.Up = up
}

func (b *basicNode) OnDown(down func()) {
	b.Down = down
}

func (b *basicNode) SetSuffix(suffix string) {
	b.Suffix = suffix
}
