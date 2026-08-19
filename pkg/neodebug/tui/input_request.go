package tui

import (
	"strings"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"github.com/Liphium/neoroute/neoschema"
	"github.com/Liphium/neoroute/pkg/neodebug/model"
)

// Idea:
// - Different nodes for the different types we have in Neo, some can share the same thing
// - WantedHeight() for getting the height (should add height of all things together basically)
// - Selected state is handled in the nodes themselves
// - All keys are handled by the currently selected node
// - The View() function takes in indicies that represent the global position of the node (in Y) -> only render what's needed
// - The selected global position can be gotten by going through the node tree
//
// Unclear:
// - Do we already send when the user presses enter? Like is that something we handle here and then just send what we get?
// 	- Might be bad for UX, but at the same time, there is no reason to press enter on any objects

type SendMsg struct {
	Cancelled bool
	Route     string
	Value     any
}

var _ keyProvider = inputRequestCreator{}

type inputRequestCreator struct {
	height   int
	width    int
	route    string
	schema   neoschema.PackedType
	rootNode SchemaNode

	// keys
	send key.Binding
	back key.Binding
}

// Children implements [keyProvider].
func (m inputRequestCreator) Children() []keyProvider {
	return []keyProvider{m.rootNode}
}

// FooterKeys implements [keyProvider].
func (m inputRequestCreator) FooterKeys() []key.Binding {
	return []key.Binding{m.send, m.back}
}

// FullKeyHelp implements [keyProvider].
func (m inputRequestCreator) FullKeyHelp() FullKeyHelp {
	return FullKeyHelp{
		Title: "Request editor",
		Keys: [][]key.Binding{
			[]key.Binding{m.send, m.back},
		},
	}
}

func newInputRequestCreator(route string, schema neoschema.PackedType, width, height int) inputRequestCreator {
	root := createNode(schema, schema.ObjectRegistry())
	root.Init()

	// Set the basic functions for up and down to select the node from the top / the bottom again
	root.OnUp(func() {
		root.SelectFromTop()
	})
	root.OnDown(func() {
		root.SelectFromBottom()
	})
	root.SelectFromTop()

	return inputRequestCreator{
		width:    width,
		height:   height,
		route:    route,
		schema:   schema,
		rootNode: root,

		// Define default keys
		send: key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "send request")),
		back: key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "cancel")),
	}
}

func (m *inputRequestCreator) SetWidth(w int) {
	m.width = w
}

func (m *inputRequestCreator) SetHeight(h int) {
	m.height = h
}

func (m inputRequestCreator) WantedHeight() int {
	return min(m.rootNode.Height(), 7) // TODO: Calculate the global height here (when smaller, also make this smaller)
}

func (m inputRequestCreator) Update(msg tea.Msg) (inputRequestCreator, tea.Cmd) {

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, m.send):

			return m, model.Plain(SendMsg{
				Route: m.route,
				Value: m.rootNode.Request(),
			})
		case key.Matches(msg, m.back):
			return m, model.Plain(SendMsg{Cancelled: true})
		}
	}

	// Update the rootNode (handles all of the things anyway from here)
	cmd := m.rootNode.Update(msg)
	return m, cmd
}

func (m inputRequestCreator) View() (*tea.Cursor, string) {
	// Manage the scrolling here (always centered around main selected thingy)

	maxShown := m.height
	selected := m.rootNode.Selected()
	n := m.rootNode.Height()

	var start, end int
	if n <= maxShown {
		start, end = 0, n // all visible
	} else {
		half := maxShown / 2
		start = selected - half
		start = max(0, min(start, n-maxShown)) // clamp → pin start/end
		end = start + maxShown                 // inclusive, window size = maxShown
	}

	// Get the view and adjust it to what's actually visible
	cursor, v := m.rootNode.View()
	lines := strings.Split(v, "\n")
	v = strings.Join(lines[start:end], "\n")
	if cursor != nil {
		cursor.Y += selected - start - 1
	}

	return cursor, v
}
