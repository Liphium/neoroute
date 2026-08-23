package tui

import (
	"fmt"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
)

// TODO: Use up and down keybinds without vim during editing

var _ SchemaNode = &ValueNode[any]{}

type ValueNode[T any] struct {
	basicNode
	convert    func(s string) (T, error)
	value      T
	focused    bool
	keyHandled bool
	prefix     string
	// This suffix here is different from the one in basicNode in that it applies to the input field instead of to what's after the entire node (also after the error label)
	suffix string
	input  textinput.Model

	// Keys
	up     key.Binding
	down   key.Binding
	clear  key.Binding
	insert key.Binding
	finish key.Binding
}

// Request implements SchemaNode.
func (v *ValueNode[T]) Request() any {
	return v.value
}

// Children implements SchemaNode.
func (v *ValueNode[T]) Children() []keyProvider {
	return []keyProvider{}
}

// FooterKeys implements SchemaNode.
func (v *ValueNode[T]) FooterKeys() []key.Binding {
	if v.input.Focused() {
		return []key.Binding{v.finish, v.up, v.down}
	}
	return []key.Binding{v.insert, v.up, v.down}
}

// KeyHandled implements [SchemaNode].
func (v *ValueNode[T]) KeyHandled() bool {
	return v.keyHandled
}

// FullKeyHelp implements SchemaNode.
func (v *ValueNode[T]) FullKeyHelp() FullKeyHelp {
	return FullKeyHelp{
		Title: "Currently selected value",
		Keys: [][]key.Binding{
			v.FooterKeys(),
			[]key.Binding{v.clear},
		},
	}
}

// Init implements [SchemaNode].
func (v *ValueNode[T]) Init() {
	v.input = textinput.New()
	v.input.SetVirtualCursor(false)
	v.input.Prompt = v.prefix
	v.input.Placeholder = ""
	v.input.SetValue(fmt.Sprintf("%v", v.value))

	styles := textinput.DefaultDarkStyles()
	styles.Blurred.Prompt, styles.Focused.Prompt = highlightStyle.Bold(true), highlightStyle.Bold(true)
	styles.Blurred.Text, styles.Focused.Text = highlightStyle, highlightStyle
	v.input.SetStyles(styles)

	// Define all of the keys
	v.up = key.NewBinding(standardUpKey, key.WithHelp("↑", "up"))
	v.down = key.NewBinding(standardDownKey, key.WithHelp("↓", "down"))
	v.clear = key.NewBinding(key.WithKeys("c"), key.WithHelp("c", "clear input"))
	v.insert = key.NewBinding(key.WithKeys("i", "enter"), key.WithHelp("i", "insert value"))
	v.finish = key.NewBinding(key.WithKeys("esc", "enter"), key.WithHelp("esc", "finish"))

	v.input.Validate = func(s string) error {
		var err error
		v.value, err = v.convert(s)
		return err
	}
}

// Height implements [SchemaNode].
func (v *ValueNode[T]) Height() int {
	if v.input.Err != nil {
		return 2
	} else {
		return 1
	}
}

// SelectFromTop implements [SchemaNode].
func (v *ValueNode[T]) SelectFromTop() {
	v.focused = true
}

// SelectFromBottom implements [SchemaNode].
func (v *ValueNode[T]) SelectFromBottom() {
	v.focused = true
}

// Unselect implements [SchemaNode].
func (v *ValueNode[T]) Unselect() {
	v.focused = false
	v.input.Blur()
}

// Selected implements [SchemaNode].
func (v *ValueNode[T]) Selected() int {
	if v.focused {
		return 1
	}
	return 0
}

// Update implements [SchemaNode].
func (v *ValueNode[T]) Update(msg tea.Msg) tea.Cmd {
	v.keyHandled = false
	if !v.focused {
		return nil
	}

	// Handle priority keybinds (up, down, finish)
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, v.up):
			v.keyHandled = true
			v.focused = false
			v.input.Blur()
			v.GoUp()
			return nil

		case key.Matches(msg, v.down):
			v.keyHandled = true
			v.focused = false
			v.input.Blur()
			v.GoDown()
			return nil

		case v.input.Focused() && key.Matches(msg, v.finish):
			v.keyHandled = true
			v.input.Blur()
			return nil
		}
	}

	if !v.input.Focused() {
		// Handle keybinds during normal focus
		switch msg := msg.(type) {
		case tea.KeyPressMsg:
			switch {
			case key.Matches(msg, v.insert):
				v.keyHandled = true
				v.input.Focus()
				return nil

			case key.Matches(msg, v.clear):
				v.keyHandled = true
				v.input.SetValue("")
				return nil
			}
		}
	} else {
		// Let the input take over when we're in insert mode
		v.keyHandled = true
		var cmd tea.Cmd
		v.input, cmd = v.input.Update(msg)
		return cmd
	}

	return nil
}

// View implements [SchemaNode].
func (v *ValueNode[T]) View() (*tea.Cursor, string) {
	style := secondaryTextStyle
	if v.focused {
		style = highlightStyle.Bold(true)
	}

	textView := secondaryTextStyle.Render(v.prefix + v.input.Value())
	if v.input.Focused() {
		textView = v.input.View()
	} else if v.focused {
		textView = style.Render(v.prefix) + highlightStyle.Render(v.input.Value())
		if v.input.Value() == "" && v.prefix == "" && v.suffix == "" {
			textView = highlightStyle.Render("/* no value */")
		}
	}

	view := textView + style.Render(v.suffix) + v.basicNode.Suffix
	if v.input.Err != nil {
		view += "\n" + errorStyle.Render(v.input.Err.Error())
	}

	return v.input.Cursor(), view
}
