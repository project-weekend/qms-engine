package config

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"runtime"
	"strings"

	"github.com/project-weekend/qms-engine/server/config"
)

func NewLogger(appCfg *config.Config) *slog.Logger {
	level := mapLogLevel(appCfg.Logger.LogLevel)
	opts := &slog.HandlerOptions{
		Level: level,
	}

	var handler slog.Handler
	// Read format from config and create appropriate handler
	format := appCfg.Logger.LogFormat
	switch {
	case format == "json":
		handler = slog.NewJSONHandler(os.Stdout, opts)
	case strings.HasPrefix(format, "cgls"):
		handler = NewCGLSHandler(os.Stdout, opts)
	default:
		handler = slog.NewTextHandler(os.Stdout, opts)
	}

	return slog.New(handler)
}

func mapLogLevel(logLevel int) slog.Level {
	switch logLevel {
	case 0, 1, 2, 3, 4, 5:
		return slog.LevelDebug
	case 6:
		return slog.LevelInfo
	case 7, 8:
		return slog.LevelWarn
	default:
		return slog.LevelError
	}
}

// CGLSHandler is a custom slog handler that formats logs in CGLS format
type CGLSHandler struct {
	opts  *slog.HandlerOptions
	w     io.Writer
	attrs []slog.Attr
	group string
}

// NewCGLSHandler creates a new CGLS format handler
func NewCGLSHandler(w io.Writer, opts *slog.HandlerOptions) *CGLSHandler {
	if opts == nil {
		opts = &slog.HandlerOptions{}
	}
	return &CGLSHandler{
		opts: opts,
		w:    w,
	}
}

// Enabled reports whether the handler handles records at the given level
func (h *CGLSHandler) Enabled(_ context.Context, level slog.Level) bool {
	minLevel := slog.LevelInfo
	if h.opts.Level != nil {
		minLevel = h.opts.Level.Level()
	}
	return level >= minLevel
}

// Handle formats and writes a log record
func (h *CGLSHandler) Handle(_ context.Context, r slog.Record) error {
	buf := make([]byte, 0, 1024)

	// Format: [YYYY-MM-DD HH:MM:SS.mmm] [LEVEL] [file:line] message key=value key=value

	// Timestamp
	timestamp := r.Time.Format("2006-01-02 15:04:05.000")
	buf = fmt.Appendf(buf, "[%s] ", timestamp)

	// Level with color coding (optional)
	levelStr := h.formatLevel(r.Level)
	buf = fmt.Appendf(buf, "[%s] ", levelStr)

	// Source location
	if h.opts.AddSource {
		fs := runtime.CallersFrames([]uintptr{r.PC})
		f, _ := fs.Next()
		if f.File != "" {
			// Get just the filename, not the full path
			file := f.File
			if idx := strings.LastIndex(file, "/"); idx >= 0 {
				file = file[idx+1:]
			}
			buf = fmt.Appendf(buf, "[%s:%d] ", file, f.Line)
		}
	}

	// Message
	buf = fmt.Appendf(buf, "%s", r.Message)

	// Add handler's attributes
	for _, attr := range h.attrs {
		buf = h.appendAttr(buf, attr)
	}

	// Add record's attributes
	r.Attrs(func(attr slog.Attr) bool {
		buf = h.appendAttr(buf, attr)
		return true
	})

	buf = append(buf, '\n')
	_, err := h.w.Write(buf)
	return err
}

// formatLevel returns a formatted level string
func (h *CGLSHandler) formatLevel(level slog.Level) string {
	switch level {
	case slog.LevelDebug:
		return "DEBUG"
	case slog.LevelInfo:
		return "INFO "
	case slog.LevelWarn:
		return "WARN "
	case slog.LevelError:
		return "ERROR"
	default:
		return level.String()
	}
}

// appendAttr appends a formatted attribute to the buffer
func (h *CGLSHandler) appendAttr(buf []byte, attr slog.Attr) []byte {
	if attr.Equal(slog.Attr{}) {
		return buf
	}

	key := attr.Key
	if h.group != "" {
		key = h.group + "." + key
	}

	buf = append(buf, ' ')
	buf = append(buf, key...)
	buf = append(buf, '=')

	value := attr.Value.String()
	// Quote the value if it contains spaces
	if strings.Contains(value, " ") {
		buf = fmt.Appendf(buf, "%q", value)
	} else {
		buf = append(buf, value...)
	}

	return buf
}

// WithAttrs returns a new handler with additional attributes
func (h *CGLSHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	newHandler := &CGLSHandler{
		opts:  h.opts,
		w:     h.w,
		attrs: make([]slog.Attr, 0, len(h.attrs)+len(attrs)),
		group: h.group,
	}
	newHandler.attrs = append(newHandler.attrs, h.attrs...)
	newHandler.attrs = append(newHandler.attrs, attrs...)
	return newHandler
}

// WithGroup returns a new handler with a group name
func (h *CGLSHandler) WithGroup(name string) slog.Handler {
	if name == "" {
		return h
	}
	newGroup := name
	if h.group != "" {
		newGroup = h.group + "." + name
	}
	return &CGLSHandler{
		opts:  h.opts,
		w:     h.w,
		attrs: h.attrs,
		group: newGroup,
	}
}
