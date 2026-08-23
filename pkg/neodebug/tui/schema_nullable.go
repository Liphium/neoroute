package tui

import (
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
)

var _ SchemaNode = &NullableNode{}

type nullableState int

type NullableNode struct {
	basicNode
	focused    bool
	keyHandled bool
	null       bool
	other      SchemaNode

	// Keys
	up         key.Binding
	down       key.Binding
	toggleNull key.Binding
}

// KeyHandled implements [SchemaNode].
func (n *NullableNode) KeyHandled() bool {
	return n.keyHandled
}

// OnUp needs to be forwarded to other as well
func (n *NullableNode) OnUp(up func()) {
	n.basicNode.OnUp(up)
	n.other.OnUp(func() {
		n.focused = false
		up()
	})
}

// OnDown needs to be forwarded to other as well
func (n *NullableNode) OnDown(down func()) {
	n.basicNode.OnDown(down)
	n.other.OnDown(func() {
		n.focused = false
		down()
	})
}

// Request implements SchemaNode.
func (n *NullableNode) Request() any {
	if n.null {
		return nil
	}
	return n.other.Request()
}

// Children implements SchemaNode.
func (n *NullableNode) Children() []keyProvider {
	return []keyProvider{}
}

// FooterKeys implements SchemaNode.
func (n *NullableNode) FooterKeys() []key.Binding {
	if n.null {
		return []key.Binding{n.toggleNull, n.up, n.down}
	}
	return []key.Binding{n.toggleNull}
}

// FullKeyHelp implements SchemaNode.
func (n *NullableNode) FullKeyHelp() FullKeyHelp {
	return FullKeyHelp{
		Title: "Nullable editing",
		Keys: [][]key.Binding{
			[]key.Binding{n.toggleNull},
		},
	}
}

// Init implements [SchemaNode].
func (n *NullableNode) Init() {
	n.up = key.NewBinding(standardUpKey, key.WithHelp(SymbolArrowUp, "up"))
	n.down = key.NewBinding(standardDownKey, key.WithHelp(SymbolArrowDown, "down"))
	n.toggleNull = key.NewBinding(key.WithKeys("n"), key.WithHelp("n", "toggle nil"))

	n.other.Init()
}

// Height implements [SchemaNode].
func (n *NullableNode) Height() int {
	if n.null {
		return 1
	} else {
		return n.other.Height()
	}
}

// SelectFromTop implements [SchemaNode].
func (n *NullableNode) SelectFromTop() {
	if n.null {
		n.focused = true
	} else {
		n.other.SelectFromTop()
	}
}

// SelectFromBottom implements [SchemaNode].
func (n *NullableNode) SelectFromBottom() {
	if n.null {
		n.focused = true
	} else {
		n.other.SelectFromBottom()
	}
}

// Selected implements [SchemaNode].
func (n *NullableNode) Selected() int {
	if n.null {
		if n.focused {
			return 1
		}
		return 0
	}
	return n.other.Selected()
}

// Unselect implements [SchemaNode].
func (n *NullableNode) Unselect() {
	if n.null {
		n.focused = false
	}
	n.other.Unselect()
}

// Update implements [SchemaNode].
func (n *NullableNode) Update(msg tea.Msg) tea.Cmd {
	n.keyHandled = false

	// If null, handle all of the navigation keys
	if n.null {
		switch msg := msg.(type) {
		case tea.KeyPressMsg:
			switch {
			case key.Matches(msg, n.up):
				n.keyHandled = true
				n.focused = false
				n.GoUp()
				return nil

			case key.Matches(msg, n.down):
				n.keyHandled = true
				n.focused = false
				n.GoDown()
				return nil
			}
		}
	}

	cmd := n.other.Update(msg)
	if n.other.KeyHandled() {
		n.keyHandled = true
		return cmd
	}

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, n.toggleNull):
			n.keyHandled = true
			n.null = !n.null
			if !n.null {
				n.focused = false
				n.other.SelectFromTop()
			} else {
				n.focused = true
			}
			return nil
		}
	}
	return cmd
}

// View implements [SchemaNode].
func (n *NullableNode) View() (*tea.Cursor, string) {
	if n.null {
		style := secondaryTextStyle
		if n.focused {
			style = highlightStyle.Bold(true)
		}
		return nil, style.Render("nil") + n.basicNode.Suffix
	}

	// If not nil, render the other node
	c, view := n.other.View()
	return c, view + n.basicNode.Suffix
}
