package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Public so caller can modify via commands (except Type/Expanded managed internally).
type Section struct {
	ID       string
	T        sectionType
	Loading  bool
	Label    string
	Data     []content
	Expanded bool
}

// Commands / Msgs

type AddSectionMsg struct{ Section Section }
type UpdateSectionMsg struct {
	ID      string
	Label   *string
	Loading *bool
	Data    *[]content
}
type ToggleSnapMsg struct{}

func AddSection(s Section) tea.Cmd {
	return func() tea.Msg { return AddSectionMsg{Section: s} }
}
func UpdateSection(id string, label *string, loading *bool, data *[]content) tea.Cmd {
	return func() tea.Msg { return UpdateSectionMsg{ID: id, Label: label, Loading: loading, Data: data} }
}

var (
	labelStyle    = lipgloss.NewStyle().Bold(true)
	selectedStyle = lipgloss.NewStyle().Background(lipgloss.Color("62")).Foreground(lipgloss.Color("230")).Bold(true)
	normalStyle   = lipgloss.NewStyle()
)

type History struct {
	snapToBottom   bool // true = snapped to bottom, selection hidden, auto-scroll
	sections       []Section
	selected       int
	width          int
	height         int
	viewport       viewport.Model
	spinner        spinner.Model
	following      bool // true = at bottom, auto-scroll on new sections
	manualExpanded map[string]bool
}

func NewHistory(w, h int) History {
	sp := spinner.New()
	sp.Spinner = spinner.Dot
	vp := viewport.New(w, h)
	return History{
		sections:       []Section{},
		selected:       0,
		width:          w,
		height:         h,
		viewport:       vp,
		spinner:        sp,
		snapToBottom:   true,
		manualExpanded: map[string]bool{},
	}
}

func (h History) Init() tea.Cmd { return h.spinner.Tick }

func (h History) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case spinner.TickMsg:
		var cmd tea.Cmd
		h.spinner, cmd = h.spinner.Update(msg)
		cmds = append(cmds, cmd)
		hasLoading := false
		for _, s := range h.sections {
			if s.Loading {
				hasLoading = true
				break
			}
		}
		if hasLoading {
			if h.snapToBottom {
				h.viewport.SetContent(h.renderContent())
				h.viewport.GotoBottom()
			} else {
				y := h.viewport.YOffset
				h.viewport.SetContent(h.renderContent())
				h.viewport.SetYOffset(y)
			}
		}

	case tea.WindowSizeMsg:
		h.width = msg.Width
		h.height = msg.Height
		h.viewport.Width = msg.Width
		h.viewport.Height = msg.Height
		h.viewport.SetContent(h.renderContent())

	case AddSectionMsg:
		// pin viewport so selected line stays at same screen row when not following
		prevOffset := -1
		if !h.snapToBottom && len(h.sections) > 0 {
			prevOffset = 0
			for i := 0; i < h.selected && i < len(h.sections); i++ {
				prevOffset++
				if h.sections[i].Expanded {
					prevOffset += len(h.sections[i].Data)
				}
			}
		}
		s := msg.Section
		s.Expanded = true
		h.sections = append(h.sections, s)
		for i := range h.sections {
			isNewest3 := i >= len(h.sections)-3
			if isNewest3 {
				h.sections[i].Expanded = true
			} else if !h.manualExpanded[h.sections[i].ID] {
				h.sections[i].Expanded = false
			}
		}
		if h.snapToBottom {
			h.selected = len(h.sections) - 1
			h.refreshViewport()
		} else {
			newOffset := 0
			for i := 0; i < h.selected && i < len(h.sections); i++ {
				newOffset++
				if h.sections[i].Expanded {
					newOffset += len(h.sections[i].Data)
				}
			}
			delta := newOffset - prevOffset
			h.viewport.SetContent(h.renderContent())
			h.viewport.SetYOffset(h.viewport.YOffset + delta)
			h.ensureVisible()
		}

	case UpdateSectionMsg:
		for i, s := range h.sections {
			if s.ID == msg.ID {
				if msg.Label != nil {
					h.sections[i].Label = *msg.Label
				}
				if msg.Loading != nil {
					h.sections[i].Loading = *msg.Loading
				}
				if msg.Data != nil {
					h.sections[i].Data = *msg.Data
				}
				break
			}
		}
		h.refreshViewport()

	case ToggleSnapMsg:
		h.toggleSnap()

	case tea.KeyMsg:
		switch msg.String() {
		case "v":
			h.toggleSnap()
		case "up", "k":
			if h.snapToBottom {
				h.snapToBottom = false
			}
			if h.selected > 0 {
				h.selected--
				h.refreshViewport()
			}
		case "down", "j":
			if h.snapToBottom {
				// already snapped → ignore
			} else if h.selected < len(h.sections)-1 {
				h.selected++
				if h.selected == len(h.sections)-1 {
					h.snapToBottom = true
				}
				h.refreshViewport()
			}
		case "left", "h":
			if h.selected >= 0 && h.selected < len(h.sections) {
				h.sections[h.selected].Expanded = false
				delete(h.manualExpanded, h.sections[h.selected].ID)
				h.refreshViewport()
			}
		case "right", "l":
			if h.selected >= 0 && h.selected < len(h.sections) {
				h.sections[h.selected].Expanded = true
				h.manualExpanded[h.sections[h.selected].ID] = true
				h.refreshViewport()
			}
		case "b":
			h.snapToBottom = true
			if len(h.sections) > 0 {
				h.selected = len(h.sections) - 1
			}
			h.refreshViewport()
		}
	}

	// let viewport handle scroll (mouse etc) - detect if user scrolled away
	oldOffset := h.viewport.YOffset
	var cmd tea.Cmd
	h.viewport, cmd = h.viewport.Update(msg)
	cmds = append(cmds, cmd)
	if h.viewport.YOffset != oldOffset {
		h.snapToBottom = h.viewport.AtBottom()
	}

	return h, tea.Batch(cmds...)
}

