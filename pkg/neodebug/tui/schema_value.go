package tui

import (
	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
)

var _ SchemaNode = &ValueNode[any]{}

type valueNodeKeyMap struct {
	Up    key.Binding
	Down  key.Binding
	Clear key.Binding
}

func defaultValueNodeKeyMap() valueNodeKeyMap {
	return valueNodeKeyMap{
		Up:    key.NewBinding(key.WithKeys("up", "k"), key.WithHelp("↑/k", "up")),
		Down:  key.NewBinding(key.WithKeys("down", "j"), key.WithHelp("↓/j", "down")),
		Clear: key.NewBinding(key.WithKeys("ctrl+backspace"), key.WithHelp("ctrl+backspace", "clear input")),
	}
}

type ValueNode[T any] struct {
	*basicSelection
	convert func(s string) (T, error)
	value   T
	prefix  string
	suffix  string
	input   textinput.Model
	keyMap  valueNodeKeyMap
}

// Init implements [SchemaNode].
func (v *ValueNode[T]) Init() {
	v.input = textinput.New()
	v.input.SetVirtualCursor(false)
	v.input.Prompt = v.prefix
	v.input.Placeholder = ""
	styles := textinput.DefaultDarkStyles()
	styles.Blurred.Prompt, styles.Focused.Prompt = secondaryTextStyle, secondaryTextStyle
	styles.Blurred.Text, styles.Focused.Text = secondaryTextStyle, textStyle
	v.input.SetStyles(styles)
	v.keyMap = defaultValueNodeKeyMap()
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
func (v *ValueNode[T]) Update(msg tea.Msg) (SchemaNode, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, v.keyMap.Up):
			v.input.Blur()
			v.GoUp()
			return v, nil

		case key.Matches(msg, v.keyMap.Down):
			v.input.Blur()
			v.GoDown()
			return v, nil

		case key.Matches(msg, v.keyMap.Clear):
			v.input.SetValue("")
			return v, nil
		}
	}

	if v.input.Focused() {
		var cmd tea.Cmd
		v.input, cmd = v.input.Update(msg)
		return v, cmd
	}

	return v, nil
}

// View implements [SchemaNode].
func (v *ValueNode[T]) View() (*tea.Cursor, string) {
	textView := secondaryTextStyle.Render(v.prefix + v.input.Value())
	if v.input.Focused() {
		textView = v.input.View()
	}

	view := textView + secondaryTextStyle.Render(v.suffix)
	if v.input.Err != nil {
		view += "\n" + errorStyle.Render(v.input.Err.Error())
	}

	return v.input.Cursor(), view
}
