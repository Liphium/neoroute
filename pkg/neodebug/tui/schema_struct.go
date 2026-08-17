package tui

import (
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

var (
	structPadding    = 2
	structChildStyle = lipgloss.NewStyle().Padding(0, 0, 0, structPadding)
)

// Wishlist:
// - Expanding and collapsing for structs

var _ SchemaNode = &StructNode{}

type StructNode struct {
	*basicSelection
	name     string
	children []StructField
}

type StructField struct {
	Name string
	Node SchemaNode
}

// Init implements [SchemaNode].
func (s *StructNode) Init() {
	for i, child := range s.children {

		// When up, go to above child or above node
		if i == 0 {
			child.Node.OnUp(s.GoUp)
		} else {
			child.Node.OnUp(s.children[i-1].Node.SelectFromBottom)
		}

		// When down, go to below child or below node
		if i == len(s.children)-1 {
			child.Node.OnDown(s.GoDown)
		} else {
			child.Node.OnDown(s.children[i+1].Node.SelectFromTop)
		}
	}
}

// SelectFromTop implements [SchemaNode].
func (s *StructNode) SelectFromTop() {
	// We don't actually have any selection state, just our children do, select the first children from the top
	s.children[0].Node.SelectFromTop()
}

// SelectFromBottom implements [SchemaNode].
func (s *StructNode) SelectFromBottom() {
	// We don't actually have any selection state, just our children do, select the first children from the bottom
	s.children[len(s.children)-1].Node.SelectFromBottom()
}

// Height implements [SchemaNode].
func (s *StructNode) Height() int {
	// Our height is the sum of the one of our children + a little bit
	sum := 0
	for _, child := range s.children {
		sum += child.Node.Height()
	}

	return sum + 2 /* Name of the struct and closing brace */
}

// Selected implements [SchemaNode].
func (s *StructNode) Selected() int {
	// For selection we need to just find the child that has Selected() != 0 and add all of the previous heights till then
	sum := 1
	for _, child := range s.children {
		sel := child.Node.Selected()
		if sel != 0 {
			return sum + sel
		}

		sum += child.Node.Height()
	}

	return 0
}

// Update implements [SchemaNode].
func (s *StructNode) Update(msg tea.Msg) (SchemaNode, tea.Cmd) {
	// Just update the child with Selected() != 0
	for _, child := range s.children {
		if child.Node.Selected() != 0 {
			newChild, cmd := child.Node.Update(msg)
			child.Node = newChild
			return s, cmd
		}
	}

	return s, nil
}

// View implements [SchemaNode].
func (s *StructNode) View() (*tea.Cursor, string) {
	// Render all of the children with the field name prefixed + some padding (also add the padding to the cursor and stuff)
	var cursor *tea.Cursor
	var b strings.Builder

	// Write the name of the struct we're editing
	b.WriteString(textStyle.Render(s.name+" {") + "\n")

	for _, child := range s.children {
		var fb strings.Builder
		sel := child.Node.Selected()

		// Create the line for the struct field
		fieldLine := unselectedStyle.Render(child.Name+":") + " "
		if sel != 0 {
			fieldLine = selectedStyle.Render(child.Name+":") + " "
		}
		fb.WriteString(fieldLine)

		// Write the actual view of the thing
		var v string
		if sel != 0 {
			cursor, v = child.Node.View()
		}
		fb.WriteString(v)

		// The field builder is rendered here to make sure the padding is applied to everything
		b.WriteString(structChildStyle.Render(fb.String()) + secondaryTextStyle.Render(",") + "\n")
	}

	// Write the closing bracket for the struct
	b.WriteString(textStyle.Render("}"))

	cursor.X += structPadding
	return cursor, b.String()
}
