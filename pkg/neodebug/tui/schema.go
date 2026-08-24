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

	// Clear the selection
	Unselect()

	// Should take a callback that is called when the upper boundary of the selection is reached and the selection should go up.
	//
	// This should also unselect this node at the same time.
	OnUp(func())

	// Should take a callback that is called when the lower boundary of the selection is reached and the seelction should go down.
	//
	// This should also unselect this node at the same time.
	OnDown(func())

	// If the last key event was handled by this node or any of its children.
	KeyHandled() bool

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

	// Should construct the request using go types.
	Request() any
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

// configureUpDownSelection adds up and down functions to a child node based on its parents' fields.
//
// snapshot is the current children in the following order: previous (i-1), current (i), next (i+1). If that would go out of bounds, just put anything, it won't be touched.
//
// This is meant for nodes that use up/down to navigate, well, up and down (like structs, slices, maps, etc.).
func configureUpDownSelection(i int, snapshot [3]SchemaNode, parent basicNode, size int) {
	// When up, go to above child or above node
	if i == 0 {
		// The method can not be put in here directly (as s.GoUp instead of wrapping it in a func), otherwise the thing is copied meaning if OnUp on the child is called after, the method is not called correctly.
		snapshot[1].OnUp(func() {
			parent.GoUp()
		})
	} else {
		snapshot[1].OnUp(snapshot[0].SelectFromBottom)
	}

	// When down, go to below child or below node
	if i == size-1 {
		// The method can not be put in here directly (as s.GoDown instead of wrapping it in a func), otherwise the thing is copied meaning if OnUp on the child is called after, the method is not called correctly.
		snapshot[1].OnDown(func() {
			parent.GoDown()
		})
	} else {
		snapshot[1].OnDown(snapshot[2].SelectFromTop)
	}
}
