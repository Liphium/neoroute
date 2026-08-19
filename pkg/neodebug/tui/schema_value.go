package tui

import (
	"fmt"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
)

var _ SchemaNode = &ValueNode[any]{}

type ValueNode[T any] struct {
	basicNode
	convert func(s string) (T, error)
	value   T
	prefix  string
	// This suffix here is different from the one in basicNode in that it applies to the input field instead of to what's after the entire node (also after the error label)
	suffix string
	input  textinput.Model

	// Keys
	up    key.Binding
	down  key.Binding
	clear key.Binding
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
	return []key.Binding{v.up, v.down}
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
	styles.Blurred.Prompt, styles.Focused.Prompt = secondaryTextStyle, secondaryTextStyle
	styles.Blurred.Text, styles.Focused.Text = secondaryTextStyle, textStyle
	v.input.SetStyles(styles)

	// Define all of the keys
	v.up = key.NewBinding(standardUpKey, key.WithHelp("↑", "up"))
	v.down = key.NewBinding(standardDownKey, key.WithHelp("↓", "down"))
	v.clear = key.NewBinding(key.WithKeys("ctrl+backspace"), key.WithHelp("ctrl+backspace", "clear input"))

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

// SelectFromBottom implements [SchemaNode].
func (v *ValueNode[T]) SelectFromBottom() {
	v.input.Focus()
}

// SelectFromTop implements [SchemaNode].
func (v *ValueNode[T]) SelectFromTop() {
	v.input.Focus()
}

// Selected implements [SchemaNode].
func (v *ValueNode[T]) Selected() int {
	if v.input.Focused() {
		return 1
	}
	return 0
}

// Update implements [SchemaNode].
func (v *ValueNode[T]) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, v.up):
			v.input.Blur()
			v.GoUp()
			return nil

		case key.Matches(msg, v.down):
			v.input.Blur()
			v.GoDown()
			return nil

		case key.Matches(msg, v.clear):
			v.input.SetValue("")
			return nil
		}
	}

	if v.input.Focused() {
		var cmd tea.Cmd
		v.input, cmd = v.input.Update(msg)
		return cmd
	}

	return nil
}

// View implements [SchemaNode].
func (v *ValueNode[T]) View() (*tea.Cursor, string) {
	textView := secondaryTextStyle.Render(v.prefix + v.input.Value())
	if v.input.Focused() {
		textView = v.input.View()
	}

	view := textView + secondaryTextStyle.Render(v.suffix) + v.basicNode.Suffix
	if v.input.Err != nil {
		view += "\n" + errorStyle.Render(v.input.Err.Error())
	}

	return v.input.Cursor(), view
}
