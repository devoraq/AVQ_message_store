// Package logger предоставляет адаптер, который выводит структурированные логи slog
// в читабельном текстовом виде.
package logger

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"log/slog"
	"strings"

	"github.com/fatih/color"
)

// PrettyHandlerOptions описывает настройки PrettyHandler и оборачивает slog.HandlerOptions.
type PrettyHandlerOptions struct {
	Opts slog.HandlerOptions
}

// PrettyHandler реализует интерфейс slog.Handler и печатает записи в компактном виде.
type PrettyHandler struct {
	l      *log.Logger
	attrs  []slog.Attr
	groups []string
}

// Enabled всегда возвращает true, предоставляя slog.Logger решать, писать ли запись.
func (h *PrettyHandler) Enabled(_ context.Context, _ slog.Level) bool {
	return true
}

// Handle форматирует запись slog с подсветкой уровня и печатает её.
func (h *PrettyHandler) Handle(_ context.Context, r slog.Record) error {
	var b strings.Builder
	b.Grow(256)

	level := r.Level.String() + ":"
	switch r.Level {
	case slog.LevelDebug:
		level = color.MagentaString(level)
	case slog.LevelInfo:
		level = color.GreenString(level)
	case slog.LevelWarn:
		level = color.YellowString(level)
	case slog.LevelError:
		level = color.RedString(level)
	}

	b.WriteString(r.Time.Format("[15:05:05.000] "))
	b.WriteString(level)
	b.WriteByte(' ')
	b.WriteString(color.WhiteString(r.Message))

	prefix := strings.Join(h.groups, ".")
	pairs := make([]attrPair, 0, len(h.attrs)+r.NumAttrs())

	for _, a := range h.attrs {
		flattenAttr(prefix, a, &pairs)
	}

	r.Attrs(func(a slog.Attr) bool {
		flattenAttr(prefix, a, &pairs)
		return true
	})

	if len(pairs) == 0 {
		h.l.Println(b.String())
		return nil
	}

	b.WriteString(" {\n")

	for _, a := range pairs {
		valueLines := formatAttrValue(a.value)
		for i, line := range valueLines {
			if i == 0 {
				b.WriteString("  ")
				b.WriteString(color.CyanString(a.key))
				b.WriteString(": ")
			} else {
				b.WriteString("    ")
			}
			b.WriteString(line)
			b.WriteByte('\n')
		}
	}

	b.WriteString("}")

	h.l.Println(b.String())
	return nil
}

// WithAttrs возвращает новый обработчик с добавленными атрибутами.
func (h *PrettyHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	if len(attrs) == 0 {
		return h
	}
	clone := h.clone()
	for _, a := range attrs {
		clone.attrs = append(clone.attrs, resolveAttr(a))
	}
	return clone
}

// WithGroup возвращает новый обработчик, который учитывает указанную группу.
func (h *PrettyHandler) WithGroup(name string) slog.Handler {
	if name == "" {
		return h
	}
	clone := h.clone()
	clone.groups = append(clone.groups, name)
	return clone
}

// NewPrettyHandler создаёт обработчик, который печатает структурированные логи slog
// в дружественном для человека формате с использованием стандартного log.Logger.
func NewPrettyHandler(
	out io.Writer,
	opts PrettyHandlerOptions,
) *PrettyHandler {
	_ = opts
	h := &PrettyHandler{
		l: log.New(out, "", 0),
	}

	return h
}

func (h *PrettyHandler) clone() *PrettyHandler {
	attrs := make([]slog.Attr, len(h.attrs))
	copy(attrs, h.attrs)
	groups := make([]string, len(h.groups))
	copy(groups, h.groups)
	return &PrettyHandler{
		l:      h.l,
		attrs:  attrs,
		groups: groups,
	}
}

type attrPair struct {
	key   string
	value any
}

func flattenAttr(prefix string, attr slog.Attr, out *[]attrPair) {
	value := attr.Value.Resolve()

	if value.Kind() == slog.KindGroup {
		nextPrefix := prefix
		if attr.Key != "" {
			if nextPrefix != "" {
				nextPrefix += "."
			}
			nextPrefix += attr.Key
		}
		for _, child := range value.Group() {
			flattenAttr(nextPrefix, child, out)
		}
		return
	}

	key := attr.Key
	if prefix != "" && key != "" {
		key = prefix + "." + key
	} else if prefix != "" && key == "" {
		key = prefix
	}
	if key == "" {
		key = "<root>"
	}
	*out = append(*out, attrPair{
		key:   key,
		value: value.Any(),
	})
}

func resolveAttr(a slog.Attr) slog.Attr {
	return slog.Attr{
		Key:   a.Key,
		Value: a.Value.Resolve(),
	}
}

func formatAttrValue(v any) []string {
	var text string
	switch val := v.(type) {
	case error:
		return formatErrorChain(val)
	default:
		text = fmt.Sprint(val)
	}
	text = strings.TrimSpace(text)
	if text == "" {
		return []string{""}
	}
	parts := strings.Split(text, "\n")
	for i := range parts {
		parts[i] = color.WhiteString(parts[i])
	}
	return parts
}

func formatErrorChain(err error) []string {
	if err == nil {
		return []string{""}
	}

	var lines []string
	current := err
	level := 0
	for current != nil {
		prefix := ""
		if level > 0 {
			prefix = strings.Repeat("  ", level-1) + color.YellowString("? ")
		}
		lines = append(lines, prefix+color.WhiteString(current.Error()))
		current = errors.Unwrap(current)
		level++
	}
	return lines
}
