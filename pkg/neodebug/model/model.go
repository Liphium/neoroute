package model

import (
	"fmt"
	"time"

	tea "charm.land/bubbletea/v2"
)

// Shared messages for avoidance of circular dependencies

type Content interface {
	Creation() time.Time
}

// Message for adding content to the TUI's history
type AddContentMsg struct{ Content }

type BasicContent struct {
	CreatedAt time.Time
}

func Now() BasicContent {
	return BasicContent{time.Now()}
}

func CreatedAt(stamp time.Time) BasicContent {
	return BasicContent{stamp}
}

func (b BasicContent) Creation() time.Time {
	return b.CreatedAt
}

type InfoContent struct {
	Message string
	BasicContent
}

func Info(format string, a ...any) tea.Msg {
	return AddContentMsg{
		Content: InfoContent{
			Message:      fmt.Sprintf(format, a...),
			BasicContent: Now(),
		},
	}
}

type ErrorContent struct {
	Message string
	BasicContent
}

func Error(format string, a ...any) tea.Msg {
	return AddContentMsg{
		Content: ErrorContent{
			Message:      fmt.Sprintf(format, a...),
			BasicContent: Now(),
		},
	}
}

type SendingContent struct {
	Route string
	Data  any
	BasicContent
}

func Sending(route string, data any) tea.Msg {
	return AddContentMsg{
		Content: SendingContent{
			Route:        route,
			Data:         data,
			BasicContent: Now(),
		},
	}
}

type EventContent struct {
	Name  string
	Event any
	BasicContent
}

func Event(name string, event any) tea.Msg {
	return AddContentMsg{
		Content: EventContent{
			Name:         name,
			Event:        event,
			BasicContent: Now(),
		},
	}
}

type ResponseContent struct {
	Route string
	Data  any
	BasicContent
}

func Response(route string, data any) tea.Msg {
	return AddContentMsg{
		Content: ResponseContent{
			Route:        route,
			Data:         data,
			BasicContent: Now(),
		},
	}
}
