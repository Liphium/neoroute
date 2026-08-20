package tui

import (
	"strings"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

var (
	structPadding    = 4
	structChildStyle = lipgloss.NewStyle().Padding(0, 0, 0, structPadding)
)

// Wishlist:
// - Expanding and collapsing for structs

var _ SchemaNode = &StructNode{}

type StructNode struct {
	basicNode
	name     string
	children []StructField
}

// Request implements SchemaNode.
func (s *StructNode) Request() any {
	req := map[string]any{}
	for _, f := range s.children {
		req[f.Name] = f.Node.Request()
	}

	return req
}

// Children implements SchemaNode.
func (s *StructNode) Children() []keyProvider {
	var sel SchemaNode
	for _, field := range s.children {
		if field.Node.Selected() != 0 {
			sel = field.Node
			break
		}
	}
	if sel == nil {
		return []keyProvider{}
	}
	return []keyProvider{sel}
}

// FooterKeys implements SchemaNode.
func (s *StructNode) FooterKeys() []key.Binding {
	return []key.Binding{}
}

// FullKeyHelp implements SchemaNode.
func (s *StructNode) FullKeyHelp() FullKeyHelp {
	return FullKeyHelp{}
}

type StructField struct {
	Name string
	Node SchemaNode
}

// Init implements [SchemaNode].
func (s *StructNode) Init() {
	for i, child := range s.children {
		child.Node.Init()
		child.Node.SetSuffix(secondaryTextStyle.Render(","))

		// Configure up down selection for the struct (we want to walk fields this way)
		configureUpDownSelection(i, [3]SchemaNode{s.children[max(i-1, 0)].Node, child.Node, s.children[min(i+1, len(s.children)-1)].Node}, s.basicNode, len(s.children))
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
func (s *StructNode) Update(msg tea.Msg) tea.Cmd {
	// Just update the child with Selected() != 0
	for _, child := range s.children {
		if child.Node.Selected() != 0 {
			cmd := child.Node.Update(msg)
			return cmd
		}
	}

	return nil
}

// View implements [SchemaNode].
func (s *StructNode) View() (*tea.Cursor, string) {
	// Render all of the children with the field name prefixed + some padding (also add the padding to the cursor and stuff)
	var cursor *tea.Cursor
	var b strings.Builder

	// Write the name of the struct we're editing
	b.WriteString(textStyle.Render(s.name+" {") + "\n")

	fieldPadding := 0
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
		c, v := child.Node.View()
		if sel != 0 {
			cursor = c
			if sel == 1 {
				fieldPadding = len(child.Name + ": ")
			}
		}
		fb.WriteString(v)

		// The field builder is rendered here to make sure the padding is applied to everything
		b.WriteString(structChildStyle.Render(fb.String()) + "\n")
	}

	// Write the closing bracket for the struct
	b.WriteString(textStyle.Render("}"))

	if cursor != nil {
		cursor.X += structPadding + fieldPadding
	}
	return cursor, b.String() + s.basicNode.Suffix
}
