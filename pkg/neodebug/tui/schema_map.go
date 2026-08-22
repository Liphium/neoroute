package tui

import (
	"fmt"
	"strings"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"github.com/Liphium/neoroute/neoschema"
)

// AI-generated prototype to get a feeling for how this will work, re-work with actual better implementation.

var _ SchemaNode = &MapNode{}

type MapEntry struct{ key, value SchemaNode }

type MapNode struct {
	basicNode
	selected                                    int // -1 = child editing, 0 = map gap, n = entry n-1
	items                                       []MapEntry
	keyType, valueType                          neoschema.PackedType
	registry                                    map[string]neoschema.PackedType
	up, down, add, remove, edit, clear, keyEdit key.Binding
}

func (m *MapNode) Init() {
	for i := range m.items {
		m.items[i].key.Init()
		m.items[i].value.Init()
	}
	m.up = key.NewBinding(standardUpKey, key.WithHelp("↑", "up"))
	m.down = key.NewBinding(standardDownKey, key.WithHelp("↓", "down"))
	m.add = key.NewBinding(key.WithKeys("a"), key.WithHelp("a", "add"))
	m.remove = key.NewBinding(key.WithKeys("r"), key.WithHelp("r", "remove"))
	m.edit = key.NewBinding(key.WithKeys("e"), key.WithHelp("e", "edit"))
	m.clear = key.NewBinding(key.WithKeys("c"), key.WithHelp("c", "clear map"))
	m.keyEdit = key.NewBinding(key.WithKeys("ctrl+a"), key.WithHelp("ctrl+a", "edit key"))
	m.rebind()
}
func (m *MapNode) Unselect() {}
func (m *MapNode) rebind() {
	for i := range m.items {
		i := i
		m.items[i].value.SetSuffix(secondaryTextStyle.Render(","))
		m.items[i].key.OnUp(func() { m.selected = i + 1 })
		m.items[i].value.OnUp(func() { m.selected = i + 1 })
		m.items[i].value.OnDown(func() {
			if i == len(m.items)-1 {
				m.selected = -1
				m.GoDown()
			} else {
				m.selected = i + 2
			}
		})
		m.items[i].key.OnDown(func() {
			if i == len(m.items)-1 {
				m.selected = -1
				m.GoDown()
			} else {
				m.selected = i + 2
			}
		})
	}
}
func (m *MapNode) KeyHandled() bool {
	return true
}
func (m *MapNode) Children() []keyProvider {
	for i := range m.items {
		if m.items[i].key.Selected() != 0 {
			return []keyProvider{m.items[i].key}
		}
		if m.items[i].value.Selected() != 0 {
			return []keyProvider{m.items[i].value}
		}
	}
	return []keyProvider{}
}
func (m *MapNode) FooterKeys() []key.Binding {
	if m.selected < 0 {
		return []key.Binding{m.up, m.down}
	}
	return []key.Binding{m.add, m.remove, m.edit, m.keyEdit, m.up, m.down}
}
func (m *MapNode) FullKeyHelp() FullKeyHelp {
	return FullKeyHelp{Title: "Map editing", Keys: [][]key.Binding{{m.up, m.down}, {m.add, m.remove, m.edit, m.keyEdit, m.clear}}}
}
func (m *MapNode) Request() any {
	out := map[string]any{}
	for _, e := range m.items {
		out[fmt.Sprint(e.key.Request())] = e.value.Request()
	}
	return out
}
func (m *MapNode) Height() int {
	n := 2
	for _, e := range m.items {
		n += e.value.Height()
	}
	return n
}
func (m *MapNode) Selected() int {
	if m.selected == 0 {
		return 1
	}
	n := 1
	for i, e := range m.items {
		if m.selected == i+1 {
			return n
		}
		if e.key.Selected() != 0 {
			return n + e.key.Selected()
		}
		if e.value.Selected() != 0 {
			return n + e.value.Selected()
		}
		n += e.value.Height()
	}
	return 0
}
func (m *MapNode) SelectFromTop()    { m.selected = 0 }
func (m *MapNode) SelectFromBottom() { m.selected = len(m.items) }
func (m *MapNode) Update(msg tea.Msg) tea.Cmd {
	if k, ok := msg.(tea.KeyPressMsg); ok {
		switch {
		case m.selected >= 0 && key.Matches(k, m.add):
			e := MapEntry{createNode(m.keyType, m.registry), createNode(m.valueType, m.registry)}
			e.key.Init()
			e.value.Init()
			m.items = append(m.items[:m.selected], append([]MapEntry{e}, m.items[m.selected:]...)...)
			m.selected++
			m.rebind()
			return nil
		case m.selected >= 1 && key.Matches(k, m.remove):
			i := m.selected - 1
			m.items = append(m.items[:i], m.items[i+1:]...)
			if m.selected > len(m.items) {
				m.selected--
			}
			m.rebind()
			return nil
		case m.selected == 0 && key.Matches(k, m.clear):
			m.items = nil
			return nil
		case m.selected >= 1 && key.Matches(k, m.edit):
			m.items[m.selected-1].value.SelectFromTop()
			m.selected = -1
			return nil
		case m.selected >= 1 && key.Matches(k, m.keyEdit):
			m.items[m.selected-1].key.SelectFromTop()
			m.selected = -1
			return nil
		case m.selected >= 0 && key.Matches(k, m.up):
			if m.selected == 0 {
				m.selected = -1
				m.GoUp()
			} else {
				m.selected--
			}
			return nil
		case m.selected >= 0 && key.Matches(k, m.down):
			if m.selected == len(m.items) {
				m.selected = -1
				m.GoDown()
			} else {
				m.selected++
			}
			return nil
		}
	}
	for _, e := range m.items {
		if e.key.Selected() != 0 {
			return e.key.Update(msg)
		}
		if e.value.Selected() != 0 {
			return e.value.Update(msg)
		}
	}
	return nil
}
func (m *MapNode) View() (*tea.Cursor, string) {
	var b strings.Builder
	var c *tea.Cursor
	b.WriteString(textStyle.Render("{") + "\n")
	if m.selected == 0 {
		b.WriteString(" " + highlightStyle.Render("/* a to add, c to clear */") + "\n")
	}
	for i, e := range m.items {
		sel := e.key.Selected() != 0 || e.value.Selected() != 0
		if m.selected == i+1 {
			b.WriteString(structChildStyle.Render(highlightStyle.Render("/* a to add, e to edit, r to remove, ctrl+a to edit key */")) + "\n")
		}
		kc, kv := e.key.View()
		vc, vv := e.value.View()
		if sel {
			c = vc
			if e.key.Selected() != 0 {
				c = kc
			}
		}
		b.WriteString(structChildStyle.Render(kv+" "+textStyle.Render(":")+" "+vv) + "\n")
	}
	b.WriteString(textStyle.Render("}"))
	return c, b.String() + m.basicNode.Suffix
}
