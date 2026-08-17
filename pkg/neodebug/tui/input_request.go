package tui

import "github.com/Liphium/neoroute/neoschema"

// Idea:
// - Different nodes for the different types we have in Neo, some can share the same thing
// - WantedHeight() for getting the height (should add height of all things together basically)
// - Selected state is handled in the nodes themselves
// - All keys are handled by the currently selected node
// - The View() function takes in indicies that represent the global position of the node (in Y) -> only render what's needed
// - The selected global position can be gotten by going through the node tree
//
// Unclear:
// - Do we already send when the user presses enter? Like is that something we handle here and then just send what we get?
// 	- Might be bad for UX, but at the same time, there is no reason to press enter on any objects

type inputRequestCreator struct {
	height int
	width  int
	schema neoschema.PackedType
}

func newInputRequestCreator(schema neoschema.PackedType) inputRequestCreator {
	return inputRequestCreator{
		schema: schema,
	}
}

func (m *inputRequestCreator) SetWidth(w int) {
	m.width = w
}

func (m *inputRequestCreator) SetHeight(h int) {
	m.height = h
}

func (m inputRequestCreator) WantedHeight() int {
	return 7 // TODO: Calculate the global height here (when smaller, also make this smaller)
}
