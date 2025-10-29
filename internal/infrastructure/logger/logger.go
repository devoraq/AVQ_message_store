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
	l *log.Logger
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

	var attrs []slog.Attr
	r.Attrs(func(a slog.Attr) bool {
		attrs = append(attrs, a)
		return true
	})

	if len(attrs) == 0 {
		h.l.Println(b.String())
		return nil
	}

	b.WriteString(" {\n")

	for _, a := range attrs {
		valueLines := formatAttrValue(a.Value.Any())
		for i, line := range valueLines {
			if i == 0 {
				b.WriteString("  ")
				b.WriteString(color.CyanString(a.Key))
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

// WithAttrs возвращает обработчик без изменений, поскольку PrettyHandler не кэширует атрибуты.
func (h *PrettyHandler) WithAttrs(_ []slog.Attr) slog.Handler {
	return h
}

// WithGroup игнорирует группировку и возвращает исходный обработчик.
func (h *PrettyHandler) WithGroup(_ string) slog.Handler {
	return h
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
