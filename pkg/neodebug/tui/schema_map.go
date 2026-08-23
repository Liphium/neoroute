package tui

import (
	"fmt"
	"strings"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/Liphium/neoroute/neoschema"
)

var _ SchemaNode = &MapNode{}

// MapEntry stores one map key together with its value node.
type MapEntry struct {
	Key, Value SchemaNode
}

// MapNode renders and edits map entries.
//
// Each entry owns two nodes. Selection tracks active entry plus whether key or value currently owns focus.
type MapNode struct {
	basicNode
	keyHandled         bool
	manageSelection    bool // When the comment (manage selection) in the beginning of the map is shown
	keyType, valueType neoschema.PackedType
	registry           map[string]neoschema.PackedType
	items              []MapEntry
	selected           int
	valueSelected      bool

	// Keys
	up     key.Binding
	down   key.Binding
	left   key.Binding
	right  key.Binding
	add    key.Binding
	remove key.Binding
	clear  key.Binding
}

// KeyHandled implements [SchemaNode].
func (m *MapNode) KeyHandled() bool { return m.keyHandled }

// Init implements [SchemaNode].
func (m *MapNode) Init() {
	m.up = key.NewBinding(standardUpKey, key.WithHelp(SymbolArrowUp, "up"))
	m.down = key.NewBinding(standardDownKey, key.WithHelp(SymbolArrowDown, "down"))
	m.left = key.NewBinding(standardLeftKey, key.WithHelp(SymbolArrowLeft, "key"))
	m.right = key.NewBinding(standardRightKey, key.WithHelp(SymbolArrowRight, "value"))
	m.add = key.NewBinding(key.WithKeys("a"), key.WithHelp("a", "add"))
	m.remove = key.NewBinding(key.WithKeys("r"), key.WithHelp("r", "remove"))
	m.clear = key.NewBinding(key.WithKeys("c"), key.WithHelp("c", "clear map"))

	for i := range m.items {
		m.items[i].Key.Init()
		m.items[i].Value.Init()
	}
	m.redoBindings()
}

// redoBindings configures vertical movement across map entries.
func (m *MapNode) redoBindings() {
	for i := range m.items {
		i := i
		m.items[i].Key.SetSuffix(secondaryTextStyle.Render(":"))
		m.items[i].Value.SetSuffix(secondaryTextStyle.Render(","))

		// We go up to the manageSelection when on the first item, otherwise to the entry above
		var up func()
		if i == 0 {
			up = func() { m.selected = -1; m.manageSelection = true }
		} else {
			up = func() { m.selectEntry(i-1, m.valueSelected) }
		}

		// We go down when we're on the last item, otherwise to the entry below
		var down func()
		if i == len(m.items)-1 {
			down = func() { m.selected = -1; m.GoDown() }
		} else {
			down = func() { m.selectEntry(i+1, m.valueSelected) }
		}

		m.items[i].Key.OnUp(up)
		m.items[i].Value.OnUp(up)
		m.items[i].Key.OnDown(down)
		m.items[i].Value.OnDown(down)
	}
}

// selectEntry clears current selection, then focuses key or value in entry i.
func (m *MapNode) selectEntry(i int, value bool) {
	for j := range m.items {
		m.items[j].Key.Unselect()
		m.items[j].Value.Unselect()
	}
	m.manageSelection = false
	m.selected, m.valueSelected = i, value
	if value {
		m.items[i].Value.SelectFromTop()
	} else {
		m.items[i].Key.SelectFromTop()
	}
}

// Children implements [SchemaNode].
func (m *MapNode) Children() []keyProvider {
	if m.selected < 0 || m.selected >= len(m.items) {
		return []keyProvider{}
	}
	if m.valueSelected {
		return []keyProvider{m.items[m.selected].Value}
	}
	return []keyProvider{m.items[m.selected].Key}
}

// FooterKeys implements [SchemaNode].
func (m *MapNode) FooterKeys() []key.Binding {
	if m.manageSelection {
		return []key.Binding{m.add, m.clear, m.up, m.down}
	}
	return []key.Binding{m.add, m.remove, m.left, m.right}
}

// FullKeyHelp implements [SchemaNode].
func (m *MapNode) FullKeyHelp() FullKeyHelp {
	return FullKeyHelp{
		Title: "Map editing",
		Keys: [][]key.Binding{
			m.FooterKeys(),
		},
	}
}

// Request implements [SchemaNode].
func (m *MapNode) Request() any {
	result := map[string]any{}
	for _, item := range m.items {
		// Message pack only accepts strings for now, so we need to customize this like this
		result[fmt.Sprintf("%v", item.Key.Request())] = item.Value.Request()
	}
	return result
}

// Height implements [SchemaNode].
func (m *MapNode) Height() int {
	if len(m.items) == 0 {
		return 1
	}

	sum := 0
	for _, item := range m.items {
		sum += (item.Key.Height() + item.Value.Height()) - 1
	}
	return sum + 2
}

// Selected implements [SchemaNode].
func (m *MapNode) Selected() int {
	if m.manageSelection {
		return 1
	} else if len(m.items) == 0 {
		return 0
	}

	sum := 1 /* first bracket */
	for _, item := range m.items {
		if selected := item.Key.Selected(); selected != 0 {
			return sum + selected
		}
		sum += item.Key.Height()

		if selected := item.Value.Selected(); selected != 0 {
			return sum + selected - 1
		}
		sum += item.Value.Height() - 1
	}
	return 0
}

