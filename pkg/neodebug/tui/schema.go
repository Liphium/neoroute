package tui

import (
	tea "charm.land/bubbletea/v2"
)

type SchemaNode interface {

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

	// Should return the height of the node together with its children
	Height() int

	// Should return the selected position, within the height it returned (0 = not selected).
	Selected() int

	// Should update the node based on the message and return a cursor in case there is a text field (or sth else) in this one.
	Update(msg tea.Msg) (SchemaNode, tea.Cmd)

	// Should render the entire view.
	View() (*tea.Cursor, string)
}

// A struct that can be embedded to store the basic functions for going up and down.
type basicSelection struct {
	Down func()
	Up   func()
}

func (b basicSelection) GoUp() {
	if b.Up == nil {
		panic("no up function set") // This is important because it should ALWAYS BE SET
	}
	b.Up()
}

func (b basicSelection) GoDown() {
	if b.Down == nil {
		panic("no up function set") // This is important because it should ALWAYS BE SET
	}
	b.Down()
}

func (b *basicSelection) OnUp(up func()) {
	b.Up = up
}

func (b *basicSelection) OnDown(down func()) {
	b.Down = down
}
