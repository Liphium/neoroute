package tui

import (
	"fmt"
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
	keyHandled      bool
	manageSelection bool // If the comment (manage selection) in the beginning of the map is shown
	items           []SchemaNode
	element         neoschema.PackedType
	registry        map[string]neoschema.PackedType

	// Keys
	up     key.Binding
	down   key.Binding
	add    key.Binding
	remove key.Binding
	clear  key.Binding
}

// KeyHandled implements [SchemaNode].
func (s *SliceNode) KeyHandled() bool {
	return s.keyHandled
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
	s.add = key.NewBinding(key.WithKeys("a"), key.WithHelp("a", "add"))
	s.remove = key.NewBinding(key.WithKeys("r"), key.WithHelp("r", "remove"))
	s.clear = key.NewBinding(key.WithKeys("c"), key.WithHelp("c", "clear slice"))
}

func (s *SliceNode) redoBindings() {
	for i, item := range s.items {
		item.SetSuffix(secondaryTextStyle.Render(","))

		// For the first item we want to select the manage selection
		if i == 0 {
			item.OnUp(func() {
				s.manageSelection = true
			})
		} else {
			item.OnUp(func() {
				s.items[i-1].SelectFromBottom()
			})
		}

		if i == len(s.items)-1 {
			// When on the last item, we should just go down to the next element straight away
			item.OnDown(func() {
				s.GoDown()
			})
		} else {
			// When going down, we should go to the next item
			item.OnDown(func() {
				s.items[i+1].SelectFromTop()
			})
		}
	}
}

// Children implements [SchemaNode].
func (s *SliceNode) Children() []keyProvider {

	// When nothing is selected, we don't need to return anything
	if s.manageSelection {
		return []keyProvider{}
	}

	// If an item is selected, we return that (currently in edit mode there)
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
	if s.manageSelection {
		return []key.Binding{s.add, s.clear, s.up, s.down}
	}
	return []key.Binding{s.add, s.remove}
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

	// Collapse list when no items are there
	if len(s.items) == 0 {
		return 1
	}

	// Our height is the sum of the one of our children + a little bit
	sum := 0
	for _, item := range s.items {
		sum += item.Height()
	}

	return sum + 2 /* Opening and closing brace of the array */
}

// Selected implements [SchemaNode].
func (s *SliceNode) Selected() int {

	// When we're in manageSelection, it's always 1
	if s.manageSelection {
		return 1
	}

	// For selection we need to just find the child that has Selected() != 0 and add all of the previous heights till then
	sum := 1 /* The open bracket that already exists */
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
	s.manageSelection = true
}

// SelectFromBottom implements [SchemaNode].
func (s *SliceNode) SelectFromBottom() {
	if len(s.items) == 0 {
		s.manageSelection = true
	} else {
		s.items[len(s.items)-1].SelectFromBottom()
	}
}

// Unselect implements [SchemaNode].
func (s *SliceNode) Unselect() {
	s.manageSelection = false
	for _, item := range s.items {
		item.Unselect()
	}
}

// selectedIndex finds the currently selected thingy or returns -1 if there isn't one.
func (s *SliceNode) selectedIndex() int {
	for i, item := range s.items {
		if item.Selected() != 0 {
			return i
		}
	}
	return -1
}

// Update implements [SchemaNode].
func (s *SliceNode) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	s.keyHandled = false

	// First try updating the child with Selected() != 0, if they handled keys, we can immediately return
	for _, item := range s.items {
		if item.Selected() != 0 {
			cmd = item.Update(msg)
			if item.KeyHandled() {
				s.keyHandled = true
				return cmd
			}
		}
	}

	selected := s.selectedIndex()
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch {
		// Handle up only for the manageSelection
		case s.manageSelection && matchesFocusedNavigation(msg, false, s.up):
			s.keyHandled = true
			s.manageSelection = false
			s.GoUp()
			return nil

		// Handle down only for the manageSelection
		case s.manageSelection && matchesFocusedNavigation(msg, false, s.down):
			s.keyHandled = true
			s.manageSelection = false

			if len(s.items) == 0 {
				s.GoDown()
				return nil
			}

			s.items[0].SelectFromTop()
			return nil

		// In the manage selection, allow clearing the slice
		case s.manageSelection && key.Matches(msg, s.clear):
			s.items = []SchemaNode{}
			s.redoBindings()
			return nil

		case (selected != -1 || s.manageSelection) && key.Matches(msg, s.add):
			s.keyHandled = true
			if s.manageSelection {
				selected = -1 // Nothing set yet
				s.manageSelection = false
			}

			// Add an item exactly at the selected index
			item := createNode(s.element, s.registry)
			s.items = append(s.items[:selected+1], append([]SchemaNode{item}, s.items[selected+1:]...)...)
			s.items[selected+1].Init()
			s.redoBindings()

			// Unselect the thing originally selected and select the thingy immediately
			if selected >= 0 {
				s.items[selected].Unselect()
			}
			s.items[selected+1].SelectFromTop()
			return nil

		case selected != -1 && key.Matches(msg, s.remove):
			s.keyHandled = true

			// Remove item at the selected index
			s.items = append(s.items[:selected], s.items[selected+1:]...)
			s.redoBindings()

			// Select the next thingy
			if len(s.items) == 0 {
				s.manageSelection = true
			} else if selected >= len(s.items)-1 {
				s.items[len(s.items)-1].SelectFromTop()
			} else {
				s.items[selected].SelectFromTop()
			}
			return nil
		}
	}

	return cmd
}

// View implements [SchemaNode].
func (s *SliceNode) View() (*tea.Cursor, string) {

	// Render all of the children with the field name prefixed + some padding (also add the padding to the cursor and stuff)
	var cursor *tea.Cursor
	var b strings.Builder

	// Compute the style for the brackets
	selected := s.selectedIndex()
	style := secondaryTextStyle
	if s.manageSelection || selected != -1 {
		style = textStyle
	}

	// Write the name of the struct we're editing
	b.WriteString(style.Render("["))

	// For collapsed list, instantly write the other bracket
	if len(s.items) == 0 {
		b.WriteString(style.Render("]"))
	}

	// Write the first selection line (in case there)
	if s.manageSelection {
		b.WriteString(" " + highlightStyle.Render("/* a to add, c to clear */"))
	}

	// Check for collapsed list
	if len(s.items) == 0 {
		return nil, b.String() + s.basicNode.Suffix
	}
	b.WriteString("\n")

	fieldPadding := 0
	for i, item := range s.items {
		var fb strings.Builder
		sel := item.Selected()

		// Write the actual view of the thing
		prefix := fmt.Sprintf("%d: ", i)
		prefixStyle := secondaryTextStyle
		c, v := item.View()
		if sel != 0 {
			cursor = c
			if sel == 1 {
				fieldPadding = len(prefix)
			}
			prefixStyle = textStyle
		}
		fb.WriteString(prefixStyle.Render(prefix) + v)

		// The field builder is rendered here to make sure the padding is applied to everything
		b.WriteString(structChildStyle.Render(fb.String()) + "\n")
	}

	// Write the closing bracket for the struct
	b.WriteString(style.Render("]"))

	if cursor != nil {
		cursor.X += structPadding + fieldPadding
	}
	return cursor, b.String() + s.basicNode.Suffix
}
