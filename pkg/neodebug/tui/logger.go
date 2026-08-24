package tui

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/Liphium/neoroute"
	"github.com/Liphium/neoroute/client"
	"github.com/Liphium/neoroute/pkg/neodebug/model"
)

// InitLogger installs a slog logger that routes to the TUI history via program.Send.
// - Ignores Debug and Info
// - Error -> model.ErrorContent (ERR, red)
// - Warn  -> model.WarnContent (WRN, yellow)
// Does NOT write to os.Stderr/Stdout so the alt-screen TUI isn't broken.
func InitLogger(program *tea.Program) {
	h := &historyHandler{program: program}
	logger := slog.New(h)
	neoroute.SetLogger(logger)
	client.SetLogger(logger)
}

type historyHandler struct {
	program *tea.Program
	attrs   []slog.Attr
	groups  []string
}

func (h *historyHandler) Enabled(_ context.Context, level slog.Level) bool {
	// Ignore Debug and Info, only Warn and above.
	return level >= slog.LevelWarn
}

func (h *historyHandler) Handle(_ context.Context, r slog.Record) error {
	if h.program == nil {
		return nil
	}

	var sb strings.Builder
	sb.WriteString(r.Message)

	// Collect handler-level attrs
	for _, a := range h.attrs {
		sb.WriteString(fmt.Sprintf(" %s=%v", a.Key, a.Value))
	}
	// Collect record attrs
	r.Attrs(func(a slog.Attr) bool {
		// Resolve value
		v := a.Value.Resolve()
		key := a.Key
		if len(h.groups) > 0 {
			key = strings.Join(append(h.groups, key), ".")
		}
		sb.WriteString(fmt.Sprintf(" %s=%v", key, v))
		return true
	})

	msg := sb.String()

	switch {
	case r.Level >= slog.LevelError:
		h.program.Send(model.Error("%s", msg))
	case r.Level >= slog.LevelWarn:
		h.program.Send(model.Warn("%s", msg))
	default:
		// Should not happen due to Enabled, but ignore.
	}
	return nil
}

func (h *historyHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	newAttrs := make([]slog.Attr, 0, len(h.attrs)+len(attrs))
	newAttrs = append(newAttrs, h.attrs...)
	newAttrs = append(newAttrs, attrs...)
	return &historyHandler{program: h.program, attrs: newAttrs, groups: h.groups}
}

func (h *historyHandler) WithGroup(name string) slog.Handler {
	newGroups := append(append([]string(nil), h.groups...), name)
	return &historyHandler{program: h.program, attrs: h.attrs, groups: newGroups}
}
