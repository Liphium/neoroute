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

type WarnContent struct {
	Message string
	BasicContent
}

func Warn(format string, a ...any) tea.Msg {
	return AddContentMsg{
		Content: WarnContent{
			Message:      fmt.Sprintf(format, a...),
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
	Route    string
	Request  any
	Response any
	BasicContent
}

func Response(route string, request any, response any) tea.Msg {
	return AddContentMsg{
		Content: ResponseContent{
			Route:        route,
			Request:      request,
			Response:     response,
			BasicContent: Now(),
		},
	}
}
