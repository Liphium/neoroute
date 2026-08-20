package tui

import (
	"strings"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"github.com/Liphium/neoroute/neoschema"
)

// Wishlist:
// - Expanding and collapsing for slices

var _ SchemaNode = &SliceNode{}

type SliceNode struct {
	basicNode
	items    []SchemaNode
	element  neoschema.PackedType
	registry map[string]neoschema.PackedType

	// Keys
	add key.Binding
}

// Init implements [SchemaNode].
func (s *SliceNode) Init() {
	for _, field := range s.items {
		field.Init()
	}
	s.redoBindings()

	s.add = key.NewBinding(key.WithKeys("ctrl+a"), key.WithHelp("ctrl+a", "add"))
}

func (s *SliceNode) redoBindings() {
	for i, item := range s.items {
		item.SetSuffix(secondaryTextStyle.Render(","))

		configureUpDownSelection(i, [3]SchemaNode{s.items[max(i-1, 0)], item, s.items[min(i+1, len(s.items)-1)]}, s.basicNode, len(s.items))
	}
}

// Children implements [SchemaNode].
func (s *SliceNode) Children() []keyProvider {
	var sel SchemaNode
	for _, item := range s.items {
		if item.Selected() != 0 {
			sel = item
			break
		}
	}
	if sel == nil {
		return []keyProvider{}
	}
	return []keyProvider{sel}
}

// FooterKeys implements [SchemaNode].
func (s *SliceNode) FooterKeys() []key.Binding {
	return []key.Binding{s.add}
}

// FullKeyHelp implements [SchemaNode].
func (s *SliceNode) FullKeyHelp() FullKeyHelp {
	return FullKeyHelp{
		Title: "Slice editing",
		Keys: [][]key.Binding{
			s.FooterKeys(),
		},
	}
}

// Request implements [SchemaNode].
func (s *SliceNode) Request() any {
	values := make([]any, len(s.items))
	for i, item := range s.items {
		values[i] = item.Request()
	}

	return values
}

// Height implements [SchemaNode].
func (s *SliceNode) Height() int {
	// Our height is the sum of the one of our children + a little bit
	sum := 0
	for _, item := range s.items {
		sum += item.Height()
	}

	return sum + 2 /* Name of the struct and closing brace */
}

// Selected implements [SchemaNode].
func (s *SliceNode) Selected() int {
	// For selection we need to just find the child that has Selected() != 0 and add all of the previous heights till then
	sum := 1
	for _, item := range s.items {
		sel := item.Selected()
		if sel != 0 {
			return sum + sel
		}

		sum += item.Height()
	}

	return 0
}

// SelectFromTop implements [SchemaNode].
func (s *SliceNode) SelectFromTop() {
	// We don't actually have any selection state, just our children do, select the first children from the top
	s.items[0].SelectFromTop()
}

// SelectFromBottom implements [SchemaNode].
func (s *SliceNode) SelectFromBottom() {
	// We don't actually have any selection state, just our children do, select the first children from the bottom
	s.items[len(s.items)-1].SelectFromBottom()
}

// Update implements [SchemaNode].
func (s *SliceNode) Update(msg tea.Msg) tea.Cmd {

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, s.add):
			s.items = append(s.items, createNode(s.element, s.registry))
			s.items[len(s.items)-1].Init()
			s.redoBindings()
			return nil
		}
	}

	// Just update the child with Selected() != 0
	for _, item := range s.items {
		if item.Selected() != 0 {
			return item.Update(msg)
		}
	}

	return nil
}

// View implements [SchemaNode].
func (s *SliceNode) View() (*tea.Cursor, string) {
	// Render all of the children with the field name prefixed + some padding (also add the padding to the cursor and stuff)
	var cursor *tea.Cursor
	var b strings.Builder

	// Write the name of the struct we're editing
	b.WriteString(textStyle.Render("[") + "\n")

	for _, item := range s.items {
		var fb strings.Builder
		sel := item.Selected()

		// Write the actual view of the thing
		c, v := item.View()
		if sel != 0 {
			cursor = c
		}
		fb.WriteString(v)

		// The field builder is rendered here to make sure the padding is applied to everything
		b.WriteString(structChildStyle.Render(fb.String()) + "\n")
	}

	// Write the closing bracket for the struct
	b.WriteString(textStyle.Render("]"))

	if cursor != nil {
		cursor.X += structPadding
	}
	return cursor, b.String() + s.basicNode.Suffix
}