// SelectFromTop implements [SchemaNode].
func (m *MapNode) SelectFromTop() {
	m.manageSelection = true
}

// SelectFromBottom implements [SchemaNode].
func (m *MapNode) SelectFromBottom() {
	if len(m.items) == 0 {
		m.manageSelection = true
		return
	}
	m.selectEntry(len(m.items)-1, m.valueSelected)
}

// Unselect implements [SchemaNode].
func (m *MapNode) Unselect() {
	m.manageSelection = false
	for i := range m.items {
		m.items[i].Key.Unselect()
		m.items[i].Value.Unselect()
	}
	m.selected = -1
}

// Update implements [SchemaNode].
func (m *MapNode) Update(msg tea.Msg) tea.Cmd {
	m.keyHandled = false

	// Handle keys and navigation for the manage selection
	if m.manageSelection {
		if msg, ok := msg.(tea.KeyPressMsg); ok {
			switch {
			case matchesFocusedNavigation(msg, false, m.up):
				m.keyHandled = true
				m.manageSelection = false
				m.GoUp()

			case matchesFocusedNavigation(msg, false, m.down):
				m.keyHandled = true
				m.manageSelection = false

				// Go down or select first entry
				if len(m.items) == 0 {
					m.GoDown()
					return nil
				}
				m.selectEntry(0, m.valueSelected)

			case key.Matches(msg, m.add):
				m.keyHandled = true
				m.addEntry(-1)

			case key.Matches(msg, m.clear):
				m.keyHandled = true
				m.items = []MapEntry{}
				m.redoBindings()
			}
		}
		return nil
	}

	if m.selected >= 0 {
		var child SchemaNode = m.items[m.selected].Key
		if m.valueSelected {
			child = m.items[m.selected].Value
		}
		cmd := child.Update(msg)
		if child.KeyHandled() {
			m.keyHandled = true
			return cmd
		}

		if msg, ok := msg.(tea.KeyPressMsg); ok {
			switch {
			// Select the key when going right and the key is currently selected
			case !m.valueSelected && key.Matches(msg, m.right):
				m.keyHandled = true
				m.selectEntry(m.selected, true)
				return nil

			// Select the key when going left and the value is currently selected
			case m.valueSelected && key.Matches(msg, m.left):
				m.keyHandled = true
				m.selectEntry(m.selected, false)
				return nil

			// Add a new item at the index
			case key.Matches(msg, m.add):
				m.keyHandled = true
				m.addEntry(m.selected)
				return nil

			// Remove the item
			case key.Matches(msg, m.remove):
				m.items = append(m.items[:m.selected], m.items[m.selected+1:]...)
				m.redoBindings()

				// Select the next item if possible, otherwise the previous one, otherwise manage selection
				if m.selected >= len(m.items) {
					m.selected = len(m.items) - 1
				}
				if m.selected < 0 {
					m.manageSelection = true
				} else {
					m.selectEntry(m.selected, m.valueSelected)
				}
				return nil
			}
		}
		return cmd
	}
	return nil
}

// addEntry creates a new entry at the given index, initializes it, and selects it.
func (m *MapNode) addEntry(index int) {
	item := MapEntry{
		Key:   createNode(m.keyType, m.registry),
		Value: createNode(m.valueType, m.registry),
	}
	m.items = append(m.items[:index+1], append([]MapEntry{item}, m.items[index+1:]...)...)
	m.items[index+1].Key.Init()
	m.items[index+1].Value.Init()
	m.redoBindings()

	// Select the thing
	m.selectEntry(index+1, false)
}

// View implements [SchemaNode].
func (m *MapNode) View() (*tea.Cursor, string) {

	// Render all of the children with the key prefixed + some padding (also add the padding to the cursor and stuff)
	var b strings.Builder
	var cursor *tea.Cursor

	// Compute bracket style
	style := secondaryTextStyle
	if m.manageSelection || m.selected >= 0 {
		style = textStyle
	}

	b.WriteString(style.Render("{"))

	// When there are no items, make sure to just render the empty object
	if len(m.items) == 0 {
		b.WriteString(style.Render("}")) // We return later due to the manageSelection
	}

	// Render manageSelection when on
	if m.manageSelection {
		b.WriteString(highlightStyle.Render(" /* a to add, c to clear */"))
	}

	// Close here when there are no items (so the manageSelection is still visible when active)
	if len(m.items) == 0 {
		return nil, b.String() + m.basicNode.Suffix
	}
	b.WriteString("\n")

	fieldPadding := 0
	for _, item := range m.items {
		c1, k := item.Key.View()
		c2, v := item.Value.View()
		if item.Key.Selected() != 0 {
			cursor = c1
		}
		if item.Value.Selected() != 0 {
			lines := strings.Split(k, "\n")
			fieldPadding = lipgloss.Width(lines[len(lines)-1]) + 1 /* space for separation */
			cursor = c2
		}

		// Write only this, the prefixes are already set in redoBindings above
		b.WriteString(structChildStyle.Render(k + " " + v))
		b.WriteByte('\n')
	}
	b.WriteString(style.Render("}"))

	if cursor != nil {
		cursor.X += structPadding + fieldPadding
	}
	b.WriteString(m.basicNode.Suffix)
	return cursor, b.String()
}
