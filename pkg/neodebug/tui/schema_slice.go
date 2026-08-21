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
	gap      int // -1 = no gap selection, the 0 gap is before the item with index 0, last gap is len(s.items)
	items    []SchemaNode
	element  neoschema.PackedType
	registry map[string]neoschema.PackedType

	// Keys
	up     key.Binding
	down   key.Binding
	add    key.Binding
	remove key.Binding
}

// Init implements [SchemaNode].
func (s *SliceNode) Init() {
	for _, field := range s.items {
		field.Init()
	}
	s.redoBindings()

	// Define the key bindings
	s.up = key.NewBinding(standardUpKey, key.WithHelp("↑", "up"))
	s.down = key.NewBinding(standardDownKey, key.WithHelp("↓", "down"))
	s.add = key.NewBinding(key.WithKeys("+"), key.WithHelp("+", "add"))
	s.remove = key.NewBinding(key.WithKeys("-"), key.WithHelp("-", "remove"))
}

func (s *SliceNode) redoBindings() {
	for i, item := range s.items {
		item.SetSuffix(secondaryTextStyle.Render(","))

		// The gap above an item is always the correct one to select
		item.OnUp(func() {
			s.selectGap(i)
		})

		// When going down, we should go to the gap below
		item.OnDown(func() {
			s.selectGap(i + 1)
		})
	}
}

func (s *SliceNode) selectGap(gap int) {
	s.gap = gap
}

// Children implements [SchemaNode].
func (s *SliceNode) Children() []keyProvider {

	// When no gap is selected, we don't need to return anything
	if s.gap != -1 {
		return []keyProvider{}
	}

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
	base := []key.Binding{s.up, s.down}
	if len(s.items) > 0 {
		base = append([]key.Binding{s.remove}, base...)
	}
	return append([]key.Binding{s.add}, base...)
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

	// Add one if we're in a gap selection
	if s.gap != -1 {
		sum += 1
	}

	return sum + 2 /* Name of the struct and closing brace */
}

// Selected implements [SchemaNode].
func (s *SliceNode) Selected() int {

	// For selection we need to just find the child that has Selected() != 0 and add all of the previous heights till then
	sum := 1 /* The open bracket that already exists */
	for i, item := range s.items {
		if s.gap != -1 && s.gap == i {
			return sum + 1
		}

		sel := item.Selected()
		if sel != 0 {
			return sum + sel
		}

		sum += item.Height()
	}

	// If we haven't found the gap so far, we selected the last one
	if s.gap != -1 {
		return sum + 1
	}

	return 0
}

// SelectFromTop implements [SchemaNode].
func (s *SliceNode) SelectFromTop() {
	s.gap = 0
}

// SelectFromBottom implements [SchemaNode].
func (s *SliceNode) SelectFromBottom() {
	s.gap = len(s.items)
}

// Update implements [SchemaNode].
func (s *SliceNode) Update(msg tea.Msg) tea.Cmd {

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch {
		case s.gap >= 0 && key.Matches(msg, s.add):

			// Add an item exactly at the gap index and select it
			item := createNode(s.element, s.registry)
			s.items = append(s.items[:s.gap], append([]SchemaNode{item}, s.items[s.gap:]...)...)
			s.items[s.gap].Init()
			s.redoBindings()

			// Select the new item
			s.items[s.gap].SelectFromTop()
			s.gap = -1

			return nil
		case s.gap >= 1 && key.Matches(msg, s.remove):

			// Remove item before the gap
			i := s.gap - 1
			s.items = append(s.items[:i], s.items[i+1:]...)
			s.redoBindings()

			// Make sure to adjust the gap now
			s.gap -= 1

			return nil

		// When we're focused we don't have any elements, meaning we're just straight going up and down essentially
		case key.Matches(msg, s.up):
			// First gap means we need to select the element above this one
			if s.gap == 0 {
				s.gap = -1
				s.GoUp()
				return nil
			}

			if s.gap != -1 {
				// If we're in gap selection and we know it's not the first one (check above), we need to select the item before the gap
				s.items[s.gap-1].SelectFromTop()
				s.gap = -1
			} else {
				// This is impossible because items should forward their things from up / down events
			}
		case key.Matches(msg, s.down):
			// Last gap means we need to select the element below
			if s.gap == len(s.items) {
				s.gap = -1
				s.GoDown()
				return nil
			}

			if s.gap != -1 {
				// If we're in the gap selection, we need to select the item below
				s.items[s.gap].SelectFromBottom()
				s.gap = -1
			} else {
				// This is impossible because items should forward their things from up / down events
			}
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

	for i, item := range s.items {
		var fb strings.Builder
		sel := item.Selected()

		// Render gap when we're in it (gap with index is always before the item with the index)
		if s.gap == i {
			fb.WriteString(s.renderGap() + "\n")
		}

		// Write the actual view of the thing
		c, v := item.View()
		if sel != 0 {
			cursor = c
		}
		fb.WriteString(v)

		// The field builder is rendered here to make sure the padding is applied to everything
		b.WriteString(structChildStyle.Render(fb.String()) + "\n")
	}

	// Render last gap in case we're in that
	if s.gap == len(s.items) {
		b.WriteString(structChildStyle.Render(s.renderGap()) + "\n")
	}

	// Write the closing bracket for the struct
	b.WriteString(textStyle.Render("]"))

	if cursor != nil {
		cursor.X += structPadding
	}
	return cursor, b.String() + s.basicNode.Suffix
}

func (s *SliceNode) renderGap() string {
	if s.gap == 0 {
		return textStyle.Bold(true).Render("Gap selection") + " " + secondaryTextStyle.Render("(+ to add item)")
	}
	return textStyle.Bold(true).Render("Gap selection") + " " + secondaryTextStyle.Render("(+ to add item, - to remove item above)")
}