func (h *History) toggleSnap() {
	if h.snapToBottom {
		h.snapToBottom = false
	} else {
		h.snapToBottom = true
		if len(h.sections) > 0 {
			h.selected = len(h.sections) - 1
		}
	}
	h.refreshViewport()
}

func (h *History) refreshViewport() {
	if h.snapToBottom {
		h.viewport.SetContent(h.renderContent())
		h.viewport.GotoBottom()
		return
	}
	y := h.viewport.YOffset
	h.viewport.SetContent(h.renderContent())
	h.viewport.SetYOffset(y)
	h.ensureVisible()
}

func (h *History) ensureVisible() {
	offset := 0
	for i, s := range h.sections {
		if i == h.selected {
			break
		}
		offset++
		if s.Expanded {
			offset += len(s.Data)
		}
	}
	if offset < h.viewport.YOffset {
		h.viewport.SetYOffset(offset)
	} else if offset >= h.viewport.YOffset+h.viewport.Height {
		h.viewport.SetYOffset(offset - h.viewport.Height + 1)
	}
}

func (h History) renderContent() string {
	var b strings.Builder
	for i, s := range h.sections {
		arrow := "▸"
		if s.Expanded {
			arrow = "▾"
		}
		spinnerStr := ""
		if s.Loading {
			spinnerStr = h.spinner.View()
		}
		line := fmt.Sprintf("%s %s%s", arrow, spinnerStr, s.Label)
		var rendered string
		if !h.snapToBottom && i == h.selected {
			rendered = selectedStyle.Render(line)
		} else {
			rendered = labelStyle.Render(line)
		}
		b.WriteString(rendered)
		b.WriteString("\n")
		if s.Expanded {
			for _, c := range s.Data {
				// only label highlighted, data always normal
				b.WriteString(normalStyle.Render("  "+c.Line) + "\n")
			}
		}
	}
	return strings.TrimSuffix(b.String(), "\n")
}

func (h History) View() string { return h.viewport.View() }
