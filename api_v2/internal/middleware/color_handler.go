package middleware

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"runtime"
	"strconv"
	"strings"
	"sync"
)

type ColorHandler struct {
	opts   *slog.HandlerOptions
	writer io.Writer
	mu     sync.Mutex
}

func NewColorHandler(w io.Writer, opts *slog.HandlerOptions) *ColorHandler {
	if opts == nil {
		opts = &slog.HandlerOptions{}
	}
	return &ColorHandler{
		opts:   opts,
		writer: w,
	}
}

func (h *ColorHandler) Enabled(_ context.Context, level slog.Level) bool {
	minLevel := slog.LevelInfo
	if h.opts.Level != nil {
		minLevel = h.opts.Level.Level()
	}
	return level >= minLevel
}

func (h *ColorHandler) Handle(ctx context.Context, r slog.Record) error {
	buf := make([]byte, 0, 1024)

	buf = append(buf, []byte(Gray)...)
	buf = append(buf, []byte(r.Time.Format("15:04:05"))...)
	buf = append(buf, []byte(Reset+" ")...)

	buf = append(buf, h.formatLevel(r.Level)...)
	buf = append(buf, []byte(" ")...)

	if h.opts.AddSource && r.PC != 0 {
		fs := runtime.CallersFrames([]uintptr{r.PC})
		f, _ := fs.Next()
		if f.File != "" {
			buf = append(buf, []byte(fmt.Sprintf("%s%s:%d%s ",
				Purple, shortenFile(f.File), f.Line, Reset))...)
		}
	}

	buf = append(buf, []byte(Bold)...)
	buf = append(buf, []byte(r.Message)...)
	buf = append(buf, []byte(Reset)...)

	if r.NumAttrs() > 0 {
		buf = append(buf, []byte(" ")...)
		first := true
		r.Attrs(func(a slog.Attr) bool {
			if !first {
				buf = append(buf, []byte(" ")...)
			}
			first = false
			buf = append(buf, h.formatAttr(a)...)
			return true
		})
	}

	buf = append(buf, '\n')

	h.mu.Lock()
	defer h.mu.Unlock()
	_, err := h.writer.Write(buf)
	return err
}

func (h *ColorHandler) WithAttrs(attrs []slog.Attr) slog.Handler {

	return h
}

func (h *ColorHandler) WithGroup(name string) slog.Handler {

	return h
}

func (h *ColorHandler) formatLevel(level slog.Level) string {
	switch level {
	case slog.LevelDebug:
		return fmt.Sprintf("%s%s%-5s%s", Bold, Gray, "DEBUG", Reset)
	case slog.LevelInfo:
		return fmt.Sprintf("%s%s%-5s%s", Bold, Green, "INFO", Reset)
	case slog.LevelWarn:
		return fmt.Sprintf("%s%s%-5s%s", Bold, Yellow, "WARN", Reset)
	case slog.LevelError:
		return fmt.Sprintf("%s%s%-5s%s", Bold, Red, "ERROR", Reset)
	default:
		return fmt.Sprintf("%s%-5s%s", Bold, level.String(), Reset)
	}
}

func (h *ColorHandler) formatAttr(a slog.Attr) []byte {
	key := a.Key
	value := a.Value

	var buf []byte

	buf = append(buf, []byte(Cyan)...)
	buf = append(buf, []byte(key)...)
	buf = append(buf, []byte(Reset+"=")...)

	switch value.Kind() {
	case slog.KindString:
		buf = append(buf, []byte(fmt.Sprintf("%q", value.String()))...)
	case slog.KindInt64:
		buf = append(buf, []byte(Yellow)...)
		buf = append(buf, []byte(strconv.FormatInt(value.Int64(), 10))...)
		buf = append(buf, []byte(Reset)...)
	case slog.KindFloat64:
		buf = append(buf, []byte(Yellow)...)
		buf = append(buf, []byte(strconv.FormatFloat(value.Float64(), 'f', -1, 64))...)
		buf = append(buf, []byte(Reset)...)
	case slog.KindBool:
		buf = append(buf, []byte(Purple)...)
		if value.Bool() {
			buf = append(buf, []byte("true")...)
		} else {
			buf = append(buf, []byte("false")...)
		}
		buf = append(buf, []byte(Reset)...)
	default:
		buf = append(buf, []byte(value.String())...)
	}

	return buf
}

func shortenFile(file string) string {
	parts := strings.Split(file, "/")
	if len(parts) > 2 {
		return strings.Join(parts[len(parts)-2:], "/")
	}
	return file
}
