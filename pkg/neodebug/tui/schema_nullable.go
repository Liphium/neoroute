package tui

import (
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
)

var _ SchemaNode = &ValueNode[any]{}

type NullableNode struct {
	basicNode
	focused bool
	null    bool
	other   SchemaNode

	// Keys
	up         key.Binding
	down       key.Binding
	toggleNull key.Binding
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
	return []key.Binding{n.toggleNull, n.up, n.down}
}

// FullKeyHelp implements SchemaNode.
func (n *NullableNode) FullKeyHelp() FullKeyHelp {
	return FullKeyHelp{
		Title: "Nullable functionality",
		Keys: [][]key.Binding{
			[]key.Binding{n.up, n.down},
			[]key.Binding{n.toggleNull},
		},
	}
}

// Init implements [SchemaNode].
func (n *NullableNode) Init() {
	n.toggleNull = key.NewBinding(key.WithKeys("ctrl+d"), key.WithHelp("ctrl+d", "toggle nil"))
	n.up = key.NewBinding(standardUpKey, key.WithHelp(SymbolArrowUp, "up"))
	n.down = key.NewBinding(standardDownKey, key.WithHelp(SymbolArrowDown, "down"))

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

// SelectFromBottom implements [SchemaNode].
func (n *NullableNode) SelectFromBottom() {
	n.focused = true
	if !n.null {
		n.other.SelectFromBottom()
	}
}

// SelectFromTop implements [SchemaNode].
func (n *NullableNode) SelectFromTop() {
	n.focused = true
	if !n.null {
		n.other.SelectFromTop()
	}
}

// Selected implements [SchemaNode].
func (n *NullableNode) Selected() int {
	if n.focused {
		if !n.null {
			return n.other.Selected()
		}
		return 1
	}
	return 0
}

// Update implements [SchemaNode].
func (n *NullableNode) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, n.toggleNull):
			n.null = !n.null
			if !n.null {
				n.other.SelectFromTop()
			}
			return nil

		case n.null && key.Matches(msg, n.up):
			n.focused = false
			n.GoUp()
			return nil

		case n.null && key.Matches(msg, n.down):
			n.focused = false
			n.GoDown()
			return nil
		}
	}

	if n.focused && !n.null {
		return n.other.Update(msg)
	}
	return nil
}

// View implements [SchemaNode].
func (n *NullableNode) View() (*tea.Cursor, string) {
	if n.null {
		style := secondaryTextStyle
		if n.focused {
			style = textStyle
		}
		return nil, style.Render("nil") + n.basicNode.Suffix
	}

	// If not nil, render the other node
	c, view := n.other.View()
	return c, view + n.basicNode.Suffix
}
